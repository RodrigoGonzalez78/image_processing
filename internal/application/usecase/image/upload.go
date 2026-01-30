package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/domain/repository"
)

// UploadUseCase maneja el caso de uso de subida de imágenes
type UploadUseCase struct {
	imageRepo   repository.ImageRepository
	fileStorage repository.FileStorage
	baseURL     string
	port        string
}

// NewUploadUseCase crea una nueva instancia de UploadUseCase
func NewUploadUseCase(
	imageRepo repository.ImageRepository,
	fileStorage repository.FileStorage,
	baseURL, port string,
) *UploadUseCase {
	return &UploadUseCase{
		imageRepo:   imageRepo,
		fileStorage: fileStorage,
		baseURL:     baseURL,
		port:        port,
	}
}

// UploadInput representa los datos de entrada para la subida
type UploadInput struct {
	FileName    string
	UserName    string
	Data        []byte
	ContentType string
	Format      string
	Width       int
	Height      int
}

// UploadOutput representa los datos de salida de la subida
type UploadOutput struct {
	URL    string
	Name   string
	Size   int64
	Format string
	Width  int
	Height int
}

// Execute ejecuta el caso de uso de subida de imagen
func (uc *UploadUseCase) Execute(ctx context.Context, input UploadInput) (*UploadOutput, error) {
	if len(input.Data) == 0 {
		return nil, errors.New("no se proporcionó datos de imagen")
	}

	// Construir path del objeto
	objectPath := fmt.Sprintf("%s/%s", input.UserName, input.FileName)

	// Subir a storage
	err := uc.fileStorage.Upload(ctx, objectPath, input.Data, input.ContentType)
	if err != nil {
		return nil, fmt.Errorf("error al subir imagen: %w", err)
	}

	// Crear entidad de imagen
	image := entity.NewImage(
		input.FileName,
		input.UserName,
		objectPath,
		input.Format,
		int64(len(input.Data)),
		input.Width,
		input.Height,
	)

	// Guardar en base de datos
	err = uc.imageRepo.Create(image)
	if err != nil {
		return nil, fmt.Errorf("error al guardar imagen en BD: %w", err)
	}

	// Construir URL de la imagen
	url := fmt.Sprintf("%s:%s/images/%s/%s", uc.baseURL, uc.port, input.UserName, input.FileName)

	return &UploadOutput{
		URL:    url,
		Name:   input.FileName,
		Size:   int64(len(input.Data)),
		Format: input.Format,
		Width:  input.Width,
		Height: input.Height,
	}, nil
}
