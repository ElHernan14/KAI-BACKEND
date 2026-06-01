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
	// AddUserXP(
	// 	ctx context.Context,
	// 	tx *gorm.DB,
	// 	userID uuid.UUID,
	// 	categoryID uuid.UUID,
	// 	xp int,
	// ) error

	// FindAttributeRelationsByCategoryID(
	// 	ctx context.Context,
	// 	categoryID uuid.UUID,
	// ) ([]xpmodel.XPAttribute, error)

	FindUserXPByCategory(
		ctx context.Context,
		userID uuid.UUID,
		categoryID uuid.UUID,
	) (*xpmodel.UserXP, error)

	UpdateUserXP(
		ctx context.Context,
		tx *gorm.DB,
		userXP *xpmodel.UserXP,
	) error

	FindXPAttributesByCategory(
		ctx context.Context,
		categoryID uuid.UUID,
	) ([]xpmodel.XPAttribute, error)
}
