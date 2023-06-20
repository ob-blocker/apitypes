package tokentype

import "github.com/ethereum/go-ethereum/signer/core/apitypes"

type TypedData struct {
	Types       apitypes.Types            `json:"types"`
	PrimaryType string                    `json:"primaryType"`
	Domain      apitypes.TypedDataMessage `json:"domain"`
	Message     apitypes.TypedDataMessage `json:"message"`
}

type APITypedData struct {
	apitypes.TypedData
}

func (atd *APITypedData) ConvertTypedData() *TypedData {
	domain := apitypes.TypedDataMessage{}
	for _, field := range atd.Types["EIP712Domain"] {
		switch field.Name {
		case "name":
			domain["name"] = atd.Domain.Name
		case "chainId":
			domain["chainId"] = atd.Domain.ChainId
		case "version":
			domain["version"] = atd.Domain.Version
		case "verifyingContract":
			domain["verifyingContract"] = atd.Domain.VerifyingContract
		case "salt":
			domain["salt"] = atd.Domain.Salt
		}
	}
	return &TypedData{
		Types:       atd.Types,
		PrimaryType: atd.PrimaryType,
		Message:     atd.Message,
		Domain:      domain,
	}
}

type PermitMessage struct {
	Token    string
	Owner    string
	Spender  string
	Value    string
	Deadline int64
}
