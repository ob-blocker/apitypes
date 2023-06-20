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

func TestBuildTokenTypePolygon(t *testing.T) {
	c, err := tokentype.NewClient(os.Getenv("httpRpcPolygon"))
	c.WithAPIKey(os.Getenv("apiKeyPolygon"))
	assert.Nil(t, err)

	for _, token := range []string{
		"0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
		"0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		"0x8f3cf7ad23cd3cadbd9735aff958023239c6a063",
		"0xdab529f40e671a1d4bf91361c21bf9f0c9712ab7",
		"0x1bfd67037b42cf73acf2047067bd4f2c47d9bfd6",
		"0x2c89bbc92bd86f8075d1decc58c7f4e0107f286b",
		"0xb33eaad8d922b1083446dc23f610c2567fb5180f",
		"0x2e1ad108ff1d8c782fcbbb89aad783ac49586756",
		"0x53e0bca35ec356bd5dddfebbd1fc0fd03fabad39",
		"0xc3c7d422809852031b44ab29eec9f1eff2a58756",
		"0x45c32fa6df82ead1e2ef74d17b76547eddfaff89",
		"0x5fe2b58c013d7601147dcdd68c143a77499f5531",
		"0x0266f4f08d82372cf0fcbccc0ff74309089c74d1",
		"0x61299774020da444af134c82fa83e3810b309991",
		"0xd6df932a45c0f255f85145f286ea0b292b21c90b",
		"0xbbba073c31bf03b8acf7c28ef0738decf3695683",
		"0xa1c57f48f0deb89f569dfbe6e2b7f46d33606fd4",
		"0x50b728d8d964fd00c2d0aad81718b71311fef68a",
		"0x172370d5cd63279efa6d502dab29171933a610af",
		"0xee327f889d5947c1dc1934bb208a1e792f953e96",
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
