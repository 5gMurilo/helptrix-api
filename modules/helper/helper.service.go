package helper

import (
	"log/slog"

	"github.com/5gMurilo/helptrix-api/core/domain"
	helperinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/helper"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
)

type HelperService struct {
	repo helperinterfaces.IHelperRepository
}

func NewHelperService(repo helperinterfaces.IHelperRepository) helperinterfaces.IHelperService {
	return &HelperService{repo: repo}
}

func (s *HelperService) Search(requesterType string, params domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "helper"),
		slog.String("operation", "Search"),
	)

	if requesterType != utils.UserTypeBusiness {
		log.Warn("helper search rejected: only business users can search helpers", slog.String("requester_type", requesterType))
		return domain.HelperListResponseDTO{}, utils.ErrBusinessOnly
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	result, err := s.repo.Search(params)
	if err != nil {
		log.Error("failed to search helpers", slog.String("error", err.Error()))
		return domain.HelperListResponseDTO{}, err
	}

	return result, nil
}
