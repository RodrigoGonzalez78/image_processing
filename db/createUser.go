package db

import (
	"fmt"

	"github.com/RodrigoGonzalez78/models"
)

func CreateUser(user models.User) error {

	isUnique, err := IsUserNameUnique(user.UserName)

	if err != nil {
		return fmt.Errorf("error al verificar usuario: %v", err)
	}

	if !isUnique {
		return fmt.Errorf("el nombre de usuario '%s' ya está en uso", user.UserName)
	}

	if err := database.Create(&user).Error; err != nil {
		return fmt.Errorf("error al crear usuario: %v", err)
	}

	return nil
}
