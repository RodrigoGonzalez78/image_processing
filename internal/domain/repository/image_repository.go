package repository

import "github.com/RodrigoGonzalez78/internal/domain/entity"

type ImageRepository interface {
	Create(image *entity.Image) error

	FindByID(id int64) (*entity.Image, error)

	FindByUser(userName string, page, limit int) ([]entity.Image, int64, error)
}
