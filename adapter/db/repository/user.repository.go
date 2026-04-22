package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	userinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/user"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userinterfaces.IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetProfile(userID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
	var user domain.User
	var reviews []domain.ReviewListResponseDTO

	result := r.db.
		Preload("Address").
		Preload("Categories").
		Preload("Services").
		Preload("Services.Category").
		First(&user, "id = ?", userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.GetProfileResponseDTO{}, utils.ErrUserNotFound
		}
		return domain.GetProfileResponseDTO{}, result.Error
	}

	// Map categories
	categories := make([]domain.ProfileCategoryDTO, 0, len(user.Categories))
	for _, c := range user.Categories {
		categories = append(categories, domain.ProfileCategoryDTO{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	// Map address
	var address *domain.ProfileAddressDTO
	if user.Address.ID != uuid.Nil {
		address = &domain.ProfileAddressDTO{
			Street:       user.Address.Street,
			Number:       user.Address.Number,
			Complement:   user.Address.Complement,
			Neighborhood: user.Address.Neighborhood,
			City:         user.Address.City,
			State:        user.Address.State,
		}
	}

	// Map services with optional in-memory filters
	services := make([]domain.ServiceResponseDTO, 0)
	for _, svc := range user.Services {
		var actuationDays []string
		if len(svc.ActuationDays) > 0 {
			_ = json.Unmarshal(svc.ActuationDays, &actuationDays)
		}
		if actuationDays == nil {
			actuationDays = []string{}
		}

		// Apply CategoryID filter
		if filters.CategoryID != nil && svc.CategoryID != *filters.CategoryID {
			continue
		}

		// Apply ActuationDays filter
		if len(filters.ActuationDays) > 0 {
			if !hasIntersection(actuationDays, filters.ActuationDays) {
				continue
			}
		}

		var photos []string
		if len(svc.Photos) > 0 {
			_ = json.Unmarshal(svc.Photos, &photos)
		}
		if photos == nil {
			photos = []string{}
		}

		services = append(services, domain.ServiceResponseDTO{
			ID:            svc.ID,
			Name:          svc.Name,
			Description:   svc.Description,
			ActuationDays: actuationDays,
			Value:         svc.Value,
			StartTime:     svc.StartTime,
			EndTime:       svc.EndTime,
			OfferSince:    svc.OfferSince,
			Category: domain.ServiceCategoryDTO{
				ID:   svc.Category.ID,
				Name: svc.Category.Name,
			},
			Photos: photos,
		})
	}

	if user.UserType == "helper" {
		err := r.db.
			Table("reviews").
			Select("reviews.*, users.name AS business_name, users.profile_picture AS business_picture").
			Joins("JOIN users ON users.id = reviews.business_id AND users.deleted_at IS NULL").
			Where("reviews.helper_id = ?", user.ID).
			Where("reviews.deleted_at IS NULL").
			Find(&reviews).Error

		if err != nil {
			return domain.GetProfileResponseDTO{}, err
		}

		dto := domain.GetProfileResponseDTO{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			Biography:      user.Biography,
			ProfilePicture: user.ProfilePicture,
			UserType:       user.UserType,
			Categories:     categories,
			Address:        address,
			Reviews:        reviews,
			Services:       services,
			CreatedAt:      user.CreatedAt,
		}

		return dto, nil
	}

	dto := domain.GetProfileResponseDTO{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		Biography:      user.Biography,
		ProfilePicture: user.ProfilePicture,
		UserType:       user.UserType,
		Categories:     categories,
		Address:        address,
		Reviews:        []domain.ReviewListResponseDTO{},
		Services:       services,
		CreatedAt:      user.CreatedAt,
	}

	return dto, nil
}

func hasIntersection(a, b []string) bool {
	set := make(map[string]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}
	for _, v := range a {
		if _, ok := set[v]; ok {
			return true
		}
	}
	return false
}

func (r *userRepository) UpdateProfile(userID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Step 1: verify user exists
	var user domain.User
	if err := tx.First(&user, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrUserNotFound
		}
		return fmt.Errorf("error finding user: %w", err)
	}

	// Step 2: coupling check — if categories are being updated, verify none of the
	// removed categories have linked active services
	if len(dto.Categories) > 0 {
		var currentUCs []domain.UserCategory
		if err := tx.Where("user_id = ?", userID).Find(&currentUCs).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error fetching current user categories: %w", err)
		}

		newCatSet := make(map[uint]struct{}, len(dto.Categories))
		for _, id := range dto.Categories {
			newCatSet[id] = struct{}{}
		}

		var removedCategoryIDs []uint
		for _, uc := range currentUCs {
			if _, stillPresent := newCatSet[uc.CategoryID]; !stillPresent {
				removedCategoryIDs = append(removedCategoryIDs, uc.CategoryID)
			}
		}

		if len(removedCategoryIDs) > 0 {
			var count int64
			if err := tx.Model(&domain.Service{}).
				Where("user_id = ? AND category_id IN ? AND deleted_at IS NULL", userID, removedCategoryIDs).
				Count(&count).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("error checking linked services: %w", err)
			}
			if count > 0 {
				tx.Rollback()
				return utils.ErrCategoryHasLinkedServices
			}
		}
	}

	// Step 3: update user scalar fields (only non-empty values)
	userUpdates := map[string]interface{}{}
	if dto.Email != "" {
		userUpdates["email"] = dto.Email
	}
	if dto.Phone != "" {
		userUpdates["phone"] = dto.Phone
	}
	if dto.Biography != "" {
		userUpdates["biography"] = dto.Biography
	}
	if dto.ProfilePicture != "" {
		userUpdates["profile_picture"] = dto.ProfilePicture
	}
	if len(userUpdates) > 0 {
		if err := tx.Model(&domain.User{}).Where("id = ?", userID).Updates(userUpdates).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating user fields: %w", err)
		}
	}

	// Step 4: update address (upsert)
	if dto.Address != nil {
		var existingAddr domain.Address
		addrErr := tx.Where("user_id = ?", userID).First(&existingAddr).Error
		if addrErr != nil && !errors.Is(addrErr, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("error fetching address: %w", addrErr)
		}

		if errors.Is(addrErr, gorm.ErrRecordNotFound) {
			newAddr := domain.Address{
				UserID:       userID,
				Street:       dto.Address.Street,
				Number:       dto.Address.Number,
				Complement:   dto.Address.Complement,
				Neighborhood: dto.Address.Neighborhood,
				ZipCode:      dto.Address.ZipCode,
				City:         dto.Address.City,
				State:        dto.Address.State,
			}
			if err := tx.Create(&newAddr).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("error creating address: %w", err)
			}
		} else {
			addrUpdates := map[string]interface{}{
				"street":       dto.Address.Street,
				"number":       dto.Address.Number,
				"complement":   dto.Address.Complement,
				"neighborhood": dto.Address.Neighborhood,
				"zip_code":     dto.Address.ZipCode,
				"city":         dto.Address.City,
				"state":        dto.Address.State,
			}
			if err := tx.Model(&domain.Address{}).Where("user_id = ?", userID).Updates(addrUpdates).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating address: %w", err)
			}
		}
	}

	// Step 5: sync user categories if provided (soft-delete removed + batch upsert)
	if len(dto.Categories) > 0 {
		uniqueCategoryIDs := make([]uint, 0, len(dto.Categories))
		seenCategoryIDs := make(map[uint]struct{}, len(dto.Categories))
		for _, categoryID := range dto.Categories {
			if _, alreadyAdded := seenCategoryIDs[categoryID]; alreadyAdded {
				continue
			}
			seenCategoryIDs[categoryID] = struct{}{}
			uniqueCategoryIDs = append(uniqueCategoryIDs, categoryID)
		}

		rm := tx.Where("user_id = ?", userID).Where("category_id NOT IN ?", uniqueCategoryIDs)
		if err := rm.Delete(&domain.UserCategory{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error soft-deleting removed user categories: %w", err)
		}

		now := time.Now()
		rows := make([]domain.UserCategory, 0, len(uniqueCategoryIDs))
		for _, categoryID := range uniqueCategoryIDs {
			rows = append(rows, domain.UserCategory{
				UserID:     userID,
				CategoryID: categoryID,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "category_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"deleted_at": gorm.Expr("NULL"),
				"updated_at": gorm.Expr("NOW()"),
			}),
		}).Create(&rows).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error upserting user categories: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing update transaction: %w", err)
	}

	return nil
}

func (r *userRepository) GetProfilePicture(userID uuid.UUID) (string, error) {
	var user domain.User

	result := r.db.Select("profile_picture").First(&user, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", utils.ErrUserNotFound
		}
		return "", fmt.Errorf("error fetching profile picture: %w", result.Error)
	}

	return user.ProfilePicture, nil
}

func (r *userRepository) UpdateProfilePicture(userID uuid.UUID, url string) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", userID).Update("profile_picture", url)
	if result.Error != nil {
		return fmt.Errorf("error updating profile picture: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) DeleteProfile(userID uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Where("id = ?", userID).Delete(&domain.User{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return utils.ErrUserNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing delete transaction: %w", err)
	}

	return nil
}
