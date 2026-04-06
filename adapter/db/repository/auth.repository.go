package repository

import (
	"errors"

	"github.com/5gMurilo/helptrix-api/core/domain"
	authinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/auth"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) authinterfaces.IAuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return domain.User{}, nil, tx.Error
	}

	var existingUser domain.User
	result := tx.Where("email = ? AND user_type = ?", dto.Email, dto.UserType).First(&existingUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return domain.User{}, nil, utils.ErrUserAlreadyRegistered
	}

	user := domain.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Document: dto.Document,
		Password: hashedPassword,
		Phone:    dto.Phone,
		UserType: dto.UserType,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return domain.User{}, nil, errors.New("error creating user")
	}

	address := domain.Address{
		UserID:       user.ID,
		Street:       dto.Address.Street,
		Number:       dto.Address.Number,
		Complement:   dto.Address.Complement,
		Neighborhood: dto.Address.Neighborhood,
		ZipCode:      dto.Address.ZipCode,
		City:         dto.Address.City,
		State:        dto.Address.State,
	}
	if err := tx.Create(&address).Error; err != nil {
		tx.Rollback()
		return domain.User{}, nil, errors.New("error creating user address")
	}

	for _, categoryID := range dto.Categories {
		uc := domain.UserCategory{
			UserID:     user.ID,
			CategoryID: categoryID,
		}
		if err := tx.Create(&uc).Error; err != nil {
			tx.Rollback()
			return domain.User{}, nil, errors.New("error to assign categories for this user")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.User{}, nil, errors.New("error committing transaction")
	}

	return user, dto.Categories, nil
}

func (r *authRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}
