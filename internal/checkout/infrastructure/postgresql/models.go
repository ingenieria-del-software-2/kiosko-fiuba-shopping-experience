package postgresql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ShippingAddressModel is the PostgreSQL representation of a shipping address
type ShippingAddressModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	Name          string    `gorm:"type:varchar(100);not null"`
	StreetAddress string    `gorm:"type:varchar(255);not null"`
	City          string    `gorm:"type:varchar(100);not null"`
	State         string    `gorm:"type:varchar(100);not null"`
	PostalCode    string    `gorm:"type:varchar(20);not null"`
	Country       string    `gorm:"type:varchar(100);not null"`
	Phone         string    `gorm:"type:varchar(20)"`
	IsDefault     bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
	UpdatedAt     time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (ShippingAddressModel) TableName() string {
	return "shipping_addresses"
}

// ShippingMethodModel is the PostgreSQL representation of a shipping method
type ShippingMethodModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string    `gorm:"type:varchar(100);not null"`
	Description   string    `gorm:"type:text"`
	Price         float64   `gorm:"type:decimal(10,2);not null"`
	EstimatedDays int       `gorm:"type:integer;not null"`
	Active        bool      `gorm:"default:true"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
	UpdatedAt     time.Time `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (ShippingMethodModel) TableName() string {
	return "shipping_methods"
}

// CheckoutModel is the PostgreSQL representation of a checkout
type CheckoutModel struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey"`
	UserID            uuid.UUID         `gorm:"type:uuid;not null"`
	CartID            uuid.UUID         `gorm:"type:uuid;not null"`
	ShippingAddressID *uuid.UUID        `gorm:"type:uuid;references:shipping_addresses(id)"`
	ShippingMethodID  *uuid.UUID        `gorm:"type:uuid;references:shipping_methods(id)"`
	PaymentMethod     string            `gorm:"type:varchar(50)"`
	Subtotal          float64           `gorm:"type:decimal(10,2);not null"`
	ShippingCost      float64           `gorm:"type:decimal(10,2);default:0"`
	Tax               float64           `gorm:"type:decimal(10,2);default:0"`
	Total             float64           `gorm:"type:decimal(10,2);not null"`
	Status            string            `gorm:"type:varchar(20);not null;default:'PENDING'"`
	Items             CheckoutItemsJSON `gorm:"type:jsonb"`
	CreatedAt         time.Time         `gorm:"not null;default:now()"`
	UpdatedAt         time.Time         `gorm:"not null;default:now()"`
	CompletedAt       *time.Time        `gorm:"type:timestamp with time zone"`
}

// TableName overrides the table name for GORM
func (CheckoutModel) TableName() string {
	return "checkouts"
}

// CheckoutItemsJSON is a custom type for storing checkout items as JSON in PostgreSQL
type CheckoutItemsJSON []*CheckoutItemJSON

// CheckoutItemJSON is the JSON representation of a checkout item
type CheckoutItemJSON struct {
	ProductID uuid.UUID `json:"productId"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
	ImageURL  string    `json:"imageUrl"`
}

// Value implements the driver.Valuer interface for CheckoutItemsJSON
func (c CheckoutItemsJSON) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface for CheckoutItemsJSON
func (c *CheckoutItemsJSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
