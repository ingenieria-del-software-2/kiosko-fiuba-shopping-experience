package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/domain/model"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/domain/repository"
)

// PostgreSQLCartRepository implements the CartRepository interface using PostgreSQL
type PostgreSQLCartRepository struct {
	db *sql.DB
}

// NewPostgreSQLCartRepository creates a new PostgreSQL repository for carts
func NewPostgreSQLCartRepository(db *sql.DB) repository.CartRepository {
	return &PostgreSQLCartRepository{
		db: db,
	}
}

// FindByID retrieves a cart by its ID
func (r *PostgreSQLCartRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Cart, error) {
	query := `
		SELECT id, user_id, items, created_at, updated_at
		FROM carts
		WHERE id = $1
	`

	var (
		cartID    uuid.UUID
		userID    uuid.UUID
		itemsJSON []byte
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cartID,
		&userID,
		&itemsJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}

	// Deserialize items from JSON
	var items []*model.CartItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, err
	}

	cart := &model.Cart{
		ID:        cartID,
		UserID:    userID,
		Items:     items,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}

	return cart, nil
}

// FindByUserID retrieves the current active cart for a user
func (r *PostgreSQLCartRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	query := `
		SELECT id, user_id, items, created_at, updated_at
		FROM carts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var (
		cartID    uuid.UUID
		dbUserID  uuid.UUID
		itemsJSON []byte
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&cartID,
		&dbUserID,
		&itemsJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}

	// Deserialize items from JSON
	var items []*model.CartItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, err
	}

	cart := &model.Cart{
		ID:        cartID,
		UserID:    dbUserID,
		Items:     items,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}

	return cart, nil
}

// Save persists a cart (creates or updates)
func (r *PostgreSQLCartRepository) Save(ctx context.Context, cart *model.Cart) error {
	// Serialize items to JSON
	itemsJSON, err := json.Marshal(cart.Items)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO carts (id, user_id, items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET items = $3, updated_at = $5
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		cart.ID,
		cart.UserID,
		itemsJSON,
		cart.CreatedAt,
		cart.UpdatedAt,
	)

	return err
}

// Delete removes a cart
func (r *PostgreSQLCartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM carts WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
