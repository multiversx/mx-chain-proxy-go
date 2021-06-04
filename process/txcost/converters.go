package txcost

import (
	"strings"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const atSep = "@"

func (tcp *transactionCostProcessor) computeShardID(addr string) (uint32, error) {
	senderBuff, err := tcp.pubKeyConverter.Decode(addr)
	if err != nil {
		return 0, err
	}

	shardID, err := tcp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return 0, err
	}

	return shardID, nil
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

func (tcp *transactionCostProcessor) maxGasLimitPerBlockBasedOnReceiverAddr(receiver string) uint64 {
	shardID, err := tcp.computeShardID(receiver)
	if err != nil {
		return tcp.maxGasLimitPerBlockShard - 1
	}

	if shardID == core.MetachainShardId {
		return tcp.maxGasLimitPerBlockMeta - 1
	}

	return tcp.maxGasLimitPerBlockShard - 1
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
