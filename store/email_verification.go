package store

import (
	"banking/core/model"
	"banking/db/pg"
	"context"
	"encoding/json"
)

type EmailVerification interface {
	Create(ctx context.Context, v *model.EmailVerification) error
	GetByEmail(ctx context.Context, email string) (*model.EmailVerification, error)
	Get(ctx context.Context, email string, code string) (*model.EmailVerification, error)
}

type emailVerification struct {
	pgcli      *pg.Client
	verifStore Verification
}

func NewEmailVerification(cli *pg.Client, vs Verification) EmailVerification {
	return &emailVerification{cli, vs}
}

func (s *emailVerification) Create(ctx context.Context, v *model.EmailVerification) error {
	data, err := json.Marshal(v.User)
	if err != nil {
		return newCoreError("marshal json", err)
	}

	err = s.verifStore.Create(ctx, s.pgcli.DB(), &verificationDto{
		Target:    v.User.Email,
		Code:      v.Code,
		ExpiresAt: v.ExpiresAt,
		Data:      data,
	})
	if err != nil {
		return newCoreError("create verification", err)
	}
	return nil
}

func (s *emailVerification) GetByEmail(ctx context.Context, email string) (*model.EmailVerification, error) {
	v, err := s.verifStore.GetByTarget(ctx, s.pgcli.DB(), email)
	if err != nil {
		return nil, newCoreError("get verification", err)
	}

	res := model.EmailVerification{
		Code:      v.Code,
		ExpiresAt: v.ExpiresAt,
	}
	if err := json.Unmarshal(v.Data, &res.User); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *emailVerification) Get(
	ctx context.Context,
	email string,
	code string,
) (*model.EmailVerification, error) {
	v, err := s.verifStore.Get(ctx, s.pgcli.DB(), email, code)
	if err != nil {
		return nil, newCoreError("get verification", err)
	}

	res := model.EmailVerification{
		Code:      v.Code,
		ExpiresAt: v.ExpiresAt,
	}
	if err := json.Unmarshal(v.Data, &res.User); err != nil {
		return nil, err
	}
	return &res, nil
}
