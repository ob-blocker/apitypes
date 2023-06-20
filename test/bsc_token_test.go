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

func TestBuildTokenTypeBSC(t *testing.T) {
	c, err := tokentype.NewClient(os.Getenv("httpRpcBSC"))
	c.WithAPIKey(os.Getenv("apiKeyBSC"))
	assert.Nil(t, err)

	for _, token := range []string{
		"0x55d398326f99059ff775485246999027b3197955",
		"0x1af3f329e8be154074d8769d1ffa4ee058b1dbc3",
		"0x7ddee176f665cd201f93eede625770e2fd911990",
		"0x4b0f1812e5df2a09796481ff14017e6005508003",
		"0x4691937a7508860f876c9c0a2a617e7d9e945d4b",
		"0xc748673057861a797275cd8a068abb95a902e8de",
		"0xb020805e0bc7f0e353d1343d67a239f417d57bbf",
		"0xf307910a4c7bbc79691fd374889b36d8531b08e3",
		"0xd41fdb03ba84762dd66a0af1a6c8540ff1ba5dfb",
		"0xc1fdbed7dac39cae2ccc0748f7a80dc446f6a594",
		"0xb0d502e938ed5f4df2e681fe6e419ff29631d62b",
		"0xa4080f1778e69467e905b8d6f72f6e441f9e9484",
		"0x0a7e7d210c45c4abba183c1d0551b53ad1756eca",
		"0x8fff93e810a2edaafc326edee51071da9d398e83",
		"0xa1faa113cbe53436df28ff0aee54275c13b40975",
		"0xfe56d5892bdffc7bf58f2e84be1b2c32d21c308b",
		"0xfe19f0b51438fd612f6fd59c1dbb3ea319f433ba",
		"0x762539b45a1dcce3d36d080f74d1aed37844b878",
		"0x039cb485212f996a9dbb85a9a75d898f94d38da6",
		"0x98f8669f6481ebb341b522fcd3663f79a3d1a6a7",
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
