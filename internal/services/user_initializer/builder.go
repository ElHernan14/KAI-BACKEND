package userinitializer

import (
	kaimodel "kai-back/internal/modules/kai/models"
	usersmodel "kai-back/internal/modules/users/models"

	"github.com/google/uuid"
)

func buildInitialConfiguration(
	userID uuid.UUID,
) *usersmodel.UserConfiguration {

	return &usersmodel.UserConfiguration{
		UserID: userID,

		NotificationsEnabled:   true,
		SoundsEnabled:          true,
		ShowStreaks:            true,
		DiscreteMode:           false,
		KaiIntensity:           "normal",
		LockWithPIN:            false,
		AllowEmotionalMessages: true,
	}
}

func buildInitialKaiState(
	userID uuid.UUID,
) *kaimodel.KaiState {
	currentMode := "SUAVE"

	return &kaimodel.KaiState{
		UserID: userID,

		CurrentState: "CURIOSO",
		CurrentStage: "cachorro",
		CurrentMode:  &currentMode,

		Energy: 100,

		BondLevel: 0,

		DaysWithoutActivity: 0,

		RecoveryMode: false,
	}
}
