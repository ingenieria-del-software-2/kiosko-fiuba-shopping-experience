package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Cart represents the Cart aggregate root in the Cart Management bounded context
type Cart struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"userId"`
	Items     []*CartItem `json:"items"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// NewCart creates a new empty cart for a user
func NewCart(userID uuid.UUID) *Cart {
	return &Cart{
		ID:        uuid.New(),
		UserID:    userID,
		Items:     make([]*CartItem, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddItem adds a product to the cart
func (c *Cart) AddItem(productID uuid.UUID, name string, price float64, quantity int, imageURL string) error {
	// Check if the item already exists in the cart
	for _, item := range c.Items {
		if item.ProductID == productID {
			// Update quantity instead of adding a new item
			newQuantity := item.Quantity + quantity
			if err := item.UpdateQuantity(newQuantity); err != nil {
				return err
			}
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	// Create a new cart item
	newItem, err := NewCartItem(productID, name, price, quantity, imageURL)
	if err != nil {
		return err
	}

	// Add the new item to the cart
	c.Items = append(c.Items, newItem)
	c.UpdatedAt = time.Now()
	return nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (c *Cart) UpdateItemQuantity(itemID uuid.UUID, quantity int) error {
	for _, item := range c.Items {
		if item.ID == itemID {
			if err := item.UpdateQuantity(quantity); err != nil {
				return err
			}
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("item not found in cart")
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(itemID uuid.UUID) error {
	for i, item := range c.Items {
		if item.ID == itemID {
			// Remove the item from the slice
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("item not found in cart")
}

// Clear empties the cart
func (c *Cart) Clear() {
	c.Items = make([]*CartItem, 0)
	c.UpdatedAt = time.Now()
}

// TotalItems returns the total number of items in the cart
func (c *Cart) TotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// Subtotal calculates the subtotal of the cart (sum of all item subtotals)
func (c *Cart) Subtotal() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.Subtotal()
	}
	return total
}

// IsEmpty checks if the cart is empty
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

// GetItem returns a cart item by ID
func (c *Cart) GetItem(itemID uuid.UUID) (*CartItem, error) {
	for _, item := range c.Items {
		if item.ID == itemID {
			return item, nil
		}
	}
	return nil, errors.New("item not found in cart")
}
