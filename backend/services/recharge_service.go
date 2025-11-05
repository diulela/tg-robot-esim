package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"

	"gorm.io/gorm"
)

// rechargeService 充值服务实现
type rechargeService struct {
	rechargeRepo        repository.RechargeOrderRepository
	walletService       WalletService
	blockchainService   BlockchainService
	notificationService NotificationService
	db                  *gorm.DB
	depositAddress      string
	minAmount           float64
	maxAmount           float64
}

// NewRechargeService 创建充值服务实例
func NewRechargeService(
	rechargeRepo repository.RechargeOrderRepository,
	walletService WalletService,
	blockchainService BlockchainService,
	notificationService NotificationService,
	db *gorm.DB,
	depositAddress string,
	minAmount, maxAmount float64,
) RechargeService {
	return &rechargeService{
		rechargeRepo:        rechargeRepo,
		walletService:       walletService,
		blockchainService:   blockchainService,
		notificationService: notificationService,
		db:                  db,
		depositAddress:      depositAddress,
		minAmount:           minAmount,
		maxAmount:           maxAmount,
	}
}

// CreateRechargeOrder 创建充值订单
func (s *rechargeService) CreateRechargeOrder(ctx context.Context, userID int64, amount string) (*models.RechargeOrder, error) {
	// 1. 验证充值金额
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil, fmt.Errorf("充值金额格式错误")
	}

	if amountFloat < s.minAmount {
		return nil, fmt.Errorf("充值金额不能低于 %.2f USDT", s.minAmount)
	}

	if amountFloat > s.maxAmount {
		return nil, fmt.Errorf("充值金额不能超过 %.2f USDT", s.maxAmount)
	}

	// 2. 生成唯一的精确金额
	exactAmount, err := s.GenerateExactAmount(ctx, amount)
	if err != nil {
		return nil, fmt.Errorf("生成精确金额失败: %w", err)
	}

	// 3. 创建充值订单
	order := &models.RechargeOrder{
		UserID:        userID,
		Amount:        amount,
		ExactAmount:   exactAmount,
		WalletAddress: s.depositAddress,
		Status:        models.RechargeStatusPending,
		ExpiresAt:     time.Now().Add(30 * time.Minute),
	}

	if err := s.rechargeRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("创建充值订单失败: %w", err)
	}

	return order, nil
}

// GetRechargeOrder 获取充值订单详情
func (s *rechargeService) GetRechargeOrder(ctx context.Context, orderNo string) (*models.RechargeOrder, error) {
	order, err := s.rechargeRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}
	return order, nil
}

// GetUserRechargeHistory 获取用户充值历史
func (s *rechargeService) GetUserRechargeHistory(ctx context.Context, userID int64, limit, offset int) ([]*models.RechargeOrder, int64, error) {
	orders, err := s.rechargeRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取充值历史失败: %w", err)
	}

	// 获取总数（简化实现，实际应该在仓储层实现）
	allOrders, err := s.rechargeRepo.GetByUserID(ctx, userID, 0, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("获取充值历史总数失败: %w", err)
	}

	return orders, int64(len(allOrders)), nil
}

// CheckRechargeStatus 手动检查充值状态
func (s *rechargeService) CheckRechargeStatus(ctx context.Context, orderNo string) (*models.RechargeOrder, error) {
	// 1. 获取订单
	order, err := s.GetRechargeOrder(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	// 2. 检查订单状态
	if order.Status != models.RechargeStatusPending {
		return order, nil
	}

	// 3. 检查是否过期
	if order.IsExpired() {
		order.Status = models.RechargeStatusExpired
		if err := s.rechargeRepo.Update(ctx, order); err != nil {
			return nil, fmt.Errorf("更新订单状态失败: %w", err)
		}
		return order, nil
	}

	// 4. 查询区块链交易
	transactions, err := s.blockchainService.GetAddressIncomingTransactions(ctx, s.depositAddress, order.ExactAmount)
	if err != nil {
		return nil, fmt.Errorf("查询区块链交易失败: %w", err)
	}

	// 5. 查找匹配的交易
	for _, tx := range transactions {
		if s.blockchainService.MatchTransactionAmount(tx.Amount, order.ExactAmount) {
			// 检查确认数
			if tx.Confirmations >= 19 {
				// 确认充值
				if err := s.ConfirmRecharge(ctx, order, tx.TxHash); err != nil {
					return nil, fmt.Errorf("确认充值失败: %w", err)
				}
				// 重新获取更新后的订单
				return s.GetRechargeOrder(ctx, orderNo)
			}
		}
	}

	return order, nil
}

// ProcessPendingRecharges 处理待支付的充值订单（定时任务调用）
func (s *rechargeService) ProcessPendingRecharges(ctx context.Context) error {
	// 1. 获取所有待支付订单
	orders, err := s.rechargeRepo.GetPendingOrders(ctx)
	if err != nil {
		return fmt.Errorf("获取待支付订单失败: %w", err)
	}

	// 2. 处理每个订单
	for _, order := range orders {
		// 检查是否过期
		if order.IsExpired() {
			order.Status = models.RechargeStatusExpired
			if err := s.rechargeRepo.Update(ctx, order); err != nil {
				// 记录错误但继续处理其他订单
				fmt.Printf("更新过期订单失败: %v\n", err)
			}
			continue
		}

		// 查询区块链交易（带重试机制）
		transactions, err := s.queryTransactionsWithRetry(ctx, s.depositAddress, order.ExactAmount, 3)
		if err != nil {
			// 记录错误但继续处理其他订单
			fmt.Printf("查询订单 %s 的区块链交易失败: %v\n", order.OrderNo, err)
			continue
		}

		// 查找匹配的交易
		for _, tx := range transactions {
			if s.blockchainService.MatchTransactionAmount(tx.Amount, order.ExactAmount) {
				// 检查确认数
				if tx.Confirmations >= 19 {
					// 确认充值
					if err := s.ConfirmRecharge(ctx, order, tx.TxHash); err != nil {
						fmt.Printf("确认订单 %s 充值失败: %v\n", order.OrderNo, err)
					}
					break
				}
			}
		}
	}

	return nil
}

// ConfirmRecharge 确认充值并更新余额，并发送 Telegram 通知
func (s *rechargeService) ConfirmRecharge(ctx context.Context, order *models.RechargeOrder, txHash string) error {
	// 使用数据库事务确保原子性
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 检查订单状态，防止重复处理
		if order.Status != models.RechargeStatusPending {
			return fmt.Errorf("订单已处理，当前状态: %s", order.Status)
		}

		// 2. 检查交易哈希是否已被使用
		var existingOrder models.RechargeOrder
		err := tx.Where("tx_hash = ? AND id != ?", txHash, order.ID).First(&existingOrder).Error
		if err == nil {
			return fmt.Errorf("交易哈希已被使用")
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("检查交易哈希失败: %w", err)
		}

		// 3. 更新订单状态
		now := time.Now()
		order.Status = models.RechargeStatusConfirmed
		order.TxHash = txHash
		order.ConfirmedAt = &now

		if err := tx.Save(order).Error; err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 4. 增加用户余额
		remark := fmt.Sprintf("充值到账，订单号: %s", order.OrderNo)
		if err := s.walletService.AddBalanceWithRemark(ctx, order.UserID, order.Amount, remark); err != nil {
			return fmt.Errorf("增加用户余额失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 5. 发送 Telegram 通知（在事务外执行，失败不影响充值确认）
	if err := s.notificationService.SendRechargeSuccessNotification(ctx, order.UserID, order.Amount, order.OrderNo); err != nil {
		// 通知发送失败不影响充值确认
		fmt.Printf("发送充值成功通知失败: %v\n", err)
	}

	return nil
}

// ExpireOldOrders 将过期订单标记为已过期
func (s *rechargeService) ExpireOldOrders(ctx context.Context) error {
	return s.rechargeRepo.ExpireOldOrders(ctx)
}

// GenerateExactAmount 生成唯一的精确金额
func (s *rechargeService) GenerateExactAmount(ctx context.Context, baseAmount string) (string, error) {
	// 解析基础金额
	baseFloat, err := strconv.ParseFloat(baseAmount, 64)
	if err != nil {
		return "", fmt.Errorf("基础金额格式错误")
	}

	// 最多尝试 100 次生成唯一的精确金额
	for i := 0; i < 100; i++ {
		// 生成 4 位随机数 (0000-9999)
		randomNum, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}

		// 构造精确金额：基础金额.随机4位数
		exactAmount := fmt.Sprintf("%.4f", baseFloat+float64(randomNum.Int64())/10000)

		// 检查是否已存在
		exists, err := s.rechargeRepo.IsExactAmountExists(ctx, exactAmount)
		if err != nil {
			return "", fmt.Errorf("检查精确金额唯一性失败: %w", err)
		}

		if !exists {
			return exactAmount, nil
		}
	}

	return "", fmt.Errorf("生成唯一精确金额失败，请稍后重试")
}

// queryTransactionsWithRetry 带重试机制的区块链交易查询
func (s *rechargeService) queryTransactionsWithRetry(ctx context.Context, address, minAmount string, maxRetries int) ([]*TransactionInfo, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		transactions, err := s.blockchainService.GetAddressIncomingTransactions(ctx, address, minAmount)
		if err == nil {
			return transactions, nil
		}

		lastErr = err

		// 如果不是最后一次重试，等待一段时间后重试
		if i < maxRetries-1 {
			// 指数退避：1秒、2秒、4秒
			waitTime := time.Duration(1<<uint(i)) * time.Second
			fmt.Printf("区块链查询失败，%v 后重试 (第 %d/%d 次): %v\n", waitTime, i+1, maxRetries, err)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
				// 继续重试
			}
		}
	}

	return nil, fmt.Errorf("重试 %d 次后仍然失败: %w", maxRetries, lastErr)
}
