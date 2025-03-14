package db

import "github.com/RodrigoGonzalez78/models"

func GetImagesByUser(userName string, page, limit int) ([]models.Image, int64, error) {
	var images []models.Image
	offset := (page - 1) * limit

	err := database.Model(&models.Image{}).
		Where("user_name = ?", userName).
		Limit(limit).
		Offset(offset).
		Find(&images).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = database.Model(&models.Image{}).
		Where("user_name = ?", userName).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}
