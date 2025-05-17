package db

import (
	"context"
	"weatherapi/internal/db/models"

	"github.com/uptrace/bun"
)

func RunMigrations(ctx context.Context, db *bun.DB) error {
	// Створюємо таблицю, якщо ще не існує
	_, err := db.NewCreateTable().
		Model((*models.Subscription)(nil)).
		IfNotExists().
		Exec(ctx)

	return err
}
