package repository

import (
	"fmt"
	"sync"

	"github.com/example/ordersvc/internal/model"
)

// CategoryRepository provides in-memory CRUD for product categories.
type CategoryRepository struct {
	mu         sync.RWMutex
	categories map[int]model.Category
	nextID     int
}

// NewCategoryRepository creates an empty CategoryRepository.
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		categories: make(map[int]model.Category),
		nextID:     1,
	}
}

// Create stores a new category and assigns an ID.
func (r *CategoryRepository) Create(cat *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cat.ID = r.nextID
	r.nextID++
	r.categories[cat.ID] = *cat
	return nil
}

// GetByID returns a category by its ID.
func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cat, ok := r.categories[id]
	if !ok {
		return nil, fmt.Errorf("category %d not found", id)
	}
	return &cat, nil
}

// ListAll returns all stored categories.
func (r *CategoryRepository) ListAll() ([]model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Category, 0, len(r.categories))
	for _, cat := range r.categories {
		result = append(result, cat)
	}
	return result, nil
}

// Delete removes a category by ID.
func (r *CategoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.categories[id]; !ok {
		return fmt.Errorf("category %d not found", id)
	}
	delete(r.categories, id)
	return nil
}

// ListByParent returns all categories with the given parent ID.
func (r *CategoryRepository) ListByParent(parentID int) ([]model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Category
	for _, cat := range r.categories {
		if cat.ParentID == parentID {
			result = append(result, cat)
		}
	}
	return result, nil
}
