package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"

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
	FreezeBalance(ctx context.Context, userID int64, amount string) error
	UnfreezeBalance(ctx context.Context, userID int64, amount string) error
	DeductBalance(ctx context.Context, userID int64, amount string) error
	AddBalance(ctx context.Context, userID int64, amount string) error
}

// walletService 钱包服务实现
type walletService struct {
	walletRepo        repository.WalletRepository
	rechargeOrderRepo repository.RechargeOrderRepository
	blockchainService BlockchainService
}

// NewWalletService 创建钱包服务实例
func NewWalletService(
	walletRepo repository.WalletRepository,
	rechargeOrderRepo repository.RechargeOrderRepository,
	blockchainService BlockchainService,
) WalletService {
	return &walletService{
		walletRepo:        walletRepo,
		rechargeOrderRepo: rechargeOrderRepo,
		blockchainService: blockchainService,
	}
}

// GetBalance 获取钱包余额
func (s *walletService) GetBalance(ctx context.Context, userID int64) (*WalletBalance, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
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
	// 获取钱包
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return &PaymentResult{
			Success: false,
			Message: "钱包不存在",
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

	// 获取钱包
	wallet, err := s.walletRepo.GetByUserID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
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

// FreezeBalance 冻结余额
func (s *walletService) FreezeBalance(ctx context.Context, userID int64, amount string) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	balance, _ := parseDecimal(wallet.Balance)
	freezeAmount, _ := parseDecimal(amount)

	if balance.Cmp(freezeAmount) < 0 {
		return errors.New("insufficient balance")
	}

	newBalance := new(big.Float).Sub(balance, freezeAmount)
	frozenBalance, _ := parseDecimal(wallet.FrozenBalance)
	newFrozenBalance := new(big.Float).Add(frozenBalance, freezeAmount)

	wallet.Balance = newBalance.Text('f', 8)
	wallet.FrozenBalance = newFrozenBalance.Text('f', 8)

	return s.walletRepo.Update(ctx, wallet)
}

// UnfreezeBalance 解冻余额
func (s *walletService) UnfreezeBalance(ctx context.Context, userID int64, amount string) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	frozenBalance, _ := parseDecimal(wallet.FrozenBalance)
	unfreezeAmount, _ := parseDecimal(amount)

	if frozenBalance.Cmp(unfreezeAmount) < 0 {
		return errors.New("insufficient frozen balance")
	}

	balance, _ := parseDecimal(wallet.Balance)
	newBalance := new(big.Float).Add(balance, unfreezeAmount)
	newFrozenBalance := new(big.Float).Sub(frozenBalance, unfreezeAmount)

	wallet.Balance = newBalance.Text('f', 8)
	wallet.FrozenBalance = newFrozenBalance.Text('f', 8)

	return s.walletRepo.Update(ctx, wallet)
}

// DeductBalance 扣除余额
func (s *walletService) DeductBalance(ctx context.Context, userID int64, amount string) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
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

// AddBalance 增加余额
func (s *walletService) AddBalance(ctx context.Context, userID int64, amount string) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	balance, _ := parseDecimal(wallet.Balance)
	addAmount, _ := parseDecimal(amount)
	newBalance := new(big.Float).Add(balance, addAmount)

	wallet.Balance = newBalance.Text('f', 8)

	return s.walletRepo.Update(ctx, wallet)
}

// getRechargeAddress 获取充值地址
func (s *walletService) getRechargeAddress() string {
	// TODO: 从配置或区块链服务获取充值地址
	// 这里暂时返回一个固定地址
	return "TYourWalletAddressHere"
}

// parseDecimal 解析 decimal 字符串为 big.Float
func parseDecimal(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	if err != nil {
		return big.NewFloat(0), err
	}
	return f, nil
}
