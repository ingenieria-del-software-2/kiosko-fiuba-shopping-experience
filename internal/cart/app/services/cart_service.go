package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/app/services/dto"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/domain/model"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/domain/repository"
)

// CartService handles operations related to shopping carts
type CartService struct {
	cartRepository repository.CartRepository
}

// NewCartService creates a new cart service
func NewCartService(cartRepository repository.CartRepository) *CartService {
	return &CartService{
		cartRepository: cartRepository,
	}
}

// CreateCart creates a new empty cart for a user
func (s *CartService) CreateCart(ctx context.Context, req *dto.CartCreateRequest) (*dto.CartDTO, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Check if user already has an active cart
	existingCart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err == nil && existingCart != nil {
		// Return existing cart
		return dto.CartFromDomain(existingCart), nil
	}

	// Create a new cart
	cart := model.NewCart(userID)
	if err := s.cartRepository.Save(ctx, cart); err != nil {
		return nil, err
	}

	return dto.CartFromDomain(cart), nil
}

// GetCart retrieves a cart by ID
func (s *CartService) GetCart(ctx context.Context, cartID string) (*dto.CartDTO, error) {
	id, err := uuid.Parse(cartID)
	if err != nil {
		return nil, errors.New("invalid cart ID format")
	}

	cart, err := s.cartRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.CartFromDomain(cart), nil
}

// AddCartItem adds a product to a cart
func (s *CartService) AddCartItem(ctx context.Context, cartID string, req *dto.CartItemRequest) (*dto.CartDTO, error) {
	id, err := uuid.Parse(cartID)
	if err != nil {
		return nil, errors.New("invalid cart ID format")
	}

	cart, err := s.cartRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	if err := cart.AddItem(productID, req.Name, req.Price, req.Quantity, req.ImageURL); err != nil {
		return nil, err
	}

	if err := s.cartRepository.Save(ctx, cart); err != nil {
		return nil, err
	}

	return dto.CartFromDomain(cart), nil
}

// UpdateCartItem updates the quantity of a cart item
func (s *CartService) UpdateCartItem(ctx context.Context, cartID string, itemID string, req *dto.CartItemUpdateRequest) (*dto.CartDTO, error) {
	cartUUID, err := uuid.Parse(cartID)
	if err != nil {
		return nil, errors.New("invalid cart ID format")
	}

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return nil, errors.New("invalid item ID format")
	}

	cart, err := s.cartRepository.FindByID(ctx, cartUUID)
	if err != nil {
		return nil, err
	}

	if err := cart.UpdateItemQuantity(itemUUID, req.Quantity); err != nil {
		return nil, err
	}

	if err := s.cartRepository.Save(ctx, cart); err != nil {
		return nil, err
	}

	return dto.CartFromDomain(cart), nil
}

// RemoveCartItem removes an item from a cart
func (s *CartService) RemoveCartItem(ctx context.Context, cartID string, itemID string) error {
	cartUUID, err := uuid.Parse(cartID)
	if err != nil {
		return errors.New("invalid cart ID format")
	}

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return errors.New("invalid item ID format")
	}

	cart, err := s.cartRepository.FindByID(ctx, cartUUID)
	if err != nil {
		return err
	}

	if err := cart.RemoveItem(itemUUID); err != nil {
		return err
	}

	return s.cartRepository.Save(ctx, cart)
}

// DeleteCart removes a cart
func (s *CartService) DeleteCart(ctx context.Context, cartID string) error {
	id, err := uuid.Parse(cartID)
	if err != nil {
		return errors.New("invalid cart ID format")
	}

	return s.cartRepository.Delete(ctx, id)
}
