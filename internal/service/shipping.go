package service

import (
	"fmt"
	"strings"

	"github.com/example/ordersvc/internal/model"
)

// ShippingRate holds the calculated shipping cost.
type ShippingRate struct {
	Method    string `json:"method"`
	CostCents int    `json:"cost_cents"`
	EstDays   int    `json:"est_days"`
}

// ShippingService calculates shipping costs and validates addresses.
type ShippingService struct {
	baseCostCents      int
	perItemCostCents   int
	internationalExtra int
}

// NewShippingService creates a ShippingService with configured rates.
func NewShippingService(baseCost, perItem, intlExtra int) *ShippingService {
	return &ShippingService{
		baseCostCents:      baseCost,
		perItemCostCents:   perItem,
		internationalExtra: intlExtra,
	}
}

// CalculateRate computes the shipping cost for an order.
func (s *ShippingService) CalculateRate(addr model.Address, itemCount int) (ShippingRate, error) {
	if err := addr.Validate(); err != nil {
		return ShippingRate{}, fmt.Errorf("invalid shipping address: %w", err)
	}
	if itemCount <= 0 {
		return ShippingRate{}, fmt.Errorf("item count must be positive")
	}

	cost := s.baseCostCents + (s.perItemCostCents * itemCount)
	estDays := 5
	method := "standard"

	if !addr.IsDomestic() {
		cost += s.internationalExtra
		estDays = 14
		method = "international"
	}

	return ShippingRate{
		Method:    method,
		CostCents: cost,
		EstDays:   estDays,
	}, nil
}

// ValidateAddress checks that a shipping address is complete and has
// a supported country code.
func (s *ShippingService) ValidateAddress(addr model.Address) error {
	if err := addr.Validate(); err != nil {
		return err
	}

	country := strings.ToUpper(strings.TrimSpace(addr.Country))
	supported := []string{"US", "CA", "GB", "DE", "FR", "AU", "JP"}
	for _, c := range supported {
		if country == c {
			return nil
		}
	}
	return fmt.Errorf("unsupported shipping country: %s", addr.Country)
}

// FormatLabel returns a formatted shipping label string.
func (s *ShippingService) FormatLabel(addr model.Address) string {
	return addr.FormatOneLine()
}
