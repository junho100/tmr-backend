package dto

import "time"

type CreateBreathingHistoryRequest struct {
	IdForLogin    string    `json:"id_for_login"`
	AverageVolume float64   `json:"average_volume"`
	Timestamp     time.Time `json:"timestamp"`
}

type CreateCueHistoryRequest struct {
	IdForLogin string    `json:"id_for_login"`
	Timestamp  time.Time `json:"timestamp"`
	TargetWord string    `json:"target_word"`
}

type StartLabRequest struct {
	LabID   string                  `json:"lab_id"`
	Type    string                  `json:"type"`
	Results []StartLabRequestResult `json:"results"`
}

type StartLabRequestResult struct {
	Word        string `json:"word"`
	WrittenWord string `json:"written_word"`
}

type CreateTestHistoryRequest struct {
	IdForLogin string                           `json:"id_for_login"`
	Type       string                           `json:"type"`
	Results    []CreateTestHistoryRequestResult `json:"results"`
}

type CreateTestHistoryRequestResult struct {
	Word      string `json:"word"`
	IsCorrect bool   `json:"is_correct"`
}

type CreatePreTestHistoryDto struct {
	Results   []CreatePreTestHistoryDtoResult
	LabTestID uint
}

type CreatePreTestHistoryDtoResult struct {
	Word      string `json:"word"`
	IsCorrect bool   `json:"is_correct"`
}

type GetTargetWordsResponse struct {
	Words []string `json:"words"`
}
