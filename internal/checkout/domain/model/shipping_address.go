package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ShippingAddress represents a shipping address entity in the Checkout Process bounded context
type ShippingAddress struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userId"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	StreetAddress string    `json:"streetAddress"`
	Apartment     string    `json:"apartment"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PostalCode    string    `json:"postalCode"`
	Country       string    `json:"country"`
	PhoneNumber   string    `json:"phoneNumber"`
	IsDefault     bool      `json:"isDefault"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// NewShippingAddress creates a new shipping address
func NewShippingAddress(
	userID uuid.UUID,
	firstName, lastName, streetAddress, apartment, city, state, postalCode, country, phoneNumber string,
	isDefault bool,
) (*ShippingAddress, error) {
	// Validate required fields
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if firstName == "" {
		return nil, errors.New("first name is required")
	}
	if lastName == "" {
		return nil, errors.New("last name is required")
	}
	if streetAddress == "" {
		return nil, errors.New("street address is required")
	}
	if city == "" {
		return nil, errors.New("city is required")
	}
	if state == "" {
		return nil, errors.New("state is required")
	}
	if postalCode == "" {
		return nil, errors.New("postal code is required")
	}
	if country == "" {
		return nil, errors.New("country is required")
	}
	if phoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	now := time.Now()
	return &ShippingAddress{
		ID:            uuid.New(),
		UserID:        userID,
		FirstName:     firstName,
		LastName:      lastName,
		StreetAddress: streetAddress,
		Apartment:     apartment,
		City:          city,
		State:         state,
		PostalCode:    postalCode,
		Country:       country,
		PhoneNumber:   phoneNumber,
		IsDefault:     isDefault,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Update updates the shipping address details
func (a *ShippingAddress) Update(
	firstName, lastName, streetAddress, apartment, city, state, postalCode, country, phoneNumber string,
	isDefault bool,
) error {
	// Validate required fields
	if firstName == "" {
		return errors.New("first name is required")
	}
	if lastName == "" {
		return errors.New("last name is required")
	}
	if streetAddress == "" {
		return errors.New("street address is required")
	}
	if city == "" {
		return errors.New("city is required")
	}
	if state == "" {
		return errors.New("state is required")
	}
	if postalCode == "" {
		return errors.New("postal code is required")
	}
	if country == "" {
		return errors.New("country is required")
	}
	if phoneNumber == "" {
		return errors.New("phone number is required")
	}

	a.FirstName = firstName
	a.LastName = lastName
	a.StreetAddress = streetAddress
	a.Apartment = apartment
	a.City = city
	a.State = state
	a.PostalCode = postalCode
	a.Country = country
	a.PhoneNumber = phoneNumber
	a.IsDefault = isDefault
	a.UpdatedAt = time.Now()

	return nil
}

// FullName returns the full name of the recipient
func (a *ShippingAddress) FullName() string {
	return a.FirstName + " " + a.LastName
}

// FormattedAddress returns a formatted address string
func (a *ShippingAddress) FormattedAddress() string {
	apartment := ""
	if a.Apartment != "" {
		apartment = ", " + a.Apartment
	}
	return a.StreetAddress + apartment + ", " + a.City + ", " + a.State + " " + a.PostalCode + ", " + a.Country
}
