package store

import (
	"banking/core/model"
	"banking/db/pg"
	"context"
)

type Account interface {
	Create(ctx context.Context, acc *model.Account) error
	GetAllByUserID(ctx context.Context, userID int) ([]*model.Account, error)
}

type account struct {
	pgcli *pg.Client
}

func NewAccount(cli *pg.Client) Account {
	return &account{cli}
}

func (s *account) Create(ctx context.Context, acc *model.Account) error {
	q := "WITH acc AS "
	q += "(INSERT INTO accounts(client_id, currency_id, balance, status) VALUES ($1,$2,$3,$4) RETURNING *) "
	q += "SELECT acc.id, acc.created_at, acc.number, curr.code "
	q += "FROM acc JOIN currencies curr ON acc.currency_id = curr.id"

	row := s.pgcli.DB().QueryRowContext(
		ctx, q,
		&acc.ClientID,
		&acc.CurrencyID,
		&acc.Balance,
		&acc.Status,
	)
	err := row.Scan(&acc.ID, &acc.CreatedAt, &acc.Number, &acc.CurrencyCode)
	return newCoreError("scan row", err)
}

func (s *account) GetAllByUserID(ctx context.Context, userID int) ([]*model.Account, error) {
	q := "SELECT a.id, a.client_id, a.currency_id, c.code, a.number, a.balance, a.status, a.created_at "
	q += "FROM accounts a JOIN currencies c ON a.currency_id = c.id "
	q += "WHERE a.client_id = $1"

	rows, err := s.pgcli.DB().QueryContext(ctx, q, userID)
	if err != nil {
		return nil, newCoreError("query rows", err)
	}
	defer rows.Close()

	accs := make([]*model.Account, 0)
	for rows.Next() {
		acc := &model.Account{}
		err := rows.Scan(
			&acc.ID,
			&acc.ClientID,
			&acc.CurrencyID,
			&acc.CurrencyCode,
			&acc.Number,
			&acc.Balance,
			&acc.Status,
			&acc.CreatedAt,
		)
		if err != nil {
			return nil, newCoreError("scan row", err)
		}
		accs = append(accs, acc)
	}
	return accs, nil
}
