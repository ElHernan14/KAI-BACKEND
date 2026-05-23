package auth

import (
	"context"

	kaimodel "kai-back/internal/modules/kai/models"
	usersmodel "kai-back/internal/modules/users/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RegisterUserWithInitialState(
	ctx context.Context,
	user *usersmodel.User,
	config *usersmodel.UserConfiguration,
	kaiState *kaimodel.KaiState,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		config.UserID = user.ID
		kaiState.UserID = user.ID

		if err := tx.Create(config).Error; err != nil {
			return err
		}

		if err := tx.Create(kaiState).Error; err != nil {
			return err
		}

		return nil
	})
}
