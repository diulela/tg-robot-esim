package repository

import (
	"context"

	"gorm.io/gorm"

	"tg-robot-sim/storage/models"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	GetByTelegramIDs(ctx context.Context, telegramIDs []int64) ([]*models.User, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) GetByTelegramIDs(ctx context.Context, telegramIDs []int64) ([]*models.User, error) {
	if len(telegramIDs) == 0 {
		return []*models.User{}, nil
	}

	var users []*models.User
	err := r.db.WithContext(ctx).
		Where("telegram_id IN ?", telegramIDs).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
