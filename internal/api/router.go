package api

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"weatherapi/internal/api/handlers"
)

func SetupRouter(db bun.IDB) *gin.Engine {
	r := gin.Default()
	r.GET("/api/weather", handlers.GetWeatherHandler())
	r.POST("/api/subscribe", handlers.SubscribeHandler(db))
	r.GET("/api/confirm/:token", handlers.ConfirmHandler(db))
	r.GET("/api/unsubscribe/:token", handlers.UnsubscribeHandler(db))

	return r
}