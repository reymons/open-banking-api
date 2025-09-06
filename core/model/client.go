package model

import (
	"banking/core/security"
	"time"
)

type Client struct {
	ID        int64
	Role      security.Role
	FirstName string
	LastName  string
	BirthDate time.Time
	Phone     string
	Email     string
	Password  string
	IsPartner bool
	CreatedAt time.Time
}
