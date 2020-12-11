package database

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type object = map[string]interface{}

func encodeQuery(query object) (bytes.Buffer, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(query); err != nil {
		return bytes.Buffer{}, fmt.Errorf("error encoding query: %w", err)
	}

	return buff, nil
}

func blockByNonceAndShardIDQuery(nonce uint64, shardID uint32) object {
	return object{
		"query": object{
			"bool": object{
				"must": []interface{}{
					object{
						"match": object{
							"nonce": fmt.Sprintf("%d", nonce),
						},
					},
					object{
						"match": object{
							"shardId": fmt.Sprintf("%d", shardID),
						},
					},
				},
			},
		},
	}
}

func blockByHashQuery(hash string) object {
	return object{
		"query": object{
			"match": object{
				"_id": hash,
			},
		},
	}
}

func txsByMiniblockHashQuery(hash string) object {
	return object{
		"query": object{
			"match": object{
				"miniBlockHash": hash,
			},
		},
	}
}
