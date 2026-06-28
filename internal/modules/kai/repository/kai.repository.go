package kai

import (
	"context"
	kaimodel "kai-back/internal/modules/kai/models"
	"time"

	"github.com/google/uuid"
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

func (r *Repository) FindKaiAttribute(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	attributeID uuid.UUID,
) (*kaimodel.KaiAttribute, error) {

	var attr kaimodel.KaiAttribute

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Where("usuario_id = ?", userID).
		Where("atributo_kai_id = ?", attributeID).
		First(&attr).
		Error

	if err != nil {
		return nil, err
	}

	return &attr, nil
}

func (r *Repository) UpdateKaiAttribute(
	ctx context.Context,
	tx *gorm.DB,
	attribute *kaimodel.KaiAttribute,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Save(attribute).
		Error
}

func (r *Repository) FindUserKaiAttributes(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) ([]kaimodel.KaiAttribute, error) {

	var attributes []kaimodel.KaiAttribute

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Preload("AttributeType").
		Where(
			"usuario_id = ?",
			userID,
		).
		Find(&attributes).
		Error

	if err != nil {
		return nil, err
	}

	return attributes, nil
}

func (r *Repository) FindKaiStateByUserID(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*kaimodel.KaiState, error) {

	var state kaimodel.KaiState

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Preload("DominantAttribute").
		Where(
			"usuario_id = ?",
			userID,
		).
		First(&state).
		Error

	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (r *Repository) UpdateKaiState(
	ctx context.Context,
	tx *gorm.DB,
	state *kaimodel.KaiState,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Save(state).
		Error
}

func (r *Repository) UpdateLastMessage(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	message string,
) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&kaimodel.KaiState{}).
		Where("usuario_id = ?", userID).
		Update(
			"ultimo_mensaje",
			message,
		).
		Error
}

func (r *Repository) UpdateLastInteraction(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	interaction time.Time,
) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).
		Model(&kaimodel.KaiState{}).
		Where("usuario_id = ?", userID).
		Update("ultima_interaccion", interaction).
		Error
}

func (r *Repository) UpdateTemporalState(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	energy int,
	state string,
	mode string,
	inactiveDays int,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&kaimodel.KaiState{}).
		Where("usuario_id = ?", userID).
		Updates(map[string]interface{}{
			"energia":            energy,
			"estado_actual":      state,
			"modo_actual":        mode,
			"dias_sin_actividad": inactiveDays,
		}).
		Error
}
