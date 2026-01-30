package gorm

import (
	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type ImageRepositoryGorm struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepositoryGorm {
	return &ImageRepositoryGorm{db: db}
}

func (r *ImageRepositoryGorm) Create(image *entity.Image) error {
	model := r.toModel(image)
	return r.db.Create(model).Error
}

func (r *ImageRepositoryGorm) FindByID(id int64) (*entity.Image, error) {
	var model models.ImageModel
	err := r.db.Where("image_id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.toEntity(&model), nil
}

func (r *ImageRepositoryGorm) FindByUser(userName string, page, limit int) ([]entity.Image, int64, error) {
	var modelsResult []models.ImageModel
	offset := (page - 1) * limit

	err := r.db.Model(&models.ImageModel{}).
		Where("user_name = ?", userName).
		Limit(limit).
		Offset(offset).
		Find(&modelsResult).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = r.db.Model(&models.ImageModel{}).
		Where("user_name = ?", userName).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	images := make([]entity.Image, len(modelsResult))
	for i, model := range modelsResult {
		images[i] = *r.toEntity(&model)
	}

	return images, total, nil
}

func (r *ImageRepositoryGorm) toModel(image *entity.Image) *models.ImageModel {
	return &models.ImageModel{
		ImageID:   image.ID,
		Name:      image.Name,
		UserName:  image.UserName,
		Path:      image.Path,
		Size:      image.Size,
		Format:    image.Format,
		Width:     image.Width,
		Height:    image.Height,
		CreatedAt: image.CreatedAt,
	}
}

func (r *ImageRepositoryGorm) toEntity(model *models.ImageModel) *entity.Image {
	return &entity.Image{
		ID:        model.ImageID,
		Name:      model.Name,
		UserName:  model.UserName,
		Path:      model.Path,
		Size:      model.Size,
		Format:    model.Format,
		Width:     model.Width,
		Height:    model.Height,
		CreatedAt: model.CreatedAt,
	}
}
