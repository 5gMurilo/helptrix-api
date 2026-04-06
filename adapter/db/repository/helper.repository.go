package repository

import (
	"fmt"

	"github.com/5gMurilo/helptrix-api/core/domain"
	helperinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/helper"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type helperRepository struct {
	db *gorm.DB
}

func NewHelperRepository(db *gorm.DB) helperinterfaces.IHelperRepository {
	return &helperRepository{db: db}
}

type helperRow struct {
	ID             string `gorm:"column:id"`
	Name           string `gorm:"column:name"`
	Biography      string `gorm:"column:biography"`
	ProfilePicture string `gorm:"column:profile_picture"`
	City           string `gorm:"column:city"`
	State          string `gorm:"column:state"`
}

type categoryRow struct {
	UserID     string `gorm:"column:user_id"`
	CategoryID uint   `gorm:"column:category_id"`
	Name       string `gorm:"column:name"`
}

func (r *helperRepository) Search(params domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
	query := r.db.Table("users").
		Select("users.id, users.name, users.biography, users.profile_picture, COALESCE(addresses.city, '') AS city, COALESCE(addresses.state, '') AS state").
		Where("users.user_type = ? AND users.deleted_at IS NULL", utils.UserTypeHelper).
		Joins("LEFT JOIN addresses ON addresses.user_id = users.id AND addresses.deleted_at IS NULL")

	if params.Name != "" {
		query = query.Where("users.name ILIKE ?", "%"+params.Name+"%")
	}

	if params.CategoryID != nil {
		query = query.
			Joins("JOIN user_categories ON user_categories.user_id = users.id AND user_categories.deleted_at IS NULL").
			Where("user_categories.category_id = ?", *params.CategoryID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return domain.HelperListResponseDTO{}, fmt.Errorf("error counting helpers: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize

	var rows []helperRow
	if err := query.Offset(offset).Limit(params.PageSize).Scan(&rows).Error; err != nil {
		return domain.HelperListResponseDTO{}, fmt.Errorf("error fetching helpers: %w", err)
	}

	if len(rows) == 0 {
		return domain.HelperListResponseDTO{
			Data:     []domain.HelperCardDTO{},
			Total:    total,
			Page:     params.Page,
			PageSize: params.PageSize,
		}, nil
	}

	userIDs := make([]string, len(rows))
	for i, row := range rows {
		userIDs[i] = row.ID
	}

	var catRows []categoryRow
	if err := r.db.Table("user_categories").
		Select("user_categories.user_id, user_categories.category_id, categories.name").
		Joins("JOIN categories ON categories.id = user_categories.category_id AND categories.deleted_at IS NULL").
		Where("user_categories.user_id IN ? AND user_categories.deleted_at IS NULL", userIDs).
		Scan(&catRows).Error; err != nil {
		return domain.HelperListResponseDTO{}, fmt.Errorf("error fetching helper categories: %w", err)
	}

	catMap := make(map[string][]domain.HelperCategoryDTO, len(rows))
	for _, cr := range catRows {
		catMap[cr.UserID] = append(catMap[cr.UserID], domain.HelperCategoryDTO{
			ID:   cr.CategoryID,
			Name: cr.Name,
		})
	}

	cards := make([]domain.HelperCardDTO, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return domain.HelperListResponseDTO{}, fmt.Errorf("error parsing helper id: %w", err)
		}

		cats := catMap[row.ID]
		if cats == nil {
			cats = []domain.HelperCategoryDTO{}
		}
		cards = append(cards, domain.HelperCardDTO{
			ID:             id,
			Name:           row.Name,
			Biography:      row.Biography,
			ProfilePicture: row.ProfilePicture,
			City:           row.City,
			State:          row.State,
			Categories:     cats,
		})
	}

	return domain.HelperListResponseDTO{
		Data:     cards,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}
