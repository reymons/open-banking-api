package model

import "time"

type Client struct {
	ID        int
	Role      int
	FirstName string
	LastName  string
	BirthDate time.Time
	Email     string
	Password  string
	IsPartner bool
	CreatedAt time.Time
}
