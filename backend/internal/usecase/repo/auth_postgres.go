package repo

import (
	"fmt"
	"single-window/internal/entity"

	"gorm.io/gorm"
)

const (
	userTableName = "subjects.user"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *authRepo {
	return &authRepo{db}
}

func (ar *authRepo) GetUserByPhone(phone string) (*entity.Claims, error) {
	var claims entity.Claims
	err := ar.db.
		Table(userTableName).
		First(&claims, "phone = ?", phone).
		Error
	if err != nil {
		return nil, fmt.Errorf("cannot find user, error: %w", err)
	}

	return &claims, nil
}
