package main

import (
	"database/sql"
	"fmt"
	"os"
	"log"
	"time"

	"github.com/joho/godotenv"
	_"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	 "github.com/uptrace/bun/driver/pgdriver"
	 "github.com/uptrace/bun/extra/bundebug"
	
	"weatherapi/internal/api"
)

const MAX_DB_CONNECTION = 25

func main() {
	// Ініціалізація bun.DB
	db := NewDB()
	// Перевірка з'єднання
	if err := db.Ping(); err != nil {
		fmt.Println("DB ping failed:", err)
		return
	}

	r := api.SetupRouter(db)
	r.Run(":8080")
}

func NewDB() *bun.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set in .env file")
	}
	
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	
		sqldb.SetMaxOpenConns(MAX_DB_CONNECTION)
		sqldb.SetMaxIdleConns(MAX_DB_CONNECTION)
		sqldb.SetConnMaxLifetime(time.Hour)

		db := bun.NewDB(sqldb, pgdialect.New())

		db.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(true),
				bundebug.WithVerbose(true),
			),
		)

		// Перевірка підключення
		if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v",  err)
		}
		
		log.Println("✅ Connected to database successfully.")
	return db
}
