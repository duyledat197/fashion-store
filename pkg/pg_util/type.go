package pg_util

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

// PgTypeTimestamptz ...
func PgTypeTimestamptz(src time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  src,
		Valid: src.IsZero(),
	}
}
