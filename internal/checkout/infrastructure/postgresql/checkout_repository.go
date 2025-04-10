package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/domain/model"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/domain/repository"
)

// PostgreSQLCheckoutRepository implements the CheckoutRepository interface using PostgreSQL
type PostgreSQLCheckoutRepository struct {
	db *sql.DB
}

// NewPostgreSQLCheckoutRepository creates a new PostgreSQL repository for checkouts
func NewPostgreSQLCheckoutRepository(db *sql.DB) repository.CheckoutRepository {
	return &PostgreSQLCheckoutRepository{
		db: db,
	}
}

// FindByID retrieves a checkout by its ID
func (r *PostgreSQLCheckoutRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Checkout, error) {
	query := `
		SELECT id, cart_id, user_id, status, items, subtotal, shipping_cost, tax, total, 
		       delivery_option, payment_method, created_at, updated_at
		FROM checkouts
		WHERE id = $1
	`

	var (
		checkoutID         uuid.UUID
		cartID             uuid.UUID
		userID             uuid.UUID
		status             string
		itemsJSON          []byte
		subtotal           float64
		shippingCost       float64
		tax                float64
		total              float64
		deliveryOptionJSON sql.NullString
		paymentMethodJSON  sql.NullString
		createdAt          sql.NullTime
		updatedAt          sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&checkoutID,
		&cartID,
		&userID,
		&status,
		&itemsJSON,
		&subtotal,
		&shippingCost,
		&tax,
		&total,
		&deliveryOptionJSON,
		&paymentMethodJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("checkout not found")
		}
		return nil, err
	}

	// Deserialize items from JSON
	var items []*model.CheckoutItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, err
	}

	// Create checkout object
	checkout := &model.Checkout{
		ID:           checkoutID,
		CartID:       cartID,
		UserID:       userID,
		Status:       model.CheckoutStatus(status),
		Items:        items,
		Subtotal:     subtotal,
		ShippingCost: shippingCost,
		Tax:          tax,
		Total:        total,
		CreatedAt:    createdAt.Time,
		UpdatedAt:    updatedAt.Time,
	}

	// Deserialize delivery option if present
	if deliveryOptionJSON.Valid {
		var deliveryOption model.DeliveryOption
		if err := json.Unmarshal([]byte(deliveryOptionJSON.String), &deliveryOption); err != nil {
			return nil, err
		}
		checkout.DeliveryOption = &deliveryOption
	}

	// Deserialize payment method if present
	if paymentMethodJSON.Valid {
		var paymentMethod model.PaymentMethod
		if err := json.Unmarshal([]byte(paymentMethodJSON.String), &paymentMethod); err != nil {
			return nil, err
		}
		checkout.PaymentMethod = &paymentMethod
	}

	return checkout, nil
}

// FindByCartID retrieves a checkout by cart ID
func (r *PostgreSQLCheckoutRepository) FindByCartID(ctx context.Context, cartID uuid.UUID) (*model.Checkout, error) {
	query := `
		SELECT id, cart_id, user_id, status, items, subtotal, shipping_cost, tax, total, 
		       delivery_option, payment_method, created_at, updated_at
		FROM checkouts
		WHERE cart_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var (
		checkoutID         uuid.UUID
		dbCartID           uuid.UUID
		userID             uuid.UUID
		status             string
		itemsJSON          []byte
		subtotal           float64
		shippingCost       float64
		tax                float64
		total              float64
		deliveryOptionJSON sql.NullString
		paymentMethodJSON  sql.NullString
		createdAt          sql.NullTime
		updatedAt          sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, cartID).Scan(
		&checkoutID,
		&dbCartID,
		&userID,
		&status,
		&itemsJSON,
		&subtotal,
		&shippingCost,
		&tax,
		&total,
		&deliveryOptionJSON,
		&paymentMethodJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("checkout not found")
		}
		return nil, err
	}

	// Deserialize items from JSON
	var items []*model.CheckoutItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, err
	}

	// Create checkout object
	checkout := &model.Checkout{
		ID:           checkoutID,
		CartID:       dbCartID,
		UserID:       userID,
		Status:       model.CheckoutStatus(status),
		Items:        items,
		Subtotal:     subtotal,
		ShippingCost: shippingCost,
		Tax:          tax,
		Total:        total,
		CreatedAt:    createdAt.Time,
		UpdatedAt:    updatedAt.Time,
	}

	// Deserialize delivery option if present
	if deliveryOptionJSON.Valid {
		var deliveryOption model.DeliveryOption
		if err := json.Unmarshal([]byte(deliveryOptionJSON.String), &deliveryOption); err != nil {
			return nil, err
		}
		checkout.DeliveryOption = &deliveryOption
	}

	// Deserialize payment method if present
	if paymentMethodJSON.Valid {
		var paymentMethod model.PaymentMethod
		if err := json.Unmarshal([]byte(paymentMethodJSON.String), &paymentMethod); err != nil {
			return nil, err
		}
		checkout.PaymentMethod = &paymentMethod
	}

	return checkout, nil
}

// FindByUserID retrieves the latest checkouts for a user
func (r *PostgreSQLCheckoutRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Checkout, error) {
	query := `
		SELECT id, cart_id, user_id, status, items, subtotal, shipping_cost, tax, total, 
		       delivery_option, payment_method, created_at, updated_at
		FROM checkouts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkouts []*model.Checkout

	for rows.Next() {
		var (
			checkoutID         uuid.UUID
			cartID             uuid.UUID
			dbUserID           uuid.UUID
			status             string
			itemsJSON          []byte
			subtotal           float64
			shippingCost       float64
			tax                float64
			total              float64
			deliveryOptionJSON sql.NullString
			paymentMethodJSON  sql.NullString
			createdAt          sql.NullTime
			updatedAt          sql.NullTime
		)

		if err := rows.Scan(
			&checkoutID,
			&cartID,
			&dbUserID,
			&status,
			&itemsJSON,
			&subtotal,
			&shippingCost,
			&tax,
			&total,
			&deliveryOptionJSON,
			&paymentMethodJSON,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		// Deserialize items from JSON
		var items []*model.CheckoutItem
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			return nil, err
		}

		// Create checkout object
		checkout := &model.Checkout{
			ID:           checkoutID,
			CartID:       cartID,
			UserID:       dbUserID,
			Status:       model.CheckoutStatus(status),
			Items:        items,
			Subtotal:     subtotal,
			ShippingCost: shippingCost,
			Tax:          tax,
			Total:        total,
			CreatedAt:    createdAt.Time,
			UpdatedAt:    updatedAt.Time,
		}

		// Deserialize delivery option if present
		if deliveryOptionJSON.Valid {
			var deliveryOption model.DeliveryOption
			if err := json.Unmarshal([]byte(deliveryOptionJSON.String), &deliveryOption); err != nil {
				return nil, err
			}
			checkout.DeliveryOption = &deliveryOption
		}

		// Deserialize payment method if present
		if paymentMethodJSON.Valid {
			var paymentMethod model.PaymentMethod
			if err := json.Unmarshal([]byte(paymentMethodJSON.String), &paymentMethod); err != nil {
				return nil, err
			}
			checkout.PaymentMethod = &paymentMethod
		}

		checkouts = append(checkouts, checkout)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkouts, nil
}

// Save persists a checkout (creates or updates)
func (r *PostgreSQLCheckoutRepository) Save(ctx context.Context, checkout *model.Checkout) error {
	// Serialize checkout items to JSON
	itemsJSON, err := json.Marshal(checkout.Items)
	if err != nil {
		return err
	}

	// Serialize delivery option to JSON if present
	var deliveryOptionJSON sql.NullString
	if checkout.DeliveryOption != nil {
		deliveryOptionBytes, err := json.Marshal(checkout.DeliveryOption)
		if err != nil {
			return err
		}
		deliveryOptionJSON = sql.NullString{
			String: string(deliveryOptionBytes),
			Valid:  true,
		}
	}

	// Serialize payment method to JSON if present
	var paymentMethodJSON sql.NullString
	if checkout.PaymentMethod != nil {
		paymentMethodBytes, err := json.Marshal(checkout.PaymentMethod)
		if err != nil {
			return err
		}
		paymentMethodJSON = sql.NullString{
			String: string(paymentMethodBytes),
			Valid:  true,
		}
	}

	query := `
		INSERT INTO checkouts (
			id, cart_id, user_id, status, items, subtotal, shipping_cost, tax, total,
			delivery_option, payment_method, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE
		SET status = $4, items = $5, subtotal = $6, shipping_cost = $7, tax = $8, total = $9,
			delivery_option = $10, payment_method = $11, updated_at = $13
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		checkout.ID,
		checkout.CartID,
		checkout.UserID,
		checkout.Status,
		itemsJSON,
		checkout.Subtotal,
		checkout.ShippingCost,
		checkout.Tax,
		checkout.Total,
		deliveryOptionJSON,
		paymentMethodJSON,
		checkout.CreatedAt,
		checkout.UpdatedAt,
	)

	return err
}
