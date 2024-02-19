// Package entity ...
package entity

import "github.com/jackc/pgx/v5/pgtype"

// Coupon ...
type Coupon struct {
	ID          pgtype.Int8
	Code        pgtype.Text
	From        pgtype.Timestamptz
	To          pgtype.Timestamptz
	Rules       []byte
	ImageURL    pgtype.Text
	Description pgtype.Text
	CreatedBy   pgtype.Int8
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}
