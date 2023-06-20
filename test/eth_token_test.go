package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ob-blocker/apitypes/tokentype"
	"github.com/stretchr/testify/assert"
)

func TestBuildTokenTypeETH(t *testing.T) {
	c, err := tokentype.NewClient(os.Getenv("httpRpc"))
	c.WithAPIKey(os.Getenv("apiKey"))
	assert.Nil(t, err)

	for _, token := range []string{
		"0x6b175474e89094c44da98b954eedeac495271d0f",
		"0x4fabb145d64652a948d72533023f6e7a623c7c53",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"0x0000000000085d4780b73119b644ae5ecd22b376",
		"0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		"0x853d955acef822db058eb8505911ed77f175b99e",
		"0x4d224452801aced8b2f0aebe155379bb5d594381",
		"0x19de6b897ed14a376dda0fe53a5420d2ac828a28",
		"0xf57e7e7c23978c3caec3c3548e3d615c346e79ff",
		"0x0f5d2fb29fb7d3cfee444a200298f468908cc942",
		"0x1a4b46696b2bb4794eb3d4c26f1c55f9170fa4c5",
		"0x056fd409e1d7a124bd7017459dfea2f387b6d5cd",
		"0xe66747a101bff2dba3697199dcce5b743b454759",
		"0xe28b3b32b6c345a34ff64674606124dd5aceca30",
		"0x39aa39c021dfbae8fac545936693ac917d5e7563",
		"0x3432b6a60d23ca0dfca7761b7ab56459d9c964d0",
		"0x3506424f91fd33084466f402d5d97f05f8e3b4af",
		"0xfcf8eda095e37a41e002e266daad7efc1579bc0a",
		"0xf629cbd94d3791c9250152bd8dfbdf380e2a3b9c",
		"0x5d3a536e4d6dbd6114cc1ead35777bab948e3643",
	} {
		data, err := c.TokenTypedData(&tokentype.PermitMessage{
			Token:    token,
			Owner:    "0xD1cc56810a3947d1D8b05448afB9889c6cFCF0F1",
			Spender:  "0x56B71565F6e7f9dE4c3217A6E5d4133bc7fc67EB",
			Value:    "1000000",
			Deadline: 1686803680,
		})
		if err != nil {
			fmt.Printf("Token: %v Got error: %v \n", token, err.Error())
			continue
		}
		body, _ := json.Marshal(data)
		var str bytes.Buffer
		_ = json.Indent(&str, body, "", "    ")
		fmt.Printf("Token: %v EIP712 TypedData: %v\n", token, str.String())
	}
}
