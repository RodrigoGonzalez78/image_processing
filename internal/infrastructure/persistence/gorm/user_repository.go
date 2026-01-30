package gorm

import (
	"github.com/RodrigoGonzalez78/internal/domain/entity"
	"github.com/RodrigoGonzalez78/internal/infrastructure/persistence/gorm/models"
	"gorm.io/gorm"
)

type UserRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) Create(user *entity.User) error {
	model := r.toModel(user)
	return r.db.Create(model).Error
}

func (r *UserRepositoryGorm) FindByUserName(userName string) (*entity.User, error) {
	var model models.UserModel
	err := r.db.Where("user_name = ?", userName).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.toEntity(&model), nil
}

func (r *UserRepositoryGorm) ExistsByUserName(userName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserModel{}).Where("user_name = ?", userName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryGorm) toModel(user *entity.User) *models.UserModel {
	return &models.UserModel{
		UserName: user.UserName,
		Password: user.Password,
	}
}

func (r *UserRepositoryGorm) toEntity(model *models.UserModel) *entity.User {
	return &entity.User{
		UserName: model.UserName,
		Password: model.Password,
	}
}
