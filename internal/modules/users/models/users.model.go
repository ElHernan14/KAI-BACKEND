package usersmodel

import (
	habitsmodel "kai-back/internal/modules/habits/models"
	kaimodel "kai-back/internal/modules/kai/models"
	xpmodel "kai-back/internal/modules/xp/models"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"column:nombre;type:varchar(100);not null" json:"nombre"`
	Email        string    `gorm:"column:email;type:varchar(150);not null;unique" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null" json:"-"`
	RegisteredAt time.Time `gorm:"column:fecha_registro;autoCreateTime" json:"fecha_registro"`
	BaseProfile  *string   `gorm:"column:perfil_base;type:varchar(30)" json:"perfil_base,omitempty"`
	KaiStage     string    `gorm:"column:etapa_kai;type:varchar(30);default:cachorro" json:"etapa_kai"`
	GlobalStreak int       `gorm:"column:racha_global;default:0" json:"racha_global"`
	InactiveDays int       `gorm:"column:dias_inactivo;default:0" json:"dias_inactivo"`

	Configuration UserConfiguration `gorm:"foreignKey:UserID"`
	KaiState      kaimodel.KaiState `gorm:"foreignKey:UserID"`

	XPCategory    []xpmodel.UserXP
	KaiAttributes []kaimodel.KaiAttribute
	UserHabits    []habitsmodel.UserHabit `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "usuarios"
}
