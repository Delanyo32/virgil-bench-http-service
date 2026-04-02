package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/ordersvc/internal/model"
)

// InventoryRepository handles product/inventory persistence.
type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository creates a new InventoryRepository.
func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// GetProduct retrieves a product by ID.
func (r *InventoryRepository) GetProduct(id int) (*model.Product, error) {
	p := &model.Product{}
	err := r.db.QueryRow(
		`SELECT id, sku, name, description, price_cents, quantity, category, created_at, updated_at
		 FROM products WHERE id = $1`,
		id,
	).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents,
		&p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return p, nil
}

// ListAll returns all products in the inventory.
func (r *InventoryRepository) ListAll() ([]model.Product, error) {
	rows, err := r.db.Query(
		`SELECT id, sku, name, description, price_cents, quantity, category, created_at, updated_at
		 FROM products ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents,
			&p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// UpdateStock sets the quantity for a product.
func (r *InventoryRepository) UpdateStock(productID, quantity int) error {
	result, err := r.db.Exec(
		`UPDATE products SET quantity = $1, updated_at = NOW() WHERE id = $2`,
		quantity, productID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product %d not found", productID)
	}
	return nil
}

// DecrementStock reduces the quantity of a product.
func (r *InventoryRepository) DecrementStock(productID, amount int) error {
	result, err := r.db.Exec(
		`UPDATE products SET quantity = quantity - $1, updated_at = NOW()
		 WHERE id = $2 AND quantity >= $1`,
		amount, productID,
	)
	if err != nil {
		return fmt.Errorf("failed to decrement stock: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}
	return nil
}

// IncrementStock increases the quantity of a product.
func (r *InventoryRepository) IncrementStock(productID, amount int) error {
	_, err := r.db.Exec(
		`UPDATE products SET quantity = quantity + $1, updated_at = NOW() WHERE id = $2`,
		amount, productID,
	)
	return err
}

// Search finds products matching a query by name or category.
func (r *InventoryRepository) Search(query string) ([]model.Product, error) {
	rows, err := r.db.Query(
		`SELECT id, sku, name, description, price_cents, quantity, category, created_at, updated_at
		 FROM products WHERE name ILIKE $1 OR category ILIKE $1`,
		"%"+query+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents,
			&p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
