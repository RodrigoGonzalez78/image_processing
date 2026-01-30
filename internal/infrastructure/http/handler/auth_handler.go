package handler

import (
	"encoding/json"
	"net/http"

	authUC "github.com/RodrigoGonzalez78/internal/application/usecase/auth"
	"github.com/RodrigoGonzalez78/internal/infrastructure/http/dto"
)

// AuthHandler maneja las peticiones de autenticación
type AuthHandler struct {
	registerUC *authUC.RegisterUseCase
	loginUC    *authUC.LoginUseCase
}

// NewAuthHandler crea una nueva instancia de AuthHandler
func NewAuthHandler(registerUC *authUC.RegisterUseCase, loginUC *authUC.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

// Register godoc
// @Summary      Registra un nuevo usuario
// @Description  Crea un nuevo usuario en la base de datos con nombre de usuario y contraseña.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body dto.RegisterRequest true "Datos del usuario a registrar"
// @Success      201
// @Failure      400 {object} dto.ErrorResponse
// @Router       /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error en los datos recibidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	input := authUC.RegisterInput{
		UserName: req.UserName,
		Password: req.Password,
	}

	if err := h.registerUC.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login godoc
// @Summary      Inicia sesión
// @Description  Autentica al usuario y devuelve un token JWT si las credenciales son correctas.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body dto.LoginRequest true "Credenciales de usuario"
// @Success      201 {object} dto.LoginResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Usuario y/o contraseña inválidas: "+err.Error(), http.StatusBadRequest)
		return
	}

	input := authUC.LoginInput{
		UserName: req.UserName,
		Password: req.Password,
	}

	output, err := h.loginUC.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.LoginResponse{Token: output.Token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
