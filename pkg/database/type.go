package database

import "database/sql"

// NullString help to transform string to [database/sql.NullString]
func NullString(str string) sql.NullString {
	var result sql.NullString
	result.Scan(str)

	return result
}

// NullInt64 help to transform int64 to [database/sql.NullInt64]
func NullInt64(val int64) sql.NullInt64 {
	var result sql.NullInt64
	result.Scan(val)

	return result
}
