package habitsdto

type DeactivateHabitResponse struct {
	HabitoUsuarioID string `json:"habito_usuario_id"`
	Activo          bool   `json:"activo"`
	Mensaje         string `json:"mensaje"`
}
