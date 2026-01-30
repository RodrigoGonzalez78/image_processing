package image_test

import (
	"context"
	"testing"

	"github.com/RodrigoGonzalez78/internal/application/usecase/image"
	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/domain/repository/mocks"
)

func TestUploadUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		input       image.UploadInput
		setupMocks  func(*mocks.MockImageRepository, *mocks.MockFileStorage)
		wantErr     bool
		expectedURL string
	}{
		{
			name: "subida exitosa",
			input: image.UploadInput{
				FileName:    "test.jpg",
				UserName:    "testuser",
				Data:        []byte("fake image data"),
				ContentType: "image/jpeg",
				Format:      "jpeg",
				Width:       800,
				Height:      600,
			},
			setupMocks: func(ir *mocks.MockImageRepository, fs *mocks.MockFileStorage) {},
			wantErr:    false,
		},
		{
			name: "datos vacíos",
			input: image.UploadInput{
				FileName:    "test.jpg",
				UserName:    "testuser",
				Data:        []byte{},
				ContentType: "image/jpeg",
				Format:      "jpeg",
				Width:       800,
				Height:      600,
			},
			setupMocks: func(ir *mocks.MockImageRepository, fs *mocks.MockFileStorage) {},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockImageRepo := mocks.NewMockImageRepository()
			mockFileStorage := mocks.NewMockFileStorage()
			tt.setupMocks(mockImageRepo, mockFileStorage)
			useCase := image.NewUploadUseCase(mockImageRepo, mockFileStorage, "http://localhost", "8080")

			// Execute
			output, err := useCase.Execute(context.Background(), tt.input)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !mockFileStorage.UploadCalled {
					t.Error("FileStorage.Upload() no fue llamado")
				}
				if !mockImageRepo.CreateCalled {
					t.Error("ImageRepository.Create() no fue llamado")
				}
				if output.Name != tt.input.FileName {
					t.Errorf("output.Name = %v, want %v", output.Name, tt.input.FileName)
				}
			}
		})
	}
}

func TestGetUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     image.GetInput
		setupMock func(*mocks.MockImageRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "obtener imagen exitoso",
			input: image.GetInput{
				ImageID:  1,
				UserName: "testuser",
			},
			setupMock: func(m *mocks.MockImageRepository) {
				m.Images[1] = &entity.Image{
					ID:       1,
					Name:     "test.jpg",
					UserName: "testuser",
					Path:     "testuser/test.jpg",
					Size:     1024,
					Format:   "jpeg",
					Width:    800,
					Height:   600,
				}
			},
			wantErr: false,
		},
		{
			name: "imagen no encontrada",
			input: image.GetInput{
				ImageID:  999,
				UserName: "testuser",
			},
			setupMock: func(m *mocks.MockImageRepository) {},
			wantErr:   true,
			errMsg:    "imagen no encontrada",
		},
		{
			name: "imagen de otro usuario",
			input: image.GetInput{
				ImageID:  1,
				UserName: "otroUsuario",
			},
			setupMock: func(m *mocks.MockImageRepository) {
				m.Images[1] = &entity.Image{
					ID:       1,
					Name:     "test.jpg",
					UserName: "testuser", // Usuario diferente
					Path:     "testuser/test.jpg",
				}
			},
			wantErr: true,
			errMsg:  "no tienes permiso para esta imagen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockImageRepository()
			tt.setupMock(mockRepo)
			useCase := image.NewGetUseCase(mockRepo, "http://localhost", "8080")

			// Execute
			output, err := useCase.Execute(tt.input)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Execute() error = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr && output == nil {
				t.Error("Execute() retornó nil")
			}
		})
	}
}

func TestListUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		input         image.ListInput
		setupMock     func(*mocks.MockImageRepository)
		wantErr       bool
		expectedTotal int64
	}{
		{
			name: "listar imágenes exitoso",
			input: image.ListInput{
				UserName: "testuser",
				Page:     1,
				Limit:    10,
			},
			setupMock: func(m *mocks.MockImageRepository) {
				m.Images[1] = &entity.Image{ID: 1, Name: "img1.jpg", UserName: "testuser"}
				m.Images[2] = &entity.Image{ID: 2, Name: "img2.jpg", UserName: "testuser"}
			},
			wantErr:       false,
			expectedTotal: 2,
		},
		{
			name: "sin imágenes",
			input: image.ListInput{
				UserName: "nuevoUsuario",
				Page:     1,
				Limit:    10,
			},
			setupMock:     func(m *mocks.MockImageRepository) {},
			wantErr:       false,
			expectedTotal: 0,
		},
		{
			name: "página por defecto cuando es 0",
			input: image.ListInput{
				UserName: "testuser",
				Page:     0,
				Limit:    0,
			},
			setupMock: func(m *mocks.MockImageRepository) {
				m.Images[1] = &entity.Image{ID: 1, Name: "img1.jpg", UserName: "testuser"}
			},
			wantErr:       false,
			expectedTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockImageRepository()
			tt.setupMock(mockRepo)
			useCase := image.NewListUseCase(mockRepo, "http://localhost", "8080")

			// Execute
			output, err := useCase.Execute(tt.input)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if output.Total != tt.expectedTotal {
					t.Errorf("Execute() Total = %v, want %v", output.Total, tt.expectedTotal)
				}
			}
		})
	}
}
