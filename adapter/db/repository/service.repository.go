package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) serviceinterfaces.IServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return domain.ServiceResponseDTO{}, tx.Error
	}

	actuationDaysJSON, err := json.Marshal(dto.ActuationDays)
	if err != nil {
		tx.Rollback()
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	photosJSON, err := json.Marshal(dto.Photos)
	if err != nil {
		tx.Rollback()
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	value, err := decimal.NewFromString(dto.Value)
	if err != nil {
		tx.Rollback()
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	service := domain.Service{
		UserID:        userID,
		CategoryID:    dto.CategoryID,
		Name:          dto.Name,
		Description:   dto.Description,
		ActuationDays: actuationDaysJSON,
		Value:         value,
		StartTime:     dto.StartTime,
		EndTime:       dto.EndTime,
		OfferSince:    dto.OfferSince,
		Photos:        photosJSON,
	}

	if err := tx.Create(&service).Error; err != nil {
		tx.Rollback()
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	if err := r.db.Preload("Category").First(&service, "id = ?", service.ID).Error; err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	var actuationDays []string
	if err := json.Unmarshal(service.ActuationDays, &actuationDays); err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
	}

	var photos []string
	if len(service.Photos) > 0 {
		if err := json.Unmarshal(service.Photos, &photos); err != nil {
			return domain.ServiceResponseDTO{}, fmt.Errorf("error creating service: %w", err)
		}
	}

	response := domain.ServiceResponseDTO{
		ID:            service.ID,
		Name:          service.Name,
		Description:   service.Description,
		ActuationDays: actuationDays,
		Value:         service.Value,
		StartTime:     service.StartTime,
		EndTime:       service.EndTime,
		OfferSince:    service.OfferSince,
		Category: domain.ServiceCategoryDTO{
			ID:   service.Category.ID,
			Name: service.Category.Name,
		},
		Photos: photos,
	}

	return response, nil
}

func (r *serviceRepository) ExistsByNameAndUser(name string, userID uuid.UUID) (bool, error) {
	result := r.db.Where("name = ? AND user_id = ? AND deleted_at IS NULL", name, userID).First(&domain.Service{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *serviceRepository) UserHasCategory(userID uuid.UUID, categoryID uint) (bool, error) {
	result := r.db.Where("user_id = ? AND category_id = ? AND deleted_at IS NULL", userID, categoryID).First(&domain.UserCategory{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func mapServiceToDTO(service domain.Service) (domain.ServiceResponseDTO, error) {
	var actuationDays []string
	if err := json.Unmarshal(service.ActuationDays, &actuationDays); err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error mapping service: %w", err)
	}

	var photos []string
	if len(service.Photos) > 0 {
		if err := json.Unmarshal(service.Photos, &photos); err != nil {
			return domain.ServiceResponseDTO{}, fmt.Errorf("error mapping service: %w", err)
		}
	}

	return domain.ServiceResponseDTO{
		ID:            service.ID,
		Name:          service.Name,
		Description:   service.Description,
		ActuationDays: actuationDays,
		Value:         service.Value,
		StartTime:     service.StartTime,
		EndTime:       service.EndTime,
		OfferSince:    service.OfferSince,
		Category: domain.ServiceCategoryDTO{
			ID:   service.Category.ID,
			Name: service.Category.Name,
		},
		Photos: photos,
	}, nil
}

func (r *serviceRepository) List(userID uuid.UUID) ([]domain.ServiceResponseDTO, error) {
	var services []domain.Service
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Preload("Category").Find(&services).Error; err != nil {
		return []domain.ServiceResponseDTO{}, fmt.Errorf("error listing services: %w", err)
	}

	result := make([]domain.ServiceResponseDTO, 0, len(services))
	for _, s := range services {
		dto, err := mapServiceToDTO(s)
		if err != nil {
			return []domain.ServiceResponseDTO{}, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (r *serviceRepository) GetByID(serviceID uuid.UUID, userID uuid.UUID) (domain.ServiceResponseDTO, error) {
	var service domain.Service
	err := r.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", serviceID, userID).Preload("Category").First(&service).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
	}
	if err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error fetching service: %w", err)
	}
	return mapServiceToDTO(service)
}

func (r *serviceRepository) Update(serviceID uuid.UUID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return domain.ServiceResponseDTO{}, tx.Error
	}

	var service domain.Service
	if err := tx.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
		}
		return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
	}

	updateMap := map[string]interface{}{}

	if dto.Name != nil {
		updateMap["name"] = *dto.Name
	}
	if dto.Description != nil {
		updateMap["description"] = *dto.Description
	}
	if dto.ActuationDays != nil {
		daysJSON, err := json.Marshal(dto.ActuationDays)
		if err != nil {
			tx.Rollback()
			return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
		}
		updateMap["actuation_days"] = daysJSON
	}
	if dto.Value != nil {
		val, err := decimal.NewFromString(*dto.Value)
		if err != nil {
			tx.Rollback()
			return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
		}
		updateMap["value"] = val
	}
	if dto.StartTime != nil {
		updateMap["start_time"] = *dto.StartTime
	}
	if dto.EndTime != nil {
		updateMap["end_time"] = *dto.EndTime
	}
	if dto.OfferSince != nil {
		updateMap["offer_since"] = *dto.OfferSince
	}
	if dto.Photos != nil {
		photosJSON, err := json.Marshal(dto.Photos)
		if err != nil {
			tx.Rollback()
			return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
		}
		updateMap["photos"] = photosJSON
	}
	if dto.CategoryID != nil {
		updateMap["category_id"] = *dto.CategoryID
	}

	if err := tx.Model(&service).Updates(updateMap).Error; err != nil {
		tx.Rollback()
		return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
	}

	if err := r.db.Preload("Category").First(&service, "id = ?", service.ID).Error; err != nil {
		return domain.ServiceResponseDTO{}, fmt.Errorf("error updating service: %w", err)
	}

	return mapServiceToDTO(service)
}

func (r *serviceRepository) Delete(serviceID uuid.UUID, userID uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Where("id = ? AND user_id = ?", serviceID, userID).Delete(&domain.Service{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting service: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return utils.ErrServiceNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}
	return nil
}

func (r *serviceRepository) ExistsByNameAndUserExcluding(name string, userID uuid.UUID, excludeID uuid.UUID) (bool, error) {
	result := r.db.Where("name = ? AND user_id = ? AND id != ? AND deleted_at IS NULL", name, userID, excludeID).First(&domain.Service{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
