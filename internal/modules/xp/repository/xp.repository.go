package xp

import (
	"context"

	xpmodel "kai-back/internal/modules/xp/models"

	"github.com/google/uuid"
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

func (r *Repository) FindUserXPByCategory(
	ctx context.Context,
	userID uuid.UUID,
	categoryID uuid.UUID,
) (*xpmodel.UserXP, error) {

	var xp xpmodel.UserXP

	err := r.db.
		WithContext(ctx).
		Where("usuario_id = ?", userID).
		Where("categoria_xp_id = ?", categoryID).
		First(&xp).
		Error

	if err != nil {
		return nil, err
	}

	return &xp, nil
}

func (r *Repository) UpdateUserXP(
	ctx context.Context,
	tx *gorm.DB,
	userXP *xpmodel.UserXP,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Save(userXP).
		Error
}

func (r *Repository) FindXPAttributesByCategory(
	ctx context.Context,
	categoryID uuid.UUID,
) ([]xpmodel.XPAttribute, error) {

	var mappings []xpmodel.XPAttribute

	err := r.db.
		WithContext(ctx).
		Where(
			"categoria_xp_id = ?",
			categoryID,
		).
		Find(&mappings).
		Error

	if err != nil {
		return nil, err
	}

	return mappings, nil
}
