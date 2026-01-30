package image

import (
	"fmt"

	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/domain/repository"
)

// ListUseCase maneja el caso de uso de listar imágenes
type ListUseCase struct {
	imageRepo repository.ImageRepository
	baseURL   string
	port      string
}

// NewListUseCase crea una nueva instancia de ListUseCase
func NewListUseCase(imageRepo repository.ImageRepository, baseURL, port string) *ListUseCase {
	return &ListUseCase{
		imageRepo: imageRepo,
		baseURL:   baseURL,
		port:      port,
	}
}

// ListInput representa los datos de entrada
type ListInput struct {
	UserName string
	Page     int
	Limit    int
}

// ImageItem representa un item de imagen en la lista
type ImageItem struct {
	ID     int64
	URL    string
	Name   string
	Size   int64
	Format string
	Width  int
	Height int
}

// ListOutput representa los datos de salida
type ListOutput struct {
	Page   int
	Limit  int
	Total  int64
	Images []ImageItem
}

// Execute ejecuta el caso de uso
func (uc *ListUseCase) Execute(input ListInput) (*ListOutput, error) {
	// Valores por defecto
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Limit < 1 {
		input.Limit = 10
	}

	images, total, err := uc.imageRepo.FindByUser(input.UserName, input.Page, input.Limit)
	if err != nil {
		return nil, err
	}

	// Convertir a output
	items := make([]ImageItem, len(images))
	for i, img := range images {
		items[i] = uc.toImageItem(&img)
	}

	return &ListOutput{
		Page:   input.Page,
		Limit:  input.Limit,
		Total:  total,
		Images: items,
	}, nil
}

func (uc *ListUseCase) toImageItem(img *entity.Image) ImageItem {
	url := fmt.Sprintf("%s:%s/images/%s/%s", uc.baseURL, uc.port, img.UserName, img.Name)
	return ImageItem{
		ID:     img.ID,
		URL:    url,
		Name:   img.Name,
		Size:   img.Size,
		Format: img.Format,
		Width:  img.Width,
		Height: img.Height,
	}
}
