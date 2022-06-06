package data

// BlockInfo defines the data structure for the block at which an resource (e.g. Account object) is fetched from the Network
type BlockInfo struct {
	Nonce    uint64 `json:"nonce"`
	Hash     []byte `json:"hash"`
	RootHash []byte `json:"rootHash"`
}
