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
	CodeDuration = time.Minute * 5
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

func getCodeExpiryDate() time.Time {
	return time.Now().Add(CodeDuration)
}

// If a user already exists, return an error.
// If not, then:
//
// 1) if there's no code yet, send it
// 2) if code exists but it's expired, send a new one
// 3) if code exists but it's not expired, do nothing and return from the func with a successa
// 4) If there's some other error, return it and do not send a code

// Sends a verification code to the user's email
// and temporarily stores user's data in the database
func (s *auth) SignUp(
	ctx context.Context,
	firstName string,
	lastName string,
	birthDate time.Time,
	email string,
	password string,
) error {
	if s.clientStore.ExistsByEmail(ctx, email) {
		return core.ErrEmailTaken
	}

	verif, err := s.emailVerifStore.GetByEmail(ctx, email)
	if err == nil && !verif.Expired() {
		return nil
	}

	if err != nil && !errors.Is(err, core.ErrResourceNotFound) {
		return fmt.Errorf("get verification by email: %w", err)
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	code, err := generateCode(CodeLength)
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}

	verif = &model.EmailVerification{
		Code:      code,
		ExpiresAt: getCodeExpiryDate(),
		User: model.EmailVerificationUser{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			BirthDate: birthDate,
			Password:  hashedPassword,
		},
	}

	if err := s.emailVerifStore.Create(ctx, verif); err != nil {
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
	if s.clientStore.ExistsByEmail(ctx, email) {
		return core.ErrEmailTaken
	}

	verif, err := s.emailVerifStore.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("get verification: %w", err)
	}

	if verif.Expired() {
		code, err := generateCode(CodeLength)
		if err != nil {
			return fmt.Errorf("generate code: %w", err)
		}

		verif.Code = code
		verif.ExpiresAt = getCodeExpiryDate()

		if err := s.emailVerifStore.Create(ctx, verif); err != nil {
			return fmt.Errorf("create verification: %w", err)
		}
	}

	return s.sendVerificationCode(ctx, email, verif.Code)
}

func (s *auth) SubmitVerification(ctx context.Context, email string, code string) (*model.Client, error) {
	verif, err := s.emailVerifStore.Get(ctx, email, code)
	if err != nil {
		if errors.Is(err, core.ErrResourceNotFound) {
			return nil, core.ErrInvalidVerificationCode
		}
		return nil, fmt.Errorf("get verification: %w", err)
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

	return cli, nil
}
