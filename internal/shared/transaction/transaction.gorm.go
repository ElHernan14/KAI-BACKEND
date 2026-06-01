package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	WithTransaction(
		ctx context.Context,
		fn func(tx *gorm.DB) error,
	) error
}

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{db: db}
}

func (m *GormTransactionManager) WithTransaction(
	ctx context.Context,
	fn func(tx *gorm.DB) error,
) error {

	return m.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			return fn(tx)
		})
}
