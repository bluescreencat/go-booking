package repository

import (
	"booking/internal/entity"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	FindUserByUsername(username string) (user *entity.User, err error)
	CreateAccount(user *entity.User) (err error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db: db}
}

func (repo *authRepository) FindUserByUsername(username string) (user *entity.User, err error) {
	tx := repo.db.First(user, "username = ?", username)
	return user, tx.Error
}

func (repo *authRepository) CreateAccount(user *entity.User) (err error) {
	tx := repo.db.Save(user)
	return tx.Error
}
