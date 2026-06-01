package messagesRepository

import (
	"context"
	messagesmodel "kai-back/internal/modules/messages/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) FindRandomMessageByAttribute(
	ctx context.Context,
	attributeID uuid.UUID,
) (*messagesmodel.KaiMessage, error) {

	var message messagesmodel.KaiMessage

	err := r.db.
		WithContext(ctx).
		Model(&messagesmodel.KaiMessage{}).
		Joins(`
			JOIN mensaje_atributos ma
				ON ma.mensaje_kai_id = mensajes_kai.id
		`).
		Where("ma.atributo_kai_id = ?", attributeID).
		Where("mensajes_kai.activo = true").
		Order("RANDOM()").
		First(&message).
		Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) CreateUserMessage(
	ctx context.Context,
	tx *gorm.DB,
	userMessage *messagesmodel.UserMessage,
) error {

	return tx.
		WithContext(ctx).
		Create(userMessage).
		Error
}
