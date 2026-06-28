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
		tx *gorm.DB,
		attributeID uuid.UUID,
	) (*messagesmodel.KaiMessage, error)

	FindRandomMessageByRules(
		ctx context.Context,
		tx *gorm.DB,
		attributeID uuid.UUID,
		attributeValue int,
		types []string,
		contexts []string,
	) (*messagesmodel.KaiMessage, error)

	FindRandomMessageByType(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
		messageType string,
	) (*messagesmodel.KaiMessage, error)

	HasUserMessages(
		ctx context.Context,
		userID uuid.UUID,
	) (bool, error)

	CreateUserMessage(
		ctx context.Context,
		tx *gorm.DB,
		userMessage *messagesmodel.UserMessage,
	) error

	FindLastUserMessage(
		ctx context.Context,
		tx *gorm.DB,
		userID uuid.UUID,
	) (*messagesmodel.UserMessage, error)
}
