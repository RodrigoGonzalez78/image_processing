package image

import (
	"errors"
	"fmt"

	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/domain/repository"
)

// GetUseCase maneja el caso de uso de obtener una imagen
type GetUseCase struct {
	imageRepo repository.ImageRepository
	baseURL   string
	port      string
}

// NewGetUseCase crea una nueva instancia de GetUseCase
func NewGetUseCase(imageRepo repository.ImageRepository, baseURL, port string) *GetUseCase {
	return &GetUseCase{
		imageRepo: imageRepo,
		baseURL:   baseURL,
		port:      port,
	}
}

// GetInput representa los datos de entrada
type GetInput struct {
	ImageID  int64
	UserName string
}

// GetOutput representa los datos de salida
type GetOutput struct {
	URL    string
	Name   string
	Size   int64
	Format string
	Width  int
	Height int
}

// Execute ejecuta el caso de uso
func (uc *GetUseCase) Execute(input GetInput) (*GetOutput, error) {
	image, err := uc.imageRepo.FindByID(input.ImageID)
	if err != nil {
		return nil, errors.New("imagen no encontrada")
	}

	// Verificar que la imagen pertenece al usuario
	if image.UserName != input.UserName {
		return nil, errors.New("no tienes permiso para esta imagen")
	}

	url := fmt.Sprintf("%s:%s/images/%s/%s", uc.baseURL, uc.port, image.UserName, image.Name)

	return &GetOutput{
		URL:    url,
		Name:   image.Name,
		Size:   image.Size,
		Format: image.Format,
		Width:  image.Width,
		Height: image.Height,
	}, nil
}

// GetImageEntity obtiene la entidad completa de la imagen
func (uc *GetUseCase) GetImageEntity(imageID int64, userName string) (*entity.Image, error) {
	image, err := uc.imageRepo.FindByID(imageID)
	if err != nil {
		return nil, errors.New("imagen no encontrada")
	}

	if image.UserName != userName {
		return nil, errors.New("no tienes permiso para esta imagen")
	}

	return image, nil
}
