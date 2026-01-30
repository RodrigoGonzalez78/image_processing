package image

import (
	"context"
	"errors"
	"image"
	"image/color"
	"io"

	"github.com/RodrigoGonzalez78/internal/domain/repository"
	"github.com/disintegration/imaging"
)

// TransformUseCase maneja el caso de uso de transformación de imágenes
type TransformUseCase struct {
	imageRepo   repository.ImageRepository
	fileStorage repository.FileStorage
}

// NewTransformUseCase crea una nueva instancia de TransformUseCase
func NewTransformUseCase(
	imageRepo repository.ImageRepository,
	fileStorage repository.FileStorage,
) *TransformUseCase {
	return &TransformUseCase{
		imageRepo:   imageRepo,
		fileStorage: fileStorage,
	}
}

// ResizeParams parámetros de redimensionado
type ResizeParams struct {
	Width  int
	Height int
}

// CropParams parámetros de recorte
type CropParams struct {
	Width  int
	Height int
	X      int
	Y      int
}

// FilterParams parámetros de filtros
type FilterParams struct {
	Grayscale bool
	Sepia     bool
}

// TransformInput representa los datos de entrada
type TransformInput struct {
	ImageID  int64
	UserName string
	Resize   ResizeParams
	Crop     CropParams
	Rotate   float64
	Format   string
	Filters  FilterParams
}

// TransformOutput representa los datos de salida
type TransformOutput struct {
	Image       image.Image
	Format      imaging.Format
	ContentType string
}

// Execute ejecuta el caso de uso
func (uc *TransformUseCase) Execute(ctx context.Context, input TransformInput) (*TransformOutput, error) {
	// Obtener metadata de la imagen
	imageData, err := uc.imageRepo.FindByID(input.ImageID)
	if err != nil {
		return nil, errors.New("imagen no encontrada")
	}

	// Verificar permisos (opcional, puede manejarse en el handler)
	if imageData.UserName != input.UserName {
		return nil, errors.New("no tienes permiso para esta imagen")
	}

	// Obtener imagen del storage
	reader, err := uc.fileStorage.Get(ctx, imageData.Path)
	if err != nil {
		return nil, errors.New("no se pudo leer la imagen desde el storage")
	}
	defer reader.Close()

	// Decodificar imagen
	srcImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, errors.New("error al decodificar la imagen")
	}

	dstImage := srcImage

	// Aplicar transformaciones
	if input.Resize.Width > 0 && input.Resize.Height > 0 {
		dstImage = imaging.Resize(dstImage, input.Resize.Width, input.Resize.Height, imaging.Lanczos)
	}

	if input.Crop.Width > 0 && input.Crop.Height > 0 {
		cropRect := image.Rect(
			input.Crop.X,
			input.Crop.Y,
			input.Crop.X+input.Crop.Width,
			input.Crop.Y+input.Crop.Height,
		)
		dstImage = imaging.Crop(dstImage, cropRect)
	}

	if input.Rotate != 0 {
		dstImage = imaging.Rotate(dstImage, input.Rotate, color.Transparent)
	}

	if input.Filters.Grayscale {
		dstImage = imaging.Grayscale(dstImage)
	}

	if input.Filters.Sepia {
		dstImage = imaging.AdjustSaturation(dstImage, -100)
		dstImage = imaging.AdjustContrast(dstImage, 10)
	}

	// Determinar formato de salida
	var imgFormat imaging.Format
	var contentType string

	switch input.Format {
	case "jpg", "jpeg":
		imgFormat = imaging.JPEG
		contentType = "image/jpeg"
	case "gif":
		imgFormat = imaging.GIF
		contentType = "image/gif"
	default:
		imgFormat = imaging.PNG
		contentType = "image/png"
	}

	return &TransformOutput{
		Image:       dstImage,
		Format:      imgFormat,
		ContentType: contentType,
	}, nil
}

// ServeImage obtiene una imagen del storage para servirla
func (uc *TransformUseCase) ServeImage(ctx context.Context, objectPath string) (io.ReadCloser, error) {
	return uc.fileStorage.Get(ctx, objectPath)
}
