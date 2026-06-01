package messagesRepository

import (
	"context"
	messagesmodel "kai-back/internal/modules/messages/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepositoryPort interface {
	// FindCompletionMessage(
	// 	ctx context.Context,
	// 	attributeTypeID uuid.UUID,
	// ) (*messagesmodel.KaiMessage, error)

	FindRandomMessageByAttribute(
		ctx context.Context,
		attributeID uuid.UUID,
	) (*messagesmodel.KaiMessage, error)

	CreateUserMessage(
		ctx context.Context,
		tx *gorm.DB,
		userMessage *messagesmodel.UserMessage,
	) error
}
