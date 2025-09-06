package rest

import (
	"time"
)

// Auth
type signInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (signInReq) Valid() map[string]string {
	return map[string]string{}
}

type signInRes struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type signUpReq struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birth_date"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
}

func (r signUpReq) Valid() map[string]string {
	return map[string]string{}
}

type signUpRes struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
