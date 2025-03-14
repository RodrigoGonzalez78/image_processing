package db

import "github.com/RodrigoGonzalez78/models"

func GetImageByID(imageID int64) (*models.Image, error) {
	var image models.Image
	err := database.Model(&models.Image{}).
		Where("image_id = ?", imageID).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}
