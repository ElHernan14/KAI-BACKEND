package authdto

type RegisterRequest struct {
	Name     string `json:"nombre" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,max=150"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Password string `json:"password" validate:"required,strong_password,max=100"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"omitempty,max=150"`
	Email      string `json:"email" validate:"omitempty,max=150"`
	Password   string `json:"password" validate:"required"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"idToken" validate:"required"`
}

type AuthUserResponse struct {
	Name     string  `json:"nombre"`
	Email    string  `json:"email"`
	Username *string `json:"username,omitempty"`
}

type AuthResponse struct {
	UserResponse AuthUserResponse `json:"usuario"`
	Token        string           `json:"token"`
}

type ValidateTokenResponse struct {
	Valid        bool             `json:"valid"`
	UserResponse AuthUserResponse `json:"usuario"`
	Token        string           `json:"token"`
}
