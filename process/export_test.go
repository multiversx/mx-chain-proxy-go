package process

// ParseTxStatusResponses -
func ParseTxStatusResponses(responses map[uint32][]string) (string, error) {
	return parseTxStatusResponses(responses)
}
