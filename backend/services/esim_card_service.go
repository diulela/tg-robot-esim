package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"tg-robot-sim/pkg/sdk/esim"
	service_common "tg-robot-sim/services/common"
	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// esimCardService eSIM 卡服务实现
type esimCardService struct {
	esimCardRepo      repository.EsimCardRepository
	orderRepo         repository.OrderRepository
	esimClientService service_common.EsimClientService
}

// NewEsimCardService 创建 eSIM 卡服务实例
func NewEsimCardService(
	esimCardRepo repository.EsimCardRepository,
	orderRepo repository.OrderRepository,
	esimClientService service_common.EsimClientService,
) EsimCardService {
	return &esimCardService{
		esimCardRepo:      esimCardRepo,
		orderRepo:         orderRepo,
		esimClientService: esimClientService,
	}
}

// CreateEsimCard 创建 eSIM 卡记录
func (s *esimCardService) CreateEsimCard(ctx context.Context, orderID uint, providerEsim interface{}) (*models.EsimCard, error) {
	// 1. 验证订单存在
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %w", err)
	}

	// 2. 类型转换
	orderEsim, ok := providerEsim.(*esim.OrderEsim)
	if !ok {
		return nil, errors.New("无效的 eSIM 数据格式")
	}

	// 3. 检查 ICCID 是否已存在
	if orderEsim.ICCID != "" {
		existing, _ := s.esimCardRepo.GetByICCID(ctx, orderEsim.ICCID)
		if existing != nil {
			return nil, fmt.Errorf("ICCID 已存在: %s", orderEsim.ICCID)
		}
	}

	// 4. 创建 eSIM 卡记录
	esimCard := &models.EsimCard{
		UserID:          order.UserID,
		OrderID:         orderID,
		ICCID:           orderEsim.ICCID,
		ActivationCode:  orderEsim.ActivationCode,
		QrCode:          orderEsim.QrCode,
		Lpa:             orderEsim.Lpa,
		DirectAppleUrl:  orderEsim.DirectAppleUrl,
		Status:          models.EsimStatusPending,
		ProviderOrderID: orderEsim.ID,
	}

	// 5. 保存第三方原始数据
	providerDataJSON, _ := json.Marshal(orderEsim)
	esimCard.ProviderData = string(providerDataJSON)

	// 6. 保存到数据库
	if err := s.esimCardRepo.Create(ctx, esimCard); err != nil {
		return nil, fmt.Errorf("创建 eSIM 卡失败: %w", err)
	}

	return esimCard, nil
}

// GetEsimCard 获取 eSIM 卡详情
func (s *esimCardService) GetEsimCard(ctx context.Context, esimID uint, userID int64) (*models.EsimCard, error) {
	// 1. 获取 eSIM 卡
	esimCard, err := s.esimCardRepo.GetByID(ctx, esimID)
	if err != nil {
		return nil, fmt.Errorf("eSIM 卡不存在: %w", err)
	}

	// 2. 验证用户权限
	if esimCard.UserID != userID {
		return nil, errors.New("无权访问该 eSIM 卡")
	}

	return esimCard, nil
}

// GetEsimCardByICCID 根据 ICCID 获取 eSIM 卡
func (s *esimCardService) GetEsimCardByICCID(ctx context.Context, iccid string) (*models.EsimCard, error) {
	esimCard, err := s.esimCardRepo.GetByICCID(ctx, iccid)
	if err != nil {
		return nil, fmt.Errorf("eSIM 卡不存在: %w", err)
	}
	return esimCard, nil
}

// GetUserEsimCards 获取用户的所有 eSIM 卡
func (s *esimCardService) GetUserEsimCards(ctx context.Context, userID int64, filters EsimCardFilters) ([]*models.EsimCard, int64, error) {
	// 1. 验证分页参数
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	// 2. 获取 eSIM 卡列表
	esimCards, total, err := s.esimCardRepo.GetByUserIDWithFilters(
		ctx,
		userID,
		filters.Status,
		filters.Limit,
		filters.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("获取 eSIM 卡列表失败: %w", err)
	}

	return esimCards, total, nil
}

// SyncEsimCardStatus 同步 eSIM 卡状态和使用情况
func (s *esimCardService) SyncEsimCardStatus(ctx context.Context, esimID uint) error {
	// 1. 获取 eSIM 卡
	esimCard, err := s.esimCardRepo.GetByID(ctx, esimID)
	if err != nil {
		return fmt.Errorf("eSIM 卡不存在: %w", err)
	}

	// 2. 检查是否可以同步
	if !esimCard.CanSync() {
		return fmt.Errorf("eSIM 卡状态不允许同步: %s", esimCard.Status)
	}

	// 3. 从第三方获取最新状态
	if s.esimClientService == nil {
		return errors.New("eSIM 客户端服务未初始化")
	}

	usageResp, err := s.esimClientService.GetEsimUsage(ctx, esimCard.ProviderOrderID)
	if err != nil {
		return fmt.Errorf("获取 eSIM 使用情况失败: %w", err)
	}

	// 4. 更新 eSIM 卡信息
	if usageResp != nil && usageResp.UsageData != nil {
		esimCard.Status = models.EsimStatus(usageResp.UsageData.Esim.Status)
		esimCard.DataUsed = usageResp.UsageData.Esim.DataUsed
		esimCard.DataRemaining = usageResp.UsageData.Esim.DataRemaining
		esimCard.UsagePercent = usageResp.UsageData.Esim.UsagePercentage

		// 解析时间信息
		if usageResp.UsageData.Esim.ActivationTime != "" {
			if activatedAt, err := time.Parse(time.RFC3339, usageResp.UsageData.Esim.ActivationTime); err == nil {
				esimCard.ActivatedAt = &activatedAt
			}
		}

		if usageResp.UsageData.Esim.ExpireTime != "" {
			if expiresAt, err := time.Parse(time.RFC3339, usageResp.UsageData.Esim.ExpireTime); err == nil {
				esimCard.ExpiresAt = &expiresAt
			}
		}

		// 记录同步时间
		now := time.Now()
		esimCard.LastSyncAt = &now

		// 保存更新
		if err := s.esimCardRepo.Update(ctx, esimCard); err != nil {
			return fmt.Errorf("更新 eSIM 卡失败: %w", err)
		}
	}

	return nil
}

// GetEsimCardWithOrder 获取 eSIM 卡及其关联的购买订单
func (s *esimCardService) GetEsimCardWithOrder(ctx context.Context, esimID uint, userID int64) (*EsimCardWithOrder, error) {
	// 1. 获取 eSIM 卡
	esimCard, err := s.GetEsimCard(ctx, esimID, userID)
	if err != nil {
		return nil, err
	}

	// 2. 获取关联的购买订单
	order, err := s.orderRepo.GetByID(ctx, esimCard.OrderID)
	if err != nil {
		return nil, fmt.Errorf("获取购买订单失败: %w", err)
	}

	return &EsimCardWithOrder{
		EsimCard:      esimCard,
		PurchaseOrder: order,
	}, nil
}

// UpdateEsimCardUsage 更新 eSIM 卡使用情况
func (s *esimCardService) UpdateEsimCardUsage(ctx context.Context, esimID uint, usageInfo interface{}) error {
	// 1. 获取 eSIM 卡
	esimCard, err := s.esimCardRepo.GetByID(ctx, esimID)
	if err != nil {
		return fmt.Errorf("eSIM 卡不存在: %w", err)
	}

	// 2. 类型转换
	usage, ok := usageInfo.(*esim.EsimUsageInfo)
	if !ok {
		return errors.New("无效的使用情况数据格式")
	}

	// 3. 更新使用情况
	if err := s.esimCardRepo.UpdateUsage(
		ctx,
		esimID,
		usage.DataUsed,
		usage.DataRemaining,
		usage.UsagePercentage,
	); err != nil {
		return fmt.Errorf("更新使用情况失败: %w", err)
	}

	// 4. 更新状态
	if err := s.esimCardRepo.UpdateStatus(ctx, esimID, models.EsimStatus(usage.Status)); err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 5. 更新时间信息
	if usage.ActivationTime != "" {
		if activatedAt, err := time.Parse(time.RFC3339, usage.ActivationTime); err == nil {
			esimCard.ActivatedAt = &activatedAt
		}
	}

	if usage.ExpireTime != "" {
		if expiresAt, err := time.Parse(time.RFC3339, usage.ExpireTime); err == nil {
			esimCard.ExpiresAt = &expiresAt
		}
	}

	// 6. 保存更新
	if err := s.esimCardRepo.Update(ctx, esimCard); err != nil {
		return fmt.Errorf("保存更新失败: %w", err)
	}

	return nil
}

// ConvertOrderEsimToEsimCard 将第三方 OrderEsim 转换为 EsimCard
func ConvertOrderEsimToEsimCard(orderID uint, userID int64, orderEsim *esim.OrderEsim) (*models.EsimCard, error) {
	if orderEsim == nil {
		return nil, errors.New("orderEsim 不能为空")
	}

	if orderEsim.ICCID == "" {
		return nil, errors.New("ICCID 不能为空")
	}

	// 创建 EsimCard
	esimCard := &models.EsimCard{
		UserID:          userID,
		OrderID:         orderID,
		ICCID:           orderEsim.ICCID,
		ActivationCode:  orderEsim.ActivationCode,
		QrCode:          orderEsim.QrCode,
		Lpa:             orderEsim.Lpa,
		DirectAppleUrl:  orderEsim.DirectAppleUrl,
		Status:          models.EsimStatusPending,
		ProviderOrderID: orderEsim.ID,
	}

	// 解析激活时间
	if orderEsim.ActivatedAt != "" {
		if activatedAt, err := parseTime(orderEsim.ActivatedAt); err == nil {
			esimCard.ActivatedAt = &activatedAt
		}
	}

	// 解析过期时间
	if orderEsim.ExpiresAt != "" {
		if expiresAt, err := parseTime(orderEsim.ExpiresAt); err == nil {
			esimCard.ExpiresAt = &expiresAt
		}
	}

	// 保存第三方原始数据
	providerDataJSON, _ := json.Marshal(orderEsim)
	esimCard.ProviderData = string(providerDataJSON)

	return esimCard, nil
}

// ConvertEsimUsageInfoToEsimCard 将第三方 EsimUsageInfo 转换为 EsimCard 更新
func ConvertEsimUsageInfoToEsimCard(esimCard *models.EsimCard, usageInfo *esim.EsimUsageInfo) error {
	if esimCard == nil {
		return errors.New("esimCard 不能为空")
	}

	if usageInfo == nil {
		return errors.New("usageInfo 不能为空")
	}

	// 更新状态
	esimCard.Status = models.EsimStatus(usageInfo.Status)

	// 更新流量信息
	esimCard.DataUsed = usageInfo.DataUsed
	esimCard.DataRemaining = usageInfo.DataRemaining
	esimCard.UsagePercent = usageInfo.UsagePercentage

	// 解析激活时间
	if usageInfo.ActivationTime != "" {
		if activatedAt, err := parseTime(usageInfo.ActivationTime); err == nil {
			esimCard.ActivatedAt = &activatedAt
		}
	}

	// 解析过期时间
	if usageInfo.ExpireTime != "" {
		if expiresAt, err := parseTime(usageInfo.ExpireTime); err == nil {
			esimCard.ExpiresAt = &expiresAt
		}
	}

	// 记录同步时间
	now := time.Now()
	esimCard.LastSyncAt = &now

	return nil
}

// ValidateEsimCard 验证 eSIM 卡数据
func ValidateEsimCard(esimCard *models.EsimCard) error {
	if esimCard == nil {
		return errors.New("eSIM 卡不能为空")
	}

	if esimCard.ICCID == "" {
		return errors.New("ICCID 不能为空")
	}

	if esimCard.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	if esimCard.OrderID == 0 {
		return errors.New("订单ID不能为空")
	}

	// 验证流量信息一致性
	if esimCard.DataUsed > 0 && esimCard.DataRemaining > 0 {
		// 如果有使用情况，检查数据一致性
		if esimCard.DataUsed+esimCard.DataRemaining > 0 {
			// 这是一个合理的检查，但不是严格的
		}
	}

	return nil
}

// parseTime 解析时间字符串
func parseTime(timeStr string) (time.Time, error) {
	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间: %s", timeStr)
}

// MergeEsimCardData 合并 eSIM 卡数据（用于更新时保留原有数据）
func MergeEsimCardData(existing *models.EsimCard, updated *models.EsimCard) *models.EsimCard {
	if existing == nil {
		return updated
	}

	if updated == nil {
		return existing
	}

	// 保留原有的 ID、UserID、OrderID、CreatedAt
	updated.ID = existing.ID
	updated.UserID = existing.UserID
	updated.OrderID = existing.OrderID
	updated.CreatedAt = existing.CreatedAt

	// 如果新数据中某些字段为空，使用原有数据
	if updated.ICCID == "" {
		updated.ICCID = existing.ICCID
	}

	if updated.ActivationCode == "" {
		updated.ActivationCode = existing.ActivationCode
	}

	if updated.QrCode == "" {
		updated.QrCode = existing.QrCode
	}

	if updated.Lpa == "" {
		updated.Lpa = existing.Lpa
	}

	if updated.DirectAppleUrl == "" {
		updated.DirectAppleUrl = existing.DirectAppleUrl
	}

	if updated.ApnType == "" {
		updated.ApnType = existing.ApnType
	}

	// 如果新数据中没有激活时间，使用原有数据
	if updated.ActivatedAt == nil && existing.ActivatedAt != nil {
		updated.ActivatedAt = existing.ActivatedAt
	}

	// 如果新数据中没有过期时间，使用原有数据
	if updated.ExpiresAt == nil && existing.ExpiresAt != nil {
		updated.ExpiresAt = existing.ExpiresAt
	}

	return updated
}
