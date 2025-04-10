package dto

import (
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/domain/model"
)

// CartItemDTO represents cart item data for API responses
type CartItemDTO struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
	ImageURL  string  `json:"imageUrl"`
}

// CartDTO represents cart data for API responses
type CartDTO struct {
	ID         string        `json:"id"`
	UserID     string        `json:"userId"`
	Items      []CartItemDTO `json:"items"`
	TotalItems int           `json:"totalItems"`
	Subtotal   float64       `json:"subtotal"`
	CreatedAt  string        `json:"createdAt"`
	UpdatedAt  string        `json:"updatedAt"`
}

// CartCreateRequest represents the request to create a new cart
type CartCreateRequest struct {
	UserID string `json:"userId" validate:"required,uuid"`
}

// CartItemRequest represents the request to add a product to a cart
type CartItemRequest struct {
	ProductID string  `json:"productId" validate:"required,uuid"`
	Name      string  `json:"name" validate:"required"`
	Price     float64 `json:"price" validate:"gte=0"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	ImageURL  string  `json:"imageUrl"`
}

// CartItemUpdateRequest represents the request to update a cart item
type CartItemUpdateRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

// FromDomain converts a cart domain model to a DTO
func CartFromDomain(cart *model.Cart) *CartDTO {
	items := make([]CartItemDTO, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = CartItemDTO{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal(),
			ImageURL:  item.ImageURL,
		}
	}

	return &CartDTO{
		ID:         cart.ID.String(),
		UserID:     cart.UserID.String(),
		Items:      items,
		TotalItems: cart.TotalItems(),
		Subtotal:   cart.Subtotal(),
		CreatedAt:  cart.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  cart.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
