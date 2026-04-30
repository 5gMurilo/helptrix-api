package repository_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/5gMurilo/helptrix-api/adapter/db/repository"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+uuid.New().String()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
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
			PRIMARY KEY (user_id, category_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create table: %v", err)
		}
	}

	db.Callback().Create().Before("gorm:create").Register("test:uuid_primary_key", func(tx *gorm.DB) {
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

func seedCategories(t *testing.T, db *gorm.DB, ids []uint) {
	t.Helper()
	for _, id := range ids {
		cat := domain.Category{ID: id, Name: "Category " + string(rune('A'+id-1)), Description: "desc"}
		if err := db.Create(&cat).Error; err != nil {
			t.Fatalf("failed to seed category %d: %v", id, err)
		}
	}
}

func helperRegisterDTO() domain.RegisterRequestDTO {
	return domain.RegisterRequestDTO{
		Name:       "John Silva",
		Email:      "john@example.com",
		Password:   "password123",
		UserType:   utils.UserTypeHelper,
		Document:   "12345678901",
		Phone:      "11999999999",
		Categories: []uint{1, 2},
		Address: domain.AddressInputDTO{
			Street:       "Flower Street",
			Number:       "100",
			Neighborhood: "Downtown",
			ZipCode:      "01001000",
			City:         "Sao Paulo",
			State:        "SP",
		},
	}
}

func businessRegisterDTO() domain.RegisterRequestDTO {
	return domain.RegisterRequestDTO{
		Name:       "Company Ltd",
		Email:      "company@example.com",
		Password:   "password123",
		UserType:   utils.UserTypeBusiness,
		Document:   "12345678000195",
		Phone:      "11888888888",
		Categories: []uint{3},
		Address: domain.AddressInputDTO{
			Street:       "Paulista Ave",
			Number:       "1000",
			Neighborhood: "Bela Vista",
			ZipCode:      "01310100",
			City:         "Sao Paulo",
			State:        "SP",
		},
	}
}

func TestAuthRepository_Register(t *testing.T) {
	t.Run("happy path: helper registration", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		seedCategories(t, db, []uint{1, 2})
		repo := repository.NewAuthRepository(db)
		dto := helperRegisterDTO()

		// Act
		user, categoryIDs, err := repo.Register(dto, "salt:hash")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if user.ID == (uuid.UUID{}) {
			t.Error("expected non-zero UUID for user ID")
		}
		if user.Email != dto.Email {
			t.Errorf("expected email %q, got %q", dto.Email, user.Email)
		}
		if len(categoryIDs) != 2 {
			t.Errorf("expected 2 category IDs, got: %d", len(categoryIDs))
		}
	})

	t.Run("happy path: business registration", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		seedCategories(t, db, []uint{3})
		repo := repository.NewAuthRepository(db)
		dto := businessRegisterDTO()

		// Act
		user, categoryIDs, err := repo.Register(dto, "salt:hash")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if user.ID == (uuid.UUID{}) {
			t.Error("expected non-zero UUID for user ID")
		}
		if user.UserType != utils.UserTypeBusiness {
			t.Errorf("expected user_type %q, got %q", utils.UserTypeBusiness, user.UserType)
		}
		if len(categoryIDs) != 1 {
			t.Errorf("expected 1 category ID, got: %d", len(categoryIDs))
		}
	})

	t.Run("duplicate user returns ErrUserAlreadyRegistered", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		seedCategories(t, db, []uint{1, 2})
		repo := repository.NewAuthRepository(db)
		dto := helperRegisterDTO()
		if _, _, err := repo.Register(dto, "salt:hash"); err != nil {
			t.Fatalf("first registration failed unexpectedly: %v", err)
		}

		// Act
		_, _, err := repo.Register(dto, "salt:hash")

		// Assert
		if err == nil {
			t.Fatal("expected error for duplicate registration, got nil")
		}
		if !errors.Is(err, utils.ErrUserAlreadyRegistered) {
			t.Errorf("expected ErrUserAlreadyRegistered, got: %v", err)
		}
	})

	t.Run("address is persisted in DB", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		seedCategories(t, db, []uint{1, 2})
		repo := repository.NewAuthRepository(db)
		dto := helperRegisterDTO()

		// Act
		user, _, err := repo.Register(dto, "salt:hash")
		if err != nil {
			t.Fatalf("registration failed: %v", err)
		}

		// Assert
		var address domain.Address
		if err := db.Where("user_id = ?", user.ID).First(&address).Error; err != nil {
			t.Fatalf("expected address in DB, got error: %v", err)
		}
		if address.Street != dto.Address.Street {
			t.Errorf("expected street %q, got %q", dto.Address.Street, address.Street)
		}
		if address.City != dto.Address.City {
			t.Errorf("expected city %q, got %q", dto.Address.City, address.City)
		}
	})

	t.Run("user_categories are persisted in DB", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		seedCategories(t, db, []uint{1, 2})
		repo := repository.NewAuthRepository(db)
		dto := helperRegisterDTO()

		// Act
		user, _, err := repo.Register(dto, "salt:hash")
		if err != nil {
			t.Fatalf("registration failed: %v", err)
		}

		// Assert
		var rows []domain.UserCategory
		if err := db.Where("user_id = ?", user.ID).Find(&rows).Error; err != nil {
			t.Fatalf("expected user_categories in DB, got error: %v", err)
		}
		if len(rows) != 2 {
			t.Errorf("expected 2 user_category rows, got: %d", len(rows))
		}
	})

	t.Run("user creation error: duplicate document", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		seedCategories(t, db, []uint{1, 2})
		existing := domain.User{
			ID:       uuid.New(),
			Name:     "Other User",
			Email:    "other@example.com",
			Document: "12345678901",
			Password: "salt:hash",
			Phone:    "11777777777",
			UserType: utils.UserTypeHelper,
		}
		if err := db.Create(&existing).Error; err != nil {
			t.Fatalf("failed to seed user with same document: %v", err)
		}
		repo := repository.NewAuthRepository(db)
		dto := helperRegisterDTO()
		dto.Email = "newuser@example.com"

		// Act
		_, _, err := repo.Register(dto, "salt:hash")

		// Assert
		if err == nil {
			t.Fatal("expected error for duplicate document, got nil")
		}
		if err.Error() != "error creating user" {
			t.Errorf("expected 'error creating user', got: %v", err)
		}
	})

	t.Run("category creation error: category does not exist in DB", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		repo := repository.NewAuthRepository(db)
		dto := helperRegisterDTO()

		// Act
		_, _, err := repo.Register(dto, "salt:hash")

		// Assert
		if err == nil {
			t.Fatal("expected error for missing categories, got nil")
		}
		if err.Error() != "error to assign categories for this user" {
			t.Errorf("expected 'error to assign categories for this user', got: %v", err)
		}
	})
}

func TestAuthRepository_FindByEmail(t *testing.T) {
	t.Run("found: returns matching user", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		repo := repository.NewAuthRepository(db)
		inserted := domain.User{
			ID:       uuid.New(),
			Name:     "John Silva",
			Email:    "john@example.com",
			Document: "12345678901",
			Password: "salt:hash",
			Phone:    "11999999999",
			UserType: utils.UserTypeHelper,
		}
		if err := db.Create(&inserted).Error; err != nil {
			t.Fatalf("failed to seed user: %v", err)
		}

		// Act
		user, err := repo.FindByEmail("john@example.com")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if user.Email != inserted.Email {
			t.Errorf("expected email %q, got %q", inserted.Email, user.Email)
		}
	})

	t.Run("not found: returns ErrUserNotFound", func(t *testing.T) {
		// Arrange
		db := setupAuthRepositoryTestDB(t)
		repo := repository.NewAuthRepository(db)

		// Act
		_, err := repo.FindByEmail("ghost@example.com")

		// Assert
		if err == nil {
			t.Fatal("expected error for unknown email, got nil")
		}
		if !errors.Is(err, utils.ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got: %v", err)
		}
	})
}
