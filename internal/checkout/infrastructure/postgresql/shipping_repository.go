package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/model"
	"github.com/kiosko-fiuba/shopping-experience/internal/checkout/domain/repository"
)

// PostgreSQLShippingRepository implements the ShippingRepository interface using PostgreSQL
type PostgreSQLShippingRepository struct {
	db *sql.DB
}

// NewPostgreSQLShippingRepository creates a new PostgreSQL repository for shipping
func NewPostgreSQLShippingRepository(db *sql.DB) repository.ShippingRepository {
	return &PostgreSQLShippingRepository{
		db: db,
	}
}

// FindAddressByID retrieves a shipping address by its ID
func (r *PostgreSQLShippingRepository) FindAddressByID(ctx context.Context, id uuid.UUID) (*model.ShippingAddress, error) {
	query := `
		SELECT id, user_id, first_name, last_name, street_address, apartment, city, state, 
		       postal_code, country, phone_number, is_default, created_at, updated_at
		FROM shipping_addresses
		WHERE id = $1
	`

	var (
		addressID     uuid.UUID
		userID        uuid.UUID
		firstName     string
		lastName      string
		streetAddress string
		apartment     sql.NullString
		city          string
		state         string
		postalCode    string
		country       string
		phoneNumber   string
		isDefault     bool
		createdAt     sql.NullTime
		updatedAt     sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&addressID,
		&userID,
		&firstName,
		&lastName,
		&streetAddress,
		&apartment,
		&city,
		&state,
		&postalCode,
		&country,
		&phoneNumber,
		&isDefault,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("shipping address not found")
		}
		return nil, err
	}

	apartmentStr := ""
	if apartment.Valid {
		apartmentStr = apartment.String
	}

	address := &model.ShippingAddress{
		ID:            addressID,
		UserID:        userID,
		FirstName:     firstName,
		LastName:      lastName,
		StreetAddress: streetAddress,
		Apartment:     apartmentStr,
		City:          city,
		State:         state,
		PostalCode:    postalCode,
		Country:       country,
		PhoneNumber:   phoneNumber,
		IsDefault:     isDefault,
		CreatedAt:     createdAt.Time,
		UpdatedAt:     updatedAt.Time,
	}

	return address, nil
}

// FindAddressesByUserID retrieves all shipping addresses for a user
func (r *PostgreSQLShippingRepository) FindAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ShippingAddress, error) {
	query := `
		SELECT id, user_id, first_name, last_name, street_address, apartment, city, state, 
		       postal_code, country, phone_number, is_default, created_at, updated_at
		FROM shipping_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*model.ShippingAddress

	for rows.Next() {
		var (
			addressID     uuid.UUID
			dbUserID      uuid.UUID
			firstName     string
			lastName      string
			streetAddress string
			apartment     sql.NullString
			city          string
			state         string
			postalCode    string
			country       string
			phoneNumber   string
			isDefault     bool
			createdAt     sql.NullTime
			updatedAt     sql.NullTime
		)

		if err := rows.Scan(
			&addressID,
			&dbUserID,
			&firstName,
			&lastName,
			&streetAddress,
			&apartment,
			&city,
			&state,
			&postalCode,
			&country,
			&phoneNumber,
			&isDefault,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		apartmentStr := ""
		if apartment.Valid {
			apartmentStr = apartment.String
		}

		address := &model.ShippingAddress{
			ID:            addressID,
			UserID:        dbUserID,
			FirstName:     firstName,
			LastName:      lastName,
			StreetAddress: streetAddress,
			Apartment:     apartmentStr,
			City:          city,
			State:         state,
			PostalCode:    postalCode,
			Country:       country,
			PhoneNumber:   phoneNumber,
			IsDefault:     isDefault,
			CreatedAt:     createdAt.Time,
			UpdatedAt:     updatedAt.Time,
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

// SaveAddress persists a shipping address (creates or updates)
func (r *PostgreSQLShippingRepository) SaveAddress(ctx context.Context, address *model.ShippingAddress) error {
	query := `
		INSERT INTO shipping_addresses (
			id, user_id, first_name, last_name, street_address, apartment, city, state,
			postal_code, country, phone_number, is_default, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (id) DO UPDATE
		SET first_name = $3, last_name = $4, street_address = $5, apartment = $6,
			city = $7, state = $8, postal_code = $9, country = $10,
			phone_number = $11, is_default = $12, updated_at = $14
	`

	apartment := sql.NullString{
		String: address.Apartment,
		Valid:  address.Apartment != "",
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		address.ID,
		address.UserID,
		address.FirstName,
		address.LastName,
		address.StreetAddress,
		apartment,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.PhoneNumber,
		address.IsDefault,
		address.CreatedAt,
		address.UpdatedAt,
	)

	return err
}

// DeleteAddress removes a shipping address
func (r *PostgreSQLShippingRepository) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shipping_addresses WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// FindMethodByID retrieves a shipping method by its ID
func (r *PostgreSQLShippingRepository) FindMethodByID(ctx context.Context, id uuid.UUID) (*model.ShippingMethod, error) {
	query := `
		SELECT id, name, description, price, estimated_delivery_days
		FROM shipping_methods
		WHERE id = $1
	`

	var (
		methodID              uuid.UUID
		name                  string
		description           string
		price                 float64
		estimatedDeliveryDays int
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&methodID,
		&name,
		&description,
		&price,
		&estimatedDeliveryDays,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("shipping method not found")
		}
		return nil, err
	}

	method := &model.ShippingMethod{
		ID:                    methodID,
		Name:                  name,
		Description:           description,
		Price:                 price,
		EstimatedDeliveryDays: estimatedDeliveryDays,
	}

	return method, nil
}

// FindAllMethods retrieves all available shipping methods
func (r *PostgreSQLShippingRepository) FindAllMethods(ctx context.Context) ([]*model.ShippingMethod, error) {
	query := `
		SELECT id, name, description, price, estimated_delivery_days
		FROM shipping_methods
		ORDER BY price ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []*model.ShippingMethod

	for rows.Next() {
		var (
			methodID              uuid.UUID
			name                  string
			description           string
			price                 float64
			estimatedDeliveryDays int
		)

		if err := rows.Scan(
			&methodID,
			&name,
			&description,
			&price,
			&estimatedDeliveryDays,
		); err != nil {
			return nil, err
		}

		method := &model.ShippingMethod{
			ID:                    methodID,
			Name:                  name,
			Description:           description,
			Price:                 price,
			EstimatedDeliveryDays: estimatedDeliveryDays,
		}

		methods = append(methods, method)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return methods, nil
}
