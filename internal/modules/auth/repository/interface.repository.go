package auth

import (
	"context"
	authmodel "kai-back/internal/modules/auth/models"
	kaimodel "kai-back/internal/modules/kai/models"
	usersmodel "kai-back/internal/modules/users/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	RegisterUserWithInitialState(
		ctx context.Context,
		user *usersmodel.User,
		config *usersmodel.UserConfiguration,
		kaiState *kaimodel.KaiState,
	) error

	FindUserByEmail(
		ctx context.Context,
		email string,
	) (*usersmodel.User, error)

	CreatePasswordResetToken(
		ctx context.Context,
		token *authmodel.PasswordResetToken,
	) error

	FindPasswordResetToken(
		ctx context.Context,
		tokenHash string,
	) (*authmodel.PasswordResetToken, error)

	MarkPasswordResetTokenUsed(
		ctx context.Context,
		tx *gorm.DB,
		tokenID uuid.UUID,
	) error

	UpdateUserPassword(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		passwordHash string,
	) error

	CreatePasswordResetCode(ctx context.Context, code *authmodel.PasswordResetCode) error

	FindValidPasswordResetCode(ctx context.Context, userID uuid.UUID, codeHash string) (*authmodel.PasswordResetCode, error)

	MarkPasswordResetCodeUsed(ctx context.Context, tx *gorm.DB, codeID uuid.UUID) error
}
