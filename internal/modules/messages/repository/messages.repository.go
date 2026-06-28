package messagesRepository

import (
	"context"
	"errors"
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
	tx *gorm.DB,
	attributeID uuid.UUID,
) (*messagesmodel.KaiMessage, error) {

	var message messagesmodel.KaiMessage

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
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

func (r *MessageRepository) FindRandomMessageByRules(
	ctx context.Context,
	tx *gorm.DB,
	attributeID uuid.UUID,
	attributeValue int,
	types []string,
	contexts []string,
) (*messagesmodel.KaiMessage, error) {

	var message messagesmodel.KaiMessage

	db := r.db

	if tx != nil {
		db = tx
	}

	query := db.
		WithContext(ctx).
		Model(&messagesmodel.KaiMessage{}).
		Joins(`
			JOIN mensaje_atributos ma
				ON ma.mensaje_kai_id = mensajes_kai.id
		`).
		Where("ma.atributo_kai_id = ?", attributeID).
		Where("ma.nivel_minimo <= ?", attributeValue).
		Where("mensajes_kai.activo = true")

	if len(types) > 0 && len(contexts) > 0 {
		query = query.Where(
			"(mensajes_kai.tipo IN ? OR mensajes_kai.contexto IN ?)",
			types,
			contexts,
		)
	} else if len(types) > 0 {
		query = query.Where("mensajes_kai.tipo IN ?", types)
	} else if len(contexts) > 0 {
		query = query.Where("mensajes_kai.contexto IN ?", contexts)
	}

	err := query.
		Order("RANDOM()").
		First(&message).
		Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) FindRandomMessageByType(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	messageType string,
) (*messagesmodel.KaiMessage, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	baseQuery := db.WithContext(ctx).
		Model(&messagesmodel.KaiMessage{}).
		Where("tipo = ? AND activo = true", messageType)

	var message messagesmodel.KaiMessage
	err := baseQuery.
		Where(`NOT EXISTS (
			SELECT 1
			FROM mensajes_usuario mu
			WHERE mu.usuario_id = ?
			AND mu.mensaje_kai_id = mensajes_kai.id
		)`, userID).
		Order("RANDOM()").
		First(&message).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = baseQuery.Order("RANDOM()").First(&message).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) HasUserMessages(
	ctx context.Context,
	userID uuid.UUID,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&messagesmodel.UserMessage{}).
		Where("usuario_id = ?", userID).
		Limit(1).
		Count(&count).
		Error

	return count > 0, err
}

func (r *MessageRepository) CreateUserMessage(
	ctx context.Context,
	tx *gorm.DB,
	userMessage *messagesmodel.UserMessage,
) error {

	db := r.db
	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Create(userMessage).
		Error
}

func (r *MessageRepository) FindLastUserMessage(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*messagesmodel.UserMessage, error) {
	var userMessage messagesmodel.UserMessage

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		Preload("KaiMessage").
		Where("usuario_id = ?", userID).
		Order("mostrado_en DESC").
		First(&userMessage).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &userMessage, nil
}
