package dto

import "time"

type HealthResponse struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpTime    int       `json:"up_time,omitempty"`
}
