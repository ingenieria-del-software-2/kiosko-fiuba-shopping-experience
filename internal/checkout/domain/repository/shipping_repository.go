package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/model"
)

// ShippingRepository defines the interface for shipping-related persistence operations
type ShippingRepository interface {
	// FindAddressByID retrieves a shipping address by its ID
	FindAddressByID(ctx context.Context, id uuid.UUID) (*model.ShippingAddress, error)

	// FindAddressesByUserID retrieves all shipping addresses for a user
	FindAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ShippingAddress, error)

	// SaveAddress persists a shipping address (creates or updates)
	SaveAddress(ctx context.Context, address *model.ShippingAddress) error

	// DeleteAddress removes a shipping address
	DeleteAddress(ctx context.Context, id uuid.UUID) error

	// FindMethodByID retrieves a shipping method by its ID
	FindMethodByID(ctx context.Context, id uuid.UUID) (*model.ShippingMethod, error)

	// FindAllMethods retrieves all available shipping methods
	FindAllMethods(ctx context.Context) ([]*model.ShippingMethod, error)
}
