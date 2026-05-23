package kai

import (
	"context"
	kaimodel "kai-back/internal/modules/kai/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAllAttributeTypes(
	ctx context.Context,
) ([]kaimodel.KaiAttributeType, error) {

	var attributeTypes []kaimodel.KaiAttributeType

	err := r.db.WithContext(ctx).
		Find(&attributeTypes).Error

	if err != nil {
		return nil, err
	}

	return attributeTypes, nil
}

func (r *Repository) CreateAttributes(
	ctx context.Context,
	tx *gorm.DB,
	attributes []kaimodel.KaiAttribute,
) error {
	db := r.db

	if tx != nil {
		db = tx
	}

	if len(attributes) == 0 {
		return nil
	}

	return db.WithContext(ctx).
		Create(&attributes).Error
}

func (r *Repository) CreateKaiState(
	ctx context.Context,
	tx *gorm.DB,
	kaiState *kaimodel.KaiState,
) error {
	return tx.WithContext(ctx).Create(kaiState).Error
}
