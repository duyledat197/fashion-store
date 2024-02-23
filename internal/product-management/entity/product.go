package entity

import (
	"database/sql"

	"github.com/lib/pq"
)

// Product represents the entity structure for products in the database.
type Product struct {
	ID          sql.NullInt64   `db:"id"`
	Name        sql.NullString  `db:"name"`
	Type        sql.NullString  `db:"type"`
	ImageURLs   pq.StringArray  `db:"image_urls"`
	Description sql.NullString  `db:"description"`
	Price       sql.NullFloat64 `db:"price"`
	CreatedBy   sql.NullInt64   `db:"created_by"`
	CreatedAt   sql.NullTime    `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}

// TableName returns the name of the database table associated with the Product entity.
func (u *Product) TableName() string {
	return "products"
}
