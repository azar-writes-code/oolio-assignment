package models

// PaginationParams is used to bind pagination query parameters from the client.
type PaginationParams struct {
	Page     int32 `form:"page" binding:"omitempty,min=1"`
	PageSize int32 `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PaginatedResponse is a generic wrapper to provide metadata about the current page and total records.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Page       int32       `json:"page"`
	PageSize   int32       `json:"page_size"`
	TotalPages int32       `json:"total_pages"`
}
