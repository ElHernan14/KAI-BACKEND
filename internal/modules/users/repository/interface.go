package users

import (
	"context"
	usersmodel "kai-back/internal/modules/users/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	FindMeByID(ctx context.Context, userID uuid.UUID) (*usersmodel.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (*usersmodel.User, error)
	FindUserByEmail(ctx context.Context, email string) (*usersmodel.User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user *usersmodel.User) error
	CreateConfiguration(ctx context.Context, tx *gorm.DB, config *usersmodel.UserConfiguration) error
	Update(ctx context.Context, user *usersmodel.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
}
