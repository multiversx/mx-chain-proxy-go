package data

import "time"

// EndpointMetrics holds statistics about the requests for a specific endpoint
type EndpointMetrics struct {
	NumRequests         uint64        `json:"num_requests"`
	NumErrors           uint64        `json:"num_errors"`
	TotalResponseTime   time.Duration `json:"total_response_time"`
	LowestResponseTime  time.Duration `json:"lowest_response_time"`
	HighestResponseTime time.Duration `json:"highest_response_time"`
}
