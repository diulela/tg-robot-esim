package services

import (
	"context"
	"fmt"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// TransactionFilters 交易筛选条件
type TransactionFilters struct {
	UserID    int64
	Status    models.TransactionStatus
	StartDate string
	EndDate   string
	Limit     int
	Offset    int
}

// TransactionStats 交易统计信息
type TransactionStats struct {
	TotalTransactions int64  `json:"total_transactions"`
	TotalAmount       string `json:"total_amount"`
	PendingCount      int64  `json:"pending_count"`
	ConfirmedCount    int64  `json:"confirmed_count"`
	FailedCount       int64  `json:"failed_count"`
}

// TransactionService 交易服务接口
type TransactionService interface {
	GetTransactionHistory(ctx context.Context, userID int64, limit, offset int) ([]*models.Transaction, error)
	GetTransactionByHash(ctx context.Context, txHash string) (*models.Transaction, error)
	GetTransactionStats(ctx context.Context, userID int64) (*TransactionStats, error)
	FilterTransactions(ctx context.Context, filters TransactionFilters) ([]*models.Transaction, error)
}

// transactionService 交易服务实现
type transactionService struct {
	transactionRepo repository.TransactionRepository
}

// NewTransactionService 创建交易服务实例
func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

// GetTransactionHistory 获取交易历史
func (s *transactionService) GetTransactionHistory(ctx context.Context, userID int64, limit, offset int) ([]*models.Transaction, error) {
	// TODO: 需要从用户获取钱包地址
	// 暂时返回空列表，实际应该根据用户的钱包地址查询
	return []*models.Transaction{}, nil
}

// GetTransactionByHash 根据交易哈希获取交易
func (s *transactionService) GetTransactionByHash(ctx context.Context, txHash string) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.GetByTxHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	return transaction, nil
}

// GetTransactionStats 获取交易统计信息
func (s *transactionService) GetTransactionStats(ctx context.Context, userID int64) (*TransactionStats, error) {
	// 获取用户的所有交易
	transactions, err := s.GetTransactionHistory(ctx, userID, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	stats := &TransactionStats{
		TotalTransactions: int64(len(transactions)),
		TotalAmount:       "0.00",
	}

	var totalAmount float64
	for _, tx := range transactions {
		switch tx.Status {
		case models.TransactionStatusPending:
			stats.PendingCount++
		case models.TransactionStatusConfirmed:
			stats.ConfirmedCount++
		case models.TransactionStatusFailed:
			stats.FailedCount++
		}

		// 累计总金额（已确认的交易）
		if tx.Status == models.TransactionStatusConfirmed {
			amount, _ := parseDecimal(tx.Amount)
			amountFloat, _ := amount.Float64()
			totalAmount += amountFloat
		}
	}

	stats.TotalAmount = fmt.Sprintf("%.2f", totalAmount)

	return stats, nil
}

// FilterTransactions 根据条件筛选交易
func (s *transactionService) FilterTransactions(ctx context.Context, filters TransactionFilters) ([]*models.Transaction, error) {
	// 获取所有交易
	transactions, err := s.GetTransactionHistory(ctx, filters.UserID, filters.Limit, filters.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// 应用筛选条件
	var filtered []*models.Transaction
	for _, tx := range transactions {
		// 按状态筛选
		if filters.Status != "" && tx.Status != filters.Status {
			continue
		}

		// 按日期范围筛选
		// TODO: 实现日期范围筛选

		filtered = append(filtered, tx)
	}

	return filtered, nil
}
