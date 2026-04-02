package model

import (
	"time"
)

// Category represents a product category with optional parent for hierarchy.
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	ParentID    int       `json:"parent_id,omitempty"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// IsRoot returns true if the category has no parent.
func (c *Category) IsRoot() bool {
	return c.ParentID == 0
}

// IsChildOf returns true if this category's parent matches the given ID.
func (c *Category) IsChildOf(parentID int) bool {
	return c.ParentID == parentID
}

// BuildBreadcrumb builds a breadcrumb path from a flat list of categories.
// It walks from the current category up through its parents.
func BuildBreadcrumb(cat Category, all []Category) []string {
	crumbs := []string{cat.Name}
	current := cat

	for current.ParentID != 0 {
		found := false
		for _, c := range all {
			if c.ID == current.ParentID {
				crumbs = append([]string{c.Name}, crumbs...)
				current = c
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	return crumbs
}
