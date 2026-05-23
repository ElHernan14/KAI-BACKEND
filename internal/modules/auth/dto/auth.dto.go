package authdto

type RegisterRequest struct {
	Name     string `json:"nombre" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,strong_password,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required"`
}

type AuthUserResponse struct {
	Name  string `json:"nombre"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
