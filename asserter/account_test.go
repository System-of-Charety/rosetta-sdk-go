// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asserter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"

	"github.com/stretchr/testify/assert"
)

func TestContainsCurrency(t *testing.T) {
	var tests = map[string]struct {
		currencies []*types.Currency
		currency   *types.Currency
		contains   bool
	}{
		"simple contains": {
			currencies: []*types.Currency{
				{
					Symbol:   "BTC",
					Decimals: 8,
				},
			},
			currency: &types.Currency{
				Symbol:   "BTC",
				Decimals: 8,
			},
			contains: true,
		},
		"complex contains": {
			currencies: []*types.Currency{
				{
					Symbol:   "BTC",
					Decimals: 8,
					Metadata: &map[string]interface{}{
						"blah": "hello",
					},
				},
			},
			currency: &types.Currency{
				Symbol:   "BTC",
				Decimals: 8,
				Metadata: &map[string]interface{}{
					"blah": "hello",
				},
			},
			contains: true,
		},
		"empty": {
			currencies: []*types.Currency{},
			currency: &types.Currency{
				Symbol:   "BTC",
				Decimals: 8,
			},
			contains: false,
		},
		"symbol mismatch": {
			currencies: []*types.Currency{
				{
					Symbol:   "ERX",
					Decimals: 8,
				},
			},
			currency: &types.Currency{
				Symbol:   "BTC",
				Decimals: 6,
			},
			contains: false,
		},
		"decimal mismatch": {
			currencies: []*types.Currency{
				{
					Symbol:   "BTC",
					Decimals: 8,
				},
			},
			currency: &types.Currency{
				Symbol:   "BTC",
				Decimals: 6,
			},
			contains: false,
		},
		"metadata mismatch": {
			currencies: []*types.Currency{
				{
					Symbol:   "BTC",
					Decimals: 8,
					Metadata: &map[string]interface{}{
						"blah": "hello",
					},
				},
			},
			currency: &types.Currency{
				Symbol:   "BTC",
				Decimals: 8,
				Metadata: &map[string]interface{}{
					"blah": "bye",
				},
			},
			contains: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			exists := containsCurrency(test.currencies, test.currency)
			assert.Equal(t, test.contains, exists)
		})
	}
}

func TestContainsAccountIdentifier(t *testing.T) {
	var tests = map[string]struct {
		identifiers []*types.AccountIdentifier
		identifier  *types.AccountIdentifier
		contains    bool
	}{
		"simple contains": {
			identifiers: []*types.AccountIdentifier{
				{
					Address: "acct1",
				},
			},
			identifier: &types.AccountIdentifier{
				Address: "acct1",
			},
			contains: true,
		},
		"complex contains": {
			identifiers: []*types.AccountIdentifier{
				{
					Address: "acct1",
					SubAccount: &types.SubAccountIdentifier{
						Address: "subacct1",
						Metadata: &map[string]interface{}{
							"blah": "hello",
						},
					},
				},
			},
			identifier: &types.AccountIdentifier{
				Address: "acct1",
				SubAccount: &types.SubAccountIdentifier{
					Address: "subacct1",
					Metadata: &map[string]interface{}{
						"blah": "hello",
					},
				},
			},
			contains: true,
		},
		"simple mismatch": {
			identifiers: []*types.AccountIdentifier{
				{
					Address: "acct1",
				},
			},
			identifier: &types.AccountIdentifier{
				Address: "acct2",
			},
			contains: false,
		},
		"empty": {
			identifiers: []*types.AccountIdentifier{},
			identifier: &types.AccountIdentifier{
				Address: "acct2",
			},
			contains: false,
		},
		"subaccount mismatch": {
			identifiers: []*types.AccountIdentifier{
				{
					Address: "acct1",
					SubAccount: &types.SubAccountIdentifier{
						Address: "subacct2",
						Metadata: &map[string]interface{}{
							"blah": "hello",
						},
					},
				},
			},
			identifier: &types.AccountIdentifier{
				Address: "acct1",
				SubAccount: &types.SubAccountIdentifier{
					Address: "subacct1",
					Metadata: &map[string]interface{}{
						"blah": "hello",
					},
				},
			},
			contains: false,
		},
		"metadata mismatch": {
			identifiers: []*types.AccountIdentifier{
				{
					Address: "acct1",
					SubAccount: &types.SubAccountIdentifier{
						Address: "subacct1",
						Metadata: &map[string]interface{}{
							"blah": "hello",
						},
					},
				},
			},
			identifier: &types.AccountIdentifier{
				Address: "acct1",
				SubAccount: &types.SubAccountIdentifier{
					Address: "subacct1",
					Metadata: &map[string]interface{}{
						"blah": "bye",
					},
				},
			},
			contains: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			exists := containsAccountIdentifier(test.identifiers, test.identifier)
			assert.Equal(t, test.contains, exists)
		})
	}
}

func TestAccoutBalance(t *testing.T) {
	validBlock := &types.BlockIdentifier{
		Index: 1000,
		Hash:  "jsakdl",
	}

	invalidBlock := &types.BlockIdentifier{
		Index: 1,
		Hash:  "",
	}

	validIdentifier := &types.AccountIdentifier{
		Address: "acct1",
	}

	invalidIdentifier := &types.AccountIdentifier{
		Address: "",
	}

	validAmount := &types.Amount{
		Value: "100",
		Currency: &types.Currency{
			Symbol:   "BTC",
			Decimals: 8,
		},
	}

	var tests = map[string]struct {
		block    *types.BlockIdentifier
		balances []*types.Balance
		err      error
	}{
		"simple balance": {
			block: validBlock,
			balances: []*types.Balance{
				{
					AccountIdentifier: validIdentifier,
					Amounts: []*types.Amount{
						validAmount,
					},
				},
			},
			err: nil,
		},
		"invalid block": {
			block: invalidBlock,
			balances: []*types.Balance{
				{
					AccountIdentifier: validIdentifier,
					Amounts: []*types.Amount{
						validAmount,
					},
				},
			},
			err: errors.New("BlockIdentifier.Hash is missing"),
		},
		"invalid account identifier": {
			block: validBlock,
			balances: []*types.Balance{
				{
					AccountIdentifier: invalidIdentifier,
					Amounts: []*types.Amount{
						validAmount,
					},
				},
			},
			err: errors.New("Account.Address is missing"),
		},
		"duplicate currency": {
			block: validBlock,
			balances: []*types.Balance{
				{
					AccountIdentifier: validIdentifier,
					Amounts: []*types.Amount{
						validAmount,
						validAmount,
					},
				},
			},
			err: fmt.Errorf("currency %+v used in balance multiple times", validAmount.Currency),
		},
		"duplicate identifier": {
			block: validBlock,
			balances: []*types.Balance{
				{
					AccountIdentifier: validIdentifier,
					Amounts: []*types.Amount{
						validAmount,
					},
				},
				{
					AccountIdentifier: validIdentifier,
					Amounts: []*types.Amount{
						validAmount,
					},
				},
			},
			err: fmt.Errorf("account identifier %+v used in balance multiple times", validIdentifier),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := AccountBalance(test.block, test.balances)
			assert.Equal(t, test.err, err)
		})
	}
}