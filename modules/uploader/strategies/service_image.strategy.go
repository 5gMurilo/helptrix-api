package strategies

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ServiceImageStrategy struct{}

func NewServiceImageStrategy() *ServiceImageStrategy {
	return &ServiceImageStrategy{}
}

func (s *ServiceImageStrategy) Upload(
	ctx context.Context,
	requesterID uuid.UUID,
	ownerID uuid.UUID,
	filename string,
	data []byte,
	contentType string,
) (string, error) {
	return "", errors.New("not implemented")
}
