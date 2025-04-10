package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/app/services/dto"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/domain/model"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/domain/repository"
)

// CheckoutService handles operations related to the checkout process
type CheckoutService struct {
	checkoutRepository repository.CheckoutRepository
	shippingRepository repository.ShippingRepository
	// External service clients would be injected here
	// productClient, inventoryClient, etc.
}

// NewCheckoutService creates a new checkout service
func NewCheckoutService(
	checkoutRepository repository.CheckoutRepository,
	shippingRepository repository.ShippingRepository,
) *CheckoutService {
	return &CheckoutService{
		checkoutRepository: checkoutRepository,
		shippingRepository: shippingRepository,
	}
}

// InitiateCheckout creates a new checkout from a cart
func (s *CheckoutService) InitiateCheckout(ctx context.Context, req *dto.CheckoutInitRequest) (*dto.CheckoutResponseDTO, error) {
	cartID, err := uuid.Parse(req.CartID)
	if err != nil {
		return nil, errors.New("invalid cart ID format")
	}

	// In a real implementation:
	// 1. Call Cart service to get cart details
	// 2. Check if cart exists and is not empty
	// 3. Check if items are in stock (via inventory service)
	// 4. Reserve inventory

	// For demonstration, we'll create a mock cart items list
	// In a real implementation, these would come from the cart service
	userID := uuid.New() // In reality, this would come from the cart or auth context
	items := []*model.CheckoutItem{
		{
			ProductID: uuid.New(),
			Name:      "Sample Product",
			Price:     29.99,
			Quantity:  2,
			Subtotal:  59.98,
			ImageURL:  "https://example.com/product.jpg",
		},
	}
	subtotal := 59.98

	// Create a new checkout
	checkout, err := model.NewCheckout(cartID, userID, items, subtotal)
	if err != nil {
		return nil, err
	}

	// Store the checkout
	if err := s.checkoutRepository.Save(ctx, checkout); err != nil {
		return nil, err
	}

	return dto.CheckoutFromDomain(checkout), nil
}

// GetCheckout retrieves a checkout by ID
func (s *CheckoutService) GetCheckout(ctx context.Context, checkoutID string) (*dto.CheckoutResponseDTO, error) {
	id, err := uuid.Parse(checkoutID)
	if err != nil {
		return nil, errors.New("invalid checkout ID format")
	}

	checkout, err := s.checkoutRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.CheckoutFromDomain(checkout), nil
}

// UpdateShipping updates the shipping details for a checkout
func (s *CheckoutService) UpdateShipping(ctx context.Context, checkoutID string, req *dto.ShippingDetailsRequest) (*dto.CheckoutResponseDTO, error) {
	id, err := uuid.Parse(checkoutID)
	if err != nil {
		return nil, errors.New("invalid checkout ID format")
	}

	checkout, err := s.checkoutRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	addressID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return nil, errors.New("invalid address ID format")
	}

	methodID, err := uuid.Parse(req.MethodID)
	if err != nil {
		return nil, errors.New("invalid shipping method ID format")
	}

	// Validate that address exists and belongs to the user
	address, err := s.shippingRepository.FindAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if address.UserID != checkout.UserID {
		return nil, errors.New("shipping address does not belong to the user")
	}

	// Validate that shipping method exists
	method, err := s.shippingRepository.FindMethodByID(ctx, methodID)
	if err != nil {
		return nil, err
	}

	// Create delivery option
	deliveryOption := model.NewDeliveryOption(addressID, methodID)

	// Update checkout with delivery option and shipping cost
	if err := checkout.SetDeliveryOption(deliveryOption, method.Price); err != nil {
		return nil, err
	}

	// Apply a mock tax rate of 10%
	checkout.CalculateTax(0.1)

	// Save the updated checkout
	if err := s.checkoutRepository.Save(ctx, checkout); err != nil {
		return nil, err
	}

	return dto.CheckoutFromDomain(checkout), nil
}

// SetPaymentMethod sets the payment method for a checkout
func (s *CheckoutService) SetPaymentMethod(ctx context.Context, checkoutID string, req *dto.PaymentMethodRequest) (*dto.CheckoutResponseDTO, error) {
	id, err := uuid.Parse(checkoutID)
	if err != nil {
		return nil, errors.New("invalid checkout ID format")
	}

	checkout, err := s.checkoutRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update checkout with payment method
	if err := checkout.SetPaymentMethod(req.PaymentType, req.PaymentDetails); err != nil {
		return nil, err
	}

	// Save the updated checkout
	if err := s.checkoutRepository.Save(ctx, checkout); err != nil {
		return nil, err
	}

	return dto.CheckoutFromDomain(checkout), nil
}

// CompleteCheckout finalizes the checkout process
func (s *CheckoutService) CompleteCheckout(ctx context.Context, checkoutID string) (*dto.CheckoutResponseDTO, error) {
	id, err := uuid.Parse(checkoutID)
	if err != nil {
		return nil, errors.New("invalid checkout ID format")
	}

	checkout, err := s.checkoutRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// In a real implementation:
	// 1. Process payment (call payment service)
	// 2. Convert reserved inventory to confirmed orders
	// 3. Create order in order management system
	// 4. Notify fulfillment service

	// Complete the checkout
	if err := checkout.Complete(); err != nil {
		return nil, err
	}

	// Save the updated checkout
	if err := s.checkoutRepository.Save(ctx, checkout); err != nil {
		return nil, err
	}

	return dto.CheckoutFromDomain(checkout), nil
}
