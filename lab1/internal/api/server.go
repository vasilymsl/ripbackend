package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("/Users/vasilymaslovsky/Documents/lab1/templates/*")
	r.Static("/static", "/Users/vasilymaslovsky/Documents/lab1/resources")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/", handler.GetOrders)
	r.GET("/credit-services", handler.GetOrders)
	r.GET("/credit/:id", handler.GetOrder)
	r.GET("/creditapplicationbasket", handler.GetApplications)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}