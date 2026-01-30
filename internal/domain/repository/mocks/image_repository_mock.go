package mocks

import (
	"errors"

	"github.com/RodrigoGonzalez78/internal/domain/entity"
)

// MockImageRepository es un mock del repositorio de imágenes para testing
type MockImageRepository struct {
	Images           map[int64]*entity.Image
	NextID           int64
	CreateError      error
	FindByIDError    error
	FindByUserError  error
	CreateCalled     bool
	FindByIDCalled   bool
	FindByUserCalled bool
}

// NewMockImageRepository crea un nuevo mock de ImageRepository
func NewMockImageRepository() *MockImageRepository {
	return &MockImageRepository{
		Images: make(map[int64]*entity.Image),
		NextID: 1,
	}
}

// Create simula la creación de una imagen
func (m *MockImageRepository) Create(image *entity.Image) error {
	m.CreateCalled = true
	if m.CreateError != nil {
		return m.CreateError
	}
	image.ID = m.NextID
	m.Images[m.NextID] = image
	m.NextID++
	return nil
}

// FindByID simula buscar una imagen por ID
func (m *MockImageRepository) FindByID(id int64) (*entity.Image, error) {
	m.FindByIDCalled = true
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}
	image, exists := m.Images[id]
	if !exists {
		return nil, errors.New("imagen no encontrada")
	}
	return image, nil
}

// FindByUser simula obtener imágenes de un usuario
func (m *MockImageRepository) FindByUser(userName string, page, limit int) ([]entity.Image, int64, error) {
	m.FindByUserCalled = true
	if m.FindByUserError != nil {
		return nil, 0, m.FindByUserError
	}

	var result []entity.Image
	for _, img := range m.Images {
		if img.UserName == userName {
			result = append(result, *img)
		}
	}

	// Simular paginación simple
	start := (page - 1) * limit
	end := start + limit
	if start > len(result) {
		return []entity.Image{}, int64(len(result)), nil
	}
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], int64(len(result)), nil
}
