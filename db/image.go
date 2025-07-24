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
