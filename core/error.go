package core

import "errors"

var (
	ErrResourceNotFound   = errors.New("resource not found")
	ErrInvalidAccess      = errors.New("invalid access")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
