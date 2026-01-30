package auth

import (
	"errors"

	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/domain/repository"
	"github.com/RodrigoGonzalez78/internal/domain/service"
)

type RegisterUseCase struct {
	userRepo        repository.UserRepository
	passwordService *service.PasswordService
}

func NewRegisterUseCase(userRepo repository.UserRepository, passwordService *service.PasswordService) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
	}
}

type RegisterInput struct {
	UserName string
	Password string
}

func (uc *RegisterUseCase) Execute(input RegisterInput) error {

	if len(input.UserName) == 0 {
		return errors.New("el nombre de usuario es requerido")
	}

	if len(input.Password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	exists, err := uc.userRepo.ExistsByUserName(input.UserName)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("el nombre de usuario ya está registrado")
	}

	hashedPassword, err := uc.passwordService.Hash(input.Password)
	if err != nil {
		return errors.New("error al encriptar la contraseña")
	}

	user := entity.NewUser(input.UserName, hashedPassword)
	return uc.userRepo.Create(user)
}
