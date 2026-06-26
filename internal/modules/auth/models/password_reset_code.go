package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetCode struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID  `gorm:"column:usuario_id;type:uuid;not null;index"`
	CodeHash  string     `gorm:"column:code_hash;type:text;not null;index"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (PasswordResetCode) TableName() string {
	return "password_reset_codes"
}
