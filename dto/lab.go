package dto

import "time"

type CreateBreathingHistoryRequest struct {
	IdForLogin    string    `json:"id_for_login"`
	AverageVolume int       `json:"average_volume"`
	Timestamp     time.Time `json:"timestamp"`
}
