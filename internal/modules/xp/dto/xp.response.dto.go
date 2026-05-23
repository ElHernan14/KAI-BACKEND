package xpdto

type UserXPSummary struct {
	CategoryID   string `json:"categoria_xp_id"`
	CategoryName string `json:"nombre"`
	Value        int    `json:"valor"`
}
