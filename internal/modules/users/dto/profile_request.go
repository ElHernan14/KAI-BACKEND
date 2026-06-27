package usersdto

type UpdateUserProfileRequest struct {
	Nombre     *string `json:"nombre" validate:"omitempty,min=2,max=100"`
	Username   *string `json:"username" validate:"omitempty,min=3,max=30"`
	PerfilBase *string `json:"perfil_base" validate:"omitempty,max=30"`
}
