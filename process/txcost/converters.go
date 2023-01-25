package txcost

import (
	"strings"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

const argsSeparator = "@"

func (tcp *transactionCostProcessor) computeShardID(addr string) (uint32, error) {
	senderBuff, err := tcp.pubKeyConverter.Decode(addr)
	if err != nil {
		return 0, err
	}

	return tcp.proc.ComputeShardId(senderBuff)
}

func (tcp *transactionCostProcessor) computeSenderAndReceiverShardID(sender, receiver string) (uint32, uint32, error) {
	senderShardID, err := tcp.computeShardID(sender)
	if err != nil {
		return 0, 0, err
	}

	receiverShardID, err := tcp.computeShardID(receiver)
	if err != nil {
		return 0, 0, err
	}

	return senderShardID, receiverShardID, nil
}

func convertSCRInTransaction(scr *data.ExtendedApiSmartContractResult, originalTx *data.Transaction) *data.Transaction {
	newDataField := removeLatestArgumentFromDataField(scr.Data)

	return &data.Transaction{
		Nonce:     scr.Nonce,
		Value:     scr.Value.String(),
		Receiver:  scr.RcvAddr,
		Sender:    scr.SndAddr,
		GasPrice:  scr.GasPrice,
		GasLimit:  scr.GasLimit,
		Data:      []byte(newDataField),
		Signature: "",
		ChainID:   originalTx.ChainID,
		Version:   originalTx.Version,
		Options:   originalTx.Options,
	}
}

func removeLatestArgumentFromDataField(dataField string) string {
	splitDataField := strings.Split(dataField, argsSeparator)
	newStr := splitDataField[:len(splitDataField)-1]
	if len(newStr) == 0 {
		return dataField
	}

	newDataField := strings.Join(newStr, argsSeparator)

	return newDataField
}
