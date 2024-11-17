package main

import (
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

	slackUtil := util.NewSlackUtil()

	handler.NewLabHandler(router, labModel, slackUtil)
	handler.NewSubjectHandler(router, subjectModel)
	handler.NewHealthCheck(router)

	router.Run(":8080")
}
