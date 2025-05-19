package handlers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"weatherapi/internal/api/handlers"
	dbpkg "weatherapi/internal/db"
	"weatherapi/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func setupRouter(db bun.IDB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/api/subscribe", handlers.SubscribeHandler(db))
	r.GET("/api/confirm/:token", handlers.ConfirmHandler(db))
	r.GET("/api/unsubscribe/:token", handlers.UnsubscribeHandler(db))

	return r
}

func initialTestDB(t *testing.T) *bun.DB {

	db := setupTestDB(t)
	ctx := context.Background()
	// 5. Міграції або створення таблиць
	if err := db.ResetModel(ctx, (*models.Subscription)(nil)); err != nil {
		log.Fatalf("failed to reset model: %v", err)
	}

	return db
}

func isDatabaseExistsError(err error) bool {
	return err != nil && ( // PostgreSQL код "42P04": duplicate_database
	err.Error() == `pq: database "weatherdb_test" already exists` ||
		err.Error() == `ERROR: database "weatherdb_test" already exists (SQLSTATE 42P04)`)
}

func TestSubscribe_Success(t *testing.T) {
	// використайте in-memory DB або мок
	db := setupTestDB(t)
	router := setupRouter(db)

	payload := map[string]string{
		"email":     "test@example.com",
		"city":      "Kyiv",
		"frequency": "daily",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Subscription successful")
}

func TestSubscribe_Conflict(t *testing.T) {
	db := setupTestDB(t)
	// вручну вставляємо підписку
	_ = insertTestSubscription(db, "test@example.com")

	router := setupRouter(db)

	payload := map[string]string{
		"email":     "test@example.com",
		"city":      "Kyiv",
		"frequency": "daily",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}
func TestConfirm_Success(t *testing.T) {
	db := setupTestDB(t)
	token := uuid.New().String()
	_ = insertTestSubscriptionWithToken(db, "test@example.com", token)

	router := setupRouter(db)

	req, _ := http.NewRequest("GET", "/api/confirm/"+token, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestUnsubscribe_Success(t *testing.T) {
	db := setupTestDB(t)
	token := uuid.New().String()
	_ = insertTestSubscriptionWithToken(db, "test@example.com", token)

	router := setupRouter(db)

	req, _ := http.NewRequest("GET", "/api/unsubscribe/"+token, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func insertTestSubscription(db bun.IDB, email string) error {
	sub := &models.Subscription{
		Email:     email,
		City:      "Kyiv",
		Frequency: "daily",
		Token:     uuid.New().String(),
		CreatedAt: time.Now(),
	}
	_, err := db.NewInsert().Model(sub).Exec(context.Background())
	return err
}

func insertTestSubscriptionWithToken(db bun.IDB, email, token string) error {
	sub := &models.Subscription{
		Email:     email,
		City:      "Kyiv",
		Frequency: "daily",
		Token:     token,
		CreatedAt: time.Now(),
	}
	_, err := db.NewInsert().Model(sub).Exec(context.Background())
	return err
}

func setupTestDB(t *testing.T) *bun.DB {
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		t.Fatal("TEST_DB_URL not set")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	// Запускаємо міграції (в твоєму пакеті dbpkg)
	if err := dbpkg.RunMigrations(ctx, db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Очистка таблиці subscriptions
	_, err := db.ExecContext(ctx, `TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE;`)
	if err != nil {
		t.Fatalf("failed to truncate subscriptions table: %v", err)
	}

	return db
}
