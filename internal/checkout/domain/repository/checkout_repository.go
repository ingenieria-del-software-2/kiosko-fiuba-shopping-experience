package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/model"
)

// CheckoutRepository defines the interface for checkout persistence operations
type CheckoutRepository interface {
	// FindByID retrieves a checkout by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Checkout, error)

	// FindByCartID retrieves a checkout by cart ID
	FindByCartID(ctx context.Context, cartID uuid.UUID) (*model.Checkout, error)

	// FindByUserID retrieves the latest checkout for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Checkout, error)

	// Save persists a checkout (creates or updates)
	Save(ctx context.Context, checkout *model.Checkout) error
}
