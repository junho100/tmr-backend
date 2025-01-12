package handler

import (
	"fmt"
	"log"
	"net/http"
	"tmr-backend/dto"
	"tmr-backend/model"
	"tmr-backend/util"

	"github.com/gin-gonic/gin"
)

type LabHandler struct {
	labModel  model.LabModel
	slackUtil util.SlackUtil
	fileUtil  util.FileUtil
}

func NewLabHandler(router *gin.Engine, labModel model.LabModel, slackUtil util.SlackUtil, fileUtil util.FileUtil) {
	labHandler := &LabHandler{
		labModel:  labModel,
		slackUtil: slackUtil,
		fileUtil:  fileUtil,
	}

	router.POST("/api/labs/breathing", labHandler.CreateBreathingHistory)
	router.POST("/api/labs/cue", labHandler.CreateCueHistory)
	router.POST("/api/labs/start-test", labHandler.StartTest)
	router.POST("/api/labs/test", labHandler.CreateTestHistory)
	router.GET("/api/labs/cue", labHandler.GetTargetWords)
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

func (h *LabHandler) StartTest(c *gin.Context) {
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

	if startLabRequest.Type == "pretest" {
		// CSV 파일 생성을 위한 데이터 준비
		csvContent := "Word,WrittenWord\n"
		for _, result := range startLabRequest.Results {
			csvContent += fmt.Sprintf("%s,%s\n", result.Word, result.WrittenWord)
		}

		// 임시 CSV 파일 생성 - fileUtil 직접 사용
		filename, err := h.fileUtil.CreateTempCSVFile(csvContent)
		if err != nil {
			log.Printf("Failed to create CSV file: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		if err := h.labModel.CreatePreTest(lab.ID, startLabRequest.Results); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		// Slack 메시지 전송
		if err := h.slackUtil.SendPreTestStartMessage(startLabRequest.LabID, filename); err != nil {
			log.Printf("Failed to send slack message: %v", err)
		}

		c.JSON(http.StatusCreated, nil)
		return
	}

	if startLabRequest.Type == "test" {
		if err := h.labModel.CreateTest(lab.ID); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		if err := h.slackUtil.SendTestStartMessage(startLabRequest.LabID); err != nil {
			log.Printf("Failed to send slack message: %v", err)
		}

		c.JSON(http.StatusCreated, nil)
		return
	}

	c.JSON(http.StatusBadRequest, nil)
}

func (h *LabHandler) CreateTestHistory(c *gin.Context) {
	var createTestHistoryRequest dto.CreateTestHistoryRequest

	if err := c.BindJSON(&createTestHistoryRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	labTest, err := h.labModel.GetLabTestByIdForLogin(createTestHistoryRequest.IdForLogin, createTestHistoryRequest.Type)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	createPreTestHistoryDtoResults := make([]dto.CreatePreTestHistoryDtoResult, len(createTestHistoryRequest.Results))
	for i, v := range createTestHistoryRequest.Results {
		createPreTestHistoryDtoResults[i] = dto.CreatePreTestHistoryDtoResult(v)
	}
	createTestHistoryDto := dto.CreatePreTestHistoryDto{
		LabTestID: labTest.ID,
		Results:   createPreTestHistoryDtoResults,
	}

	if createTestHistoryRequest.Type == "pretest" {
		selectedWords, correctCount, wrongCount, err := h.labModel.CreatePreTestHistory(createTestHistoryDto)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		// Slack 메시지 전송
		if err := h.slackUtil.SendTestResultMessage(createTestHistoryRequest.IdForLogin, correctCount, wrongCount, selectedWords); err != nil {
			// Slack 메시지 전송 실패는 클라이언트에게 에러를 반환하지 않음
			log.Printf("Failed to send slack message: %v", err)
		}

		c.JSON(http.StatusCreated, nil)
		return
	}

	if createTestHistoryRequest.Type == "test" {
		err := h.labModel.CreateTestHistory(createTestHistoryDto)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusCreated, nil)
		return
	}

	c.JSON(http.StatusBadRequest, nil)
}

func (h *LabHandler) GetTargetWords(c *gin.Context) {
	idForLogin := c.Query("id")

	lab, err := h.labModel.GetLabBySubjectIdForLogin(idForLogin)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	words, err := h.labModel.GetTargetWordsByLabId(lab.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	getTargetWordsResponse := dto.GetTargetWordsResponse{
		Words: words,
	}
	c.JSON(http.StatusOK, getTargetWordsResponse)
}
