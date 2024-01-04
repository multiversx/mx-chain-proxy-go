package process

import (
	"math/big"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const (
	esdtContractAddress   = "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u"
	initialESDTSupplyFunc = "getTokenProperties"

	networkESDTSupplyPath = "/network/esdt/supply/"
	zeroBigIntStr         = "0"
)

type esdtSupplyProcessor struct {
	baseProc    Processor
	scQueryProc SCQueryService
}

// NewESDTSupplyProcessor will create a new instance of the ESDT supply processor
func NewESDTSupplyProcessor(baseProc Processor, scQueryProc SCQueryService) (*esdtSupplyProcessor, error) {
	if check.IfNil(baseProc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(scQueryProc) {
		return nil, ErrNilSCQueryService
	}

	return &esdtSupplyProcessor{
		baseProc:    baseProc,
		scQueryProc: scQueryProc,
	}, nil
}

// GetESDTSupply will return the total supply for the provided token
func (esp *esdtSupplyProcessor) GetESDTSupply(tokenIdentifier string) (*data.ESDTSupplyResponse, error) {
	totalSupply, err := esp.getSupplyFromShards(tokenIdentifier)
	if err != nil {
		return nil, err
	}

	res := &data.ESDTSupplyResponse{
		Code: data.ReturnCodeSuccess,
	}
	if !isFungibleESDT(tokenIdentifier) {
		res.Data = *totalSupply
		makeInitialMintedNotEmpty(res)
		return res, nil
	}

	initialSupply, err := esp.getInitialSupplyFromMeta(tokenIdentifier)
	if err != nil {
		return nil, err
	}

	res.Data.InitialMinted = initialSupply.String()
	if totalSupply.RecomputedSupply {
		res.Data.Supply = totalSupply.Supply
		res.Data.Burned = zeroBigIntStr
		res.Data.Minted = zeroBigIntStr
		res.Data.RecomputedSupply = true
	} else {
		res.Data.Supply = sumStr(totalSupply.Supply, initialSupply.String())
		res.Data.Burned = totalSupply.Burned
		res.Data.Minted = totalSupply.Minted
	}

	makeInitialMintedNotEmpty(res)
	return res, nil
}

func makeInitialMintedNotEmpty(resp *data.ESDTSupplyResponse) {
	if resp.Data.InitialMinted == "" {
		resp.Data.InitialMinted = zeroBigIntStr
	}
}

func (esp *esdtSupplyProcessor) getSupplyFromShards(tokenIdentifier string) (*data.ESDTSupply, error) {
	totalSupply := &data.ESDTSupply{}
	shardIDs := esp.baseProc.GetShardIDs()
	numNodesQueried := 0
	numNodesWithRecomputedSupply := 0
	for _, shardID := range shardIDs {
		if shardID == core.MetachainShardId {
			continue
		}

		supply, err := esp.getShardSupply(tokenIdentifier, shardID)
		if err != nil {
			return nil, err
		}

		addToSupply(totalSupply, supply)
		if supply.RecomputedSupply {
			numNodesWithRecomputedSupply++
		}
		numNodesQueried++
	}

	if numNodesWithRecomputedSupply > 0 {
		totalSupply.RecomputedSupply = true
	}

	return totalSupply, nil
}

func addToSupply(dstSupply, sourceSupply *data.ESDTSupply) {
	dstSupply.Supply = sumStr(dstSupply.Supply, sourceSupply.Supply)
	dstSupply.Burned = sumStr(dstSupply.Burned, sourceSupply.Burned)
	dstSupply.Minted = sumStr(dstSupply.Minted, sourceSupply.Minted)

	return
}

func sumStr(s1, s2 string) string {
	s1Big, ok := big.NewInt(0).SetString(s1, 10)
	if !ok {
		return s2
	}
	s2Big, ok := big.NewInt(0).SetString(s2, 10)
	if !ok {
		return s1
	}

	return big.NewInt(0).Add(s1Big, s2Big).String()
}

func (esp *esdtSupplyProcessor) getInitialSupplyFromMeta(token string) (*big.Int, error) {
	scQuery := &data.SCQuery{
		ScAddress: esdtContractAddress,
		FuncName:  initialESDTSupplyFunc,
		Arguments: [][]byte{[]byte(token)},
	}

	res, _, err := esp.scQueryProc.ExecuteQuery(scQuery)
	if err != nil {
		return nil, err
	}

	if len(res.ReturnData) < 4 {
		return big.NewInt(0), nil
	}

	supplyBytes := res.ReturnData[3]
	supplyBig, ok := big.NewInt(0).SetString(string(supplyBytes), 10)
	if !ok {
		return big.NewInt(0), nil
	}

	return supplyBig, nil
}

func (esp *esdtSupplyProcessor) getShardSupply(token string, shardID uint32) (*data.ESDTSupply, error) {
	shardObservers, errObs := esp.baseProc.GetObservers(shardID, data.AvailabilityAll)
	if errObs != nil {
		return nil, errObs
	}

	apiPath := networkESDTSupplyPath + token
	for _, observer := range shardObservers {
		var responseEsdtSupply data.ESDTSupplyResponse

		_, errGet := esp.baseProc.CallGetRestEndPoint(observer.Address, apiPath, &responseEsdtSupply)
		if errGet != nil {
			log.Error("esdt supply request", "shard ID", observer.ShardId, "observer", observer.Address, "error", errGet.Error())
			continue
		}

		log.Info("esdt supply request", "shard ID", observer.ShardId, "observer", observer.Address)

		return &responseEsdtSupply.Data, nil

	}

	return nil, ErrSendingRequest
}

func isFungibleESDT(tokenIdentifier string) bool {
	splitToken := strings.Split(tokenIdentifier, "-")

	return len(splitToken) < 3
}
