package authdto

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"codigo" validate:"required,len=6"`
	Password string `json:"password" validate:"required,min=6"`
}
