package service

import (
	"fmt"
	"strings"

	"github.com/example/ordersvc/internal/model"
	"github.com/example/ordersvc/internal/repository"
)

// CategoryService handles category management business logic.
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategory validates and stores a new category.
func (s *CategoryService) CreateCategory(name, slug, description string, parentID int) (*model.Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("category name is required")
	}
	if strings.TrimSpace(slug) == "" {
		return nil, fmt.Errorf("category slug is required")
	}

	cat := &model.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
		ParentID:    parentID,
	}

	if err := s.repo.Create(cat); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return cat, nil
}

// GetCategory returns a category by ID.
func (s *CategoryService) GetCategory(id int) (*model.Category, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return cat, nil
}

// ListRootCategories returns all top-level categories (no parent).
func (s *CategoryService) ListRootCategories() ([]model.Category, error) {
	return s.repo.ListByParent(0)
}

// ListChildren returns all child categories for the given parent.
func (s *CategoryService) ListChildren(parentID int) ([]model.Category, error) {
	return s.repo.ListByParent(parentID)
}

// DeleteCategory removes a category by ID.
func (s *CategoryService) DeleteCategory(id int) error {
	children, err := s.repo.ListByParent(id)
	if err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete category with %d children", len(children))
	}
	return s.repo.Delete(id)
}
