package xp

import (
	"context"

	xpmodel "kai-back/internal/modules/xp/models"

	"gorm.io/gorm"
)

type XpRepository interface {
	FindAllCategories(ctx context.Context) ([]xpmodel.XPCategory, error)
	CreateUserXP(
		ctx context.Context,
		tx *gorm.DB,
		userXP []xpmodel.UserXP,
	) error
}
