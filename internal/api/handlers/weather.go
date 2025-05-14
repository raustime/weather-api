package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

func GetWeatherHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
			return
		}
		// Dummy data for now
		c.JSON(http.StatusOK, WeatherResponse{
			Temperature: 20.5,
			Humidity:    60,
			Description: "Sunny",
		})
	}
}