package models

import (
	"time"
	"github.com/uptrace/bun"
)

type Subscription struct {
	bun.BaseModel `bun:"table:subscriptions"`

	ID          int64     `bun:",pk,autoincrement"`
	Email       string    `bun:",unique,notnull"`
	City        string    `bun:",notnull"`
	Frequency   string    `bun:",notnull"`
	Confirmed   bool      `bun:",notnull,default:false"`
	Token       string    `bun:",notnull"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	ConfirmedAt time.Time `bun:",nullzero"`
}
