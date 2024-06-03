package dto

type PaginationMetadata struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalDocs   int  `json:"total_docs"`
	TotalPages  int  `json:"total_pages"`
	HasNextPage bool `json:"has_next_page"`
}

type Pagination[T any] struct {
	Code     int                `json:"code"`
	Data     T                  `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}
