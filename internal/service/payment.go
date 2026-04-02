package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// PaymentResult represents the outcome of a payment attempt.
type PaymentResult struct {
	TransactionID string
	Success       bool
	ErrorMessage  string
}

// PaymentService handles payment processing via external gateway.
type PaymentService struct {
	gatewayURL string
	client     *http.Client
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(gatewayURL string) *PaymentService {
	return &PaymentService{
		gatewayURL: gatewayURL,
		client: &http.Client{
			Timeout: 30 * time.Second, // FLAW: magic number for timeout
		},
	}
}

// ProcessPayment sends a payment request to the external gateway.
// FLAW: sync-blocking-in-async -- this function makes a blocking HTTP call
// and is called from the request handler goroutine. If the gateway is slow,
// the goroutine blocks for up to 30 seconds, consuming resources.
func (s *PaymentService) ProcessPayment(orderID, amountCents int) (*PaymentResult, error) {
	payload := fmt.Sprintf(`{"order_id": %d, "amount_cents": %d}`, orderID, amountCents)

	// FLAW: sync blocking call in request handler path
	resp, err := s.client.Post(
		s.gatewayURL+"/v1/charge",
		"application/json",
		strings.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("payment request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payment response: %w", err)
	}

	var result struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
		Error         string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse payment response: %w", err)
	}

	return &PaymentResult{
		TransactionID: result.TransactionID,
		Success:       result.Status == "success",
		ErrorMessage:  result.Error,
	}, nil
}

// RefundPayment issues a refund for a previous payment.
// FLAW: sync-blocking-in-async -- same blocking pattern as ProcessPayment.
func (s *PaymentService) RefundPayment(transactionID string, amountCents int) error {
	payload := fmt.Sprintf(`{"transaction_id": "%s", "amount_cents": %d}`,
		transactionID, amountCents)

	resp, err := s.client.Post(
		s.gatewayURL+"/v1/refund",
		"application/json",
		strings.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("refund request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("refund failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ValidatePaymentMethod checks if a payment method is valid.
// Dead code -- not called from any service or handler.
func (s *PaymentService) ValidatePaymentMethod(methodID string) (bool, error) {
	resp, err := s.client.Get(s.gatewayURL + "/v1/methods/" + methodID)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// retryPayment attempts a payment with exponential backoff.
// Dead code -- not used anywhere.
func retryPayment(fn func() error, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			backoff := time.Duration(1<<uint(i)) * time.Second
			log.Printf("payment attempt %d failed, retrying in %v: %v", i+1, backoff, err)
			time.Sleep(backoff)
			continue
		}
		return nil
	}
	return fmt.Errorf("payment failed after %d retries: %w", maxRetries, lastErr)
}
