package domain

import "time"

type User struct {
	ID       int64
	Email    string
	Password string

	Name     string
	Birthday time.Time
	Profile  string
}
