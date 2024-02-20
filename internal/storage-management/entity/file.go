package entity

import (
	"database/sql"
)

type File struct {
	ID        sql.NullString
	FileName  sql.NullString
	MimeType  sql.NullString
	Size      sql.NullInt64
	CreatedBy sql.NullInt64
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func (t *File) TableName() string {
	return "files"
}
