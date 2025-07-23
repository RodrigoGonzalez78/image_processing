package models

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
