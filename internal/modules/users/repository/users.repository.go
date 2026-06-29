package users

import (
	"context"
	"errors"
	usersmodel "kai-back/internal/modules/users/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) FindUserSyncState(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*usersmodel.User, error) {
	var user usersmodel.User
	db := r.db
	if tx != nil {
		db = tx
	}

	query := db.WithContext(ctx).
		Select("id", "dias_inactivo", "habit_records_synced_at", "activity_sync_at").
		Where("id = ?", userID)
	if tx != nil {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateHabitRecordsSyncedAt(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	syncedAt time.Time,
) error {
	return r.updateSyncTimestamp(ctx, tx, userID, "habit_records_synced_at", syncedAt)
}

func (r *Repository) UpdateActivitySyncAt(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	syncedAt time.Time,
) error {
	return r.updateSyncTimestamp(ctx, tx, userID, "activity_sync_at", syncedAt)
}

func (r *Repository) updateSyncTimestamp(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	column string,
	syncedAt time.Time,
) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).
		Model(&usersmodel.User{}).
		Where("id = ?", userID).
		Update(column, syncedAt).
		Error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (*usersmodel.User, error) {

	var user usersmodel.User

	err := r.db.
		WithContext(ctx).
		First(&user, "id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindMeByID(
	ctx context.Context,
	userID uuid.UUID,
) (*usersmodel.User, error) {

	var user usersmodel.User

	err := r.db.
		WithContext(ctx).
		Preload("Configuration").
		Preload("KaiState").
		Preload("KaiState.DominantAttribute").
		Preload("XPCategory.Category").
		Preload("KaiAttributes.AttributeType").
		Preload("UserHabits").
		Preload("UserHabits.HabitCatalog").
		Preload(
			"UserHabits.HabitRecords",
			"fecha = CURRENT_DATE",
		).
		First(&user, "id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*usersmodel.User, error) {
	var user usersmodel.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (*usersmodel.User, error) {
	var user usersmodel.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByEmailOrUsername(ctx context.Context, identifier string) (*usersmodel.User, error) {
	var user usersmodel.User

	err := r.db.
		WithContext(ctx).
		Where("email = ? OR username = ?", identifier, identifier).
		First(&user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateUser(
	ctx context.Context,
	tx *gorm.DB,
	user *usersmodel.User,
) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *Repository) CreateConfiguration(
	ctx context.Context,
	tx *gorm.DB,
	config *usersmodel.UserConfiguration,
) error {
	return tx.WithContext(ctx).Create(config).Error
}

func (r *Repository) Update(
	ctx context.Context,
	user *usersmodel.User,
) error {
	return r.db.WithContext(ctx).
		Save(user).
		Error
}

func (r *Repository) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	passwordHash string,
) error {

	return r.db.
		WithContext(ctx).
		Model(&usersmodel.User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).
		Error
}

func (r *Repository) UpdateActivityStats(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	globalStreak int,
	inactiveDays int,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&usersmodel.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"racha_global":  globalStreak,
			"dias_inactivo": inactiveDays,
		}).
		Error
}

func (r *Repository) FindUserProfileByID(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*usersmodel.User, error) {

	var user usersmodel.User

	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Preload("Configuration").
		Where("id = ?", userID).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUserProfile(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	updates map[string]interface{},
) error {

	if len(updates) == 0 {
		return nil
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&usersmodel.User{}).
		Where("id = ?", userID).
		Updates(updates).
		Error
}

func (r *Repository) UpdateUserConfiguration(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	updates map[string]interface{},
) error {

	if len(updates) == 0 {
		return nil
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&usersmodel.UserConfiguration{}).
		Where("usuario_id = ?", userID).
		Updates(updates).
		Error
}
