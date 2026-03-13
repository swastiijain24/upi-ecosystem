package models

import "time"

type Account struct {
	ID         string
	Name       string
	Phone      string
	Balance    int64
	Created_at time.Time
}

