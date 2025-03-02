package db

import "github.com/RodrigoGonzalez78/models"

func GetImagesByUserName(userName string) ([]models.Image, error) {

	var images []models.Image

	err := database.Model(&models.Image{}).Where("user_name = ?", userName).Find(&images).Error

	if err != nil {
		return nil, err
	}

	return images, nil
}
