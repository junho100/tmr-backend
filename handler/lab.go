package handler

import (
	"net/http"
	"tmr-backend/dto"
	"tmr-backend/model"

	"github.com/gin-gonic/gin"
)

type LabHandler struct {
	labModel model.LabModel
}

func NewLabHandler(router *gin.Engine, labModel model.LabModel) {
	labHandler := &LabHandler{
		labModel: labModel,
	}

	router.POST("/api/labs/breathing", labHandler.CreateBreathingHistory)
	router.POST("/api/labs/cue", labHandler.CreateCueHistory)
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
