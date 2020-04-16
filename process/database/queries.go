package database

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
)

func encodeQuery(query map[string]interface{}) (bytes.Buffer, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(query); err != nil {
		return bytes.Buffer{}, fmt.Errorf("error encoding query: %w", err)
	}

	return buff, nil
}

func txsByAddrQuery(addr string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"sender": addr,
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"receiver": addr,
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"timestamp": map[string]interface{}{
				"order": "desc",
			},
		},
	}
}

func latestBlockQuery() map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"shardId": fmt.Sprintf("%d", core.MetachainShardId),
			},
		},
		"sort": map[string]interface{}{
			"nonce": map[string]interface{}{
				"order": "desc",
			},
		},
	}
}
