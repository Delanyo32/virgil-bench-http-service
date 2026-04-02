package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/example/ordersvc/internal/model"
	"github.com/example/ordersvc/internal/repository"
)

// OrderService handles order business logic.
type OrderService struct {
	orderRepo     *repository.OrderRepository
	inventoryRepo *repository.InventoryRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(
	orderRepo *repository.OrderRepository,
	inventoryRepo *repository.InventoryRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CreateOrder persists a new order and reserves inventory.
// FLAW: N+1 query pattern -- fetches each product price individually
// in a loop instead of batch-fetching all prices at once.
func (s *OrderService) CreateOrder(order *model.Order) error {
	// Calculate total by fetching each product individually
	var totalCents int
	for _, item := range order.Items {
		// FLAW: N+1 query -- each iteration hits the database
		product, err := s.inventoryRepo.GetProduct(item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}
		item.PriceCents = product.PriceCents
		totalCents += product.PriceCents * item.Quantity
	}
	order.TotalCents = totalCents

	// Reserve stock for each item individually
	for _, item := range order.Items {
		// FLAW: N+1 query -- another loop of individual DB calls
		err := s.inventoryRepo.DecrementStock(item.ProductID, item.Quantity)
		if err != nil {
			// No rollback of previously decremented items
			return fmt.Errorf("failed to reserve stock for product %d: %w", item.ProductID, err)
		}
	}

	// Persist order
	if err := s.orderRepo.Create(order); err != nil {
		return fmt.Errorf("failed to persist order: %w", err)
	}

	return nil
}

// GetOrder retrieves an order by ID.
func (s *OrderService) GetOrder(id int) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// FLAW: N+1 query -- fetches order items and then each product separately
	items, err := s.orderRepo.FindItemsByOrderID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load order items: %w", err)
	}
	for i := range items {
		product, err := s.inventoryRepo.GetProduct(items[i].ProductID)
		if err != nil {
			log.Printf("failed to load product %d for order item: %v", items[i].ProductID, err)
			continue
		}
		items[i].ProductName = product.Name
	}
	order.Items = items

	return order, nil
}

// ListOrders returns paginated orders for a user.
func (s *OrderService) ListOrders(userID, page, limit int) ([]model.Order, error) {
	offset := (page - 1) * limit
	orders, err := s.orderRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// FLAW: N+1 query -- for each order, fetch its items separately
	for i := range orders {
		items, err := s.orderRepo.FindItemsByOrderID(orders[i].ID)
		if err != nil {
			log.Printf("failed to load items for order %d: %v", orders[i].ID, err)
			continue
		}
		orders[i].Items = items
	}

	return orders, nil
}

// UpdateStatus changes the order status string.
func (s *OrderService) UpdateStatus(orderID int, status string) error {
	return s.orderRepo.UpdateStatus(orderID, status)
}

// CancelOrder marks an order as cancelled and restores inventory.
func (s *OrderService) CancelOrder(orderID int) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	if order.Status != "pending" && order.Status != "confirmed" {
		return fmt.Errorf("cannot cancel order in status: %s", order.Status)
	}

	items, err := s.orderRepo.FindItemsByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to load order items: %w", err)
	}

	// Restore stock for each item
	for _, item := range items {
		err := s.inventoryRepo.IncrementStock(item.ProductID, item.Quantity)
		if err != nil {
			_ = err // FLAW: error swallowed -- partial stock restore on failure
		}
	}

	return s.orderRepo.UpdateStatus(orderID, "cancelled")
}

// GetOrderStats returns basic statistics about orders.
func (s *OrderService) GetOrderStats(userID int) (map[string]interface{}, error) {
	orders, err := s.orderRepo.FindByUserID(userID, 10000, 0) // FLAW: magic number, no real pagination
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_orders":     len(orders),
		"total_spent":      0,
		"pending_count":    0,
		"confirmed_count":  0,
		"cancelled_count":  0,
	}

	for _, o := range orders {
		stats["total_spent"] = stats["total_spent"].(int) + o.TotalCents
		switch o.Status {
		case "pending":
			stats["pending_count"] = stats["pending_count"].(int) + 1
		case "confirmed":
			stats["confirmed_count"] = stats["confirmed_count"].(int) + 1
		case "cancelled":
			stats["cancelled_count"] = stats["cancelled_count"].(int) + 1
		}
	}

	return stats, nil
}

// ProcessExpiredOrders cancels orders older than the given duration.
// Dead code -- not called from any handler or worker.
func (s *OrderService) ProcessExpiredOrders(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)
	orders, err := s.orderRepo.FindPendingBefore(cutoff)
	if err != nil {
		return err
	}
	for _, order := range orders {
		_ = s.CancelOrder(order.ID)
	}
	return nil
}

// unused helper -- this is dead code
func calculateDiscount(totalCents int, userTier string) int {
	switch userTier {
	case "gold":
		return totalCents * 10 / 100
	case "silver":
		return totalCents * 5 / 100
	default:
		return 0
	}
}

// Ensure sql import is used
var _ = sql.ErrNoRows
