package handler

import (
	"fmt"
	"net/http"
	"tmr-backend/dto"
	"tmr-backend/model"

	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	subjectModel model.SubjectModel
}

func NewSubjectHandler(router *gin.Engine, subjectModel model.SubjectModel) {
	subjectHandler := &SubjectHandler{
		subjectModel: subjectModel,
	}

	router.POST("/api/subjects", subjectHandler.CreateSubject)
	router.GET("/api/subjects/check", subjectHandler.CheckSubjectExists)
}

func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var createSubjectRequest dto.CreateSubjectRequest
	var err error
	var idForLogin string

	if err = c.BindJSON(&createSubjectRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	if idForLogin, err = h.subjectModel.CreateSubject(createSubjectRequest.Age, createSubjectRequest.EnglishLevel, createSubjectRequest.Detail); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusCreated, &dto.CreateSubjectResponse{
		IdForLogin: idForLogin,
	})
}

func (h *SubjectHandler) CheckSubjectExists(c *gin.Context) {
	idForLogin := c.Query("id")
	fmt.Println(idForLogin)

	if _, err := h.subjectModel.FindSubjectByIdForLogin(idForLogin); err == nil {
		c.JSON(http.StatusOK, &dto.CheckSubjectExistsResponse{
			IsExists: true,
		})
		return
	}

	c.JSON(http.StatusOK, &dto.CheckSubjectExistsResponse{
		IsExists: false,
	})
}
