package service_test

import (
	"testing"

	"github.com/RodrigoGonzalez78/internal/domain/service"
)

func TestPasswordService_Hash(t *testing.T) {
	passwordService := service.NewPasswordService()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "hash contraseña válida",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "hash contraseña corta",
			password: "abc",
			wantErr:  false,
		},
		{
			name:     "hash contraseña vacía",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := passwordService.Hash(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Error("Hash() retornó un hash vacío")
			}
			if !tt.wantErr && hash == tt.password {
				t.Error("Hash() no debe retornar la contraseña sin hashear")
			}
		})
	}
}

func TestPasswordService_Verify(t *testing.T) {
	passwordService := service.NewPasswordService()

	// Primero creamos un hash válido
	password := "miContraseñaSegura123"
	hashedPassword, _ := passwordService.Hash(password)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		want           bool
	}{
		{
			name:           "contraseña correcta",
			hashedPassword: hashedPassword,
			password:       password,
			want:           true,
		},
		{
			name:           "contraseña incorrecta",
			hashedPassword: hashedPassword,
			password:       "otraContraseña",
			want:           false,
		},
		{
			name:           "hash vacío",
			hashedPassword: "",
			password:       password,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := passwordService.Verify(tt.hashedPassword, tt.password)
			if got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenService_GenerateAndValidate(t *testing.T) {
	tokenService := service.NewTokenService("mi-secreto-super-seguro", 24*60*60*1000000000) // 24h

	tests := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{
			name:     "generar token para usuario válido",
			userName: "testuser",
			wantErr:  false,
		},
		{
			name:     "generar token para usuario con espacios",
			userName: "test user",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generar token
			token, err := tokenService.Generate(tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("Generate() retornó un token vacío")
			}

			// Validar token
			claims, err := tokenService.Validate("Bearer " + token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claims.UserName != tt.userName {
				t.Errorf("Validate() UserName = %v, want %v", claims.UserName, tt.userName)
			}
		})
	}
}

func TestTokenService_ValidateInvalidToken(t *testing.T) {
	tokenService := service.NewTokenService("mi-secreto-super-seguro", 24*60*60*1000000000)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "token sin prefijo Bearer",
			token:   "algún-token-inválido",
			wantErr: true,
		},
		{
			name:    "token vacío",
			token:   "",
			wantErr: true,
		},
		{
			name:    "token con Bearer pero inválido",
			token:   "Bearer token-inválido",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tokenService.Validate(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
