package process

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
)

var errTimeIsOut = errors.New("time is out")

type networkConfigResponseData struct {
	Config struct {
		NumShards uint32 `json:"erd_num_shards_without_meta"`
	} `json:"config"`
}

type networkConfigResponse struct {
	Data  networkConfigResponseData `json:"data"`
	Error string                    `json:"error"`
	Code  string                    `json:"code"`
}

// ArgNumShardsProcessor is the DTO used to create a new instance of numShardsProcessor
type ArgNumShardsProcessor struct {
	HttpClient                    HttpClient
	Observers                     []string
	TimeBetweenNodesRequestsInSec int
	NumShardsTimeoutInSec         int
	RequestTimeoutInSec           int
}

type numShardsProcessor struct {
	observers                []string
	httpClient               HttpClient
	timeBetweenNodesRequests time.Duration
	numShardsTimeout         time.Duration
	requestTimeout           time.Duration
}

// NewNumShardsProcessor returns a new instance of numShardsProcessor
func NewNumShardsProcessor(args ArgNumShardsProcessor) (*numShardsProcessor, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &numShardsProcessor{
		observers:                args.Observers,
		httpClient:               args.HttpClient,
		timeBetweenNodesRequests: time.Second * time.Duration(args.TimeBetweenNodesRequestsInSec),
		numShardsTimeout:         time.Second * time.Duration(args.NumShardsTimeoutInSec),
		requestTimeout:           time.Second * time.Duration(args.RequestTimeoutInSec),
	}, nil
}

func checkArgs(args ArgNumShardsProcessor) error {
	if check.IfNilReflect(args.HttpClient) {
		return ErrNilHttpClient
	}
	if len(args.Observers) == 0 {
		return fmt.Errorf("%w for Observers, empty list provided", core.ErrInvalidValue)
	}
	if args.TimeBetweenNodesRequestsInSec == 0 {
		return fmt.Errorf("%w for TimeBetweenNodesRequestsInSec, %d provided", core.ErrInvalidValue, args.TimeBetweenNodesRequestsInSec)
	}
	if args.NumShardsTimeoutInSec == 0 {
		return fmt.Errorf("%w for NumShardsTimeoutInSec, %d provided", core.ErrInvalidValue, args.NumShardsTimeoutInSec)
	}
	if args.RequestTimeoutInSec == 0 {
		return fmt.Errorf("%w for RequestTimeoutInSec, %d provided", core.ErrInvalidValue, args.RequestTimeoutInSec)
	}

	return nil
}

// GetNetworkNumShards tries to get the number of shards from the network
func (processor *numShardsProcessor) GetNetworkNumShards(ctx context.Context) (uint32, error) {
	log.Info("getting the number of shards from observers...")

	waitNodeTicker := time.NewTicker(processor.timeBetweenNodesRequests)
	for {
		select {
		case <-waitNodeTicker.C:
			for _, observerAddress := range processor.observers {
				numShards, httpStatus := processor.tryGetnumShardsFromObserver(observerAddress)
				if httpStatus == http.StatusOK {
					log.Info("fetched the number of shards", "shards", numShards)
					return numShards, nil
				}
			}
		case <-time.After(processor.numShardsTimeout):
			return 0, fmt.Errorf("%w, no observer online", errTimeIsOut)
		case <-ctx.Done():
			log.Debug("closing the getNetworkNumShards loop due to context done...")
			return 0, errTimeIsOut
		}
	}
}

func (processor *numShardsProcessor) tryGetnumShardsFromObserver(observerAddress string) (uint32, int) {
	ctx, cancel := context.WithTimeout(context.Background(), processor.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, observerAddress+NetworkConfigPath, nil)
	if err != nil {
		return 0, http.StatusNotFound
	}

	resp, err := processor.httpClient.Do(req)
	if err != nil {
		return 0, http.StatusNotFound
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			log.LogIfError(resp.Body.Close())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, resp.StatusCode
	}

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	var response networkConfigResponse
	err = json.Unmarshal(responseBodyBytes, &response)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	return response.Data.Config.NumShards, resp.StatusCode
}
