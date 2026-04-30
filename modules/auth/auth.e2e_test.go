package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	pasetoauth "github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/adapter/db/repository"
	"github.com/5gMurilo/helptrix-api/core/domain"
	authmodule "github.com/5gMurilo/helptrix-api/modules/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const e2eSymmetricKey = "3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b"

func setupE2EDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+uuid.New().String()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			document TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			phone TEXT NOT NULL,
			user_type TEXT NOT NULL,
			biography TEXT,
			profile_picture TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS addresses (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			street TEXT NOT NULL,
			number TEXT NOT NULL,
			complement TEXT,
			neighborhood TEXT NOT NULL,
			zip_code TEXT,
			city TEXT NOT NULL,
			state TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS user_categories (
			user_id TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			PRIMARY KEY (user_id, category_id)
		)`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create table: %v", err)
		}
	}

	db.Callback().Create().Before("gorm:create").Register("e2e:uuid_primary_key", func(tx *gorm.DB) {
		if tx.Statement == nil || tx.Statement.ReflectValue.Kind() != reflect.Struct {
			return
		}
		idField := tx.Statement.ReflectValue.FieldByName("ID")
		if idField.IsValid() && idField.Type() == reflect.TypeOf(uuid.UUID{}) {
			if idField.Interface().(uuid.UUID) == (uuid.UUID{}) {
				idField.Set(reflect.ValueOf(uuid.New()))
			}
		}
	})

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	return db
}

func seedE2ECategories(t *testing.T, db *gorm.DB, ids []uint) {
	t.Helper()
	for _, id := range ids {
		cat := domain.Category{ID: id, Name: "Category " + string(rune('A'+id-1)), Description: "desc"}
		if err := db.Create(&cat).Error; err != nil {
			t.Fatalf("failed to seed category %d: %v", id, err)
		}
	}
}

func setupE2ERouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()

	db := setupE2EDB(t)
	seedE2ECategories(t, db, []uint{1, 2, 3})

	maker, err := pasetoauth.NewPasetoMaker(e2eSymmetricKey)
	if err != nil {
		t.Fatalf("failed to create token maker: %v", err)
	}

	repo := repository.NewAuthRepository(db)
	svc := authmodule.NewAuthService(repo, maker)
	ctrl := authmodule.NewAuthController(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/register", ctrl.Register)
	r.POST("/auth/login", ctrl.Login)

	return r, db
}

func validE2ERegisterBody(t *testing.T) []byte {
	t.Helper()
	body := map[string]interface{}{
		"name":       "John Silva",
		"email":      "john@example.com",
		"password":   "password123",
		"user_type":  "helper",
		"document":   "12345678901",
		"phone":      "11999999999",
		"categories": []uint{1, 2},
		"address": map[string]interface{}{
			"street":       "Flower Street",
			"number":       "100",
			"neighborhood": "Downtown",
			"city":         "Sao Paulo",
			"state":        "SP",
			"zip_code":     "01001000",
		},
	}
	b, _ := json.Marshal(body)
	return b
}

func doPost(t *testing.T, router *gin.Engine, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

func TestAuthE2E_Register(t *testing.T) {
	t.Run("201: valid helper registration", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)

		// Act
		w := doPost(t, router, "/auth/register", validE2ERegisterBody(t))

		// Assert
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got: %d — body: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if _, ok := resp["id"]; !ok {
			t.Error("response must contain key 'id'")
		}
		if resp["user_type"] != "helper" {
			t.Errorf("expected user_type 'helper', got: %v", resp["user_type"])
		}
	})

	t.Run("201: valid business registration", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		body := map[string]interface{}{
			"name":       "Company Ltd",
			"email":      "company@example.com",
			"password":   "password123",
			"user_type":  "business",
			"document":   "12345678000195",
			"phone":      "11888888888",
			"categories": []uint{3},
			"address": map[string]interface{}{
				"street":       "Paulista Ave",
				"number":       "1000",
				"neighborhood": "Bela Vista",
				"city":         "Sao Paulo",
				"state":        "SP",
				"zip_code":     "01310100",
			},
		}
		b, _ := json.Marshal(body)

		// Act
		w := doPost(t, router, "/auth/register", b)

		// Assert
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got: %d — body: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp["user_type"] != "business" {
			t.Errorf("expected user_type 'business', got: %v", resp["user_type"])
		}
	})

	t.Run("409: duplicate email and user_type", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		doPost(t, router, "/auth/register", validE2ERegisterBody(t))

		// Act
		w := doPost(t, router, "/auth/register", validE2ERegisterBody(t))

		// Assert
		if w.Code != http.StatusConflict {
			t.Fatalf("expected status 409, got: %d — body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("400: missing required field", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		body, _ := json.Marshal(map[string]interface{}{"email": "john@example.com"})

		// Act
		w := doPost(t, router, "/auth/register", body)

		// Assert
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got: %d", w.Code)
		}
	})

	t.Run("400: password shorter than 6 characters", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		body := map[string]interface{}{
			"name":       "John Silva",
			"email":      "john@example.com",
			"password":   "abc",
			"user_type":  "helper",
			"document":   "12345678901",
			"phone":      "11999999999",
			"categories": []uint{1},
			"address": map[string]interface{}{
				"street": "Flower Street", "number": "100",
				"neighborhood": "Downtown", "city": "Sao Paulo",
				"state": "SP", "zip_code": "01001000",
			},
		}
		b, _ := json.Marshal(body)

		// Act
		w := doPost(t, router, "/auth/register", b)

		// Assert
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got: %d", w.Code)
		}
	})
}

func TestAuthE2E_Login(t *testing.T) {
	t.Run("201: valid login returns token and id", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		doPost(t, router, "/auth/register", validE2ERegisterBody(t))
		loginBody, _ := json.Marshal(map[string]string{
			"email": "john@example.com", "password": "password123",
		})

		// Act
		w := doPost(t, router, "/auth/login", loginBody)

		// Assert
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got: %d — body: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		token, ok := resp["token"]
		if !ok {
			t.Fatal("response must contain key 'token'")
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
		if _, ok := resp["id"]; !ok {
			t.Error("response must contain key 'id'")
		}
	})

	t.Run("401: wrong password", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		doPost(t, router, "/auth/register", validE2ERegisterBody(t))
		loginBody, _ := json.Marshal(map[string]string{
			"email": "john@example.com", "password": "wrongpassword",
		})

		// Act
		w := doPost(t, router, "/auth/login", loginBody)

		// Assert
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got: %d", w.Code)
		}
	})

	t.Run("401: email not found", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		loginBody, _ := json.Marshal(map[string]string{
			"email": "ghost@example.com", "password": "password123",
		})

		// Act
		w := doPost(t, router, "/auth/login", loginBody)

		// Assert
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got: %d", w.Code)
		}
	})

	t.Run("400: invalid request body", func(t *testing.T) {
		// Arrange
		router, _ := setupE2ERouter(t)
		body, _ := json.Marshal(map[string]string{})

		// Act
		w := doPost(t, router, "/auth/login", body)

		// Assert
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got: %d", w.Code)
		}
	})
}
