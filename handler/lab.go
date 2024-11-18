package handler

import (
	"net/http"
	"tmr-backend/dto"
	"tmr-backend/model"
	"tmr-backend/util"

	"github.com/gin-gonic/gin"
)

type LabHandler struct {
	labModel  model.LabModel
	slackUtil util.SlackUtil
}

func NewLabHandler(router *gin.Engine, labModel model.LabModel, slackUtil util.SlackUtil) {
	labHandler := &LabHandler{
		labModel:  labModel,
		slackUtil: slackUtil,
	}

	router.POST("/api/labs/breathing", labHandler.CreateBreathingHistory)
	router.POST("/api/labs/cue", labHandler.CreateCueHistory)
	router.POST("/api/labs/start-test", labHandler.StartLab)
	router.POST("/api/labs/test", labHandler.CreateTestHistory)
}

func (h *LabHandler) CreateBreathingHistory(c *gin.Context) {
	var createBreathingHistoryRequest dto.CreateBreathingHistoryRequest

	if err := c.BindJSON(&createBreathingHistoryRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	if err := h.labModel.CreateBreathingHistory(createBreathingHistoryRequest.IdForLogin, createBreathingHistoryRequest.AverageVolume, createBreathingHistoryRequest.Timestamp); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *LabHandler) CreateCueHistory(c *gin.Context) {
	var createCueHistoryRequest dto.CreateCueHistoryRequest

	if err := c.BindJSON(&createCueHistoryRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	if err := h.labModel.CreateCueHistory(createCueHistoryRequest.IdForLogin, createCueHistoryRequest.Timestamp, createCueHistoryRequest.TargetWord); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *LabHandler) StartLab(c *gin.Context) {
	var startLabRequest dto.StartLabRequest

	if err := c.BindJSON(&startLabRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	lab, err := h.labModel.GetLabBySubjectIdForLogin(startLabRequest.LabID)

	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := h.labModel.CreatePreTest(lab.ID); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := h.slackUtil.SendTestStartMessage(startLabRequest.LabID); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *LabHandler) CreateTestHistory(c *gin.Context) {
	var createTestHistoryRequest dto.CreateTestHistoryRequest

	if err := c.BindJSON(&createTestHistoryRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	labTest, err := h.labModel.GetLabTestByIdForLogin(createTestHistoryRequest.IdForLogin)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	createTestHistoryDtoResults := make([]dto.CreateTestHistoryDtoResult, len(createTestHistoryRequest.Results))
	for i, v := range createTestHistoryRequest.Results {
		createTestHistoryDtoResults[i] = dto.CreateTestHistoryDtoResult(v)
	}
	createTestHistoryDto := dto.CreateTestHistoryDto{
		LabTestID: labTest.ID,
		Results:   createTestHistoryDtoResults,
	}

	if err := h.labModel.CreateTestHistory(createTestHistoryDto); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
