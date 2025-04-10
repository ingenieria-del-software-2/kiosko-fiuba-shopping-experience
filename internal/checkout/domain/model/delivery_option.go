package model

import (
	"github.com/google/uuid"
)

// DeliveryOption represents a value object for delivery options in the Checkout Process bounded context
type DeliveryOption struct {
	ShippingAddressID uuid.UUID `json:"shippingAddressId"`
	ShippingMethodID  uuid.UUID `json:"shippingMethodId"`
}

// NewDeliveryOption creates a new delivery option
func NewDeliveryOption(shippingAddressID, shippingMethodID uuid.UUID) *DeliveryOption {
	return &DeliveryOption{
		ShippingAddressID: shippingAddressID,
		ShippingMethodID:  shippingMethodID,
	}
}
