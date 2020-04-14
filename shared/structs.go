package shared

// SCQuery represents a prepared query for executing a function of the smart contract
type SCQuery struct {
	ScAddress string
	FuncName  string
	Arguments [][]byte
}
