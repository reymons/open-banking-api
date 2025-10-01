package store

import (
	"banking/core/model"
	"banking/db/pg"
	"context"
)

type Client interface {
	GetByEmail(ctx context.Context, email string) (*model.Client, error)
	Create(ctx context.Context, cli *model.Client) error
}

type client struct {
	pgcli *pg.Client
}

func NewClient(cli *pg.Client) Client {
	return &client{cli}
}

func (s *client) GetByEmail(ctx context.Context, email string) (*model.Client, error) {
	row := s.pgcli.DB().QueryRowContext(
		ctx,
		"SELECT id, role, first_name, last_name, birth_date, email, password, is_partner, created_at FROM clients WHERE email = $1",
		email,
	)
	var c model.Client
	err := row.Scan(
		&c.ID,
		&c.Role,
		&c.FirstName,
		&c.LastName,
		&c.BirthDate,
		&c.Email,
		&c.Password,
		&c.IsPartner,
		&c.CreatedAt,
	)
	if err != nil {
		return nil, newCoreError("scan client row", err)
	}
	return &c, nil
}

func (s *client) Create(ctx context.Context, cli *model.Client) error {
	row := s.pgcli.DB().QueryRowContext(
		ctx,
		"INSERT INTO clients(role, first_name, last_name, birth_date, email, password, is_partner) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at",
		cli.Role,
		cli.FirstName,
		cli.LastName,
		cli.BirthDate,
		cli.Email,
		cli.Password,
		cli.IsPartner,
	)
	if err := row.Scan(&cli.ID, &cli.CreatedAt); err != nil {
		return newCoreError("scan client row", err)
	}
	return nil
}

func (s *client) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	row := s.pgcli.DB().QueryRowContext(ctx, "SELECT count(1) FROM clients WHERE email = $1", email)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, newCoreError("scan client row", err)
	}
	return count > 0, nil
}
