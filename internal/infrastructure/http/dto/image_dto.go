package dto

type ErrorResponse struct {
	Message string `json:"message" example:"Descripción del error"`
}

type UploadedImageDetail struct {
	URL    string `json:"url" example:"http://localhost:8080/images/username/imagen123.jpg"`
	Name   string `json:"name" example:"imagen123.jpg"`
	Size   int64  `json:"size" example:"204800"`
	Format string `json:"format" example:"jpeg"`
	Width  int    `json:"width" example:"1920"`
	Height int    `json:"height" example:"1080"`
}

type UploadResponse struct {
	Message string              `json:"message" example:"Imagen subida exitosamente"`
	Image   UploadedImageDetail `json:"image"`
}

type ImageDetailResponse struct {
	URL    string `json:"url" example:"http://localhost:8080/images/rodrick/imagen123.jpg"`
	Name   string `json:"name" example:"imagen123.jpg"`
	Size   int64  `json:"size" example:"204800"`
	Format string `json:"format" example:"jpeg"`
	Width  int    `json:"width" example:"1920"`
	Height int    `json:"height" example:"1080"`
}

type ImageItem struct {
	ID     int64  `json:"id"`
	URL    string `json:"url"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PaginatedImagesResponse struct {
	Page   int         `json:"page"`
	Limit  int         `json:"limit"`
	Total  int64       `json:"total"`
	Images []ImageItem `json:"images"`
}

type TransformationRequest struct {
	Transformations struct {
		Resize struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"resize"`
		Crop struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			X      int `json:"x"`
			Y      int `json:"y"`
		} `json:"crop"`
		Rotate  float64 `json:"rotate"`
		Format  string  `json:"format"`
		Filters struct {
			Grayscale bool `json:"grayscale"`
			Sepia     bool `json:"sepia"`
		} `json:"filters"`
	} `json:"transformations"`
}
