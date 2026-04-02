package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/example/ordersvc/internal/model"
)

// OrderRepository handles order persistence.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create inserts a new order and its items into the database.
func (r *OrderRepository) Create(order *model.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = tx.QueryRow(
		`INSERT INTO orders (user_id, status, total_cents, shipping_address, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		order.UserID, order.Status, order.TotalCents, order.ShippingAddress, order.CreatedAt,
	).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	for _, item := range order.Items {
		_, err := tx.Exec(
			`INSERT INTO order_items (order_id, product_id, quantity, price_cents)
			 VALUES ($1, $2, $3, $4)`,
			order.ID, item.ProductID, item.Quantity, item.PriceCents,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an order by its ID.
func (r *OrderRepository) FindByID(id int) (*model.Order, error) {
	order := &model.Order{}
	err := r.db.QueryRow(
		`SELECT id, user_id, status, total_cents, shipping_address, created_at, updated_at
		 FROM orders WHERE id = $1`,
		id,
	).Scan(
		&order.ID, &order.UserID, &order.Status, &order.TotalCents,
		&order.ShippingAddress, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// FindByUserID returns orders for a user with pagination.
func (r *OrderRepository) FindByUserID(userID, limit, offset int) ([]model.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, status, total_cents, shipping_address, created_at, updated_at
		 FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Status, &o.TotalCents,
			&o.ShippingAddress, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// FindItemsByOrderID returns all items for an order.
func (r *OrderRepository) FindItemsByOrderID(orderID int) ([]model.OrderItem, error) {
	rows, err := r.db.Query(
		`SELECT id, order_id, product_id, quantity, price_cents
		 FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.PriceCents,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// UpdateStatus changes the status of an order.
func (r *OrderRepository) UpdateStatus(orderID int, status string) error {
	_, err := r.db.Exec(
		`UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), orderID,
	)
	return err
}

// FindPendingBefore returns orders with "pending" status created before the cutoff.
func (r *OrderRepository) FindPendingBefore(cutoff time.Time) ([]model.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, status, total_cents, shipping_address, created_at, updated_at
		 FROM orders WHERE status = 'pending' AND created_at < $1`,
		cutoff,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Status, &o.TotalCents,
			&o.ShippingAddress, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
