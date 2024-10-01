package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/tomaszkoziara/codenotarybe/accounting"
)

const (
	TypeSending   = "sending"
	TypeReceiving = "receiving"
)

type Type string

type AccountingInfo struct {
	AccountNumber string  `json:"accountNumber"`
	AccountName   string  `json:"accountName"`
	IBAN          string  `json:"iban"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Type          Type    `json:"type"`
}

type AccountinInfoList []AccountingInfo

func CreateStoreAccountingInfo(accountingService *accounting.Accounting) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		var accountingInfo AccountingInfo
		err := json.NewDecoder(request.Body).Decode(&accountingInfo)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(fmt.Sprintf("couldn't decode request body: %v", err)))
			return
		}

		id, err := accountingService.StoreAccountingInfo(ctx, accounting.AccountingInfo{
			AccountNumber: accountingInfo.AccountNumber,
			AccountName:   accountingInfo.AccountName,
			IBAN:          accountingInfo.IBAN,
			Address:       accountingInfo.Address,
			Amount:        accountingInfo.Amount,
			Type:          accounting.Type(accountingInfo.Type),
		})
		if err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte(fmt.Sprintf("invalid request payload: %v", ve)))
				return
			}
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
			return
		}

		response.WriteHeader(http.StatusCreated)
		response.Write([]byte(fmt.Sprintf(`{"id": "%v"}`, id)))
		return
	}
}

func CreateGetListAccountingInfo(accountingService *accounting.Accounting) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		query := request.URL.Query()

		pageStr := query.Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(fmt.Sprintf("couldn't decode page: %v", err)))
			return
		}

		pageSizeStr := query.Get("pageSize")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(fmt.Sprintf("couldn't decode pageSize: %v", err)))
			return
		}

		accountName := query.Get("accountName")
		if accountName == "" {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("accountName is missing"))
			return
		}

		list, err := accountingService.ListAccountingInfo(ctx, accountName, page, pageSize)
		if err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte(fmt.Sprintf("invalid request payload: %v", ve)))
				return
			}
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
			return
		}

		accountingInfoList := make([]AccountingInfo, 0, len(list))
		for _, accountingInfo := range list {
			accountingInfoList = append(accountingInfoList, AccountingInfo{
				AccountNumber: accountingInfo.AccountNumber,
				AccountName:   accountingInfo.AccountName,
				IBAN:          accountingInfo.IBAN,
				Address:       accountingInfo.Address,
				Amount:        accountingInfo.Amount,
				Type:          Type(accountingInfo.Type),
			})
		}

		respBody, err := json.Marshal(accountingInfoList)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
			return
		}

		response.WriteHeader(http.StatusOK)
		response.Header().Add("Content-Type", "application/json")
		response.Write(respBody)
	}
}
