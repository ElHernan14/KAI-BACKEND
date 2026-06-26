package auth

import (
	"context"
	"time"

	authmodel "kai-back/internal/modules/auth/models"
	kaimodel "kai-back/internal/modules/kai/models"
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

func (r *Repository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*usersmodel.User, error) {

	var user usersmodel.User

	err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreatePasswordResetToken(
	ctx context.Context,
	token *authmodel.PasswordResetToken,
) error {

	return r.db.
		WithContext(ctx).
		Create(token).
		Error
}

func (r *Repository) FindPasswordResetToken(
	ctx context.Context,
	tokenHash string,
) (*authmodel.PasswordResetToken, error) {

	var token authmodel.PasswordResetToken

	err := r.db.
		WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&token).
		Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *Repository) MarkPasswordResetTokenUsed(
	ctx context.Context,
	tx *gorm.DB,
	tokenID uuid.UUID,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	now := time.Now()

	return db.
		WithContext(ctx).
		Model(&authmodel.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Update("used_at", now).
		Error
}

func (r *Repository) UpdateUserPassword(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	passwordHash string,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&usersmodel.User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).
		Error
}

func (r *Repository) FindValidPasswordResetCode(
	ctx context.Context,
	userID uuid.UUID,
	codeHash string,
) (*authmodel.PasswordResetCode, error) {

	var code authmodel.PasswordResetCode

	err := r.db.
		WithContext(ctx).
		Where("usuario_id = ?", userID).
		Where("code_hash = ?", codeHash).
		Where("used_at IS NULL").
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		First(&code).
		Error

	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (r *Repository) MarkPasswordResetCodeUsed(
	ctx context.Context,
	tx *gorm.DB,
	codeID uuid.UUID,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	now := time.Now()

	return db.
		WithContext(ctx).
		Model(&authmodel.PasswordResetCode{}).
		Where("id = ?", codeID).
		Update("used_at", now).
		Error
}

func (r *Repository) CreatePasswordResetCode(
	ctx context.Context,
	code *authmodel.PasswordResetCode,
) error {

	return r.db.
		WithContext(ctx).
		Create(code).
		Error
}
