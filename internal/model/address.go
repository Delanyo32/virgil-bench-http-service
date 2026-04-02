package model

import (
	"fmt"
	"strings"
)

// Address represents a shipping or billing address.
type Address struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// Validate checks that required address fields are populated.
func (a *Address) Validate() error {
	if strings.TrimSpace(a.Line1) == "" {
		return fmt.Errorf("address line1 is required")
	}
	if strings.TrimSpace(a.City) == "" {
		return fmt.Errorf("address city is required")
	}
	if strings.TrimSpace(a.State) == "" {
		return fmt.Errorf("address state is required")
	}
	if strings.TrimSpace(a.PostalCode) == "" {
		return fmt.Errorf("address postal code is required")
	}
	if strings.TrimSpace(a.Country) == "" {
		return fmt.Errorf("address country is required")
	}
	return nil
}

// FormatOneLine returns the address as a single comma-separated string.
func (a *Address) FormatOneLine() string {
	parts := []string{a.Line1}
	if a.Line2 != "" {
		parts = append(parts, a.Line2)
	}
	parts = append(parts, a.City, a.State, a.PostalCode, a.Country)
	return strings.Join(parts, ", ")
}

// IsDomestic returns true if the country is "US".
func (a *Address) IsDomestic() bool {
	return strings.EqualFold(a.Country, "US")
}
