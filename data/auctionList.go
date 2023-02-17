package data

// AuctionNode holds data needed for a node in auction to respond to API calls
type AuctionNode struct {
	BlsKey    string `json:"blsKey"`
	Qualified bool   `json:"qualified"`
}

// AuctionListValidatorAPIResponse holds the data needed for an auction node validator for responding to API calls
type AuctionListValidatorAPIResponse struct {
	Owner          string         `json:"owner"`
	NumStakedNodes int64          `json:"numStakedNodes"`
	TotalTopUp     string         `json:"totalTopUp"`
	TopUpPerNode   string         `json:"topUpPerNode"`
	QualifiedTopUp string         `json:"qualifiedTopUp"`
	AuctionList    []*AuctionNode `json:"auctionList"`
}
