package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"weatherapi/internal/models"
)

func SubscribeHandler(db bun.IDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form struct {
			Email     string `form:"email" binding:"required,email"`
			City      string `form:"city" binding:"required"`
			Frequency string `form:"frequency" binding:"required,oneof=hourly daily"`
		}
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Check if already subscribed
		var existing models.Subscription
		err := db.NewSelect().Model(&existing).Where("email = ?", form.Email).Scan(c)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already subscribed"})
			return
		}

		token := uuid.New().String()
		sub := &models.Subscription{
			Email:     form.Email,
			City:      form.City,
			Frequency: form.Frequency,
			Token:     token,
			CreatedAt: time.Now(),
		}

		if _, err := db.NewInsert().Model(sub).Exec(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not subscribe"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})
	}
}

func ConfirmHandler(db bun.IDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		var sub models.Subscription
		err := db.NewSelect().Model(&sub).Where("token = ?", token).Scan(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		sub.Confirmed = true
		sub.ConfirmedAt = time.Now()
		_, err = db.NewUpdate().Model(&sub).WherePK().Exec(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
	}
}

func UnsubscribeHandler(db bun.IDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		_, err := db.NewDelete().Model((*models.Subscription)(nil)).Where("token = ?", token).Exec(c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
	}
}