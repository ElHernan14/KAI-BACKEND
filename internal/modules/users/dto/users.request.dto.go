package usersdto

type UpdateMeRequest struct {
	Name        string  `json:"nombre" validate:"required,min=2,max=100"`
	ProfileBase *string `json:"perfil_base,omitempty" validate:"omitempty,max=30"`
}
