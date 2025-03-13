package db

import (
	"fmt"

	"github.com/RodrigoGonzalez78/models"
)

func CreateImage(image models.Image) error {

	isUnique, err := IsUserNameUnique(image.UserName)

	if err != nil {
		return fmt.Errorf("error al verificar usuario: %v", err)
	}

	if isUnique {
		return fmt.Errorf("no se puede crear la imagen: el usuario '%s' no existe", image.UserName)
	}

	if err := database.Create(&image).Error; err != nil {
		return fmt.Errorf("error al crear imagen: %v", err)
	}

	return nil
}
