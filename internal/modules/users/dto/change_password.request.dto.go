package usersdto

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,strong_password,max=100"`
}
