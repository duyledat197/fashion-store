package postgresclient

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// PgTypeInt8 ...
func PgTypeInt8(src int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: src,
		Valid: src != 0,
	}
}

// PgTypeText ...
func PgTypeText(src string) pgtype.Text {
	return pgtype.Text{
		String: src,
		Valid:  src != "",
	}
}

// PgTypeTimestamp ...
func PgTypeTimestamp(src time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  src,
		Valid: src.IsZero(),
	}
}
