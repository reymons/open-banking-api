package service

import (
	"banking/core"
	"banking/core/model"
	"banking/store"
	"banking/util"
	"context"
	"fmt"
	"time"
)

type Auth interface {
	SignIn(
		ctx context.Context,
		email string,
		password string,
	) (*model.Client, error)
	SignUp(
		ctx context.Context,
		firstName string,
		lastName string,
		birthDate time.Time,
		email string,
		phone string,
		password string,
	) (*model.Client, error)
}

type auth struct {
	clientStore store.Client
}

func NewAuth(cs store.Client) Auth {
	return &auth{cs}
}

func (s *auth) SignIn(
	ctx context.Context,
	email string,
	password string,
) (*model.Client, error) {
	cli, err := s.clientStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get client by email: %w", core.ErrInvalidCredentials)
	}
	if !util.VerifyPassword(password, cli.Password) {
		return nil, fmt.Errorf("verify password: %w", core.ErrInvalidCredentials)
	}
	return cli, nil
}

func (s *auth) SignUp(
	ctx context.Context,
	firstName string,
	lastName string,
	birthDate time.Time,
	email string,
	phone string,
	password string,
) (*model.Client, error) {
	hashed, err := util.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	cli := &model.Client{
		Role:      core.RoleClient,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Email:     email,
		Phone:     phone,
		Password:  hashed,
	}
	if err := s.clientStore.Create(ctx, cli); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return cli, nil
}
