package models

type PaginatedImagesResponse struct {
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
	Total  int64   `json:"total"`
	Images []Image `json:"images"`
}
