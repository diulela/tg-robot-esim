package services

import (
	"context"
	"fmt"
	"time"

	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// SyncResult 同步结果
type SyncResult struct {
	OrderID     uint   `json:"order_id"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	NewStatus   string `json:"new_status,omitempty"`
	SyncAttempt int    `json:"sync_attempt"`
}

// SyncTaskStatus 同步任务状态
type SyncTaskStatus struct {
	OrderID      uint       `json:"order_id"`
	IsRunning    bool       `json:"is_running"`
	LastSyncAt   *time.Time `json:"last_sync_at"`
	NextSyncAt   *time.Time `json:"next_sync_at"`
	SyncAttempts int        `json:"sync_attempts"`
	LastError    string     `json:"last_error,omitempty"`
}

// OrderSyncService 定义订单同步服务接口
type OrderSyncService interface {
	// StartSyncTask 启动订单同步任务
	StartSyncTask(ctx context.Context, orderID uint) error

	// StopSyncTask 停止订单同步任务
	StopSyncTask(orderID uint) error

	// SyncOrderStatus 手动同步订单状态
	SyncOrderStatus(ctx context.Context, orderID uint) (*SyncResult, error)

	// GetSyncStatus 获取同步状态
	GetSyncStatus(orderID uint) (*SyncTaskStatus, error)

	// ProcessPendingOrders 处理所有待处理订单（定时任务）
	ProcessPendingOrders(ctx context.Context) error
}

// orderSyncService 订单同步服务实现
type orderSyncService struct {
	orderRepo       repository.OrderRepository
	orderService    OrderService
	esimClient      *esim.Client
	syncInterval    time.Duration
	maxSyncAttempts int
}

// NewOrderSyncService 创建订单同步服务实例
func NewOrderSyncService(
	orderRepo repository.OrderRepository,
	orderService OrderService,
	esimClient *esim.Client,
) OrderSyncService {
	return &orderSyncService{
		orderRepo:       orderRepo,
		orderService:    orderService,
		esimClient:      esimClient,
		syncInterval:    10 * time.Second, // 默认10秒同步间隔
		maxSyncAttempts: 100,              // 最大同步尝试次数
	}
}

// StartSyncTask 启动订单同步任务
func (s *orderSyncService) StartSyncTask(ctx context.Context, orderID uint) error {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("订单不存在: %w", err)
	}

	// 检查订单是否需要同步
	if order.Status != models.OrderStatusProcessing {
		return fmt.Errorf("订单状态不需要同步: %s", order.Status)
	}

	if order.ProviderOrderID == "" {
		return fmt.Errorf("订单缺少第三方订单ID")
	}

	// 设置下次同步时间
	nextSyncAt := time.Now().Add(s.syncInterval)
	return s.orderRepo.UpdateSyncInfo(ctx, orderID, order.SyncAttempts, &nextSyncAt)
}

// StopSyncTask 停止订单同步任务
func (s *orderSyncService) StopSyncTask(orderID uint) error {
	// 简单实现：将下次同步时间设为 nil
	ctx := context.Background()
	return s.orderRepo.UpdateSyncInfo(ctx, orderID, 0, nil)
}

// SyncOrderStatus 手动同步订单状态
func (s *orderSyncService) SyncOrderStatus(ctx context.Context, orderID uint) (*SyncResult, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %w", err)
	}

	result := &SyncResult{
		OrderID:     orderID,
		SyncAttempt: order.SyncAttempts + 1,
	}

	// 检查订单状态
	if order.Status != models.OrderStatusProcessing {
		result.Success = false
		result.Message = fmt.Sprintf("订单状态不需要同步: %s", string(order.Status))
		return result, nil
	}

	if order.ProviderOrderID == "" {
		result.Success = false
		result.Message = "订单缺少第三方订单ID"
		return result, nil
	}

	// 调用第三方 API 查询订单状态
	providerOrderIDInt := 0 // 需要将字符串转换为整数
	fmt.Sscanf(order.ProviderOrderID, "%d", &providerOrderIDInt)

	providerOrder, err := s.esimClient.GetOrder(providerOrderIDInt)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("查询第三方订单失败: %v", err)

		// 更新同步尝试次数
		nextSyncAt := time.Now().Add(s.syncInterval)
		s.orderRepo.UpdateSyncInfo(ctx, orderID, result.SyncAttempt, &nextSyncAt)

		return result, nil
	}

	// 检查是否成功解析订单详情
	if providerOrder.OrderDetail == nil {
		result.Success = false
		result.Message = "第三方订单数据解析失败"
		return result, nil
	}

	// 处理第三方订单状态
	switch providerOrder.OrderDetail.Status {
	case esim.OrderStatusCompleted:
		// 订单完成
		providerOrderData := &ProviderOrderData{
			OrderID:     providerOrder.OrderDetail.ID,
			OrderNumber: providerOrder.OrderDetail.OrderNumber,
			Status:      string(providerOrder.OrderDetail.Status),
			OrderItems:  convertOrderItems(providerOrder.OrderDetail.OrderItems),
			Esims:       convertEsims(providerOrder.OrderDetail.Esims),
		}

		err = s.orderService.ProcessOrderCompletion(ctx, orderID, providerOrderData)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("处理订单完成失败: %v", err)
		} else {
			result.Success = true
			result.Message = "订单处理完成"
			result.NewStatus = string(models.OrderStatusCompleted)
		}

	case esim.OrderStatusCancelled, esim.OrderStatusFailed:
		// 订单失败
		reason := fmt.Sprintf("第三方订单状态: %s", providerOrder.OrderDetail.Status)
		err = s.orderService.ProcessOrderFailure(ctx, orderID, reason)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("处理订单失败失败: %v", err)
		} else {
			result.Success = true
			result.Message = "订单已标记为失败"
			result.NewStatus = string(models.OrderStatusFailed)
		}

	case esim.OrderStatusPending, esim.OrderStatusPaid, esim.OrderStatusProcessing:
		// 订单仍在处理中，继续等待
		result.Success = true
		result.Message = fmt.Sprintf("订单仍在处理中，第三方状态: %s", providerOrder.OrderDetail.Status)

		// 设置下次同步时间
		nextSyncAt := time.Now().Add(s.syncInterval)
		s.orderRepo.UpdateSyncInfo(ctx, orderID, result.SyncAttempt, &nextSyncAt)

	default:
		result.Success = false
		result.Message = fmt.Sprintf("未知的第三方订单状态: %s", providerOrder.OrderDetail.Status)
	}

	return result, nil
}

// GetSyncStatus 获取同步状态
func (s *orderSyncService) GetSyncStatus(orderID uint) (*SyncTaskStatus, error) {
	ctx := context.Background()
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %w", err)
	}

	status := &SyncTaskStatus{
		OrderID:      orderID,
		IsRunning:    order.Status == models.OrderStatusProcessing && order.NextSyncAt != nil,
		LastSyncAt:   order.LastSyncAt,
		NextSyncAt:   order.NextSyncAt,
		SyncAttempts: order.SyncAttempts,
	}

	return status, nil
}

// ProcessPendingOrders 处理所有待处理订单（定时任务）
func (s *orderSyncService) ProcessPendingOrders(ctx context.Context) error {
	// 获取需要同步的订单
	orders, err := s.orderRepo.GetPendingSyncOrders(ctx, 50) // 每次处理50个订单
	if err != nil {
		return fmt.Errorf("获取待同步订单失败: %w", err)
	}

	fmt.Printf("开始处理 %d 个待同步订单\n", len(orders))

	for _, order := range orders {
		// 检查是否超过最大尝试次数
		if order.SyncAttempts >= s.maxSyncAttempts {
			// 超过最大尝试次数，标记为失败
			s.orderService.ProcessOrderFailure(ctx, order.ID, "同步超时，超过最大尝试次数")
			continue
		}

		// 同步订单状态
		result, err := s.SyncOrderStatus(ctx, order.ID)
		if err != nil {
			fmt.Printf("同步订单 %d 失败: %v\n", order.ID, err)
			continue
		}

		if result.Success {
			fmt.Printf("订单 %d 同步成功: %s\n", order.ID, result.Message)
		} else {
			fmt.Printf("订单 %d 同步失败: %s\n", order.ID, result.Message)
		}
	}

	return nil
}

// convertOrderItems 转换订单项数据格式
func convertOrderItems(items []esim.OrderItem) []OrderItemDetail {
	var result []OrderItemDetail
	for _, item := range items {
		result = append(result, OrderItemDetail{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
			DataSize:    item.DataSize,
			ValidDays:   item.ValidDays,
		})
	}
	return result
}

// convertEsims 转换 eSIM 数据格式
func convertEsims(esims []esim.OrderEsim) []EsimDetail {
	var result []EsimDetail
	for _, esim := range esims {
		result = append(result, EsimDetail{
			ID:                esim.ID,
			ICCID:             esim.ICCID,
			Status:            esim.Status,
			HasActivationCode: esim.HasActivationCode,
			HasQrCode:         esim.HasQrCode,
		})
	}
	return result
}
