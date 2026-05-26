package domain

const DefaultPageSize = 20

type PaginationParams struct {
	Page     int
	PageSize int
}
