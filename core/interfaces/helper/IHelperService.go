package helperinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type IHelperService interface {
	Search(requesterType string, params domain.HelperSearchParams) (domain.HelperListResponseDTO, error)
}
