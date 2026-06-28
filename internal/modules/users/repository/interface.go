package users

import (
	"context"
	usersmodel "kai-back/internal/modules/users/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	FindMeByID(ctx context.Context, userID uuid.UUID) (*usersmodel.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (*usersmodel.User, error)
	FindUserByEmail(ctx context.Context, email string) (*usersmodel.User, error)
	FindUserByUsername(ctx context.Context, username string) (*usersmodel.User, error)
	FindUserByEmailOrUsername(ctx context.Context, identifier string) (*usersmodel.User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user *usersmodel.User) error
	CreateConfiguration(ctx context.Context, tx *gorm.DB, config *usersmodel.UserConfiguration) error
	Update(ctx context.Context, user *usersmodel.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateActivityStats(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		globalStreak int,
		inactiveDays int,
	) error
	FindUserProfileByID(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) (*usersmodel.User, error)

	UpdateUserProfile(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		updates map[string]interface{},
	) error

	UpdateUserConfiguration(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		updates map[string]interface{},
	) error

	FindUserSyncState(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) (*usersmodel.User, error)

	UpdateHabitRecordsSyncedAt(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		syncedAt time.Time,
	) error

	UpdateActivitySyncAt(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		syncedAt time.Time,
	) error
}
