package db

import (
	"github.com/RodrigoGonzalez78/models"
)

func IsUserNameUnique(userName string) (bool, error) {

	var count int64

	err := database.Model(&models.User{}).Where("user_name = ?", userName).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}
