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
