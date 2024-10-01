package accounting

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/tomaszkoziara/codenotarybe/store"
)

type Type string

const (
	TypeSending   = "sending"
	TypeReceiving = "receiving"
)

type AccountingInfo struct {
	AccountNumber string  `validate:"required"`
	AccountName   string  `validate:"required"`
	IBAN          string  `validate:"required"`
	Address       string  `validate:"required"`
	Amount        float64 `validate:"required"`
	Type          Type    `validate:"required,oneof=sending receiving"`
}

type Accounting struct {
	db store.Store
}

func (a *Accounting) StoreAccountingInfo(ctx context.Context, accountingInfo AccountingInfo) (string, error) {
	validate := validator.New()
	if err := validate.Struct(accountingInfo); err != nil {
		return "", err
	}

	ai := store.AccountingInfo{
		AccountNumber: accountingInfo.AccountNumber,
		AccountName:   accountingInfo.AccountName,
		IBAN:          accountingInfo.IBAN,
		Address:       accountingInfo.Address,
		Amount:        accountingInfo.Amount,
		Type:          store.Type(accountingInfo.Type),
	}
	id, err := a.db.StoreAccountingInfo(ctx, ai)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *Accounting) ListAccountingInfo(ctx context.Context, accountName string, page int, pageSize int) ([]AccountingInfo, error) {
	accountingInfoStoreList, err := a.db.ListAccountingInfo(ctx, accountName, page, pageSize)
	if err != nil {
		return nil, err
	}

	accountingInfoList := make([]AccountingInfo, 0, len(accountingInfoStoreList))
	for _, ai := range accountingInfoStoreList {
		accountingInfoList = append(accountingInfoList, AccountingInfo{
			AccountNumber: ai.AccountNumber,
			AccountName:   ai.AccountName,
			IBAN:          ai.IBAN,
			Address:       ai.Address,
			Amount:        ai.Amount,
			Type:          Type(ai.Type),
		})
	}

	return accountingInfoList, nil
}

func New(db store.Store) *Accounting {
	return &Accounting{
		db: db,
	}
}
