package pg_util

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

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

// NullFloat64 help to transform float64 to [database/sql.NullInt64]
func NullFloat64(val float64) sql.NullFloat64 {
	var result sql.NullFloat64
	result.Scan(val)

	return result
}

// NullTime help to transform int64 to [database/sql.NullInt64]
func NullTime(val time.Time) sql.NullTime {
	var result sql.NullTime
	result.Scan(val)

	return result
}

// StringArray ...
func StringArray(val []string) pq.StringArray {
	var result pq.StringArray
	result.Scan(val)

	return result
}

// StringArrayValue ...
func StringArrayValue(val pq.StringArray) []string {
	result := make([]string, 0, len(val))

	for _, v := range val {
		result = append(result, v)
	}

	return result
}
