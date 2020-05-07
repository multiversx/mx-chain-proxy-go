package data

import "time"

// DatabaseTransaction is used for transaction route
type DatabaseTransaction struct {
	Hash          string        `json:"hash"`
	MBHash        string        `json:"miniBlockHash"`
	Nonce         uint64        `json:"nonce"`
	Round         uint64        `json:"round"`
	Value         string        `json:"value"`
	Receiver      string        `json:"receiver"`
	Sender        string        `json:"sender"`
	ReceiverShard uint32        `json:"receiverShard"`
	SenderShard   uint32        `json:"senderShard"`
	GasPrice      uint64        `json:"gasPrice"`
	GasLimit      uint64        `json:"gasLimit"`
	Data          string        `json:"data"`
	Signature     string        `json:"signature"`
	Timestamp     time.Duration `json:"timestamp"`
	Status        string        `json:"status"`
	Fee           string        `json:"fee"`
}
