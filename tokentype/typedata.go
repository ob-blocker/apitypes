package tokentype

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type TokenInfo struct {
	DomainTypehash  string
	DomainSeparator string
	PermitTypehash  string
	Name            string
	Version         string
	Nonce           *big.Int
	Token           string
}

func (c *Client) TokenTypedData(req *PermitMessage) (*TypedData, error) {
	contract, err := c.NewContract(req.Token)
	if err != nil || !contract.HasMethod("permit") {
		return nil, fmt.Errorf("contract has no permit")
	}

	info := &TokenInfo{
		Token: contract.Address(),
	}
	info.DomainTypehash, _ = contract.DOMAIN_TYPEHASH()
	info.Name, _ = contract.Name()
	info.Version, _ = contract.Version()

	// build eip712domain info
	base_domain := apitypes.TypedDataDomain{
		Name:              info.Name,
		ChainId:           (*math.HexOrDecimal256)(c.chainId),
		Version:           info.Version,
		VerifyingContract: info.Token,
		Salt:              hexutil.Encode(common.LeftPadBytes(c.chainId.Bytes(), 32)),
	}
	domain_type, ok := DomainTypedDataMap[info.DomainTypehash]
	if !ok {
		var err1 error
		info.DomainSeparator, err1 = contract.DOMAIN_SEPARATOR()
		if err1 != nil {
			return nil, err1
		}
		c.logger.Debugf("guessDomainInfoFromSeparator chain: %v , token: %v , info: %+v", c.chainId, info.Token, info)
		for _, guess_type := range DomainTypedDataMap {
			guess_domain_sep := getDomainSeparator(guess_type, base_domain)
			if guess_domain_sep == info.DomainSeparator {
				domain_type = guess_type
				break
			}
		}
	}
	if len(domain_type) == 0 {
		return nil, fmt.Errorf("build eip712domain failed")
	}

	// build permit info
	info.PermitTypehash, _ = contract.PERMIT_TYPEHASH()
	info.Nonce, _ = contract.Nonces(req.Owner)
	var permit_message apitypes.TypedDataMessage
	var permit_type []apitypes.Type
	if info.PermitTypehash == PERMIT_DAI_LIKE {
		permit_type = PERMIT_DAI_LIKE_TYPES
		permit_message = apitypes.TypedDataMessage{
			"holder":  req.Owner,
			"spender": req.Spender,
			"nonce":   info.Nonce,
			"expiry":  req.Deadline,
			"allowed": true,
		}
	} else { // use PERMIT_COMMON
		permit_type = PERMIT_COMMON_TYPES
		permit_message = apitypes.TypedDataMessage{
			"owner":    req.Owner,
			"spender":  req.Spender,
			"value":    req.Value,
			"nonce":    info.Nonce,
			"deadline": req.Deadline,
		}
	}

	ret := APITypedData{TypedData: apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": domain_type,
			"Permit":       permit_type,
		},
		PrimaryType: "Permit",
		Domain:      base_domain,
		Message:     permit_message,
	}}
	return ret.ConvertTypedData(), nil
}

func getDomainSeparator(domain_type []apitypes.Type, standard_domain apitypes.TypedDataDomain) string {
	domain := apitypes.TypedDataDomain{}
	for _, field := range domain_type {
		switch field.Name {
		case "name":
			domain.Name = standard_domain.Name
		case "chainId":
			domain.ChainId = standard_domain.ChainId
		case "version":
			domain.Version = standard_domain.Version
		case "verifyingContract":
			domain.VerifyingContract = standard_domain.VerifyingContract
		case "salt":
			domain.Salt = standard_domain.Salt
		}
	}
	td := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": domain_type,
		},
		Domain: domain,
	}
	sep, err := td.HashStruct("EIP712Domain", domain.Map())
	if err != nil {
		return ""
	}
	return hexutil.Encode(sep)
}
