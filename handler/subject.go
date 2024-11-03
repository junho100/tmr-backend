package handler

import (
	"tmr-backend/model"

	"github.com/gin-gonic/gin"
)

type SubjectHandler interface {
}

type subjectHandler struct {
	subjectModel model.SubjectModel
}

func NewSubjectHandler(router *gin.Engine, subjectModel model.SubjectModel) SubjectHandler {
	return &subjectHandler{
		subjectModel: subjectModel,
	}
}
