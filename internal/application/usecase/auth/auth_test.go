package auth_test

import (
	"testing"

	"github.com/RodrigoGonzalez78/internal/application/usecase/auth"
	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/domain/repository/mocks"
	"github.com/RodrigoGonzalez78/internal/domain/service"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     auth.RegisterInput
		setupMock func(*mocks.MockUserRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "registro exitoso",
			input: auth.RegisterInput{
				UserName: "nuevoUsuario",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {},
			wantErr:   false,
		},
		{
			name: "nombre de usuario vacío",
			input: auth.RegisterInput{
				UserName: "",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {},
			wantErr:   true,
			errMsg:    "el nombre de usuario es requerido",
		},
		{
			name: "contraseña muy corta",
			input: auth.RegisterInput{
				UserName: "usuario",
				Password: "1234567",
			},
			setupMock: func(m *mocks.MockUserRepository) {},
			wantErr:   true,
			errMsg:    "la contraseña debe tener al menos 8 caracteres",
		},
		{
			name: "usuario ya existe",
			input: auth.RegisterInput{
				UserName: "usuarioExistente",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.Users["usuarioExistente"] = nil // Simular que ya existe
			},
			wantErr: true,
			errMsg:  "el nombre de usuario ya está registrado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUserRepository()
			tt.setupMock(mockRepo)
			passwordService := service.NewPasswordService()
			useCase := auth.NewRegisterUseCase(mockRepo, passwordService)

			// Execute
			err := useCase.Execute(tt.input)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Execute() error message = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr {
				if !mockRepo.CreateCalled {
					t.Error("Create() no fue llamado")
				}
				if mockRepo.Users[tt.input.UserName] == nil {
					t.Error("Usuario no fue guardado en el repositorio")
				}
			}
		})
	}
}

func TestLoginUseCase_Execute(t *testing.T) {
	// Setup: crear usuario con contraseña hasheada
	passwordService := service.NewPasswordService()
	hashedPassword, _ := passwordService.Hash("password123")

	tests := []struct {
		name      string
		input     auth.LoginInput
		setupMock func(*mocks.MockUserRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "login exitoso",
			input: auth.LoginInput{
				UserName: "testuser",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.Users["testuser"] = &entity.User{
					UserName: "testuser",
					Password: hashedPassword,
				}
			},
			wantErr: false,
		},
		{
			name: "nombre de usuario vacío",
			input: auth.LoginInput{
				UserName: "",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {},
			wantErr:   true,
			errMsg:    "el nombre de usuario es requerido",
		},
		{
			name: "contraseña vacía",
			input: auth.LoginInput{
				UserName: "testuser",
				Password: "",
			},
			setupMock: func(m *mocks.MockUserRepository) {},
			wantErr:   true,
			errMsg:    "la contraseña es requerida",
		},
		{
			name: "usuario no encontrado",
			input: auth.LoginInput{
				UserName: "noexiste",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {},
			wantErr:   true,
			errMsg:    "usuario no encontrado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUserRepository()
			tt.setupMock(mockRepo)
			tokenService := service.NewTokenService("test-secret", 24*60*60*1000000000)
			useCase := auth.NewLoginUseCase(mockRepo, passwordService, tokenService)

			// Execute
			output, err := useCase.Execute(tt.input)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Execute() error message = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr && output.Token == "" {
				t.Error("Execute() retornó un token vacío")
			}
		})
	}
}
