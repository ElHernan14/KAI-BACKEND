package xp

import (
	"context"
	"database/sql"

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
	tx *gorm.DB,
	userID uuid.UUID,
	categoryID uuid.UUID,
) (*xpmodel.UserXP, error) {

	var xp xpmodel.UserXP

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
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

func (r *Repository) FindTotalUserXP(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (int, error) {

	var total sql.NullInt64

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Model(&xpmodel.UserXP{}).
		Select("COALESCE(SUM(valor), 0)").
		Where("usuario_id = ?", userID).
		Scan(&total).
		Error

	if err != nil {
		return 0, err
	}

	return int(total.Int64), nil
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
	tx *gorm.DB,
	categoryID uuid.UUID,
) ([]xpmodel.XPAttribute, error) {

	var mappings []xpmodel.XPAttribute

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
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

func (r *Repository) FindUserXPByUserID(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) ([]xpmodel.UserXP, error) {
	var userXP []xpmodel.UserXP

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Preload("Category").
		Where("usuario_id = ?", userID).
		Find(&userXP).
		Error

	if err != nil {
		return nil, err
	}

	return userXP, nil
}
