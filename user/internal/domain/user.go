package domain

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}
