package dto

import "time"

type CreateBreathingHistoryRequest struct {
	IdForLogin    string    `json:"id_for_login"`
	AverageVolume int       `json:"average_volume"`
	Timestamp     time.Time `json:"timestamp"`
}

type CreateCueHistoryRequest struct {
	IdForLogin string    `json:"id_for_login"`
	Timestamp  time.Time `json:"timestamp"`
	TargetWord string    `json:"target_word"`
}

type StartLabRequest struct {
	LabID string `json:"lab_id"`
}

type CreateTestHistoryRequest struct {
	IdForLogin string                           `json:"id_for_login"`
	Results    []CreateTestHistoryRequestResult `json:"results"`
}

type CreateTestHistoryRequestResult struct {
	Word      string `json:"word"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateTestHistoryDto struct {
	Results   []CreateTestHistoryDtoResult
	LabTestID uint
}

type CreateTestHistoryDtoResult struct {
	Word      string `json:"word"`
	IsCorrect bool   `json:"is_correct"`
}

type GetTargetWordsResponse struct {
	Words []string `json:"words"`
}
