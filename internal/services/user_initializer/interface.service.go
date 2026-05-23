package userinitializer

import (
	"context"

	usersmodel "kai-back/internal/modules/users/models"
)

type Service interface {
	InitializeNewUser(
		ctx context.Context,
		user *usersmodel.User,
	) error
}
