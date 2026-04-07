package strategies

import (
	"context"
	"errors"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

type stubServiceRepository struct {
	getByIDFn func(serviceID, userID uuid.UUID) (domain.ServiceResponseDTO, error)
	updateFn  func(serviceID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error)

	updateCalled bool
	lastDTO      domain.UpdateServiceRequestDTO
}

func (s *stubServiceRepository) Create(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	return domain.ServiceResponseDTO{}, nil
}

func (s *stubServiceRepository) ExistsByNameAndUser(name string, userID uuid.UUID) (bool, error) {
	return false, nil
}

func (s *stubServiceRepository) ExistsByNameAndUserExcluding(name string, userID uuid.UUID, excludeID uuid.UUID) (bool, error) {
	return false, nil
}

func (s *stubServiceRepository) UserHasCategory(userID uuid.UUID, categoryID uint) (bool, error) {
	return false, nil
}

func (s *stubServiceRepository) List(userID uuid.UUID) ([]domain.ServiceResponseDTO, error) {
	return nil, nil
}

func (s *stubServiceRepository) GetByID(serviceID, userID uuid.UUID) (domain.ServiceResponseDTO, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(serviceID, userID)
	}
	return domain.ServiceResponseDTO{}, nil
}

func (s *stubServiceRepository) Update(serviceID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	s.updateCalled = true
	s.lastDTO = dto
	if s.updateFn != nil {
		return s.updateFn(serviceID, userID, dto)
	}
	return domain.ServiceResponseDTO{}, nil
}

func (s *stubServiceRepository) Delete(serviceID, userID uuid.UUID) error {
	return nil
}

func TestServiceImageStrategy_Upload(t *testing.T) {
	bucketName := "my-bucket"
	serviceID := uuid.New()
	requesterID := uuid.New()
	newURL := "https://storage.googleapis.com/my-bucket/service-images/" + serviceID.String() + "/x.jpg"

	getErr := errors.New("get failed")
	uploadErr := errors.New("upload failed")
	updateErr := errors.New("update failed")

	tests := []struct {
		name           string
		repo           *stubServiceRepository
		storage        *stubStorageService
		wantErr        error
		wantUpload     bool
		wantUpdate     bool
		wantDeleteRoll bool
		wantPhotosLen  int
	}{
		{
			name: "success appends photo",
			repo: &stubServiceRepository{
				getByIDFn: func(sid, uid uuid.UUID) (domain.ServiceResponseDTO, error) {
					if sid != serviceID || uid != requesterID {
						t.Fatalf("unexpected ids")
					}
					return domain.ServiceResponseDTO{ID: serviceID, Photos: []string{"https://old"}}, nil
				},
			},
			storage:       &stubStorageService{},
			wantErr:       nil,
			wantUpload:    true,
			wantUpdate:    true,
			wantPhotosLen: 2,
		},
		{
			name: "success with no prior photos",
			repo: &stubServiceRepository{
				getByIDFn: func(sid, uid uuid.UUID) (domain.ServiceResponseDTO, error) {
					return domain.ServiceResponseDTO{ID: serviceID, Photos: nil}, nil
				},
			},
			storage:       &stubStorageService{},
			wantErr:       nil,
			wantUpload:    true,
			wantUpdate:    true,
			wantPhotosLen: 1,
		},
		{
			name: "GetByID error",
			repo: &stubServiceRepository{
				getByIDFn: func(uuid.UUID, uuid.UUID) (domain.ServiceResponseDTO, error) {
					return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
				},
			},
			storage:    &stubStorageService{},
			wantErr:    utils.ErrServiceNotFound,
			wantUpload: false,
			wantUpdate: false,
		},
		{
			name: "upload fails",
			repo: &stubServiceRepository{
				getByIDFn: func(uuid.UUID, uuid.UUID) (domain.ServiceResponseDTO, error) {
					return domain.ServiceResponseDTO{ID: serviceID}, nil
				},
			},
			storage: &stubStorageService{
				uploadFn: func(context.Context, string, string, string, []byte, string) (string, error) {
					return "", uploadErr
				},
			},
			wantErr:    uploadErr,
			wantUpload: true,
			wantUpdate: false,
		},
		{
			name: "update fails rolls back object",
			repo: &stubServiceRepository{
				getByIDFn: func(uuid.UUID, uuid.UUID) (domain.ServiceResponseDTO, error) {
					return domain.ServiceResponseDTO{ID: serviceID}, nil
				},
				updateFn: func(uuid.UUID, uuid.UUID, domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
					return domain.ServiceResponseDTO{}, updateErr
				},
			},
			storage: &stubStorageService{
				uploadFn: func(context.Context, string, string, string, []byte, string) (string, error) {
					return newURL, nil
				},
			},
			wantErr:        updateErr,
			wantUpload:     true,
			wantUpdate:     true,
			wantDeleteRoll: true,
		},
		{
			name: "get fails with generic error",
			repo: &stubServiceRepository{
				getByIDFn: func(uuid.UUID, uuid.UUID) (domain.ServiceResponseDTO, error) {
					return domain.ServiceResponseDTO{}, getErr
				},
			},
			storage:    &stubStorageService{},
			wantErr:    getErr,
			wantUpload: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			strategy := NewServiceImageStrategy(tc.storage, tc.repo, bucketName)
			url, err := strategy.Upload(context.Background(), requesterID, serviceID, "a.png", []byte("x"), "image/png")

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error: want %v, got %v", tc.wantErr, err)
			}
			if tc.wantUpload != tc.storage.uploadCalled {
				t.Fatalf("uploadCalled: want %v, got %v", tc.wantUpload, tc.storage.uploadCalled)
			}
			if tc.wantUpdate != tc.repo.updateCalled {
				t.Fatalf("updateCalled: want %v, got %v", tc.wantUpdate, tc.repo.updateCalled)
			}
			if tc.wantDeleteRoll && !tc.storage.deleteCalled {
				t.Fatal("expected rollback delete")
			}
			if tc.wantPhotosLen > 0 && len(tc.repo.lastDTO.Photos) != tc.wantPhotosLen {
				t.Fatalf("photos len: want %d, got %d", tc.wantPhotosLen, len(tc.repo.lastDTO.Photos))
			}
			if tc.wantErr == nil && url == "" {
				t.Fatal("expected non-empty url")
			}
		})
	}
}
