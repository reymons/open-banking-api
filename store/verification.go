package store

import (
	"banking/db/pg"
	"context"
	"time"
)

type verificationDto struct {
	Target    string
	Data      []byte
	Code      string
	ExpiresAt time.Time
}

type Verification interface {
	Create(ctx context.Context, db pg.DB, v *verificationDto) error
	Delete(ctx context.Context, db pg.DB, target string, code string) error
	GetByTarget(ctx context.Context, db pg.DB, target string) (*verificationDto, error)
}

type verification struct{}

func NewVerification() Verification {
	return &verification{}
}

func (*verification) Create(ctx context.Context, db pg.DB, v *verificationDto) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO verifications(target, data, code, expires_at) VALUES ($1, $2, $3, $4)",
		v.Target, v.Data, v.Code, v.ExpiresAt,
	)
	return newCoreError("insert", err)
}

func (*verification) Delete(ctx context.Context, db pg.DB, target string, code string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM verifications WHERE target = $1 AND code = $2", target, code)
	return newCoreError("delete", err)
}

func (*verification) GetByTarget(
	ctx context.Context,
	db pg.DB,
	target string,
) (*verificationDto, error) {
	row := db.QueryRowContext(
		ctx,
		"SELECT target, data, code, expires_at FROM verifications WHERE target = $1",
		target,
	)
	var v verificationDto
	err := row.Scan(&v.Target, &v.Data, &v.Code, &v.ExpiresAt)
	if err != nil {
		return nil, newCoreError("scan verification", err)
	}
	return &v, nil
}
