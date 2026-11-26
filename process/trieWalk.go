package process

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/availabilityCommon"
)

const maxContinuationObserverAttempts = 3

type trieWalkRequest struct {
	getObservers     func(data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	path             string
	options          common.AccountQueryOptions
	page             common.TrieWalkPage
	legacyUpstream   bool
	retryOnNodeError bool
	withKeys         bool
}

type trieWalkEnvelope struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
	Code  data.ReturnCode `json:"code"`
}

func fetchPage(ctx context.Context, proc Processor, req trieWalkRequest) (*data.GenericAPIResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path := common.BuildUrlWithAccountQueryOptions(req.path, req.options)
	if req.withKeys {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		query := u.Query()
		query.Set(common.UrlParameterWithKeys, "true")
		u.RawQuery = query.Encode()
		path = u.String()
	}
	method := http.MethodGet
	var body []byte
	if !req.legacyUpstream {
		method = http.MethodPost
		var numKeys *uint64
		if req.page.NumKeys.HasValue {
			value := req.page.NumKeys.Value
			numKeys = &value
		}
		var err error
		body, err = json.Marshal(struct {
			NumKeys       *uint64         `json:"numKeys,omitempty"`
			IteratorState json.RawMessage `json:"iteratorState,omitempty"`
		}{numKeys, req.page.IteratorState})
		if err != nil {
			return nil, err
		}
	}
	provider := availabilityCommon.AvailabilityProvider{}
	availability := []data.ObserverDataAvailabilityType{provider.AvailabilityForAccountQueryOptions(req.options)}
	continuation := req.page.IsContinuation()
	if continuation {
		availability = []data.ObserverDataAvailabilityType{data.AvailabilityAll, data.AvailabilityRecent}
	}
	seen := make(map[string]struct{})
	attempts, unavailable, allBusy := 0, false, true
	lastMessage := ""
	lookupFailed := false
	for _, pool := range availability {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if continuation && attempts == maxContinuationObserverAttempts {
			break
		}
		observers, err := req.getObservers(pool)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if err != nil {
			if !continuation {
				return nil, err
			}
			lookupFailed = true
			lastMessage = err.Error()
			continue
		}
		lookupFailed = false // a successful fallback recovers a failed All lookup
		for _, observer := range observers {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			if continuation && attempts == maxContinuationObserverAttempts {
				break
			}
			if observer == nil || observer.Address == "" {
				continue
			}
			if _, ok := seen[observer.Address]; ok {
				continue
			}
			seen[observer.Address] = struct{}{}
			attempts++
			status, raw, callErr := proc.CallTrieWalkEndpoint(ctx, method, observer.Address, path, body)
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			allBusy = allBusy && status == http.StatusTooManyRequests && callErr == nil
			if callErr != nil {
				lastMessage = callErr.Error()
				if req.legacyUpstream && !req.retryOnNodeError && (status == 200 || status == 400 || status == 500) {
					return nil, &data.ObserverAPIError{StatusCode: 500, Message: lastMessage}
				}
				continue
			}
			envelope := trieWalkEnvelope{}
			decodeErr := json.Unmarshal(raw, &envelope)
			if status == http.StatusOK {
				fields := map[string]json.RawMessage{}
				if decodeErr == nil && envelope.Error == "" && json.Unmarshal(envelope.Data, &fields) == nil && fields != nil {
					checkpoint := bytes.TrimSpace(fields[common.UrlParameterIteratorState])
					if len(checkpoint) > 0 && !bytes.Equal(checkpoint, []byte("null")) && !req.page.Paginated {
						return nil, &data.ObserverAPIError{StatusCode: 500, Message: ErrTrieCollectionLimit.Error()}
					}
					return &data.GenericAPIResponse{Data: envelope.Data, Error: envelope.Error, Code: envelope.Code}, nil
				}
				lastMessage = "invalid observer page response"
				if envelope.Error != "" {
					lastMessage = envelope.Error
				}
				if req.legacyUpstream && !req.retryOnNodeError {
					return nil, &data.ObserverAPIError{StatusCode: 500, Message: lastMessage}
				}
				continue
			}
			lastMessage = envelope.Error
			if lastMessage == "" {
				lastMessage = http.StatusText(status)
			}
			if status == http.StatusBadRequest && decodeErr == nil && envelope.Code == data.ReturnCodeIteratorStateUnavailable && continuation {
				unavailable = true
				continue
			}
			if req.legacyUpstream {
				if req.retryOnNodeError || (status != 400 && status != 500) {
					continue
				}
			} else if retryableTrieWalkStatus(status) {
				continue
			}
			if status < 400 {
				status = http.StatusInternalServerError
			}
			return nil, &data.ObserverAPIError{StatusCode: status, Message: lastMessage}
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if unavailable {
		return nil, &data.ObserverAPIError{StatusCode: 400, Message: ErrCheckpointCannotContinue.Error()}
	}
	if attempts > 0 && allBusy && !lookupFailed {
		return nil, &data.ObserverAPIError{StatusCode: 429, Message: http.StatusText(429)}
	}
	return nil, WrapObserversError(lastMessage)
}

func retryableTrieWalkStatus(status int) bool {
	if status >= 300 && status < 400 {
		return true
	}
	switch status {
	case 408, 404, 405, 429, 503:
		return true
	default:
		return false
	}
}

// Decode errors stay local to the typed account flow rather than producing an empty account.
func decodeAccountPage(response *data.GenericAPIResponse) (*data.AccountModel, error) {
	raw, ok := response.Data.(json.RawMessage)
	if !ok {
		return nil, errors.New("invalid observer account page")
	}
	account := &data.AccountModel{}
	if err := json.Unmarshal(raw, account); err != nil {
		return nil, err
	}
	return account, nil
}
