package auth

import (
	"context"
	"errors"

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

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*usersmodel.User, error) {
	var user usersmodel.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
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
