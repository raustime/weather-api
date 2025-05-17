package db

import (
	"context"

	"github.com/uptrace/bun"
)

type Subscription struct {
	bun.BaseModel `bun:"table:subscriptions"`

	ID        int64  `bun:",pk,autoincrement"`
	Email     string `bun:",notnull,unique"`
	City      string `bun:",notnull"`
	Frequency string `bun:",notnull"` // "hourly" or "daily"
	Confirmed bool   `bun:",default:false"`
	Token     string `bun:",notnull,unique"`
}

func RunMigrations(ctx context.Context, db *bun.DB) error {
	// Створюємо таблицю, якщо ще не існує
	_, err := db.NewCreateTable().
		Model((*Subscription)(nil)).
		IfNotExists().
		Exec(ctx)

	return err
}
