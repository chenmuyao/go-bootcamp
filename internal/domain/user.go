package domain

import "time"

type User struct {
	Birthday time.Time
	Email    string
	Password string
	Phone    string
	Name     string
	Profile  string
	ID       int64
}
