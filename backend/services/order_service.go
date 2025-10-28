package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// OrderService 订单服务接口
type OrderService interface {
	CreateOrder(ctx context.Context, userID int64, productID int) (*models.Order, error)
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	GetOrderByOrderNo(ctx context.Context, orderNo string) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID int64, limit, offset int) ([]*models.Order, error)
	PayOrder(ctx context.Context, orderNo string) error
	CompleteOrder(ctx context.Context, orderNo string) error
	CancelOrder(ctx context.Context, orderNo string) error
	GetOrderStats(ctx context.Context, userID int64) (*OrderStats, error)
}

// OrderStats 订单统计信息
type OrderStats struct {
	TotalOrders     int64  `json:"total_orders"`
	PendingOrders   int64  `json:"pending_orders"`
	CompletedOrders int64  `json:"completed_orders"`
	TotalAmount     string `json:"total_amount"`
}

// orderService 订单服务实现
type orderService struct {
	orderRepo     repository.OrderRepository
	productRepo   repository.ProductRepository
	walletService WalletService
}

// NewOrderService 创建订单服务实例
func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	walletService WalletService,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		walletService: walletService,
	}
}

// CreateOrder 创建订单
func (s *orderService) CreateOrder(ctx context.Context, userID int64, productID int) (*models.Order, error) {
	// 获取产品信息
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// 检查产品状态
	if product.Status != "active" {
		return nil, errors.New("product is not available")
	}

	// 创建订单
	order := &models.Order{
		UserID:      userID,
		ProductID:   productID,
		ProductName: product.Name,
		Amount:      fmt.Sprintf("%.2f", product.Price),
		Status:      models.OrderStatusPending,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrderByID 根据ID获取订单
func (s *orderService) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	return order, nil
}

// GetOrderByOrderNo 根据订单号获取订单
func (s *orderService) GetOrderByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	return order, nil
}

// GetUserOrders 获取用户订单列表
func (s *orderService) GetUserOrders(ctx context.Context, userID int64, limit, offset int) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return orders, nil
}

// PayOrder 支付订单
func (s *orderService) PayOrder(ctx context.Context, orderNo string) error {
	// 获取订单
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending {
		return errors.New("order cannot be paid")
	}

	// 处理支付
	result, err := s.walletService.ProcessPayment(ctx, order.UserID, order.ProductID, order.Amount)
	if err != nil || !result.Success {
		return fmt.Errorf("payment failed: %w", err)
	}

	// 更新订单状态
	now := time.Now()
	order.Status = models.OrderStatusPaid
	order.PaidAt = &now

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// CompleteOrder 完成订单
func (s *orderService) CompleteOrder(ctx context.Context, orderNo string) error {
	// 获取订单
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPaid {
		return errors.New("order cannot be completed")
	}

	// 更新订单状态
	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.CompletedAt = &now

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// CancelOrder 取消订单
func (s *orderService) CancelOrder(ctx context.Context, orderNo string) error {
	// 获取订单
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending {
		return errors.New("order cannot be cancelled")
	}

	// 更新订单状态
	order.Status = models.OrderStatusCancelled

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// GetOrderStats 获取订单统计信息
func (s *orderService) GetOrderStats(ctx context.Context, userID int64) (*OrderStats, error) {
	// 获取所有订单
	orders, err := s.orderRepo.GetByUserID(ctx, userID, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	stats := &OrderStats{
		TotalOrders: int64(len(orders)),
		TotalAmount: "0.00",
	}

	var totalAmount float64
	for _, order := range orders {
		switch order.Status {
		case models.OrderStatusPending:
			stats.PendingOrders++
		case models.OrderStatusCompleted:
			stats.CompletedOrders++
		}

		// 累计总金额（已支付和已完成的订单）
		if order.Status == models.OrderStatusPaid || order.Status == models.OrderStatusCompleted {
			amount, _ := parseDecimal(order.Amount)
			amountFloat, _ := amount.Float64()
			totalAmount += amountFloat
		}
	}

	stats.TotalAmount = fmt.Sprintf("%.2f", totalAmount)

	return stats, nil
}
