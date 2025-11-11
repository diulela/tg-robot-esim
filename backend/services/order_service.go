package services

import (
	"context"
	"encoding/json"
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

	// eSIM 订单处理相关方法
	// CreateEsimOrder 创建 eSIM 商品购买订单（含余额冻结）
	CreateEsimOrder(ctx context.Context, req *CreateEsimOrderRequest) (*EsimOrderResponse, error)

	// GetOrderWithDetail 获取包含详情的订单信息
	GetOrderWithDetail(ctx context.Context, orderID uint, userID int64) (*OrderWithDetail, error)

	// ProcessOrderCompletion 处理订单完成（确认扣费）
	ProcessOrderCompletion(ctx context.Context, orderID uint, providerOrderData *ProviderOrderData) error

	// ProcessOrderFailure 处理订单失败（退还冻结金额）
	ProcessOrderFailure(ctx context.Context, orderID uint, reason string) error

	// UpdateOrderSyncInfo 更新订单同步信息
	UpdateOrderSyncInfo(ctx context.Context, orderID uint, syncAttempts int, nextSyncAt *time.Time) error

	// GetUserOrdersWithFilters 根据筛选条件获取用户订单列表
	GetUserOrdersWithFilters(ctx context.Context, userID int64, status models.OrderStatus, limit, offset int) ([]*models.Order, int64, error)
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
	orderRepo       repository.OrderRepository
	orderDetailRepo repository.OrderDetailRepository
	productRepo     repository.ProductRepository
	walletService   WalletService
}

// NewOrderService 创建订单服务实例
func NewOrderService(
	orderRepo repository.OrderRepository,
	orderDetailRepo repository.OrderDetailRepository,
	productRepo repository.ProductRepository,
	walletService WalletService,
) OrderService {
	return &orderService{
		orderRepo:       orderRepo,
		orderDetailRepo: orderDetailRepo,
		productRepo:     productRepo,
		walletService:   walletService,
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

// CreateEsimOrder 创建 eSIM 商品购买订单（含余额冻结）
func (s *orderService) CreateEsimOrder(ctx context.Context, req *CreateEsimOrderRequest) (*EsimOrderResponse, error) {
	// 1. 验证输入参数
	if req.UserID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	if req.ProductID == 0 {
		return nil, errors.New("产品ID不能为空")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("购买数量必须大于0")
	}

	// 2. 获取产品信息
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("产品不存在: %w", err)
	}

	// 检查产品状态
	if product.Status != "active" {
		return nil, errors.New("产品暂不可用")
	}

	// 3. 计算订单金额
	unitPrice := product.Price
	totalAmount := unitPrice * float64(req.Quantity)
	totalAmountStr := fmt.Sprintf("%.4f", totalAmount)

	// 验证前端传入的金额是否正确
	if req.TotalAmount != totalAmountStr {
		return nil, fmt.Errorf("订单金额不匹配，期望: %s，实际: %s", totalAmountStr, req.TotalAmount)
	}

	// 4. 检查用户余额是否充足
	hasSufficient, err := s.walletService.HasSufficientBalance(ctx, req.UserID, req.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("检查余额失败: %w", err)
	}
	if !hasSufficient {
		return nil, errors.New("余额不足，请先充值")
	}

	// 5. 创建订单记录
	order := &models.Order{
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		UnitPrice:   fmt.Sprintf("%.4f", unitPrice),
		Amount:      req.TotalAmount,
		Status:      models.OrderStatusProcessing, // 直接设为处理中状态
		Remark:      req.Remark,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 6. 冻结用户余额
	err = s.walletService.FreezeBalance(
		ctx,
		req.UserID,
		req.TotalAmount,
		order.OrderNo,
		fmt.Sprintf("eSIM订单支付 - 订单号: %s", order.OrderNo),
	)
	if err != nil {
		// 冻结失败，需要删除订单或标记为失败
		s.orderRepo.UpdateStatus(ctx, order.ID, models.OrderStatusFailed)
		return nil, fmt.Errorf("冻结余额失败: %w", err)
	}

	// 7. TODO: 调用第三方 eSIM API 创建订单
	// 这里暂时模拟，实际应该调用 eSIM SDK
	// providerOrderID, err := s.createProviderOrder(ctx, order, product)
	// if err != nil {
	//     // 创建第三方订单失败，解冻余额
	//     s.walletService.UnfreezeBalance(ctx, req.UserID, req.TotalAmount, order.OrderNo, "订单创建失败退款")
	//     s.orderRepo.UpdateStatus(ctx, order.ID, models.OrderStatusFailed)
	//     return nil, fmt.Errorf("创建第三方订单失败: %w", err)
	// }
	//
	// // 更新订单的第三方订单ID
	// order.ProviderOrderID = providerOrderID
	// s.orderRepo.Update(ctx, order)

	// 8. 返回订单响应
	return &EsimOrderResponse{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		Status:      order.Status,
		TotalAmount: order.Amount,
		CreatedAt:   order.CreatedAt,
	}, nil
}

// GetOrderWithDetail 获取包含详情的订单信息
func (s *orderService) GetOrderWithDetail(ctx context.Context, orderID uint, userID int64) (*OrderWithDetail, error) {
	// 获取订单基本信息
	order, err := s.orderRepo.GetByIDWithDetail(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %w", err)
	}

	// 验证用户权限
	if order.UserID != userID {
		return nil, errors.New("无权访问该订单")
	}

	// 构建响应对象
	result := &OrderWithDetail{
		OrderID:         order.ID,
		OrderNo:         order.OrderNo,
		UserID:          order.UserID,
		ProductID:       order.ProductID,
		ProductName:     order.ProductName,
		Quantity:        order.Quantity,
		UnitPrice:       order.UnitPrice,
		TotalAmount:     order.Amount,
		Status:          order.Status,
		ProviderOrderID: order.ProviderOrderID,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		CompletedAt:     order.CompletedAt,
	}

	// 手动获取订单详情
	orderDetail, err := s.orderDetailRepo.GetByOrderID(ctx, order.ID)
	if err == nil && orderDetail != nil {
		// 解析 OrderItems
		if orderDetail.OrderItems != "" {
			var orderItems []OrderItemDetail
			if err := json.Unmarshal([]byte(orderDetail.OrderItems), &orderItems); err == nil {
				result.OrderItems = orderItems
			}
		}

		// 解析 Esims
		if orderDetail.Esims != "" {
			var esims []EsimDetail
			if err := json.Unmarshal([]byte(orderDetail.Esims), &esims); err == nil {
				result.Esims = esims
			}
		}
	}

	return result, nil
}

// ProcessOrderCompletion 处理订单完成（确认扣费）
func (s *orderService) ProcessOrderCompletion(ctx context.Context, orderID uint, providerOrderData *ProviderOrderData) error {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("订单不存在: %w", err)
	}

	// 检查订单状态
	if order.Status != models.OrderStatusProcessing {
		return fmt.Errorf("订单状态不正确，当前状态: %s", order.Status)
	}

	// 确认冻结金额的支付
	err = s.walletService.ConfirmFrozenPayment(
		ctx,
		order.UserID,
		order.Amount,
		order.OrderNo,
		fmt.Sprintf("eSIM订单支付完成 - 订单号: %s", order.OrderNo),
	)
	if err != nil {
		return fmt.Errorf("确认支付失败: %w", err)
	}

	// 更新订单状态为已完成
	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.CompletedAt = &now
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 保存订单详情
	if providerOrderData != nil {
		err = s.saveOrderDetail(ctx, orderID, providerOrderData)
		if err != nil {
			// 保存详情失败不应该影响主流程，只记录日志
			fmt.Printf("Warning: failed to save order detail: %v\n", err)
		}
	}

	return nil
}

// ProcessOrderFailure 处理订单失败（退还冻结金额）
func (s *orderService) ProcessOrderFailure(ctx context.Context, orderID uint, reason string) error {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("订单不存在: %w", err)
	}

	// 检查订单状态
	if order.Status != models.OrderStatusProcessing {
		return fmt.Errorf("订单状态不正确，当前状态: %s", order.Status)
	}

	// 解冻余额（退还给用户）
	err = s.walletService.UnfreezeBalance(
		ctx,
		order.UserID,
		order.Amount,
		order.OrderNo,
		fmt.Sprintf("eSIM订单失败退款 - 订单号: %s, 原因: %s", order.OrderNo, reason),
	)
	if err != nil {
		return fmt.Errorf("退还余额失败: %w", err)
	}

	// 更新订单状态为失败
	order.Status = models.OrderStatusFailed
	order.Remark = fmt.Sprintf("%s\n失败原因: %s", order.Remark, reason)
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	return nil
}

// UpdateOrderSyncInfo 更新订单同步信息
func (s *orderService) UpdateOrderSyncInfo(ctx context.Context, orderID uint, syncAttempts int, nextSyncAt *time.Time) error {
	return s.orderRepo.UpdateSyncInfo(ctx, orderID, syncAttempts, nextSyncAt)
}

// saveOrderDetail 保存订单详情（私有方法）
func (s *orderService) saveOrderDetail(ctx context.Context, orderID uint, providerOrderData *ProviderOrderData) error {
	// 序列化数据
	providerDataJSON, _ := json.Marshal(providerOrderData)
	orderItemsJSON, _ := json.Marshal(providerOrderData.OrderItems)
	esimsJSON, _ := json.Marshal(providerOrderData.Esims)

	// 创建或更新订单详情
	orderDetail := &models.OrderDetail{
		OrderID:      orderID,
		ProviderData: string(providerDataJSON),
		OrderItems:   string(orderItemsJSON),
		Esims:        string(esimsJSON),
	}

	return s.orderDetailRepo.CreateOrUpdate(ctx, orderDetail)
}

// GetUserOrdersWithFilters 根据筛选条件获取用户订单列表
func (s *orderService) GetUserOrdersWithFilters(ctx context.Context, userID int64, status models.OrderStatus, limit, offset int) ([]*models.Order, int64, error) {
	return s.orderRepo.GetByUserIDWithFilters(ctx, userID, status, limit, offset)
}
