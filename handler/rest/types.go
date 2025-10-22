package rest

import "time"

// Auth
type signInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (signInReq) Valid() map[string]string {
	// TODO: validate
	return map[string]string{}
}

type signInRes struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type signUpReq struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

func (r signUpReq) Valid() map[string]string {
	// TODO: validate
	return map[string]string{}
}

type submitVerificationReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (submitVerificationReq) Valid() map[string]string {
	// TODO: validate
	return map[string]string{}
}

type submitVerificationRes struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
}

type sendVerificationCodeReq struct {
	Email string `json:email"`
}

func (sendVerificationCodeReq) Valid() map[string]string {
	// TODO: validate
	return map[string]string{}
}

// User
type getUserMeRes struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

// Account
type requestAccReq struct {
	Currency string `json:"currency"`
}

func (requestAccReq) Valid() map[string]string {
	// TODO: validate
	return map[string]string{}
}

type accRes struct {
	ID       int     `json:"id"`
	Number   string  `json:"number"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

type getAllAccsRes struct {
	Accounts []accRes `json:"accounts"`
}
