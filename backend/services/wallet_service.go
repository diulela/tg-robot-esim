package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// WalletBalance 钱包余额信息
type WalletBalance struct {
	Balance       string `json:"balance"`
	FrozenBalance string `json:"frozen_balance"`
	TotalIncome   string `json:"total_income"`
	TotalExpense  string `json:"total_expense"`
}

// PaymentResult 支付结果
type PaymentResult struct {
	Success bool   `json:"success"`
	OrderID uint   `json:"order_id"`
	Message string `json:"message"`
}

// WalletService 钱包服务接口
type WalletService interface {
	GetBalance(ctx context.Context, userID int64) (*WalletBalance, error)
	CreateRechargeOrder(ctx context.Context, userID int64, amount string) (*models.RechargeOrder, error)
	ProcessPayment(ctx context.Context, userID int64, productID int, amount string) (*PaymentResult, error)
	ProcessRecharge(ctx context.Context, txHash string, amount string) error
	DeductBalance(ctx context.Context, userID int64, amount string) error

	// 新增方法用于充值功能
	GetWallet(ctx context.Context, userID int64) (*models.Wallet, error)
	GetOrCreateWallet(ctx context.Context, userID int64) (*models.Wallet, error)
	AddBalanceWithRemark(ctx context.Context, userID int64, amount string, remark string) error

	// eSIM 订单处理相关方法
	// FreezeBalance 冻结余额（不记录 wallet_history，仅内部状态变更）
	FreezeBalance(ctx context.Context, userID int64, amount string, relatedID string, description string) error

	// UnfreezeBalance 解冻余额（退还到可用余额）
	// 会创建 wallet_history 记录：type=refund, status=completed
	UnfreezeBalance(ctx context.Context, userID int64, amount string, relatedID string, description string) error

	// ConfirmFrozenPayment 确认冻结金额的支付（从冻结余额扣除）
	// 会创建 wallet_history 记录：type=payment, status=completed
	ConfirmFrozenPayment(ctx context.Context, userID int64, amount string, relatedID string, description string) error

	// AddBalance 增加余额（充值等）
	AddBalance(ctx context.Context, userID int64, amount string, relatedID string, description string) error

	// HasSufficientBalance 检查是否有足够的可用余额
	HasSufficientBalance(ctx context.Context, userID int64, amount string) (bool, error)

	// GetFrozenBalance 获取冻结余额
	GetFrozenBalance(ctx context.Context, userID int64) (string, error)
}

// walletService 钱包服务实现
type walletService struct {
	walletRepo           repository.WalletRepository
	rechargeOrderRepo    repository.RechargeOrderRepository
	blockchainService    BlockchainService
	walletHistoryService WalletHistoryService
}

// NewWalletService 创建钱包服务实例
func NewWalletService(
	walletRepo repository.WalletRepository,
	rechargeOrderRepo repository.RechargeOrderRepository,
	blockchainService BlockchainService,
	walletHistoryService WalletHistoryService,
) WalletService {
	return &walletService{
		walletRepo:           walletRepo,
		rechargeOrderRepo:    rechargeOrderRepo,
		blockchainService:    blockchainService,
		walletHistoryService: walletHistoryService,
	}
}

// GetBalance 获取钱包余额
func (s *walletService) GetBalance(ctx context.Context, userID int64) (*WalletBalance, error) {
	// 使用 GetOrCreate 确保钱包存在
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create wallet: %w", err)
	}

	return &WalletBalance{
		Balance:       wallet.Balance,
		FrozenBalance: wallet.FrozenBalance,
		TotalIncome:   wallet.TotalIncome,
		TotalExpense:  wallet.TotalExpense,
	}, nil
}

// CreateRechargeOrder 创建充值订单
func (s *walletService) CreateRechargeOrder(ctx context.Context, userID int64, amount string) (*models.RechargeOrder, error) {
	// 验证金额
	amountFloat, err := parseDecimal(amount)
	if err != nil || amountFloat.Cmp(big.NewFloat(0)) <= 0 {
		return nil, errors.New("invalid amount")
	}

	// 获取充值地址（从配置或区块链服务）
	walletAddress := s.getRechargeAddress()

	// 创建充值订单
	order := &models.RechargeOrder{
		UserID:        userID,
		Amount:        amount,
		WalletAddress: walletAddress,
		Status:        models.RechargeStatusPending,
	}

	if err := s.rechargeOrderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create recharge order: %w", err)
	}

	return order, nil
}

// ProcessPayment 处理支付
func (s *walletService) ProcessPayment(ctx context.Context, userID int64, productID int, amount string) (*PaymentResult, error) {
	// 获取或创建钱包
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return &PaymentResult{
			Success: false,
			Message: "获取钱包失败",
		}, err
	}

	// 检查余额是否充足
	balance, _ := parseDecimal(wallet.Balance)
	payAmount, _ := parseDecimal(amount)

	if balance.Cmp(payAmount) < 0 {
		return &PaymentResult{
			Success: false,
			Message: "余额不足",
		}, errors.New("insufficient balance")
	}

	// 扣除余额
	newBalance := new(big.Float).Sub(balance, payAmount)
	wallet.Balance = newBalance.Text('f', 8)

	// 更新总支出
	totalExpense, _ := parseDecimal(wallet.TotalExpense)
	newTotalExpense := new(big.Float).Add(totalExpense, payAmount)
	wallet.TotalExpense = newTotalExpense.Text('f', 8)

	// 保存钱包
	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return &PaymentResult{
			Success: false,
			Message: "支付失败",
		}, err
	}

	return &PaymentResult{
		Success: true,
		Message: "支付成功",
	}, nil
}

// ProcessRecharge 处理充值（区块链确认后调用）
func (s *walletService) ProcessRecharge(ctx context.Context, txHash string, amount string) error {
	// 查找充值订单
	order, err := s.rechargeOrderRepo.GetByTxHash(ctx, txHash)
	if err != nil {
		return fmt.Errorf("recharge order not found: %w", err)
	}

	// 检查订单状态
	if order.Status != models.RechargeStatusPending {
		return errors.New("order already processed")
	}

	// 获取或创建钱包
	wallet, err := s.walletRepo.GetOrCreate(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get or create wallet: %w", err)
	}

	// 增加余额
	balance, _ := parseDecimal(wallet.Balance)
	rechargeAmount, _ := parseDecimal(amount)
	newBalance := new(big.Float).Add(balance, rechargeAmount)
	wallet.Balance = newBalance.Text('f', 8)

	// 更新总收入
	totalIncome, _ := parseDecimal(wallet.TotalIncome)
	newTotalIncome := new(big.Float).Add(totalIncome, rechargeAmount)
	wallet.TotalIncome = newTotalIncome.Text('f', 8)

	// 保存钱包
	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	// 更新充值订单状态
	order.Status = models.RechargeStatusConfirmed
	if err := s.rechargeOrderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update recharge order: %w", err)
	}

	return nil
}

// FreezeBalance 冻结余额（不记录 wallet_history，仅内部状态变更）
func (s *walletService) FreezeBalance(ctx context.Context, userID int64, amount string, relatedID string, description string) error {
	// 验证金额格式
	freezeAmount, err := parseDecimal(amount)
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	if freezeAmount.Cmp(big.NewFloat(0)) <= 0 {
		return errors.New("amount must be positive")
	}

	// 使用原子操作：减少可用余额，增加冻结余额
	negativeAmount := new(big.Float).Neg(freezeAmount)
	balanceDelta := negativeAmount.Text('f', 8) // 减少可用余额
	frozenDelta := amount                       // 增加冻结余额

	return s.walletRepo.UpdateBalanceAtomic(ctx, userID, balanceDelta, frozenDelta)
}

// UnfreezeBalance 解冻余额（退还到可用余额）
// 会创建 wallet_history 记录：type=refund, status=completed
func (s *walletService) UnfreezeBalance(ctx context.Context, userID int64, amount string, relatedID string, description string) error {
	// 验证金额格式
	unfreezeAmount, err := parseDecimal(amount)
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	if unfreezeAmount.Cmp(big.NewFloat(0)) <= 0 {
		return errors.New("amount must be positive")
	}

	// 使用原子操作：增加可用余额，减少冻结余额
	negativeAmount := new(big.Float).Neg(unfreezeAmount)
	balanceDelta := amount                     // 增加可用余额
	frozenDelta := negativeAmount.Text('f', 8) // 减少冻结余额

	// 获取操作前的余额信息
	walletBefore, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get wallet before unfreeze: %w", err)
	}

	err = s.walletRepo.UpdateBalanceAtomic(ctx, userID, balanceDelta, frozenDelta)
	if err != nil {
		return err
	}

	// 获取操作后的余额信息
	walletAfter, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get wallet after unfreeze: %w", err)
	}

	// 记录 wallet_history（type=refund）
	if s.walletHistoryService != nil {
		err = s.walletHistoryService.CreateRefundRecord(
			ctx,
			userID,
			amount,
			walletBefore.Balance,
			walletAfter.Balance,
			relatedID,
			description,
		)
		if err != nil {
			// 记录历史失败不应该影响主要操作，只记录日志
			fmt.Printf("Warning: failed to create refund history record: %v\n", err)
		}
	}

	return nil
}

// DeductBalance 扣除余额
func (s *walletService) DeductBalance(ctx context.Context, userID int64, amount string) error {
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	balance, _ := parseDecimal(wallet.Balance)
	deductAmount, _ := parseDecimal(amount)

	if balance.Cmp(deductAmount) < 0 {
		return errors.New("insufficient balance")
	}

	newBalance := new(big.Float).Sub(balance, deductAmount)
	wallet.Balance = newBalance.Text('f', 8)

	return s.walletRepo.Update(ctx, wallet)
}

// AddBalance 增加余额（新签名，支持 relatedID 和 description）
func (s *walletService) AddBalance(ctx context.Context, userID int64, amount string, relatedID string, description string) error {
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	balance, _ := parseDecimal(wallet.Balance)
	addAmount, _ := parseDecimal(amount)
	newBalance := new(big.Float).Add(balance, addAmount)

	wallet.Balance = newBalance.Text('f', 8)

	// TODO: 记录钱包历史（使用 relatedID 和 description）
	// 这里可以调用 WalletHistoryService 来记录历史

	return s.walletRepo.Update(ctx, wallet)
}

// getRechargeAddress 获取充值地址
func (s *walletService) getRechargeAddress() string {
	// TODO: 从配置或区块链服务获取充值地址
	// 这里暂时返回一个固定地址
	return "TYourWalletAddressHere"
}

// GetWallet 获取用户钱包
func (s *walletService) GetWallet(ctx context.Context, userID int64) (*models.Wallet, error) {
	// 使用 GetOrCreate 确保钱包存在
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户钱包失败: %w", err)
	}
	return wallet, nil
}

// GetOrCreateWallet 获取或创建用户钱包
func (s *walletService) GetOrCreateWallet(ctx context.Context, userID int64) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取或创建用户钱包失败: %w", err)
	}
	return wallet, nil
}

// AddBalanceWithRemark 增加余额（带备注）- 带重试机制处理数据库锁定
func (s *walletService) AddBalanceWithRemark(ctx context.Context, userID int64, amount string, remark string) error {
	maxRetries := 3
	var err error

	for i := 0; i < maxRetries; i++ {
		err = s.addBalanceWithRemarkOnce(ctx, userID, amount, remark)
		if err == nil {
			return nil // 成功，退出重试循环
		}

		// 如果是数据库锁定错误，进行重试
		if i < maxRetries-1 && (err.Error() == "database is locked" || err.Error() == "database locked") {
			waitTime := time.Duration(100*(i+1)) * time.Millisecond // 100ms, 200ms, 300ms
			fmt.Printf("钱包操作数据库锁定，%v 后重试 (第 %d/%d 次): %v\n", waitTime, i+1, maxRetries, err)
			time.Sleep(waitTime)
			continue
		}

		// 其他错误或最后一次重试失败，直接返回
		break
	}

	return err
}

// addBalanceWithRemarkOnce 单次增加余额操作
func (s *walletService) addBalanceWithRemarkOnce(ctx context.Context, userID int64, amount string, remark string) error {
	// 获取或创建钱包
	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return err
	}

	// 解析金额
	balance, _ := parseDecimal(wallet.Balance)
	addAmount, err := parseDecimal(amount)
	if err != nil {
		return fmt.Errorf("金额格式错误: %w", err)
	}

	// 计算新余额
	newBalance := new(big.Float).Add(balance, addAmount)
	wallet.Balance = newBalance.Text('f', 8)

	// 更新总收入
	totalIncome, _ := parseDecimal(wallet.TotalIncome)
	newTotalIncome := new(big.Float).Add(totalIncome, addAmount)
	wallet.TotalIncome = newTotalIncome.Text('f', 8)

	// 保存钱包
	if err := s.walletRepo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("更新钱包余额失败: %w", err)
	}

	// TODO: 记录交易日志（包含备注）
	// 这里可以添加交易记录的逻辑

	return nil
}

// ConfirmFrozenPayment 确认冻结金额的支付（从冻结余额扣除）
// 会创建 wallet_history 记录：type=payment, status=completed
func (s *walletService) ConfirmFrozenPayment(ctx context.Context, userID int64, amount string, relatedID string, description string) error {
	// 验证金额格式
	payAmount, err := parseDecimal(amount)
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	if payAmount.Cmp(big.NewFloat(0)) <= 0 {
		return errors.New("amount must be positive")
	}

	// 使用原子操作：减少冻结余额
	negativeAmount := new(big.Float).Neg(payAmount)
	balanceDelta := "0"                        // 可用余额不变
	frozenDelta := negativeAmount.Text('f', 8) // 减少冻结余额

	err = s.walletRepo.UpdateBalanceAtomic(ctx, userID, balanceDelta, frozenDelta)
	if err != nil {
		return err
	}

	// 更新总支出（需要单独处理，因为 UpdateBalanceAtomic 不处理统计字段）
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get wallet for expense update: %w", err)
	}

	totalExpense, _ := parseDecimal(wallet.TotalExpense)
	newTotalExpense := new(big.Float).Add(totalExpense, payAmount)
	wallet.TotalExpense = newTotalExpense.Text('f', 8)

	// 获取操作前的余额信息（用于历史记录）
	balanceBefore := wallet.Balance // 可用余额在这个操作中没有变化

	err = s.walletRepo.Update(ctx, wallet)
	if err != nil {
		return fmt.Errorf("failed to update total expense: %w", err)
	}

	// 记录 wallet_history（type=payment）
	if s.walletHistoryService != nil {
		err = s.walletHistoryService.CreatePaymentRecord(
			ctx,
			userID,
			amount,
			balanceBefore,
			wallet.Balance, // 可用余额没有变化，但记录当前值
			relatedID,
			description,
		)
		if err != nil {
			// 记录历史失败不应该影响主要操作，只记录日志
			fmt.Printf("Warning: failed to create payment history record: %v\n", err)
		}
	}

	return nil
}

// HasSufficientBalance 检查是否有足够的可用余额
func (s *walletService) HasSufficientBalance(ctx context.Context, userID int64, amount string) (bool, error) {
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return false, err
	}

	balance, _ := parseDecimal(wallet.Balance)
	requiredAmount, _ := parseDecimal(amount)

	return balance.Cmp(requiredAmount) >= 0, nil
}

// GetFrozenBalance 获取冻结余额
func (s *walletService) GetFrozenBalance(ctx context.Context, userID int64) (string, error) {
	wallet, err := s.walletRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return "0", err
	}

	return wallet.FrozenBalance, nil
}

// parseDecimal 解析 decimal 字符串为 big.Float
func parseDecimal(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	if err != nil {
		return big.NewFloat(0), err
	}
	return f, nil
}
