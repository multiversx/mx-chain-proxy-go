package data

const (
	FungibleTokens     = "fungible-tokens"
	SemiFungibleTokens = "semi-fungible-tokens"
	NonFungibleTokens  = "non-fungible-tokens"
)

// ValidTokenTypes holds a slice containing the valid esdt token types
var ValidTokenTypes = []string{FungibleTokens, SemiFungibleTokens, NonFungibleTokens}

// ESDTSupplyResponse is a response holding esdt supply
type ESDTSupplyResponse struct {
	Data  ESDTSupply `json:"data"`
	Error string     `json:"error"`
	Code  ReturnCode `json:"code"`
}

// ESDTSupply is a DTO holding esdt supply
type ESDTSupply struct {
	Supply           string `json:"supply"`
	Minted           string `json:"minted"`
	Burned           string `json:"burned"`
	InitialMinted    string `json:"initialMinted"`
	RecomputedSupply bool   `json:"recomputedSupply"`
}

// IsValidEsdtPath returns true if the provided path is a valid esdt token type
func IsValidEsdtPath(path string) bool {
	for _, tokenType := range ValidTokenTypes {
		if tokenType == path {
			return true
		}
	}

	return false
}
