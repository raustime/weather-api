package db

import (
	"context"
	"weatherapi/internal/db/models"

	"github.com/uptrace/bun"
)

func RunMigrations(ctx context.Context, db *bun.DB) error {
	// Увага: тільки на час розробки!
	_, _ = db.NewDropTable().
		Model((*models.Subscription)(nil)).
		IfExists().
		Cascade().
		Exec(ctx)

	_, err := db.NewCreateTable().
		Model((*models.Subscription)(nil)).
		IfNotExists().
		Exec(ctx)

	return err
}
