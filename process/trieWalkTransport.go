package process

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

// CallTrieWalkEndpoint preserves the status and body, including error envelopes.
// Its POST redirect policy is local so checkpoints cannot be replayed as GETs.
func (bp *BaseProcessor) CallTrieWalkEndpoint(ctx context.Context, method, address, path string, body []byte) (int, []byte, error) {
	if err := ctx.Err(); err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, address+path, bytes.NewReader(body))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	req.Header.Set("Accept", "application/json")
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Multiversx Proxy / 1.0.0 <Posting to nodes>")
	} else {
		req.Header.Set("User-Agent", "Multiversx Proxy / 1.0.0 <Requesting data from nodes>")
	}
	mutHttpClient.RLock()
	client := *bp.httpClient
	mutHttpClient.RUnlock()
	if method == http.MethodPost {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}
		bp.triggerNodesSyncCheck(address)
		if isTimeoutError(err) {
			return http.StatusRequestTimeout, nil, err
		}
		return http.StatusNotFound, nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(resp.Body)
	if ctx.Err() != nil {
		return resp.StatusCode, nil, ctx.Err()
	}
	return resp.StatusCode, raw, err
}
