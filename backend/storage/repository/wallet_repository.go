package repository

import (
	"context"
	"fmt"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// WalletRepository 钱包仓储接口
type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByUserID(ctx context.Context, userID int64) (*models.Wallet, error)
	GetOrCreate(ctx context.Context, userID int64) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	UpdateBalance(ctx context.Context, userID int64, balance, frozenBalance string) error
	Delete(ctx context.Context, id uint) error
}

// walletRepository 钱包仓储实现
type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository 创建钱包仓储实例
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

// Create 创建钱包
func (r *walletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

// GetByUserID 根据用户ID获取钱包
func (r *walletRepository) GetByUserID(ctx context.Context, userID int64) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Update 更新钱包
func (r *walletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

// UpdateBalance 更新余额
func (r *walletRepository) UpdateBalance(ctx context.Context, userID int64, balance, frozenBalance string) error {
	return r.db.WithContext(ctx).Model(&models.Wallet{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"balance":        balance,
			"frozen_balance": frozenBalance,
		}).Error
}

// Delete 删除钱包
func (r *walletRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Wallet{}, id).Error
}

// GetOrCreate 获取或创建钱包
func (r *walletRepository) GetOrCreate(ctx context.Context, userID int64) (*models.Wallet, error) {
	wallet, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return wallet, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// 创建新钱包
	wallet = &models.Wallet{
		UserID:        userID,
		Balance:       "0",
		FrozenBalance: "0",
		TotalIncome:   "0",
		TotalExpense:  "0",
	}

	if err := r.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}
