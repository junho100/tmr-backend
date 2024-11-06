package main

import (
	"tmr-backend/config"
	"tmr-backend/handler"
	"tmr-backend/model"
	"tmr-backend/router"
)

func main() {
	router := router.NewRouter()
	config := config.NewConfig()

	subjectModel := model.NewSubjectModel(config.DB)
	labModel := model.NewLabModel(config.DB)

	handler.NewLabHandler(router, labModel)
	handler.NewSubjectHandler(router, subjectModel)
	handler.NewHealthCheck(router)

	router.Run(":8080")
}
