package auth

import (
	"errors"

	"github.com/RodrigoGonzalez78/internal/domain/repository"
	"github.com/RodrigoGonzalez78/internal/domain/service"
)

type LoginUseCase struct {
	userRepo        repository.UserRepository
	passwordService *service.PasswordService
	tokenService    *service.TokenService
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	passwordService *service.PasswordService,
	tokenService *service.TokenService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

type LoginInput struct {
	UserName string
	Password string
}

type LoginOutput struct {
	Token string
}

func (uc *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {

	if len(input.UserName) == 0 {
		return nil, errors.New("el nombre de usuario es requerido")
	}

	if len(input.Password) == 0 {
		return nil, errors.New("la contraseña es requerida")
	}

	user, err := uc.userRepo.FindByUserName(input.UserName)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	if !uc.passwordService.Verify(user.Password, input.Password) {
		return nil, errors.New("contraseña inválida")
	}

	token, err := uc.tokenService.Generate(input.UserName)
	if err != nil {
		return nil, errors.New("error al generar token")
	}

	return &LoginOutput{Token: token}, nil
}
