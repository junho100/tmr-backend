package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	tempDir string
}

func NewFileHandler(router *gin.Engine) {
	fileHandler := &FileHandler{
		tempDir: "./temp_files",
	}

	router.GET("/api/files/:filename", fileHandler.DownloadFile)
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(h.tempDir, filename)

	// 파일 존재 여부 확인
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Content-Disposition 헤더 설정
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.File(filePath)
}
