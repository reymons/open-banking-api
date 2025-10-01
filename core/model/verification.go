package model

import "time"

type EmailVerificationUser struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	BirthDate time.Time `json:"birthDate"`
}

type EmailVerification struct {
	ExpiresAt time.Time
	Code      string
	User      EmailVerificationUser
}

func (v *EmailVerification) Expired() bool {
	return time.Now().After(v.ExpiresAt)
}
