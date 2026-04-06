package strategies

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

// --- stubs ---

type stubStorageService struct {
	uploadFn func(ctx context.Context, folder, ownerID, filename string, data []byte, contentType string) (string, error)
	deleteFn func(ctx context.Context, objectPath string) error

	uploadCalled bool
	deleteCalled bool
	deleteArg    string
}

func (s *stubStorageService) UploadFile(ctx context.Context, folder, ownerID, filename string, data []byte, contentType string) (string, error) {
	s.uploadCalled = true
	if s.uploadFn != nil {
		return s.uploadFn(ctx, folder, ownerID, filename, data, contentType)
	}
	return fmt.Sprintf("https://storage.googleapis.com/bucket/%s/%s/%s", folder, ownerID, filename), nil
}

func (s *stubStorageService) DeleteFile(ctx context.Context, objectPath string) error {
	s.deleteCalled = true
	s.deleteArg = objectPath
	if s.deleteFn != nil {
		return s.deleteFn(ctx, objectPath)
	}
	return nil
}

type stubUserRepository struct {
	getProfilePictureFn func(userID uuid.UUID) (string, error)
	updatePictureFn     func(userID uuid.UUID, url string) error

	updateCalled bool
}

func (r *stubUserRepository) GetProfile(userID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
	return domain.GetProfileResponseDTO{}, nil
}

func (r *stubUserRepository) UpdateProfile(userID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
	return nil
}

func (r *stubUserRepository) DeleteProfile(userID uuid.UUID) error { return nil }

func (r *stubUserRepository) GetProfilePicture(userID uuid.UUID) (string, error) {
	if r.getProfilePictureFn != nil {
		return r.getProfilePictureFn(userID)
	}
	return "", nil
}

func (r *stubUserRepository) UpdateProfilePicture(userID uuid.UUID, url string) error {
	r.updateCalled = true
	if r.updatePictureFn != nil {
		return r.updatePictureFn(userID, url)
	}
	return nil
}

// --- tests ---

func TestProfileImageStrategy_Upload(t *testing.T) {
	bucketName := "my-bucket"
	ownerID := uuid.New()
	requesterID := ownerID
	differentID := uuid.New()

	existingURL := fmt.Sprintf("https://storage.googleapis.com/%s/profile-images/%s/old.jpg", bucketName, ownerID)
	expectedObjectPath := fmt.Sprintf("profile-images/%s/old.jpg", ownerID)

	uploadErr := errors.New("upload failed")
	deleteErr := errors.New("delete failed")
	updateErr := errors.New("update failed")
	getUserErr := errors.New("get user failed")

	tests := []struct {
		name          string
		requesterID   uuid.UUID
		ownerID       uuid.UUID
		storageSvc    *stubStorageService
		userRepo      *stubUserRepository
		wantErr       error
		wantDeleteArg string
		wantDelete    bool
		wantUpload    bool
		wantUpdate    bool
		wantURL       bool
	}{
		{
			name:        "no existing image — skip delete, upload and update succeed",
			requesterID: requesterID,
			ownerID:     ownerID,
			storageSvc:  &stubStorageService{},
			userRepo: &stubUserRepository{
				getProfilePictureFn: func(_ uuid.UUID) (string, error) { return "", nil },
			},
			wantErr:    nil,
			wantDelete: false,
			wantUpload: true,
			wantUpdate: true,
			wantURL:    true,
		},
		{
			name:        "existing image — delete called with correct path, then upload and update",
			requesterID: requesterID,
			ownerID:     ownerID,
			storageSvc:  &stubStorageService{},
			userRepo: &stubUserRepository{
				getProfilePictureFn: func(_ uuid.UUID) (string, error) { return existingURL, nil },
			},
			wantErr:       nil,
			wantDelete:    true,
			wantDeleteArg: expectedObjectPath,
			wantUpload:    true,
			wantUpdate:    true,
			wantURL:       true,
		},
		{
			name:        "requester is not owner — ErrNotOwner, no storage or repo calls",
			requesterID: differentID,
			ownerID:     ownerID,
			storageSvc:  &stubStorageService{},
			userRepo:    &stubUserRepository{},
			wantErr:     utils.ErrNotOwner,
			wantDelete:  false,
			wantUpload:  false,
			wantUpdate:  false,
		},
		{
			name:        "GetProfilePicture fails — error propagated",
			requesterID: requesterID,
			ownerID:     ownerID,
			storageSvc:  &stubStorageService{},
			userRepo: &stubUserRepository{
				getProfilePictureFn: func(_ uuid.UUID) (string, error) { return "", getUserErr },
			},
			wantErr:    getUserErr,
			wantDelete: false,
			wantUpload: false,
			wantUpdate: false,
		},
		{
			name:        "delete fails — error propagated",
			requesterID: requesterID,
			ownerID:     ownerID,
			storageSvc: &stubStorageService{
				deleteFn: func(_ context.Context, _ string) error { return deleteErr },
			},
			userRepo: &stubUserRepository{
				getProfilePictureFn: func(_ uuid.UUID) (string, error) { return existingURL, nil },
			},
			wantErr:    deleteErr,
			wantDelete: true,
			wantUpload: false,
			wantUpdate: false,
		},
		{
			name:        "upload fails — error propagated",
			requesterID: requesterID,
			ownerID:     ownerID,
			storageSvc: &stubStorageService{
				uploadFn: func(_ context.Context, _, _, _ string, _ []byte, _ string) (string, error) {
					return "", uploadErr
				},
			},
			userRepo: &stubUserRepository{
				getProfilePictureFn: func(_ uuid.UUID) (string, error) { return "", nil },
			},
			wantErr:    uploadErr,
			wantDelete: false,
			wantUpload: true,
			wantUpdate: false,
		},
		{
			name:        "UpdateProfilePicture fails — error propagated",
			requesterID: requesterID,
			ownerID:     ownerID,
			storageSvc:  &stubStorageService{},
			userRepo: &stubUserRepository{
				getProfilePictureFn: func(_ uuid.UUID) (string, error) { return "", nil },
				updatePictureFn:     func(_ uuid.UUID, _ string) error { return updateErr },
			},
			wantErr:    updateErr,
			wantDelete: false,
			wantUpload: true,
			wantUpdate: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			strategy := NewProfileImageStrategy(tc.storageSvc, tc.userRepo, bucketName)

			url, err := strategy.Upload(
				context.Background(),
				tc.requesterID,
				tc.ownerID,
				"photo.jpg",
				[]byte("data"),
				"image/jpeg",
			)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}

			if tc.wantDelete != tc.storageSvc.deleteCalled {
				t.Errorf("expected deleteCalled=%v, got %v", tc.wantDelete, tc.storageSvc.deleteCalled)
			}

			if tc.wantDeleteArg != "" && tc.storageSvc.deleteArg != tc.wantDeleteArg {
				t.Errorf("expected deleteArg=%q, got %q", tc.wantDeleteArg, tc.storageSvc.deleteArg)
			}

			if tc.wantUpload != tc.storageSvc.uploadCalled {
				t.Errorf("expected uploadCalled=%v, got %v", tc.wantUpload, tc.storageSvc.uploadCalled)
			}

			if tc.wantUpdate != tc.userRepo.updateCalled {
				t.Errorf("expected updateCalled=%v, got %v", tc.wantUpdate, tc.userRepo.updateCalled)
			}

			if tc.wantURL && url == "" {
				t.Error("expected non-empty URL, got empty string")
			}

			if !tc.wantURL && url != "" {
				t.Errorf("expected empty URL, got %q", url)
			}
		})
	}
}
