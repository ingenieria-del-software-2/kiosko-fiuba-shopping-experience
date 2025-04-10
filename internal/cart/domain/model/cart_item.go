package model

import (
	"errors"

	"github.com/google/uuid"
)

// CartItem represents a value object for items in a shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productId"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	ImageURL  string    `json:"imageUrl"`
}

// Subtotal calculates the subtotal for this cart item (price * quantity)
func (i *CartItem) Subtotal() float64 {
	return i.Price * float64(i.Quantity)
}

// NewCartItem creates a new cart item
func NewCartItem(productID uuid.UUID, name string, price float64, quantity int, imageURL string) (*CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	return &CartItem{
		ID:        uuid.New(),
		ProductID: productID,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
		ImageURL:  imageURL,
	}, nil
}

// UpdateQuantity updates the quantity of the cart item
func (i *CartItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	i.Quantity = quantity
	return nil
}
