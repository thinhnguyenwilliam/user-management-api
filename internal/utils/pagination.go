// user-management-api/internal/utils/pagination.go
package utils

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	totalPages := int32(0)
	if totalRecords > 0 {
		totalPages = (totalRecords + limit - 1) / limit
	}

	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

type PaginationResponse struct {
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

func NewPaginationResponse(data any, page, limit, totalRecords int32) PaginationResponse {
	return PaginationResponse{
		Data:       data,
		Pagination: NewPagination(page, limit, totalRecords),
	}
}
