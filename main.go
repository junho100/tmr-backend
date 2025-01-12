package main

import (
	"os"
	"tmr-backend/config"
	"tmr-backend/handler"
	"tmr-backend/model"
	"tmr-backend/router"
	"tmr-backend/util"
)

func main() {
	router := router.NewRouter()
	config := config.NewConfig()

	subjectModel := model.NewSubjectModel(config.DB)
	labModel := model.NewLabModel(config.DB)

	fileUtil := util.NewFileUtil()
	fileUtil.StartCleanupRoutine()

	baseURL := os.Getenv("BASE_URL")
	slackUtil := util.NewSlackUtil(fileUtil, baseURL)

	handler.NewLabHandler(router, labModel, slackUtil, fileUtil)
	handler.NewFileHandler(router)
	handler.NewSubjectHandler(router, subjectModel)
	handler.NewHealthCheck(router)

	router.Run(":8080")
}
