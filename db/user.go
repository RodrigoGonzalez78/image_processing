package db

import (
	"fmt"

	"github.com/RodrigoGonzalez78/models"
	"gorm.io/gorm"
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

func IsUserNameUnique(userName string) (bool, error) {

	var count int64

	err := database.Model(&models.User{}).Where("user_name = ?", userName).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func GetUserByUserName(userName string) (*models.User, error) {
	var user models.User

	err := database.Model(&models.User{}).Where("user_name = ?", userName).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {

			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
