package mocks

import (
	"context"
	"io"
	"strings"
)

// MockFileStorage es un mock del storage de archivos para testing
type MockFileStorage struct {
	Files        map[string][]byte
	UploadError  error
	GetError     error
	UploadCalled bool
	GetCalled    bool
}

// NewMockFileStorage crea un nuevo mock de FileStorage
func NewMockFileStorage() *MockFileStorage {
	return &MockFileStorage{
		Files: make(map[string][]byte),
	}
}

// Upload simula subir un archivo
func (m *MockFileStorage) Upload(ctx context.Context, path string, data []byte, contentType string) error {
	m.UploadCalled = true
	if m.UploadError != nil {
		return m.UploadError
	}
	m.Files[path] = data
	return nil
}

// Get simula obtener un archivo
func (m *MockFileStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	m.GetCalled = true
	if m.GetError != nil {
		return nil, m.GetError
	}
	data, exists := m.Files[path]
	if !exists {
		return nil, m.GetError
	}
	return io.NopCloser(strings.NewReader(string(data))), nil
}
