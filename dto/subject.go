package dto

type CreateSubjectRequest struct {
	Age          uint   `json:"age"`
	EnglishLevel string `json:"english_level"`
	Detail       string `json:"detail"`
}

type CreateSubjectResponse struct {
	IdForLogin string `json:"id_for_login"`
}
