package db

import (
	"fmt"

	"github.com/RodrigoGonzalez78/models"
)

func CreateImage(name, userName string) error {

	isUnique, err := IsUserNameUnique(userName)
	if err != nil {
		return fmt.Errorf("error al verificar usuario: %v", err)
	}

	if isUnique {
		return fmt.Errorf("no se puede crear la imagen: el usuario '%s' no existe", userName)
	}

	image := models.Image{Name: name, UserName: userName}
	if err := database.Create(&image).Error; err != nil {
		return fmt.Errorf("error al crear imagen: %v", err)
	}

	return nil
}
