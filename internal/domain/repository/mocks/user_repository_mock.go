package mocks

import (
	"errors"

	"github.com/RodrigoGonzalez78/internal/domain/entity"
)

// MockUserRepository es un mock del repositorio de usuarios para testing
type MockUserRepository struct {
	Users        map[string]*entity.User
	CreateError  error
	FindError    error
	ExistsError  error
	CreateCalled bool
	FindCalled   bool
	ExistsCalled bool
}

// NewMockUserRepository crea un nuevo mock de UserRepository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[string]*entity.User),
	}
}

// Create simula la creación de un usuario
func (m *MockUserRepository) Create(user *entity.User) error {
	m.CreateCalled = true
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Users[user.UserName] = user
	return nil
}

// FindByUserName simula buscar un usuario por nombre
func (m *MockUserRepository) FindByUserName(userName string) (*entity.User, error) {
	m.FindCalled = true
	if m.FindError != nil {
		return nil, m.FindError
	}
	user, exists := m.Users[userName]
	if !exists {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

// ExistsByUserName simula verificar si existe un usuario
func (m *MockUserRepository) ExistsByUserName(userName string) (bool, error) {
	m.ExistsCalled = true
	if m.ExistsError != nil {
		return false, m.ExistsError
	}
	_, exists := m.Users[userName]
	return exists, nil
}
