package utils

import (
	"fmt"
	"strconv"
	"strings"

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

func RupeesToPaise(amount string) (int64, error) {
	parts := strings.Split(amount, ".")

	if len(parts) == 1 {
		rupees, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid rupee value: %w", err)
		}
		return rupees * 100, nil
	}

	if len(parts) == 2 {
		rupees, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid rupee value: %w", err)
		}

		paisePart := parts[1]
		if len(paisePart) == 1 {
			paisePart += "0"
		}
		if len(paisePart) > 2 {
			return 0, fmt.Errorf("max 2 decimal places allowed")
		}

		paise, _ := strconv.ParseInt(paisePart, 10, 64)
		return rupees*100 + paise, nil
	}

	return 0, fmt.Errorf("invalid amount format")
}
