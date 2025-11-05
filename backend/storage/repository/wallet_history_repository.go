package repository

import (
	"context"
	"time"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// WalletHistoryRepository 钱包历史仓储接口
type WalletHistoryRepository interface {
	Create(ctx context.Context, history *models.WalletHistory) error
	GetByID(ctx context.Context, id uint) (*models.WalletHistory, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.WalletHistory, error)
	GetByUserIDAndType(ctx context.Context, userID int64, historyType models.WalletHistoryType, limit, offset int) ([]*models.WalletHistory, error)
	GetByUserIDAndStatus(ctx context.Context, userID int64, status models.WalletHistoryStatus, limit, offset int) ([]*models.WalletHistory, error)
	GetByUserIDAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time, limit, offset int) ([]*models.WalletHistory, error)
	GetByRelated(ctx context.Context, relatedType, relatedID string) (*models.WalletHistory, error)
	Update(ctx context.Context, history *models.WalletHistory) error
	UpdateStatus(ctx context.Context, id uint, status models.WalletHistoryStatus) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	GetStatsByUserID(ctx context.Context, userID int64) (map[string]interface{}, error)
}

// walletHistoryRepository 钱包历史仓储实现
type walletHistoryRepository struct {
	db *gorm.DB
}

// NewWalletHistoryRepository 创建钱包历史仓储实例
func NewWalletHistoryRepository(db *gorm.DB) WalletHistoryRepository {
	return &walletHistoryRepository{db: db}
}

// Create 创建钱包历史记录
func (r *walletHistoryRepository) Create(ctx context.Context, history *models.WalletHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetByID 根据ID获取钱包历史记录
func (r *walletHistoryRepository) GetByID(ctx context.Context, id uint) (*models.WalletHistory, error) {
	var history models.WalletHistory
	err := r.db.WithContext(ctx).First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// GetByUserID 根据用户ID获取钱包历史记录
func (r *walletHistoryRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.WalletHistory, error) {
	var histories []*models.WalletHistory
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetByUserIDAndType 根据用户ID和类型获取钱包历史记录
func (r *walletHistoryRepository) GetByUserIDAndType(ctx context.Context, userID int64, historyType models.WalletHistoryType, limit, offset int) ([]*models.WalletHistory, error) {
	var histories []*models.WalletHistory
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, historyType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetByUserIDAndStatus 根据用户ID和状态获取钱包历史记录
func (r *walletHistoryRepository) GetByUserIDAndStatus(ctx context.Context, userID int64, status models.WalletHistoryStatus, limit, offset int) ([]*models.WalletHistory, error) {
	var histories []*models.WalletHistory
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetByUserIDAndDateRange 根据用户ID和日期范围获取钱包历史记录
func (r *walletHistoryRepository) GetByUserIDAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time, limit, offset int) ([]*models.WalletHistory, error) {
	var histories []*models.WalletHistory
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&histories).Error
	return histories, err
}

// GetByRelated 根据关联信息获取钱包历史记录
func (r *walletHistoryRepository) GetByRelated(ctx context.Context, relatedType, relatedID string) (*models.WalletHistory, error) {
	var history models.WalletHistory
	err := r.db.WithContext(ctx).
		Where("related_type = ? AND related_id = ?", relatedType, relatedID).
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// Update 更新钱包历史记录
func (r *walletHistoryRepository) Update(ctx context.Context, history *models.WalletHistory) error {
	return r.db.WithContext(ctx).Save(history).Error
}

// UpdateStatus 更新钱包历史记录状态
func (r *walletHistoryRepository) UpdateStatus(ctx context.Context, id uint, status models.WalletHistoryStatus) error {
	return r.db.WithContext(ctx).Model(&models.WalletHistory{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除钱包历史记录
func (r *walletHistoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.WalletHistory{}, id).Error
}

// CountByUserID 统计用户钱包历史记录数量
func (r *walletHistoryRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.WalletHistory{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetStatsByUserID 获取用户钱包历史统计
func (r *walletHistoryRepository) GetStatsByUserID(ctx context.Context, userID int64) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"total_income":     "0.00",
		"total_expense":    "0.00",
		"pending_amount":   "0.00",
		"completed_amount": "0.00",
	}

	// 统计总记录数
	count, err := r.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats["total_records"] = count

	// 统计收入总额（充值、退款、解冻）
	var totalIncome string
	err = r.db.WithContext(ctx).Model(&models.WalletHistory{}).
		Select("COALESCE(SUM(CAST(amount AS DECIMAL(10,2))), 0)").
		Where("user_id = ? AND type IN (?, ?, ?) AND status = ?",
			userID,
			models.WalletHistoryTypeRecharge,
			models.WalletHistoryTypeRefund,
			models.WalletHistoryTypeUnfreeze,
			models.WalletHistoryStatusCompleted).
		Scan(&totalIncome).Error
	if err == nil {
		stats["total_income"] = totalIncome
	}

	// 统计支出总额（支付、冻结）
	var totalExpense string
	err = r.db.WithContext(ctx).Model(&models.WalletHistory{}).
		Select("COALESCE(ABS(SUM(CAST(amount AS DECIMAL(10,2)))), 0)").
		Where("user_id = ? AND type IN (?, ?) AND status = ?",
			userID,
			models.WalletHistoryTypePayment,
			models.WalletHistoryTypeFreeze,
			models.WalletHistoryStatusCompleted).
		Scan(&totalExpense).Error
	if err == nil {
		stats["total_expense"] = totalExpense
	}

	// 统计处理中金额
	var pendingAmount string
	err = r.db.WithContext(ctx).Model(&models.WalletHistory{}).
		Select("COALESCE(SUM(CAST(amount AS DECIMAL(10,2))), 0)").
		Where("user_id = ? AND status = ?", userID, models.WalletHistoryStatusPending).
		Scan(&pendingAmount).Error
	if err == nil {
		stats["pending_amount"] = pendingAmount
	}

	// 统计已完成金额
	var completedAmount string
	err = r.db.WithContext(ctx).Model(&models.WalletHistory{}).
		Select("COALESCE(SUM(CAST(amount AS DECIMAL(10,2))), 0)").
		Where("user_id = ? AND status = ?", userID, models.WalletHistoryStatusCompleted).
		Scan(&completedAmount).Error
	if err == nil {
		stats["completed_amount"] = completedAmount
	}

	return stats, nil
}
