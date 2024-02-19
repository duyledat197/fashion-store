package entity

import "github.com/jackc/pgx/v5/pgtype"

type File struct {
	ID        pgtype.Text
	FileName  pgtype.Text
	MimeType  pgtype.Text
	Size      pgtype.Int8
	CreatedBy pgtype.Int8
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (t *File) TableName() string {
	return "files"
}
