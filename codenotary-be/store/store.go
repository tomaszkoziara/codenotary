package store

import "context"

type Store interface {
	StoreAccountingInfo(ctx context.Context, accInfo AccountingInfo) (string, error)
	GetAccountingInfo(ctx context.Context, id string) (AccountingInfo, error)
	ListAccountingInfo(ctx context.Context, accountName string, page int, pageSize int) ([]AccountingInfo, error)
}

type Type string

const (
	TypeSending   = "sending"
	TypeReceiving = "receiving"
)

type AccountingInfo struct {
	AccountNumber string  `json:"accountNumber"`
	AccountName   string  `json:"accountName"`
	IBAN          string  `json:"iban"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Type          Type    `json:"type"`
}

type AccountingInfoVersion struct {
	Document AccountingInfo `json:"document"`
	Query    Query          `json:"query"`
}

type Query struct {
}
