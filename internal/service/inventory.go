package service

import (
	"fmt"
	"log"

	"github.com/example/ordersvc/internal/model"
	"github.com/example/ordersvc/internal/repository"
)

// InventoryService handles inventory business logic.
type InventoryService struct {
	repo *repository.InventoryRepository
}

// NewInventoryService creates a new InventoryService.
func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// CheckStock verifies that the requested quantity is available.
func (s *InventoryService) CheckStock(productID, quantity int) (bool, error) {
	product, err := s.repo.GetProduct(productID)
	if err != nil {
		return false, fmt.Errorf("failed to get product: %w", err)
	}
	return product.Quantity >= quantity, nil
}

// UpdateStock sets the stock level for a product.
func (s *InventoryService) UpdateStock(productID, quantity int) error {
	return s.repo.UpdateStock(productID, quantity)
}

// ListProducts returns all products in inventory.
func (s *InventoryService) ListProducts() ([]model.Product, error) {
	return s.repo.ListAll()
}

// GetProduct returns a single product by ID.
func (s *InventoryService) GetProduct(productID int) (*model.Product, error) {
	return s.repo.GetProduct(productID)
}

// BulkUpdateStock updates stock for multiple products.
func (s *InventoryService) BulkUpdateStock(updates map[int]int) error {
	for productID, quantity := range updates {
		if err := s.repo.UpdateStock(productID, quantity); err != nil {
			log.Printf("failed to update stock for product %d: %v", productID, err)
			continue
		}
	}
	return nil
}

// SearchProducts searches for products by name or category.
func (s *InventoryService) SearchProducts(query string) ([]model.Product, error) {
	return s.repo.Search(query)
}

// GetLowStockProducts returns products with stock below threshold.
func (s *InventoryService) GetLowStockProducts(threshold int) ([]model.Product, error) {
	products, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	var lowStock []model.Product
	for _, p := range products {
		if p.Quantity < threshold {
			lowStock = append(lowStock, p)
		}
	}
	return lowStock, nil
}
