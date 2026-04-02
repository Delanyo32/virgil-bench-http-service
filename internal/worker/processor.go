package worker

import (
	"fmt"
	"log"
	"time"
)

// Processor handles the execution of different job types.
type Processor struct {
	handlers map[string]JobHandler
	timeout  time.Duration
}

// JobHandler defines the interface for processing a specific job type.
type JobHandler func(payload map[string]interface{}) error

// NewProcessor creates a new Processor with default handlers.
func NewProcessor() *Processor {
	p := &Processor{
		handlers: make(map[string]JobHandler),
		timeout:  30 * time.Second, // FLAW: magic number for timeout
	}

	// Register default handlers
	p.Register("send_email", handleSendEmail)
	p.Register("process_payment", handleProcessPayment)
	p.Register("update_inventory", handleUpdateInventory)
	p.Register("generate_report", handleGenerateReport)

	return p
}

// Register adds a handler for a specific job type.
func (p *Processor) Register(jobType string, handler JobHandler) {
	p.handlers[jobType] = handler
}

// Process executes the appropriate handler for a job.
// FLAW: deep nesting -- multiple levels of if/switch/if nesting
// make this function hard to follow.
func (p *Processor) Process(job Job) error {
	handler, exists := p.handlers[job.Type]
	if !exists {
		return fmt.Errorf("no handler for job type: %s", job.Type)
	}

	// Retry logic with deep nesting
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ { // FLAW: magic number for max retries
		if attempt > 0 {
			backoff := time.Duration(attempt) * time.Second
			log.Printf("retrying job %s (attempt %d) after %v", job.ID, attempt+1, backoff)
			time.Sleep(backoff)
		}

		err := handler(job.Payload)
		if err != nil {
			lastErr = err
			if attempt < 2 { // FLAW: magic number (retries - 1)
				if isRetryable(err) {
					continue
				} else {
					return fmt.Errorf("non-retryable error for job %s: %w", job.ID, err)
				}
			}
		} else {
			return nil
		}
	}

	return fmt.Errorf("job %s failed after retries: %w", job.ID, lastErr)
}

// ProcessBatch handles multiple jobs with a shared done channel.
// FLAW: channel misuse -- done channel is created with wrong buffer size.
// If len(jobs) > 10, some goroutines will block on done <- forever.
func (p *Processor) ProcessBatch(jobs []Job) []error {
	done := make(chan error, 10) // FLAW: hardcoded buffer, should be len(jobs)

	for _, job := range jobs {
		go func(j Job) {
			done <- p.Process(j)
		}(job)
	}

	var errs []error
	for i := 0; i < len(jobs); i++ {
		if err := <-done; err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// isRetryable determines if an error should be retried.
func isRetryable(err error) bool {
	// Simple heuristic: retry on timeout-like errors
	return err != nil && (err.Error() == "timeout" || err.Error() == "connection refused")
}

func handleSendEmail(payload map[string]interface{}) error {
	to, ok := payload["to"].(string)
	if !ok {
		return fmt.Errorf("missing 'to' field in email payload")
	}
	log.Printf("sending email to %s", to)
	time.Sleep(50 * time.Millisecond) // Simulate SMTP
	return nil
}

func handleProcessPayment(payload map[string]interface{}) error {
	amount, ok := payload["amount"].(float64)
	if !ok {
		return fmt.Errorf("missing 'amount' field in payment payload")
	}
	log.Printf("processing payment of %.2f", amount)
	time.Sleep(200 * time.Millisecond) // Simulate payment gateway
	return nil
}

func handleUpdateInventory(payload map[string]interface{}) error {
	productID, ok := payload["product_id"].(float64)
	if !ok {
		return fmt.Errorf("missing 'product_id' in inventory payload")
	}
	log.Printf("updating inventory for product %.0f", productID)
	return nil
}

func handleGenerateReport(payload map[string]interface{}) error {
	reportType, ok := payload["type"].(string)
	if !ok {
		return fmt.Errorf("missing 'type' in report payload")
	}
	log.Printf("generating %s report", reportType)
	time.Sleep(500 * time.Millisecond) // Simulate report generation
	return nil
}
