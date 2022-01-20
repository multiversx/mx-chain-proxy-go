package data

import "time"

type EndpointMetrics struct {
	NumRequests         uint64        `json:"num_requests"`
	NumErrors           uint64        `json:"num_errors"`
	TotalResponseTime   time.Duration `json:"total_response_time"`
	LowestResponseTime  time.Duration `json:"lowest_response_time"`
	HighestResponseTime time.Duration `json:"highest_response_time"`
}

type EndpointMetricsData struct {
}
