package groups

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/common"
)

type trieWalkBody struct {
	NumKeys        *uint64         `json:"numKeys"`
	IteratorState  json.RawMessage `json:"iteratorState"`
	OnFinalBlock   *bool           `json:"onFinalBlock"`
	OnStartOfEpoch *uint32         `json:"onStartOfEpoch"`
	BlockNonce     *uint64         `json:"blockNonce"`
	BlockHash      string          `json:"blockHash"`
	BlockRootHash  string          `json:"blockRootHash"`
	HintEpoch      *uint32         `json:"hintEpoch"`
}

func parseTrieWalkRequest(c *gin.Context, address string, network bool) (common.AccountQueryOptions, common.TrieWalkPage, error) {
	options, err := parseAccountQueryOptionsForRoute(c, address, !network)
	if err != nil {
		return options, common.TrieWalkPage{}, err
	}
	page := common.TrieWalkPage{Paginated: c.Request.Method == http.MethodPost}
	page.NumKeys, err = parseUint64UrlParam(c, common.UrlParameterNumKeys)
	if err != nil {
		return options, page, fmt.Errorf("invalid numKeys: %w", err)
	}
	encoded := c.Query(common.UrlParameterIteratorState)
	page.Paginated = page.Paginated || page.NumKeys.HasValue || encoded != ""
	if encoded != "" {
		page.IteratorState, err = base64.RawURLEncoding.DecodeString(encoded)
		if err != nil {
			return options, page, fmt.Errorf("invalid iteratorState: %w", err)
		}
		if !isCheckpointObject(page.IteratorState) {
			return options, page, fmt.Errorf("invalid iteratorState: expected a JSON object")
		}
	}
	if c.Request.Method != http.MethodPost || c.Request.Body == nil {
		return options, page, nil
	}
	body := trieWalkBody{}
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err == io.EOF {
		return options, page, nil
	}
	if err != nil {
		return options, page, fmt.Errorf("invalid paging body: %w", err)
	}
	if body.NumKeys != nil {
		page.NumKeys = core.OptionalUint64{HasValue: true, Value: *body.NumKeys}
	}
	if len(body.IteratorState) > 0 && !bytes.Equal(bytes.TrimSpace(body.IteratorState), []byte("null")) {
		if !isCheckpointObject(body.IteratorState) {
			return options, page, fmt.Errorf("invalid iteratorState: expected a JSON object")
		}
		page.IteratorState = body.IteratorState
	}
	if body.OnFinalBlock != nil {
		options.OnFinalBlock = *body.OnFinalBlock
	}
	if body.OnStartOfEpoch != nil {
		options.OnStartOfEpoch = core.OptionalUint32{HasValue: true, Value: *body.OnStartOfEpoch}
	}
	if body.BlockNonce != nil {
		options.BlockNonce = core.OptionalUint64{HasValue: true, Value: *body.BlockNonce}
	}
	if body.HintEpoch != nil {
		options.HintEpoch = core.OptionalUint32{HasValue: true, Value: *body.HintEpoch}
	}
	if body.BlockHash != "" {
		options.BlockHash, err = hex.DecodeString(body.BlockHash)
		if err != nil {
			return options, page, fmt.Errorf("invalid blockHash: %w", err)
		}
	}
	if body.BlockRootHash != "" {
		options.BlockRootHash, err = hex.DecodeString(body.BlockRootHash)
		if err != nil {
			return options, page, fmt.Errorf("invalid blockRootHash: %w", err)
		}
	}
	return options, page, nil
}

func isCheckpointObject(raw json.RawMessage) bool {
	raw = bytes.TrimSpace(raw)
	return len(raw) > 0 && raw[0] == '{' && json.Valid(raw)
}
