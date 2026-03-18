package helpers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

func StringToPgUUID(s string) (pgtype.UUID) {
    u, err := uuid.Parse(s)
    if err != nil {
        return pgtype.UUID{}
    }

    return pgtype.UUID{
        Bytes: u,
        Valid: true,
    }
}