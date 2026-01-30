package repository

import "github.com/RodrigoGonzalez78/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error

	FindByUserName(userName string) (*entity.User, error)

	ExistsByUserName(userName string) (bool, error)
}
