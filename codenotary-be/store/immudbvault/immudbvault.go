package immudbvault

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tomaszkoziara/codenotarybe/store"
)

const (
	timeout    = 5 * time.Second
	retryCount = 3

	vaultBaseURLTemplate = "https://vault.immudb.io/ics/api/v1/ledger/%s/collection/%s"
	storeDocumentPath    = "/document"
	searchDocumentsPath  = "/documents/search"
)

type immudbVault struct {
	vaultBaseURL string
	apiKey       string
	client       *resty.Client
}

type createDocRespBody struct {
	TransactionID string `json:"transactionId"`
	DocumentID    string `json:"documentId"`
}

type searchDocRespBody struct {
	Revisions []revision `json:"revisions"`
}

type revision struct {
	Document store.AccountingInfo `json:"document"`
}

type searchDocReqBody struct {
	Query   query `json:"query"`
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
}

type query struct {
	Expressions []expression `json:"expressions"`
}

type expression struct {
	FieldComparisons []fieldComparison `json:"fieldComparisons"`
}

type fieldComparison struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func (im *immudbVault) StoreAccountingInfo(ctx context.Context, accInfo store.AccountingInfo) (string, error) {
	createDocRespBody := new(createDocRespBody)
	resp, err := im.client.
		R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("X-API-Key", im.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(accInfo).
		SetResult(createDocRespBody).
		Put(im.vaultBaseURL + storeDocumentPath)
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", fmt.Errorf("vault returned code: %v, resp: %v", resp.StatusCode(), string(resp.Body()))
	}

	return createDocRespBody.DocumentID, nil
}
func (im *immudbVault) GetAccountingInfo(ctx context.Context, id string) (store.AccountingInfo, error) {
	panic("not implemented")
}
func (im *immudbVault) ListAccountingInfo(ctx context.Context, accountName string, page int, pageSize int) ([]store.AccountingInfo, error) {
	searchDocRespBody := new(searchDocRespBody)
	resp, err := im.client.
		R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("X-API-Key", im.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(&searchDocReqBody{
			Query: query{
				Expressions: []expression{
					{
						FieldComparisons: []fieldComparison{
							{
								Field:    "accountName",
								Operator: "EQ",
								Value:    accountName,
							},
						},
					},
				},
			},
			Page:    page,
			PerPage: pageSize,
		}).
		SetResult(searchDocRespBody).
		Post(im.vaultBaseURL + searchDocumentsPath)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("vault returned code: %v, resp: %v", resp.StatusCode(), string(resp.Body()))
	}

	docs := make([]store.AccountingInfo, 0, len(searchDocRespBody.Revisions))
	for _, doc := range searchDocRespBody.Revisions {
		docs = append(docs, doc.Document)
	}

	return docs, nil
}

func New(ledger string, collection string, apiKey string) *immudbVault {
	client := resty.New()
	client.SetTimeout(timeout)
	client.SetRetryCount(retryCount)

	return &immudbVault{
		// for simplicity I'm not considering that I could connect to different ledgers/collection
		vaultBaseURL: fmt.Sprintf(vaultBaseURLTemplate, ledger, collection),
		client:       client,
		apiKey:       apiKey,
	}
}
