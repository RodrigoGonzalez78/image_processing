package db

import (
	"github.com/RodrigoGonzalez78/models"
	"gorm.io/gorm"
)

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
