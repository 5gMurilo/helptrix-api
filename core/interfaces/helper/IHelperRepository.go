package helperinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type IHelperRepository interface {
	Search(params domain.HelperSearchParams) (domain.HelperListResponseDTO, error)
}
