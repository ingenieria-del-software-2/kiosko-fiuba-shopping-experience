package dto

import (
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/model"
)

// CheckoutItemDTO represents an item in a checkout
type CheckoutItemDTO struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
	ImageURL  string  `json:"imageUrl"`
}

// DeliveryOptionDTO represents shipping details for a checkout
type DeliveryOptionDTO struct {
	ShippingAddressID string `json:"shippingAddressId"`
	ShippingMethodID  string `json:"shippingMethodId"`
}

// PaymentMethodDTO represents payment method details
type PaymentMethodDTO struct {
	PaymentType    string                 `json:"paymentType"`
	PaymentDetails map[string]interface{} `json:"paymentDetails"`
}

// CheckoutResponseDTO represents checkout data for API responses
type CheckoutResponseDTO struct {
	ID           string             `json:"id"`
	CartID       string             `json:"cartId"`
	Status       string             `json:"status"`
	Items        []CheckoutItemDTO  `json:"items"`
	Subtotal     float64            `json:"subtotal"`
	ShippingCost float64            `json:"shippingCost"`
	Tax          float64            `json:"tax"`
	Total        float64            `json:"total"`
	Delivery     *DeliveryOptionDTO `json:"delivery,omitempty"`
	Payment      *PaymentMethodDTO  `json:"payment,omitempty"`
	CreatedAt    string             `json:"createdAt"`
	UpdatedAt    string             `json:"updatedAt"`
}

// CheckoutInitRequest represents the request to initialize a checkout
type CheckoutInitRequest struct {
	CartID string `json:"cartId" validate:"required,uuid"`
}

// ShippingDetailsRequest represents the request to update shipping details
type ShippingDetailsRequest struct {
	AddressID string `json:"addressId" validate:"required,uuid"`
	MethodID  string `json:"shippingMethodId" validate:"required,uuid"`
}

// PaymentMethodRequest represents the request to set a payment method
type PaymentMethodRequest struct {
	PaymentType    string                 `json:"paymentType" validate:"required"`
	PaymentDetails map[string]interface{} `json:"paymentDetails" validate:"required"`
}

// ShippingAddressDTO represents a shipping address for API responses
type ShippingAddressDTO struct {
	ID            string `json:"id"`
	UserID        string `json:"userId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	StreetAddress string `json:"streetAddress"`
	Apartment     string `json:"apartment,omitempty"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
	PhoneNumber   string `json:"phoneNumber"`
	IsDefault     bool   `json:"isDefault"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// ShippingAddressRequest represents the request to create or update a shipping address
type ShippingAddressRequest struct {
	UserID        string `json:"userId" validate:"required,uuid"`
	FirstName     string `json:"firstName" validate:"required"`
	LastName      string `json:"lastName" validate:"required"`
	StreetAddress string `json:"streetAddress" validate:"required"`
	Apartment     string `json:"apartment"`
	City          string `json:"city" validate:"required"`
	State         string `json:"state" validate:"required"`
	PostalCode    string `json:"postalCode" validate:"required"`
	Country       string `json:"country" validate:"required"`
	PhoneNumber   string `json:"phoneNumber" validate:"required"`
	IsDefault     bool   `json:"isDefault"`
}

// ShippingMethodDTO represents a shipping method for API responses
type ShippingMethodDTO struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	Price                 float64 `json:"price"`
	EstimatedDeliveryDays int     `json:"estimatedDeliveryDays"`
}

// FromDomain converts a checkout domain model to a DTO
func CheckoutFromDomain(checkout *model.Checkout) *CheckoutResponseDTO {
	items := make([]CheckoutItemDTO, len(checkout.Items))
	for i, item := range checkout.Items {
		items[i] = CheckoutItemDTO{
			ProductID: item.ProductID.String(),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			ImageURL:  item.ImageURL,
		}
	}

	result := &CheckoutResponseDTO{
		ID:           checkout.ID.String(),
		CartID:       checkout.CartID.String(),
		Status:       string(checkout.Status),
		Items:        items,
		Subtotal:     checkout.Subtotal,
		ShippingCost: checkout.ShippingCost,
		Tax:          checkout.Tax,
		Total:        checkout.Total,
		CreatedAt:    checkout.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    checkout.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if checkout.DeliveryOption != nil {
		result.Delivery = &DeliveryOptionDTO{
			ShippingAddressID: checkout.DeliveryOption.ShippingAddressID.String(),
			ShippingMethodID:  checkout.DeliveryOption.ShippingMethodID.String(),
		}
	}

	if checkout.PaymentMethod != nil {
		result.Payment = &PaymentMethodDTO{
			PaymentType:    checkout.PaymentMethod.PaymentType,
			PaymentDetails: checkout.PaymentMethod.PaymentDetails,
		}
	}

	return result
}

// ShippingAddressFromDomain converts a shipping address domain model to a DTO
func ShippingAddressFromDomain(address *model.ShippingAddress) *ShippingAddressDTO {
	return &ShippingAddressDTO{
		ID:            address.ID.String(),
		UserID:        address.UserID.String(),
		FirstName:     address.FirstName,
		LastName:      address.LastName,
		StreetAddress: address.StreetAddress,
		Apartment:     address.Apartment,
		City:          address.City,
		State:         address.State,
		PostalCode:    address.PostalCode,
		Country:       address.Country,
		PhoneNumber:   address.PhoneNumber,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     address.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ShippingMethodFromDomain converts a shipping method domain model to a DTO
func ShippingMethodFromDomain(method *model.ShippingMethod) *ShippingMethodDTO {
	return &ShippingMethodDTO{
		ID:                    method.ID.String(),
		Name:                  method.Name,
		Description:           method.Description,
		Price:                 method.Price,
		EstimatedDeliveryDays: method.EstimatedDeliveryDays,
	}
}
