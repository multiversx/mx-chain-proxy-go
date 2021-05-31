package txcost

import (
	"strings"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const atSep = "@"

func (tcp *transactionCostProcessor) computeSenderAndReceiverShardID(sender, receiver string) (uint32, uint32, error) {
	senderBuff, err := tcp.pubKeyConverter.Decode(sender)
	if err != nil {
		return 0, 0, err
	}

	senderShardID, err := tcp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return 0, 0, err
	}

	receiverBuff, err := tcp.pubKeyConverter.Decode(receiver)
	if err != nil {
		return 0, 0, err
	}

	receiverShardID, err := tcp.proc.ComputeShardId(receiverBuff)
	if err != nil {
		return 0, 0, err
	}

	return senderShardID, receiverShardID, nil
}

func convertSCRInTransaction(scr *data.ApiSmartContractResultExtended, originalTx *data.Transaction) *data.Transaction {
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
	splitDataField := strings.Split(dataField, atSep)
	newStr := splitDataField[:len(splitDataField)-1]
	newDataField := strings.Join(newStr, atSep)

	return newDataField
}
