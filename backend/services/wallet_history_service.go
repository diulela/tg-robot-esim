package services

import (
	"context"
	"fmt"
	"time"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// walletHistoryService 钱包历史服务实现
type walletHistoryService struct {
	walletHistoryRepo repository.WalletHistoryRepository
}

// NewWalletHistoryService 创建钱包历史服务实例
func NewWalletHistoryService(
	walletHistoryRepo repository.WalletHistoryRepository,
) WalletHistoryService {
	return &walletHistoryService{
		walletHistoryRepo: walletHistoryRepo,
	}
}

// GetWalletHistory 获取钱包历史记录
func (s *walletHistoryService) GetWalletHistory(ctx context.Context, userID int64, filters WalletHistoryFilters) ([]*models.WalletHistory, int64, error) {
	var histories []*models.WalletHistory
	var err error

	// 根据筛选条件选择不同的查询方法
	if filters.Type != "" && filters.Status != "" {
		// 同时按类型和状态筛选（需要自定义查询）
		histories, err = s.getHistoriesByMultipleFilters(ctx, userID, filters)
	} else if filters.Type != "" {
		// 按类型筛选
		histories, err = s.walletHistoryRepo.GetByUserIDAndType(ctx, userID, filters.Type, filters.Limit, filters.Offset)
	} else if filters.Status != "" {
		// 按状态筛选
		histories, err = s.walletHistoryRepo.GetByUserIDAndStatus(ctx, userID, filters.Status, filters.Limit, filters.Offset)
	} else if filters.StartDate != "" || filters.EndDate != "" {
		// 按日期范围筛选
		startDate, endDate := s.parseDateRange(filters.StartDate, filters.EndDate)
		histories, err = s.walletHistoryRepo.GetByUserIDAndDateRange(ctx, userID, startDate, endDate, filters.Limit, filters.Offset)
	} else {
		// 获取所有记录
		histories, err = s.walletHistoryRepo.GetByUserID(ctx, userID, filters.Limit, filters.Offset)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("获取钱包历史失败: %w", err)
	}

	// 获取总数
	total, err := s.walletHistoryRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("获取记录总数失败: %w", err)
	}

	return histories, total, nil
}

// GetWalletHistoryStats 获取钱包历史统计
func (s *walletHistoryService) GetWalletHistoryStats(ctx context.Context, userID int64) (*WalletHistoryStats, error) {
	statsMap, err := s.walletHistoryRepo.GetStatsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取钱包统计失败: %w", err)
	}

	// 转换为结构体
	stats := &WalletHistoryStats{
		TotalRecords:    statsMap["total_records"].(int64),
		TotalIncome:     statsMap["total_income"].(string),
		TotalExpense:    statsMap["total_expense"].(string),
		PendingAmount:   statsMap["pending_amount"].(string),
		CompletedAmount: statsMap["completed_amount"].(string),
	}

	return stats, nil
}

// GetHistoryRecord 获取单条历史记录详情
func (s *walletHistoryService) GetHistoryRecord(ctx context.Context, recordID uint, userID int64) (*models.WalletHistory, error) {
	history, err := s.walletHistoryRepo.GetByID(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("历史记录不存在: %w", err)
	}

	// 验证用户权限
	if history.UserID != userID {
		return nil, fmt.Errorf("无权访问该记录")
	}

	return history, nil
}

// CreateRechargeRecord 创建充值记录
func (s *walletHistoryService) CreateRechargeRecord(ctx context.Context, userID int64, amount, balanceBefore, balanceAfter string, relatedID, txHash string) error {
	history := &models.WalletHistory{
		UserID:        userID,
		Type:          models.WalletHistoryTypeRecharge,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        models.WalletHistoryStatusCompleted,
		Description:   fmt.Sprintf("USDT充值 - 订单号: %s", relatedID),
		RelatedType:   "recharge_order",
		RelatedID:     relatedID,
		TxHash:        txHash,
	}

	if err := s.walletHistoryRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("创建充值记录失败: %w", err)
	}

	return nil
}

// CreatePaymentRecord 创建支付记录
func (s *walletHistoryService) CreatePaymentRecord(ctx context.Context, userID int64, amount, balanceBefore, balanceAfter string, relatedID, description string) error {
	history := &models.WalletHistory{
		UserID:        userID,
		Type:          models.WalletHistoryTypePayment,
		Amount:        fmt.Sprintf("-%s", amount), // 支出为负数
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        models.WalletHistoryStatusCompleted,
		Description:   description,
		RelatedType:   "order",
		RelatedID:     relatedID,
	}

	if err := s.walletHistoryRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("创建支付记录失败: %w", err)
	}

	return nil
}

// UpdateRecordStatus 更新记录状态
func (s *walletHistoryService) UpdateRecordStatus(ctx context.Context, recordID uint, status models.WalletHistoryStatus) error {
	if err := s.walletHistoryRepo.UpdateStatus(ctx, recordID, status); err != nil {
		return fmt.Errorf("更新记录状态失败: %w", err)
	}

	return nil
}

// getHistoriesByMultipleFilters 根据多个筛选条件获取历史记录（简化实现）
func (s *walletHistoryService) getHistoriesByMultipleFilters(ctx context.Context, userID int64, filters WalletHistoryFilters) ([]*models.WalletHistory, error) {
	// 这里简化实现，先获取所有记录再筛选
	// 实际生产环境中应该在数据库层面进行复合查询
	allHistories, err := s.walletHistoryRepo.GetByUserID(ctx, userID, 0, 0)
	if err != nil {
		return nil, err
	}

	var filtered []*models.WalletHistory
	for _, history := range allHistories {
		// 按类型筛选
		if filters.Type != "" && history.Type != filters.Type {
			continue
		}

		// 按状态筛选
		if filters.Status != "" && history.Status != filters.Status {
			continue
		}

		filtered = append(filtered, history)
	}

	// 应用分页
	start := filters.Offset
	end := start + filters.Limit
	if start > len(filtered) {
		return []*models.WalletHistory{}, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// parseDateRange 解析日期范围
func (s *walletHistoryService) parseDateRange(startDateStr, endDateStr string) (time.Time, time.Time) {
	var startDate, endDate time.Time

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(24 * time.Hour) // 包含整天
		}
	} else {
		endDate = time.Now()
	}

	return startDate, endDate
}

// CreateRefundRecord 创建退款记录（用于订单失败时的退款）
func (s *walletHistoryService) CreateRefundRecord(ctx context.Context, userID int64, amount, balanceBefore, balanceAfter string, relatedID, description string) error {
	history := &models.WalletHistory{
		UserID:        userID,
		Type:          models.WalletHistoryTypeRefund,
		Amount:        amount, // 退款为正数
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        models.WalletHistoryStatusCompleted,
		Description:   description,
		RelatedType:   "order",
		RelatedID:     relatedID,
	}

	if err := s.walletHistoryRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("创建退款记录失败: %w", err)
	}

	return nil
}
