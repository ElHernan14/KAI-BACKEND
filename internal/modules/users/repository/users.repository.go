package users

import (
	"context"
	"errors"
	usersmodel "kai-back/internal/modules/users/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (*usersmodel.User, error) {

	var user usersmodel.User

	err := r.db.
		WithContext(ctx).
		First(&user, "id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindMeByID(
	ctx context.Context,
	userID uuid.UUID,
) (*usersmodel.User, error) {

	var user usersmodel.User

	err := r.db.
		WithContext(ctx).
		Preload("Configuration").
		Preload("KaiState").
		First(&user, "id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
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

func (r *Repository) CreateUser(
	ctx context.Context,
	tx *gorm.DB,
	user *usersmodel.User,
) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *Repository) CreateConfiguration(
	ctx context.Context,
	tx *gorm.DB,
	config *usersmodel.UserConfiguration,
) error {
	return tx.WithContext(ctx).Create(config).Error
}

func (r *Repository) Update(
	ctx context.Context,
	user *usersmodel.User,
) error {
	return r.db.WithContext(ctx).
		Save(user).
		Error
}
