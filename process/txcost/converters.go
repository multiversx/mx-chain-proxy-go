package txcost

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

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

func (tcp *transactionCostProcessor) convertExtendedSCRInProtocolSCR(scr *data.ExtendedApiSmartContractResult) (*smartContractResult.SmartContractResult, error) {
	rcvAddr, err := tcp.pubKeyConverter.Decode(scr.RcvAddr)
	if err != nil {
		return nil, err
	}
	sndAddr, err := tcp.pubKeyConverter.Decode(scr.SndAddr)
	if err != nil {
		return nil, err
	}
	originalTxHashDecoded, err := hex.DecodeString(scr.OriginalTxHash)
	if err != nil {
		return nil, err
	}
	prevTxHashDecoded, err := hex.DecodeString(scr.PrevTxHash)
	if err != nil {
		return nil, err
	}

	var originalSender []byte
	if len(scr.OriginalSender) > 0 {
		originalSender, err = tcp.pubKeyConverter.Decode(scr.OriginalSender)
		if err != nil {
			return nil, err
		}
	}
	var relayer []byte
	if len(scr.RelayerAddr) > 0 {
		relayer, err = tcp.pubKeyConverter.Decode(scr.RelayerAddr)
		if err != nil {
			return nil, err
		}
	}

	return &smartContractResult.SmartContractResult{
		Nonce:          scr.Nonce,
		Value:          scr.Value,
		RcvAddr:        rcvAddr,
		SndAddr:        sndAddr,
		RelayerAddr:    relayer,
		RelayedValue:   scr.RelayedValue,
		Code:           []byte(scr.Code),
		Data:           []byte(scr.Data),
		PrevTxHash:     prevTxHashDecoded,
		OriginalTxHash: originalTxHashDecoded,
		GasLimit:       scr.GasLimit,
		GasPrice:       scr.GasPrice,
		CallType:       scr.CallType,
		CodeMetadata:   []byte(scr.CodeMetadata),
		ReturnMessage:  []byte(scr.ReturnMessage),
		OriginalSender: originalSender,
	}, nil
}
