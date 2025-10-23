package repository

import (
	"context"
	"tg-robot-sim/models"
	"time"

	"gorm.io/gorm"
)

// UserSessionRepository 用户会话仓库接口
type UserSessionRepository interface {
	Create(ctx context.Context, session *models.UserSession) error
	GetByUserID(ctx context.Context, userID int64) (*models.UserSession, error)
	Update(ctx context.Context, session *models.UserSession) error
	Delete(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context, timeout time.Duration) error
	List(ctx context.Context, limit, offset int) ([]*models.UserSession, error)
}

// userSessionRepository 用户会话仓库实现
type userSessionRepository struct {
	db *gorm.DB
}

// NewUserSessionRepository 创建用户会话仓库
func NewUserSessionRepository(db *gorm.DB) UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *userSessionRepository) GetByUserID(ctx context.Context, userID int64) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *userSessionRepository) Update(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *userSessionRepository) Delete(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

func (r *userSessionRepository) DeleteExpired(ctx context.Context, timeout time.Duration) error {
	expiredTime := time.Now().Add(-timeout)
	return r.db.WithContext(ctx).Where("last_active < ?", expiredTime).Delete(&models.UserSession{}).Error
}

func (r *userSessionRepository) List(ctx context.Context, limit, offset int) ([]*models.UserSession, error) {
	var sessions []*models.UserSession
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&sessions).Error
	return sessions, err
}
