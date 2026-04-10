package utils

import (
	
	"github.com/jackc/pgx/v5/pgtype"
)

func StringtoUUID(id string) pgtype.UUID {
	var uuid pgtype.UUID
	uuid.Scan(id)
	return uuid
}

func ToPGText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}


