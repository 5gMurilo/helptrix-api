package helper

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	helperinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/helper"
	"github.com/5gMurilo/helptrix-api/core/utils"
)

type HelperService struct {
	repo helperinterfaces.IHelperRepository
}

func NewHelperService(repo helperinterfaces.IHelperRepository) helperinterfaces.IHelperService {
	return &HelperService{repo: repo}
}

func (s *HelperService) Search(requesterType string, params domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
	if requesterType != utils.UserTypeBusiness {
		return domain.HelperListResponseDTO{}, utils.ErrBusinessOnly
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	return s.repo.Search(params)
}
