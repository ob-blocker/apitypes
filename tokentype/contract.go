package tokentype

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ob-blocker/apitypes/utils"
)

const DEFAULT_ABI = "[{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"_nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

type Contract struct {
	address  common.Address
	abi      *abi.ABI
	contract *bind.BoundContract
}

type ContractABIResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (client *Client) requestABI(url string, contractAddress string) (*abi.ABI, error) {
	var resp ContractABIResp
	httpUrl := fmt.Sprintf("%s?module=contract&action=getabi&address=%v&apikey=%v", url, contractAddress, client.apiKey)
	if err := utils.HTTPGetJSON(httpUrl, &resp); err != nil {
		return nil, fmt.Errorf("request abi: %v failed: %v", httpUrl, err.Error())
	}

	if resp.Status == "1" && len(resp.Result) != 0 {
		contractAbi, abiErr := abi.JSON(strings.NewReader(resp.Result))
		if abiErr != nil {
			return nil, fmt.Errorf("get abi failed: %v", abiErr.Error())
		}
		return &contractAbi, nil
	}
	return nil, fmt.Errorf("request abi failed: %v", resp.Message)
}

func (client *Client) getDefaultContract(address common.Address) (*Contract, error) {
	contractAbi, err := abi.JSON(strings.NewReader(DEFAULT_ABI))
	if err != nil {
		return nil, err
	}
	return &Contract{
		address:  address,
		abi:      &contractAbi,
		contract: bind.NewBoundContract(address, contractAbi, client.ethClient, nil, nil),
	}, nil
}

func (client *Client) NewContract(address string) (*Contract, error) {
	addr := common.HexToAddress(address)
	url, ok := APIMap[client.chainId.String()]
	if !ok {
		return client.getDefaultContract(addr)
	}

	contractAbi, err := client.requestABI(url, address)
	if err != nil {
		client.logger.WithError(err).Error("requestGetABI")
		return client.getDefaultContract(addr)
	}

	contract := &Contract{
		address:  addr,
		abi:      contractAbi,
		contract: bind.NewBoundContract(addr, *contractAbi, client.ethClient, nil, nil),
	}

	// check is proxy
	if !contract.HasMethod("implementation") {
		return contract, nil
	}

	// get implement address
	impl, err := contract.Implementation()
	if err != nil {
		client.logger.WithError(err).Error("get Implementation")
		return contract, err
	}

	// fetch implement abi
	implAbi, err := client.requestABI(url, impl.String())
	if err != nil {
		client.logger.WithError(err).Error("requestImplABI")
		return contract, err
	}

	return &Contract{
		address:  addr,
		abi:      implAbi,
		contract: bind.NewBoundContract(addr, *implAbi, client.ethClient, nil, nil),
	}, nil
}

func (c *Contract) Address() string {
	return c.address.String()
}

func (c *Contract) HasMethod(method string) bool {
	_, found := c.abi.Methods[method]
	return found
}

func (c *Contract) Implementation() (common.Address, error) {
	var out []interface{}
	err := c.contract.Call(nil, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

func (c *Contract) getTYPEHASH(method string) (string, error) {
	if !c.HasMethod(method) {
		return "", fmt.Errorf("contract has no method %v", method)
	}
	var out []interface{}
	err := c.contract.Call(nil, &out, method)

	if err != nil {
		return "", err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return hexutil.Encode(out0[:]), nil
}

func (c *Contract) DOMAIN_SEPARATOR() (string, error) {
	for _, method := range []string{
		"DOMAIN_SEPARATOR",
		"getDomainSeperator",
	} {
		sep, err := c.getTYPEHASH(method)
		if err == nil {
			return sep, err
		}
	}
	return "", fmt.Errorf("can not get DOMAIN_SEPARATOR")
}

func (c *Contract) DOMAIN_TYPEHASH() (string, error) {
	return c.getTYPEHASH("DOMAIN_TYPEHASH")
}

func (c *Contract) PERMIT_TYPEHASH() (string, error) {
	return c.getTYPEHASH("PERMIT_TYPEHASH")
}

func (c *Contract) getNonces(method string, sender string) (*big.Int, error) {
	if !c.HasMethod(method) {
		return nil, fmt.Errorf("contract has no method %v", method)
	}
	var out []interface{}
	err := c.contract.Call(nil, &out, method, common.HexToAddress(sender))

	if err != nil {
		return *new(*big.Int), err
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}

func (c *Contract) Nonces(sender string) (*big.Int, error) {
	for _, method := range []string{
		"nonces",
		"_nonces", // eth.aaveToken
	} {
		n, err := c.getNonces(method, sender)
		if err == nil {
			return n, err
		}
	}
	return nil, fmt.Errorf("can not get Nonces")
}

func (c *Contract) getString(method string) (string, error) {
	if !c.HasMethod(method) {
		return "", fmt.Errorf("contract has no method %v", method)
	}

	var out []interface{}
	err := c.contract.Call(nil, &out, method)

	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

func (c *Contract) Name() (string, error) {
	return c.getString("name")
}

func (c *Contract) Decimals() (uint8, error) {
	var out []interface{}
	err := c.contract.Call(nil, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

func (c *Contract) Version() (string, error) {
	for _, method := range []string{
		"version",
		"EIP712_VERSION", // polygon.USDC
	} {
		r, err := c.getString(method)
		if err == nil {
			return r, err
		}
	}
	return "", fmt.Errorf("can not get Version")
}
