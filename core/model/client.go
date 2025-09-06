package model

import "time"

type Client struct {
	ID        int
	Role      int
	FirstName string
	LastName  string
	BirthDate time.Time
	Phone     string
	Email     string
	Password  string
	IsPartner bool
	CreatedAt time.Time
}
