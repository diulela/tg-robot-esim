package repository

import (
	"context"
	"tg-robot-sim/storage/models"
	"time"

	"gorm.io/gorm"
)

// TransactionRepository 交易仓库接口
type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByTxHash(ctx context.Context, txHash string) (*models.Transaction, error)
	GetByID(ctx context.Context, id uint) (*models.Transaction, error)
	Update(ctx context.Context, tx *models.Transaction) error
	Delete(ctx context.Context, id uint) error
	GetByAddress(ctx context.Context, address string, limit int) ([]*models.Transaction, error)
	GetPendingTransactions(ctx context.Context) ([]*models.Transaction, error)
	UpdateStatus(ctx context.Context, txHash string, status models.TransactionStatus, confirmations int) error
}

// transactionRepository 交易仓库实现
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository 创建交易仓库
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) GetByTxHash(ctx context.Context, txHash string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.WithContext(ctx).Where("tx_hash = ?", txHash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.WithContext(ctx).First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *transactionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Transaction{}, id).Error
}

func (r *transactionRepository) GetByAddress(ctx context.Context, address string, limit int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := r.db.WithContext(ctx).
		Where("from_address = ? OR to_address = ?", address, address).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetPendingTransactions(ctx context.Context) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := r.db.WithContext(ctx).
		Where("status = ?", models.TransactionStatusPending).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, txHash string, status models.TransactionStatus, confirmations int) error {
	return r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("tx_hash = ?", txHash).
		Updates(map[string]interface{}{
			"status":        status,
			"confirmations": confirmations,
			"updated_at":    time.Now(),
		}).Error
}
