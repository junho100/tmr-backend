package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHealthCheck(router *gin.Engine) {
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})
}
