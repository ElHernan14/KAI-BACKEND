package xp

import (
	"context"

	xpmodel "kai-back/internal/modules/xp/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindAllCategories(
	ctx context.Context,
) ([]xpmodel.XPCategory, error) {

	var categories []xpmodel.XPCategory

	err := r.db.WithContext(ctx).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) CreateUserXP(
	ctx context.Context,
	tx *gorm.DB,
	userXP []xpmodel.UserXP,
) error {
	db := r.db

	if tx != nil {
		db = tx
	}

	if len(userXP) == 0 {
		return nil
	}

	return db.WithContext(ctx).
		Create(&userXP).Error
}
