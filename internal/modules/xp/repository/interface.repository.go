package xp

import (
	"context"

	xpmodel "kai-back/internal/modules/xp/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type XpRepository interface {
	FindAllCategories(ctx context.Context) ([]xpmodel.XPCategory, error)
	CreateUserXP(
		ctx context.Context,
		tx *gorm.DB,
		userXP []xpmodel.UserXP,
	) error

	FindUserXPByCategory(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		categoryID uuid.UUID,
	) (*xpmodel.UserXP, error)

	FindTotalUserXP(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) (int, error)

	UpdateUserXP(
		ctx context.Context,
		tx *gorm.DB,
		userXP *xpmodel.UserXP,
	) error

	FindXPAttributesByCategory(
		ctx context.Context,
		tx *gorm.DB,
		categoryID uuid.UUID,
	) ([]xpmodel.XPAttribute, error)

	FindUserXPByUserID(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) ([]xpmodel.UserXP, error)
}
