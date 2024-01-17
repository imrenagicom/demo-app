package proto

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromSQLNullTime(t sql.NullTime) *timestamppb.Timestamp {
	if !t.Valid {
		return nil
	}
	return timestamppb.New(t.Time)
}
