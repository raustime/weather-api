package api

import (
	"weatherapi/internal/api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func SetupRouter(db bun.IDB) *gin.Engine {
	r := gin.Default()

	// API роутери
	r.GET("/api/weather", handlers.GetWeatherHandler())
	r.POST("/api/subscribe", handlers.SubscribeHandler(db))
	r.GET("/api/confirm", handlers.InvalidConfirmHandler())
	r.GET("/api/confirm/*tokenPath", handlers.ConfirmHandler(db))
	r.GET("/api/unsubscribe", handlers.InvalidUnsubscribeHandler())
	r.GET("/api/unsubscribe/*tokenPath", handlers.UnsubscribeHandler(db))

	return r
}
