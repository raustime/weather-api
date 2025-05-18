package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"weatherapi/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"weatherapi/internal/mailer"
)

var ErrAlreadySubscribed = errors.New("email already subscribed")

func SubscribeHandler(db bun.IDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form struct {
			Email     string `form:"email" binding:"required,email"`
			City      string `form:"city" binding:"required"`
			Frequency string `form:"frequency" binding:"required,oneof=hourly daily"`
		}
		if err := c.ShouldBind(&form); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Check if already subscribed
		var existing models.Subscription
		err := db.NewSelect().Model(&existing).Where("email = ?", form.Email).Scan(c)
		if err == nil {
			c.AbortWithStatus(http.StatusConflict)
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
			log.Printf("DB insert error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not subscribe"})
			return
		}

		// Send confirmation email
		if err := mailer.SendConfirmationEmail(form.Email, token); err != nil {
			// Підписку вже створено, але повідомляємо про помилку з email
			log.Printf("Failed to send confirmation email to %s: %v", form.Email, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}

func InvalidConfirmHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func ConfirmHandler(db bun.IDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenWithSlash := c.Param("tokenPath")
		token := strings.TrimPrefix(tokenWithSlash, "/")

		if _, err := uuid.Parse(token); err != nil {
			c.AbortWithStatus(http.StatusBadRequest) // 400: Invalid token format
			return
		}

		var sub models.Subscription
		err := db.NewSelect().Model(&sub).Where("token = ?", token).Scan(c)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		sub.Confirmed = true
		sub.ConfirmedAt = time.Now()
		_, err = db.NewUpdate().Model(&sub).WherePK().Exec(c)
		if err != nil {
			log.Printf("Failed to confirm token %s: %v", token, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	}
}

func InvalidUnsubscribeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func UnsubscribeHandler(db bun.IDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := c.Param("token")
		tokenWithSlash := c.Param("tokenPath")
		token := strings.TrimPrefix(tokenWithSlash, "/")

		if _, err := uuid.Parse(token); err != nil {
			c.AbortWithStatus(http.StatusBadRequest) // 400: Invalid token format
			return
		}

		res, err := db.NewDelete().Model((*models.Subscription)(nil)).Where("token = ?", token).Exec(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if count == 0 {
			c.AbortWithStatus(http.StatusNotFound) // 404 Token not found
			return
		}

		c.Status(http.StatusOK)
	}
}
