package service

import (
	"banking/core"
	"banking/core/model"
	"banking/store"
	"banking/util"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var (
	CodeDuration = time.Minute * 10
	CodeLength   = 6
)

func generateCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("Code length should be a positive number")
	}
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = strconv.Itoa(rand.Intn(10))[0]
	}
	return string(code), nil
}

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
		password string,
	) error
	SubmitVerification(ctx context.Context, email string, code string) (*model.Client, error)
	SendVerificationCode(ctx context.Context, email string) error
}

type auth struct {
	clientStore     store.Client
	emailVerifStore store.EmailVerification
	emailService    EmailService
}

func NewAuth(cs store.Client, vs store.EmailVerification, es EmailService) Auth {
	return &auth{cs, vs, es}
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

// Sends a verification code to the user's email
func (s *auth) SignUp(
	ctx context.Context,
	firstName string,
	lastName string,
	birthDate time.Time,
	email string,
	password string,
) error {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	code, err := generateCode(CodeLength)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}

	verif := model.EmailVerification{
		Code:      code,
		ExpiresAt: time.Now().Add(CodeDuration),
		User: model.EmailVerificationUser{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			BirthDate: birthDate,
			Password:  hashedPassword,
		},
	}

	if err := s.emailVerifStore.Create(ctx, &verif); err != nil {
		return fmt.Errorf("create verification: %w", err)
	}

	if err := s.sendVerificationCode(ctx, email, code); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

func (s *auth) sendVerificationCode(ctx context.Context, email string, code string) error {
	return s.emailService.SendMessage(
		ctx,
		"noreply@reymons.net",
		[]string{email},
		"Open Banking - Email Address Verification",
		fmt.Sprintf("Your verification code is <strong>%s</strong>", code),
	)
}

func (s *auth) SendVerificationCode(ctx context.Context, email string) error {
	verif, err := s.emailVerifStore.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("get verification: %w", err)
	}
	return s.sendVerificationCode(ctx, email, verif.Code)
}

func (s *auth) SubmitVerification(ctx context.Context, email string, code string) (*model.Client, error) {
	verif, err := s.emailVerifStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get verification: %w", err)
	}

	if verif.Code != code {
		return nil, core.ErrInvalidVerificationCode
	}

	if verif.Expired() {
		return nil, core.ErrVerificationCodeExpired
	}

	cli := &model.Client{
		Role:      core.RoleClient,
		FirstName: verif.User.FirstName,
		LastName:  verif.User.LastName,
		BirthDate: verif.User.BirthDate,
		Email:     verif.User.Email,
		Password:  verif.User.Password,
	}
	if err := s.clientStore.Create(ctx, cli); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	if err := s.emailVerifStore.Delete(ctx, email, code); err != nil {
		return nil, fmt.Errorf("delete verification: %w", err)
	}

	return cli, nil
}
