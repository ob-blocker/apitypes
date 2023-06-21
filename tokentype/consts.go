package tokentype

import "github.com/ethereum/go-ethereum/signer/core/apitypes"

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000000000000000000000000000"

// permit types
var (
	// Permit(address holder,address spender,uint256 nonce,uint256 expiry,bool allowed)
	PERMIT_DAI_LIKE = "0xea2aa0a1be11a07ed86d755c93467f4f82362b452371d1ba94d1715123511acb"
	// Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)
	PERMIT_COMMON = "0x6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c9"

	PERMIT_DAI_LIKE_TYPES = []apitypes.Type{
		{Name: "holder", Type: "address"},
		{Name: "spender", Type: "address"},
		{Name: "nonce", Type: "uint256"},
		{Name: "expiry", Type: "uint256"},
		{Name: "allowed", Type: "bool"},
	}
	PERMIT_COMMON_TYPES = []apitypes.Type{
		{Name: "owner", Type: "address"},
		{Name: "spender", Type: "address"},
		{Name: "value", Type: "uint256"},
		{Name: "nonce", Type: "uint256"},
		{Name: "deadline", Type: "uint256"},
	}
)

// domain types
var (
	// EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)
	DOMAIN_COMMON = "0x8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f"
	// EIP712Domain(string name,uint256 chainId,address verifyingContract)
	DOMAIN_UNI_LIKE = "0x8cad95687ba82c2ce50e74f7b754645e5117c3a5bec8151c0726d5857980a866"
	// EIP712Domain(uint256 chainId,address verifyingContract)
	DOMAIN_ONLY_CHAIN = "0x47e79534a245952e8b16893a336b85a3d9ea9fa8c573f3d803afb92a79469218"
	// EIP712Domain(string name,string version,address verifyingContract,bytes32 salt)
	DOMAIN_SALT_LIKE = "0x36c25de3e541d5d970f66e4210d728721220fff5c077cc6cd008b3a0c62adab7"

	DOMAIN_COMMON_TYPES = []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}
	DOMAIN_UNI_LIKE_TYPES = []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}
	DOMAIN_ONLY_CHAIN_TYPES = []apitypes.Type{
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}
	DOMAIN_SALT_LIKE_TYPES = []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "verifyingContract", Type: "address"},
		{Name: "salt", Type: "bytes32"},
	}
)

// domain maps
var DomainTypedDataMap = map[string][]apitypes.Type{
	DOMAIN_COMMON:     DOMAIN_COMMON_TYPES,
	DOMAIN_UNI_LIKE:   DOMAIN_UNI_LIKE_TYPES,
	DOMAIN_ONLY_CHAIN: DOMAIN_ONLY_CHAIN_TYPES,
	DOMAIN_SALT_LIKE:  DOMAIN_SALT_LIKE_TYPES,
}

// api maps
/*	for request ABI
	https://api.etherscan.io/api
	   ?module=contract
	   &action=getabi
	   &address=0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413
	   &apikey=YourApiKeyToken
*/
var APIMap = map[string]string{
	"1":     "https://api.etherscan.io/api",            // eth mainnet
	"5":     "https://api-goerli.etherscan.io/api",     // goerli
	"56":    "https://api.bscscan.com/api",             // bsc mainnet
	"97":    "https://api-testnet.bscscan.com/api",     // bsc testnet
	"137":   "https://api.polygonscan.com/api",         // polygon mainnet
	"80001": "https://api-testnet.polygonscan.com/api", // polygon testnet
}

var IMPLEMENTATION_SLOTS = []string{
	"0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc", // bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)
	"0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3", // keccak256("org.zeppelinos.proxy.implementation"))
}
