package service

import (
	"banking/core"
	"banking/core/model"
	"banking/store"
	"context"
)

type Account interface {
	Request(ctx context.Context, currencyID int, userID int, role int) error
	GetAllByUserID(ctx context.Context, userID int, role int) ([]*model.Account, error)
}

type account struct {
	permService  Perm
	accountStore store.Account
}

func NewAccount(ps Perm, as store.Account) Account {
	return &account{ps, as}
}

func (s *account) Request(ctx context.Context, currencyID int, userID int, role int) error {
	if !s.permService.Can(ctx, role, core.PermRequestAccount) {
		return core.ErrInvalidAccess
	}

	acc := &model.Account{
		ClientID:   userID,
		CurrencyID: currencyID,
		Balance:    0,
		Status:     core.AccountStatusInactive,
	}

	return s.accountStore.Create(ctx, acc)
}

func (s *account) GetAllByUserID(ctx context.Context, userID int, role int) ([]*model.Account, error) {
	if !s.permService.Can(ctx, role, core.PermViewAccount) {
		return nil, core.ErrInvalidAccess
	}
	return s.accountStore.GetAllByUserID(ctx, userID)
}
