package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CheckoutStatus represents the status of a checkout process
type CheckoutStatus string

const (
	CheckoutStatusInitiated        CheckoutStatus = "INITIATED"
	CheckoutStatusShippingSelected CheckoutStatus = "SHIPPING_SELECTED"
	CheckoutStatusPaymentSelected  CheckoutStatus = "PAYMENT_SELECTED"
	CheckoutStatusCompleted        CheckoutStatus = "COMPLETED"
	CheckoutStatusCancelled        CheckoutStatus = "CANCELLED"
)

// CheckoutItem represents an item in the checkout
type CheckoutItem struct {
	ProductID uuid.UUID `json:"productId"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
	ImageURL  string    `json:"imageUrl"`
}

// PaymentMethod represents payment method details
type PaymentMethod struct {
	PaymentType    string                 `json:"paymentType"`
	PaymentDetails map[string]interface{} `json:"paymentDetails"`
}

// Checkout represents the Checkout aggregate root in the Checkout Process bounded context
type Checkout struct {
	ID             uuid.UUID       `json:"id"`
	CartID         uuid.UUID       `json:"cartId"`
	UserID         uuid.UUID       `json:"userId"`
	Status         CheckoutStatus  `json:"status"`
	Items          []*CheckoutItem `json:"items"`
	Subtotal       float64         `json:"subtotal"`
	ShippingCost   float64         `json:"shippingCost"`
	Tax            float64         `json:"tax"`
	Total          float64         `json:"total"`
	DeliveryOption *DeliveryOption `json:"deliveryOption"`
	PaymentMethod  *PaymentMethod  `json:"paymentMethod"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// NewCheckout creates a new checkout from a cart
func NewCheckout(cartID, userID uuid.UUID, items []*CheckoutItem, subtotal float64) (*Checkout, error) {
	if cartID == uuid.Nil {
		return nil, errors.New("cart ID is required")
	}
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("checkout must have at least one item")
	}

	now := time.Now()
	return &Checkout{
		ID:           uuid.New(),
		CartID:       cartID,
		UserID:       userID,
		Status:       CheckoutStatusInitiated,
		Items:        items,
		Subtotal:     subtotal,
		ShippingCost: 0,
		Tax:          0,
		Total:        subtotal, // Initially just the subtotal
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// SetDeliveryOption sets the delivery option and updates the shipping cost
func (c *Checkout) SetDeliveryOption(deliveryOption *DeliveryOption, shippingCost float64) error {
	if c.Status == CheckoutStatusCancelled {
		return errors.New("cannot update a cancelled checkout")
	}

	if deliveryOption == nil {
		return errors.New("delivery option cannot be nil")
	}

	if deliveryOption.ShippingAddressID == uuid.Nil {
		return errors.New("shipping address ID is required")
	}

	if deliveryOption.ShippingMethodID == uuid.Nil {
		return errors.New("shipping method ID is required")
	}

	c.DeliveryOption = deliveryOption
	c.ShippingCost = shippingCost
	c.Status = CheckoutStatusShippingSelected
	c.UpdateTotal()
	c.UpdatedAt = time.Now()

	return nil
}

// SetPaymentMethod sets the payment method for the checkout
func (c *Checkout) SetPaymentMethod(paymentType string, paymentDetails map[string]interface{}) error {
	if c.Status == CheckoutStatusCancelled {
		return errors.New("cannot update a cancelled checkout")
	}

	if c.Status == CheckoutStatusInitiated {
		return errors.New("shipping option must be selected before payment")
	}

	if paymentType == "" {
		return errors.New("payment type is required")
	}

	if paymentDetails == nil {
		return errors.New("payment details cannot be nil")
	}

	c.PaymentMethod = &PaymentMethod{
		PaymentType:    paymentType,
		PaymentDetails: paymentDetails,
	}
	c.Status = CheckoutStatusPaymentSelected
	c.UpdatedAt = time.Now()

	return nil
}

// Complete marks the checkout as completed
func (c *Checkout) Complete() error {
	if c.Status == CheckoutStatusCancelled {
		return errors.New("cannot complete a cancelled checkout")
	}

	if c.Status != CheckoutStatusPaymentSelected {
		return errors.New("payment method must be selected before completing checkout")
	}

	c.Status = CheckoutStatusCompleted
	c.UpdatedAt = time.Now()

	return nil
}

// Cancel marks the checkout as cancelled
func (c *Checkout) Cancel() {
	if c.Status != CheckoutStatusCompleted {
		c.Status = CheckoutStatusCancelled
		c.UpdatedAt = time.Now()
	}
}

// CalculateTax calculates the tax amount based on the subtotal and shipping cost
func (c *Checkout) CalculateTax(taxRate float64) {
	c.Tax = (c.Subtotal + c.ShippingCost) * taxRate
	c.UpdateTotal()
	c.UpdatedAt = time.Now()
}

// UpdateTotal updates the total amount
func (c *Checkout) UpdateTotal() {
	c.Total = c.Subtotal + c.ShippingCost + c.Tax
}

// IsCompleted returns true if the checkout is completed
func (c *Checkout) IsCompleted() bool {
	return c.Status == CheckoutStatusCompleted
}

// IsCancelled returns true if the checkout is cancelled
func (c *Checkout) IsCancelled() bool {
	return c.Status == CheckoutStatusCancelled
}
