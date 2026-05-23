package kai

import (
	"context"

	kaimodel "kai-back/internal/modules/kai/models"

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
}
