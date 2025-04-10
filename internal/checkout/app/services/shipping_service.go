package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/app/services/dto"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/model"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/repository"
)

// ShippingService handles operations related to shipping addresses and methods
type ShippingService struct {
	shippingRepository repository.ShippingRepository
}

// NewShippingService creates a new shipping service
func NewShippingService(shippingRepository repository.ShippingRepository) *ShippingService {
	return &ShippingService{
		shippingRepository: shippingRepository,
	}
}

// AddShippingAddress adds a new shipping address for a user
func (s *ShippingService) AddShippingAddress(ctx context.Context, req *dto.ShippingAddressRequest) (*dto.ShippingAddressDTO, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Create a new shipping address
	address, err := model.NewShippingAddress(
		userID,
		req.FirstName,
		req.LastName,
		req.StreetAddress,
		req.Apartment,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.PhoneNumber,
		req.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	// If this is set as default, unset any existing default address
	if req.IsDefault {
		if err := s.unsetDefaultAddresses(ctx, userID); err != nil {
			return nil, err
		}
	}

	// Save the address
	if err := s.shippingRepository.SaveAddress(ctx, address); err != nil {
		return nil, err
	}

	return dto.ShippingAddressFromDomain(address), nil
}

// GetShippingAddress retrieves a shipping address by ID
func (s *ShippingService) GetShippingAddress(ctx context.Context, addressID string) (*dto.ShippingAddressDTO, error) {
	id, err := uuid.Parse(addressID)
	if err != nil {
		return nil, errors.New("invalid address ID format")
	}

	address, err := s.shippingRepository.FindAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ShippingAddressFromDomain(address), nil
}

// GetUserShippingAddresses retrieves all shipping addresses for a user
func (s *ShippingService) GetUserShippingAddresses(ctx context.Context, userID string) ([]*dto.ShippingAddressDTO, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	addresses, err := s.shippingRepository.FindAddressesByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ShippingAddressDTO, len(addresses))
	for i, address := range addresses {
		result[i] = dto.ShippingAddressFromDomain(address)
	}

	return result, nil
}

// UpdateShippingAddress updates a shipping address
func (s *ShippingService) UpdateShippingAddress(ctx context.Context, addressID string, req *dto.ShippingAddressRequest) (*dto.ShippingAddressDTO, error) {
	id, err := uuid.Parse(addressID)
	if err != nil {
		return nil, errors.New("invalid address ID format")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get the existing address
	address, err := s.shippingRepository.FindAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Ensure the address belongs to the user
	if address.UserID != userID {
		return nil, errors.New("shipping address does not belong to the user")
	}

	// If this is being set as default, unset any existing default address
	if req.IsDefault && !address.IsDefault {
		if err := s.unsetDefaultAddresses(ctx, userID); err != nil {
			return nil, err
		}
	}

	// Update the address
	if err := address.Update(
		req.FirstName,
		req.LastName,
		req.StreetAddress,
		req.Apartment,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.PhoneNumber,
		req.IsDefault,
	); err != nil {
		return nil, err
	}

	// Save the updated address
	if err := s.shippingRepository.SaveAddress(ctx, address); err != nil {
		return nil, err
	}

	return dto.ShippingAddressFromDomain(address), nil
}

// DeleteShippingAddress deletes a shipping address
func (s *ShippingService) DeleteShippingAddress(ctx context.Context, addressID string) error {
	id, err := uuid.Parse(addressID)
	if err != nil {
		return errors.New("invalid address ID format")
	}

	return s.shippingRepository.DeleteAddress(ctx, id)
}

// GetShippingMethods retrieves all available shipping methods
func (s *ShippingService) GetShippingMethods(ctx context.Context) ([]*dto.ShippingMethodDTO, error) {
	methods, err := s.shippingRepository.FindAllMethods(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ShippingMethodDTO, len(methods))
	for i, method := range methods {
		result[i] = dto.ShippingMethodFromDomain(method)
	}

	return result, nil
}

// unsetDefaultAddresses unsets the default flag on all addresses for a user
func (s *ShippingService) unsetDefaultAddresses(ctx context.Context, userID uuid.UUID) error {
	addresses, err := s.shippingRepository.FindAddressesByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, address := range addresses {
		if address.IsDefault {
			address.IsDefault = false
			if err := s.shippingRepository.SaveAddress(ctx, address); err != nil {
				return err
			}
		}
	}

	return nil
}
