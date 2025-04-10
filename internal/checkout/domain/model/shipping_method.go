package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ShippingMethod represents a shipping method entity in the Checkout Process bounded context
type ShippingMethod struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	Price                 float64   `json:"price"`
	EstimatedDeliveryDays int       `json:"estimatedDeliveryDays"`
}

// NewShippingMethod creates a new shipping method
func NewShippingMethod(name, description string, price float64, estimatedDeliveryDays int) (*ShippingMethod, error) {
	// Validate required fields
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if estimatedDeliveryDays <= 0 {
		return nil, errors.New("estimated delivery days must be positive")
	}

	return &ShippingMethod{
		ID:                    uuid.New(),
		Name:                  name,
		Description:           description,
		Price:                 price,
		EstimatedDeliveryDays: estimatedDeliveryDays,
	}, nil
}

// Update updates the shipping method details
func (m *ShippingMethod) Update(name, description string, price float64, estimatedDeliveryDays int) error {
	// Validate required fields
	if name == "" {
		return errors.New("name is required")
	}
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	if estimatedDeliveryDays <= 0 {
		return errors.New("estimated delivery days must be positive")
	}

	m.Name = name
	m.Description = description
	m.Price = price
	m.EstimatedDeliveryDays = estimatedDeliveryDays

	return nil
}

// DisplayName returns a formatted display name with delivery estimate
func (m *ShippingMethod) DisplayName() string {
	return m.Name + " (" + m.DeliveryEstimate() + ")"
}

// DeliveryEstimate returns a human-readable delivery estimate
func (m *ShippingMethod) DeliveryEstimate() string {
	if m.EstimatedDeliveryDays == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", m.EstimatedDeliveryDays)
}
