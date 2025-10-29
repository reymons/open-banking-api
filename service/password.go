package service

import (
	"banking/config"
	"banking/core"
	"banking/core/model"
	"banking/db/pg"
	"banking/store"
	"banking/util"
	"context"
	"fmt"
	"time"
)

var (
	resetTokenDuration = time.Minute * 10
	resetTokenLength   = 20
)

type PasswordService interface {
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type passwordService struct {
	pgcli             *pg.Client
	emailService      EmailService
	resetPswdReqStore store.ResetPasswordReqStore
	clientStore       store.Client
	appCfg            *config.InternalConfig
}

func NewPassword(
	pgcli *pg.Client,
	es EmailService,
	rps store.ResetPasswordReqStore,
	cs store.Client,
	cfg *config.InternalConfig,
) PasswordService {
	return &passwordService{pgcli, es, rps, cs, cfg}
}

// Sends a reset-password link to the specified email
// if a user by that email exists
func (s *passwordService) RequestPasswordReset(ctx context.Context, email string) error {
	token, err := util.Base64URLString(resetTokenLength)
	if err != nil {
		return fmt.Errorf("get base64url string: %w", err)
	}
	hashedToken := util.HashString(token)

	client, err := s.clientStore.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	resetPswd := &model.ResetPasswordReq{
		ClientID:  client.ID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(resetTokenDuration),
	}
	err = s.resetPswdReqStore.Create(ctx, resetPswd)
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}

	link := fmt.Sprintf("%s/reset-password?token=%s", s.appCfg.GetMainWebsiteUrl(), token)

	s.emailService.SendMessage(
		ctx,
		[]string{email},
		"Reset Password",
		fmt.Sprintf(`Click <a href="%s">here</a> to <strong>reset</strong> your password`, link),
	)

	return nil
}

// Verifies that the reset-password token is valid
// and sets a new user password
func (s *passwordService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	hashedToken := util.HashString(token)
	resetPswd, err := s.resetPswdReqStore.GetByToken(ctx, hashedToken)
	if err != nil {
		return core.ErrInvalidToken
	}

	if resetPswd.Expired() {
		return core.ErrTokenExpired
	}

	hashedPswd, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.pgcli.ExecTx(ctx, func(tx *pg.Tx) error {
		err := s.clientStore.SetPassword(ctx, tx, resetPswd.ClientID, hashedPswd)
		if err != nil {
			return fmt.Errorf("set new password: %w", err)
		}
		return s.resetPswdReqStore.DeleteByClientID(ctx, tx, resetPswd.ClientID)
	})
}
