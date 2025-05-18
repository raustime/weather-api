package jobs

import (
	"context"
	"log"
	"os"
	"time"

	"weatherapi/internal/db/models"
	"weatherapi/internal/mailer"
	openweatherapi "weatherapi/internal/openweatherapi"

	"github.com/uptrace/bun"
)

func StartWeatherNotificationLoop(db *bun.DB) {
	go func() {
		for {
			now := time.Now()

			// Надсилати щогодини (00 хвилин)
			if now.Minute() == 0 {
				sendUpdates(db, "hourly")

				// О 8:00 — щоденні
				if now.Hour() == 8 {
					sendUpdates(db, "daily")
				}
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}
func sendUpdates(db *bun.DB, frequency string) {

	baseURL := os.Getenv("APP_BASE_URL")
	var subs []models.Subscription
	err := db.NewSelect().
		Model(&subs).
		Where("confirmed = TRUE").
		Where("frequency = ?", frequency).
		Scan(context.Background())
	if err != nil {
		log.Printf("Failed to fetch subscriptions: %v", err)
		return
	}

	for _, sub := range subs {
		weather, err := openweatherapi.FetchWeather(sub.City)
		if err != nil {
			log.Printf("Weather fetch error for %s: %v", sub.City, err)
			continue
		}

		err = mailer.SendWeatherEmail(sub.Email, sub.City, weather, baseURL, sub.Token)
		if err != nil {
			log.Printf("Failed to send weather email to %s: %v", sub.Email, err)
		}
	}
}
