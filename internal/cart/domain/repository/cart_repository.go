package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kiosko-fiuba/shopping-experience/internal/cart/domain/model"
)

// CartRepository defines the interface for cart persistence operations
type CartRepository interface {
	// FindByID retrieves a cart by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Cart, error)

	// FindByUserID retrieves the current active cart for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Cart, error)

	// Save persists a cart (creates or updates)
	Save(ctx context.Context, cart *model.Cart) error

	// Delete removes a cart
	Delete(ctx context.Context, id uuid.UUID) error
}
