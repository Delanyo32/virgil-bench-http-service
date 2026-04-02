package model

import (
	"time"
)

// Order represents a customer order.
type Order struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	Status          string      `json:"status"`
	TotalCents      int         `json:"total_cents"`
	ShippingAddress string      `json:"shipping_address"`
	Items           []OrderItem `json:"items,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem represents a single line item in an order.
type OrderItem struct {
	ID          int    `json:"id"`
	OrderID     int    `json:"order_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
	Quantity    int    `json:"quantity"`
	PriceCents  int    `json:"price_cents"`
}

// OrderStatus constants for order lifecycle.
const (
	OrderStatusPending       = "pending"
	OrderStatusConfirmed     = "confirmed"
	OrderStatusProcessing    = "processing"
	OrderStatusShipped       = "shipped"
	OrderStatusDelivered     = "delivered"
	OrderStatusCancelled     = "cancelled"
	OrderStatusPaymentFailed = "payment_failed"
)

// IsTerminal returns true if the order is in a final state.
func (o *Order) IsTerminal() bool {
	return o.Status == OrderStatusDelivered ||
		o.Status == OrderStatusCancelled ||
		o.Status == OrderStatusPaymentFailed
}

// ItemCount returns the total number of items in the order.
func (o *Order) ItemCount() int {
	total := 0
	for _, item := range o.Items {
		total += item.Quantity
	}
	return total
}
