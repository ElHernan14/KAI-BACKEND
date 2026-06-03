package kai

import (
	"context"

	kaimodel "kai-back/internal/modules/kai/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KaiRepository interface {
	FindAllAttributeTypes(
		ctx context.Context,
	) ([]kaimodel.KaiAttributeType, error)
	CreateAttributes(
		ctx context.Context,
		tx *gorm.DB,
		attributes []kaimodel.KaiAttribute,
	) error
	CreateKaiState(
		ctx context.Context,
		tx *gorm.DB,
		kaiState *kaimodel.KaiState,
	) error
	// AddAttributeValue(
	// 	ctx context.Context,
	// 	userID uuid.UUID,
	// 	attributeTypeID uuid.UUID,
	// 	value int,
	// ) error

	// FindDominantAttribute(
	// 	ctx context.Context,
	// 	userID uuid.UUID,
	// ) (*kaimodel.KaiAttribute, error)

	FindKaiStateByUserID(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) (*kaimodel.KaiState, error)

	UpdateKaiState(
		ctx context.Context,
		tx *gorm.DB,
		kaiState *kaimodel.KaiState,
	) error

	FindKaiAttribute(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		attributeID uuid.UUID,
	) (*kaimodel.KaiAttribute, error)

	UpdateKaiAttribute(
		ctx context.Context,
		tx *gorm.DB,
		attribute *kaimodel.KaiAttribute,
	) error

	FindUserKaiAttributes(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) ([]kaimodel.KaiAttribute, error)

	UpdateLastMessage(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		message string,
	) error
}
