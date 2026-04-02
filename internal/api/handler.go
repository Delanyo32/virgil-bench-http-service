package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/example/ordersvc/internal/model"
	"github.com/example/ordersvc/internal/service"
)

// Handler holds references to all service dependencies.
type Handler struct {
	orderSvc        *service.OrderService
	inventorySvc    *service.InventoryService
	paymentSvc      *service.PaymentService
	notificationSvc *service.NotificationService
}

// NewHandler creates a new Handler with all service dependencies.
func NewHandler(
	orderSvc *service.OrderService,
	inventorySvc *service.InventoryService,
	paymentSvc *service.PaymentService,
	notificationSvc *service.NotificationService,
) *Handler {
	return &Handler{
		orderSvc:        orderSvc,
		inventorySvc:    inventorySvc,
		paymentSvc:      paymentSvc,
		notificationSvc: notificationSvc,
	}
}

// CreateOrder handles the full order creation flow.
// FLAW: god function -- this handler does auth extraction, input parsing,
// validation, business logic orchestration, payment processing,
// notification sending, and response formatting all in one function.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Auth extraction (should be middleware)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		WriteError(w, http.StatusUnauthorized, "missing authorization header")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		WriteError(w, http.StatusUnauthorized, "invalid authorization format")
		return
	}
	token := parts[1]
	userID, err := validateToken(token)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	// Input parsing
	var req struct {
		Items []struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"items"`
		ShippingAddress string `json:"shipping_address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validation (should be separate)
	if len(req.Items) == 0 {
		WriteError(w, http.StatusBadRequest, "order must have at least one item")
		return
	}
	if req.ShippingAddress == "" {
		WriteError(w, http.StatusBadRequest, "shipping address is required")
		return
	}
	for _, item := range req.Items {
		if item.ProductID <= 0 {
			WriteError(w, http.StatusBadRequest, "invalid product id")
			return
		}
		if item.Quantity <= 0 || item.Quantity > 100 { // FLAW: magic number
			WriteError(w, http.StatusBadRequest, "invalid quantity")
			return
		}
	}

	// Stock check (should be service layer only)
	for _, item := range req.Items {
		available, err := h.inventorySvc.CheckStock(item.ProductID, item.Quantity)
		if err != nil {
			log.Printf("stock check failed for product %d: %v", item.ProductID, err)
			WriteError(w, http.StatusInternalServerError, "stock check failed")
			return
		}
		if !available {
			WriteError(w, http.StatusConflict,
				fmt.Sprintf("insufficient stock for product %d", item.ProductID))
			return
		}
	}

	// Build order model
	order := &model.Order{
		UserID:          userID,
		Status:          "pending", // FLAW: stringly-typed status
		ShippingAddress: req.ShippingAddress,
		CreatedAt:       time.Now(),
	}

	for _, item := range req.Items {
		order.Items = append(order.Items, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Create order
	if err := h.orderSvc.CreateOrder(order); err != nil {
		log.Printf("failed to create order: %v", err)
		WriteError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	// Process payment (should be async)
	paymentResult, err := h.paymentSvc.ProcessPayment(order.ID, order.TotalCents)
	if err != nil {
		log.Printf("payment failed for order %d: %v", order.ID, err)
		_ = h.orderSvc.UpdateStatus(order.ID, "payment_failed") // FLAW: stringly-typed
		WriteError(w, http.StatusPaymentRequired, "payment processing failed")
		return
	}

	// Update order status
	if paymentResult.Success {
		_ = h.orderSvc.UpdateStatus(order.ID, "confirmed") // FLAW: stringly-typed
	} else {
		_ = h.orderSvc.UpdateStatus(order.ID, "payment_failed") // FLAW: stringly-typed
	}

	// Send notification (should be async)
	_ = h.notificationSvc.SendOrderConfirmation(userID, order.ID) // FLAW: error swallowed

	duration := time.Since(startTime)
	log.Printf("order %d created in %v", order.ID, duration)

	// FLAW: duplicate response construction (same pattern in GetOrder, ListOrders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      order.ID,
		"status":  order.Status,
		"total":   order.TotalCents,
		"message": "order created successfully",
	})
}

// GetOrder retrieves a single order by ID.
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.orderSvc.GetOrder(id)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteError(w, http.StatusNotFound, "order not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "failed to get order")
		return
	}

	// FLAW: duplicate response construction
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      order.ID,
		"status":  order.Status,
		"total":   order.TotalCents,
		"items":   order.Items,
		"message": "success",
	})
}

// ListOrders returns orders for a user with basic pagination.
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit := 20 // FLAW: magic number for page size

	orders, err := h.orderSvc.ListOrders(userID, page, limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to list orders")
		return
	}

	// FLAW: duplicate response construction
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orders":  orders,
		"page":    page,
		"limit":   limit,
		"message": "success",
	})
}

// validateToken is a stub token validator.
func validateToken(token string) (int, error) {
	if len(token) < 10 {
		return 0, fmt.Errorf("token too short")
	}
	return 1, nil
}
