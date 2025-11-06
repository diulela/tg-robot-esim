package repository

import (
	"context"
	"fmt"
	"math/big"

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

	// 新增方法用于原子操作
	// UpdateBalanceAtomic 原子性更新余额（带乐观锁）
	UpdateBalanceAtomic(ctx context.Context, userID int64, balanceDelta, frozenDelta string) error
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

// UpdateBalanceAtomic 原子性更新余额（带乐观锁）
func (r *walletRepository) UpdateBalanceAtomic(ctx context.Context, userID int64, balanceDelta, frozenDelta string) error {
	// 使用数据库级别的原子操作来更新余额
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先锁定记录
		var wallet models.Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("user_id = ?", userID).
			First(&wallet).Error; err != nil {
			return fmt.Errorf("failed to lock wallet: %w", err)
		}

		// 解析当前余额
		balance, err := parseDecimal(wallet.Balance)
		if err != nil {
			return fmt.Errorf("invalid balance format: %w", err)
		}

		frozenBalance, err := parseDecimal(wallet.FrozenBalance)
		if err != nil {
			return fmt.Errorf("invalid frozen balance format: %w", err)
		}

		// 解析变化量
		balanceChange, err := parseDecimal(balanceDelta)
		if err != nil {
			return fmt.Errorf("invalid balance delta format: %w", err)
		}

		frozenChange, err := parseDecimal(frozenDelta)
		if err != nil {
			return fmt.Errorf("invalid frozen delta format: %w", err)
		}

		// 计算新余额
		newBalance := new(big.Float).Add(balance, balanceChange)
		newFrozenBalance := new(big.Float).Add(frozenBalance, frozenChange)

		// 检查余额不能为负
		if newBalance.Cmp(big.NewFloat(0)) < 0 {
			return fmt.Errorf("insufficient balance")
		}

		if newFrozenBalance.Cmp(big.NewFloat(0)) < 0 {
			return fmt.Errorf("insufficient frozen balance")
		}

		// 更新余额
		wallet.Balance = newBalance.Text('f', 8)
		wallet.FrozenBalance = newFrozenBalance.Text('f', 8)

		return tx.Save(&wallet).Error
	})
}

// parseDecimal 解析 decimal 字符串为 big.Float
func parseDecimal(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	if err != nil {
		return big.NewFloat(0), err
	}
	return f, nil
}
