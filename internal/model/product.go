package model

import (
	"fmt"
	"time"
)

// Product represents an inventory item.
type Product struct {
	ID          int       `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceCents  int       `json:"price_cents"`
	Quantity    int       `json:"quantity"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsInStock returns true if the product has available quantity.
func (p *Product) IsInStock() bool {
	return p.Quantity > 0
}

// PriceFormatted returns the price as a formatted dollar string.
func (p *Product) PriceFormatted() string {
	dollars := p.PriceCents / 100
	cents := p.PriceCents % 100
	return formatPrice(dollars, cents)
}

// formatPrice formats a dollar and cents value as a string.
func formatPrice(dollars, cents int) string {
	if cents < 10 {
		return fmt.Sprintf("$%d.0%d", dollars, cents)
	}
	return fmt.Sprintf("$%d.%d", dollars, cents)
}

// ensure imports are used
var _ = fmt.Sprintf
