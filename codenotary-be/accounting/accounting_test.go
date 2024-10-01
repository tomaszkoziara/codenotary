package accounting_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomaszkoziara/codenotarybe/accounting"
	"github.com/tomaszkoziara/codenotarybe/store"
	"github.com/tomaszkoziara/codenotarybe/store/immudbvaultmock"
)

func TestCreateAccountingInfoHappyPath(t *testing.T) {
	ctx := context.Background()
	accountingInfoTest := accounting.AccountingInfo{
		AccountNumber: "accounting-number",
		AccountName:   "account-name",
		IBAN:          "iban",
		Address:       "address",
		Amount:        100.99,
		Type:          accounting.TypeReceiving,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := immudbvaultmock.NewMockStore(ctrl)
	db.EXPECT().StoreAccountingInfo(ctx, gomock.Eq(store.AccountingInfo{
		AccountNumber: "accounting-number",
		AccountName:   "account-name",
		IBAN:          "iban",
		Address:       "address",
		Amount:        100.99,
		Type:          store.TypeReceiving,
	})).Return("123", nil)
	accountinService := accounting.New(db)

	_, err := accountinService.StoreAccountingInfo(ctx, accountingInfoTest)
	assert.NoError(t, err)
}

func TestCreateAccountingInfoValidationErrors(t *testing.T) {
	testCases := []struct {
		name        string
		input       accounting.AccountingInfo
		expectedMsg string
	}{
		{
			name: "missing account number",
			input: accounting.AccountingInfo{
				AccountName: "account-name",
				IBAN:        "iban",
				Address:     "address",
				Amount:      100.99,
				Type:        accounting.TypeReceiving,
			},
			expectedMsg: "Field validation for 'AccountNumber' failed on the 'required' tag",
		},
		{
			name: "missing account name",
			input: accounting.AccountingInfo{
				AccountNumber: "account-number",
				IBAN:          "iban",
				Address:       "address",
				Amount:        100.99,
				Type:          accounting.TypeReceiving,
			},
			expectedMsg: "Field validation for 'AccountName' failed on the 'required' tag",
		},
		{
			name: "missing iban",
			input: accounting.AccountingInfo{
				AccountNumber: "account-number",
				AccountName:   "account-name",
				Address:       "address",
				Amount:        100.99,
				Type:          accounting.TypeReceiving,
			},
			expectedMsg: "Field validation for 'IBAN' failed on the 'required' tag",
		},
		{
			name: "missing address",
			input: accounting.AccountingInfo{
				AccountNumber: "account-number",
				AccountName:   "account-name",
				IBAN:          "iban",
				Amount:        100.99,
				Type:          accounting.TypeReceiving,
			},
			expectedMsg: "Field validation for 'Address' failed on the 'required' tag",
		},
		{
			name: "missing amount",
			input: accounting.AccountingInfo{
				AccountNumber: "account-number",
				AccountName:   "account-name",
				IBAN:          "iban",
				Address:       "address",
				Type:          accounting.TypeReceiving,
			},
			expectedMsg: "Field validation for 'Amount' failed on the 'required' tag",
		},
		{
			name: "missing type",
			input: accounting.AccountingInfo{
				AccountNumber: "account-number",
				AccountName:   "account-name",
				IBAN:          "iban",
				Address:       "address",
				Amount:        100.99,
			},
			expectedMsg: "Field validation for 'Type' failed on the 'required' tag",
		},
		{
			name: "wrong type",
			input: accounting.AccountingInfo{
				AccountNumber: "account-number",
				AccountName:   "account-name",
				IBAN:          "iban",
				Address:       "address",
				Amount:        100.99,
				Type:          "some other type",
			},
			expectedMsg: "Field validation for 'Type' failed on the 'oneof' tag",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := immudbvaultmock.NewMockStore(ctrl)
			accountinService := accounting.New(db)

			_, err := accountinService.StoreAccountingInfo(context.Background(), tc.input)
			assert.Error(t, err)
			assert.Containsf(t, err.Error(), tc.expectedMsg, "unexpected error: %s")
		})
	}
}

func TestListAccountingInfoHappyPath(t *testing.T) {
	ctx := context.Background()
	testAccountingInfoList := []store.AccountingInfo{
		{
			AccountNumber: "1",
			AccountName:   "mario",
			IBAN:          "iban1",
			Address:       "address1",
			Amount:        10,
			Type:          accounting.TypeReceiving,
		},
		{
			AccountNumber: "2",
			AccountName:   "luigi",
			IBAN:          "iban2",
			Address:       "address2",
			Amount:        20,
			Type:          accounting.TypeSending,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := immudbvaultmock.NewMockStore(ctrl)
	db.EXPECT().ListAccountingInfo(ctx, "test-account", 5, 10).
		Return(testAccountingInfoList, nil)
	accountinService := accounting.New(db)

	list, err := accountinService.ListAccountingInfo(ctx, "test-account", 5, 10)
	assert.NoError(t, err)
	assert.Len(t, list, len(testAccountingInfoList))

	for i, it := range testAccountingInfoList {
		assert.Equal(t, it.AccountNumber, list[i].AccountNumber)
		assert.Equal(t, it.AccountName, list[i].AccountName)
		assert.Equal(t, it.IBAN, list[i].IBAN)
		assert.Equal(t, it.Address, list[i].Address)
		assert.Equal(t, it.Amount, list[i].Amount)
		assert.Equal(t, string(it.Type), string(list[i].Type))
	}
}

// TODO: add more test for errors
