package store

import (
	"banking/core/model"
	"banking/db/pg"
	"context"
)

type ResetPasswordReqStore interface {
	Create(ctx context.Context, m *model.ResetPasswordReq) error
	GetByClientID(ctx context.Context, id int) (*model.ResetPasswordReq, error)
	GetByToken(ctx context.Context, token string) (*model.ResetPasswordReq, error)
	DeleteByClientID(ctx context.Context, db pg.DB, id int) error
}

type resetPasswordReqStore struct {
	pgcli *pg.Client
}

func NewResetPasswordReq(cli *pg.Client) ResetPasswordReqStore {
	return &resetPasswordReqStore{cli}
}

func (s *resetPasswordReqStore) GetByClientID(ctx context.Context, id int) (*model.ResetPasswordReq, error) {
	row := s.pgcli.DB().QueryRowContext(
		ctx,
		"SELECT client_id, token, expires_at FROM reset_passwords WHERE client_id = $1",
		id,
	)

	var m model.ResetPasswordReq
	err := row.Scan(&m.ClientID, &m.Token, &m.ExpiresAt)
	if err != nil {
		return nil, newCoreError("scan rows", err)
	}
	return &m, nil
}

func (s *resetPasswordReqStore) GetByToken(ctx context.Context, token string) (*model.ResetPasswordReq, error) {
	row := s.pgcli.DB().QueryRowContext(
		ctx,
		"SELECT client_id, token, expires_at FROM reset_passwords WHERE token = $1",
		token,
	)

	var m model.ResetPasswordReq
	err := row.Scan(&m.ClientID, &m.Token, &m.ExpiresAt)
	if err != nil {
		return nil, newCoreError("scan rows", err)
	}
	return &m, nil
}

func (s *resetPasswordReqStore) Create(ctx context.Context, m *model.ResetPasswordReq) error {
	_, err := s.pgcli.DB().ExecContext(
		ctx,
		"INSERT INTO reset_passwords(client_id, token, expires_at) VALUES ($1,$2,$3)",
		m.ClientID,
		m.Token,
		m.ExpiresAt,
	)
	if err != nil {
		return newCoreError("create reset password", err)
	}
	return nil
}

func (s *resetPasswordReqStore) DeleteByClientID(ctx context.Context, db pg.DB, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM reset_passwords WHERE client_id = $1", id)
	return newCoreError("delete", err)
}
