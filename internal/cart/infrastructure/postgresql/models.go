package postgresql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// CartModel is the PostgreSQL representation of a cart
type CartModel struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID     `gorm:"type:uuid;not null"`
	Items     CartItemsJSON `gorm:"type:jsonb"`
	CreatedAt time.Time     `gorm:"not null;default:now()"`
	UpdatedAt time.Time     `gorm:"not null;default:now()"`
}

// TableName overrides the table name for GORM
func (CartModel) TableName() string {
	return "carts"
}

// CartItemsJSON is a custom type for storing cart items as JSON in PostgreSQL
type CartItemsJSON []*CartItemJSON

// CartItemJSON is the JSON representation of a cart item
type CartItemJSON struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productId"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	ImageURL  string    `json:"imageUrl"`
}

// Value implements the driver.Valuer interface for CartItemsJSON
func (c CartItemsJSON) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface for CartItemsJSON
func (c *CartItemsJSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
