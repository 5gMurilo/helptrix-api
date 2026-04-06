package uploader

import (
	"context"
	"errors"
	"testing"

	uploaderinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/uploader"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

// --- stub ---

type stubStrategy struct {
	uploadFn func(ctx context.Context, requesterID, ownerID uuid.UUID, filename string, data []byte, contentType string) (string, error)
}

func (s *stubStrategy) Upload(ctx context.Context, requesterID, ownerID uuid.UUID, filename string, data []byte, contentType string) (string, error) {
	if s.uploadFn != nil {
		return s.uploadFn(ctx, requesterID, ownerID, filename, data, contentType)
	}
	return "https://storage.example.com/image.jpg", nil
}

// toStrategyMap wraps stubStrategy pointers in the interface map type expected by UploaderService.
func toStrategyMap(m map[string]*stubStrategy) map[string]uploaderinterfaces.IImageUploadStrategy {
	result := make(map[string]uploaderinterfaces.IImageUploadStrategy, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// --- tests ---

func TestUploaderService_Upload(t *testing.T) {
	requesterID := uuid.New()
	ownerID := uuid.New()
	strategyErr := errors.New("strategy error")

	successStrategy := &stubStrategy{}
	failStrategy := &stubStrategy{
		uploadFn: func(_ context.Context, _, _ uuid.UUID, _ string, _ []byte, _ string) (string, error) {
			return "", strategyErr
		},
	}

	tests := []struct {
		name       string
		imageType  string
		strategies map[string]*stubStrategy
		wantErr    error
		wantURL    bool
	}{
		{
			name:      "unknown imageType returns ErrInvalidImageType",
			imageType: "unknown-type",
			strategies: map[string]*stubStrategy{
				"profile-images": successStrategy,
			},
			wantErr: utils.ErrInvalidImageType,
			wantURL: false,
		},
		{
			name:      "valid imageType delegates to strategy and returns URL",
			imageType: "profile-images",
			strategies: map[string]*stubStrategy{
				"profile-images": successStrategy,
			},
			wantErr: nil,
			wantURL: true,
		},
		{
			name:      "strategy error is propagated",
			imageType: "profile-images",
			strategies: map[string]*stubStrategy{
				"profile-images": failStrategy,
			},
			wantErr: strategyErr,
			wantURL: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &UploaderService{strategies: toStrategyMap(tc.strategies)}

			url, err := svc.Upload(
				context.Background(),
				tc.imageType,
				requesterID,
				ownerID,
				"photo.jpg",
				[]byte("data"),
				"image/jpeg",
			)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
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
