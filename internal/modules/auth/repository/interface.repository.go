package auth

import (
	"context"
	kaimodel "kai-back/internal/modules/kai/models"
	usersmodel "kai-back/internal/modules/users/models"
)

type AuthRepository interface {
	RegisterUserWithInitialState(
		ctx context.Context,
		user *usersmodel.User,
		config *usersmodel.UserConfiguration,
		kaiState *kaimodel.KaiState,
	) error
}
