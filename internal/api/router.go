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
	r.GET("/api/confirm/:token", handlers.ConfirmHandler(db))
	r.GET("/api/unsubscribe/:token", handlers.UnsubscribeHandler(db))

	// Статика
	r.Static("/static", "./static/assets")   // CSS, JS, картинки
	r.StaticFile("/", "./static/index.html") // Головна сторінка - форма

	// Для SPA підтримки, якщо треба (можна видалити, якщо не SPA)
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}
