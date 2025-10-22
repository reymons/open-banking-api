package core

import "errors"

var (
	ErrResourceNotFound        = errors.New("resource not found")
	ErrInvalidAccess           = errors.New("invalid access")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrVerificationCodeExpired = errors.New("verification code has expired")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrEmailTaken              = errors.New("email is already taken")
)
