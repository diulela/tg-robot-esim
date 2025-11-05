package services

import (
	"context"
	"fmt"
	"log"
)

// mockNotificationService 模拟通知服务实现（用于测试）
type mockNotificationService struct{}

// NewMockNotificationService 创建模拟通知服务
func NewMockNotificationService() NotificationService {
	return &mockNotificationService{}
}

// SendMessage 发送消息
func (m *mockNotificationService) SendMessage(ctx context.Context, userID int64, message string) error {
	log.Printf("[MOCK NOTIFICATION] SendMessage to user %d: %s", userID, message)
	return nil
}

// SendMenuMessage 发送菜单消息
func (m *mockNotificationService) SendMenuMessage(ctx context.Context, userID int64, response *MenuResponse) error {
	log.Printf("[MOCK NOTIFICATION] SendMenuMessage to user %d: %s", userID, response.Text)
	return nil
}

// EditMessage 编辑消息
func (m *mockNotificationService) EditMessage(ctx context.Context, userID int64, messageID int, newText string) error {
	log.Printf("[MOCK NOTIFICATION] EditMessage for user %d, message %d: %s", userID, messageID, newText)
	return nil
}

// SendTransactionNotification 发送交易通知
func (m *mockNotificationService) SendTransactionNotification(ctx context.Context, userID int64, txInfo *TransactionInfo) error {
	log.Printf("[MOCK NOTIFICATION] SendTransactionNotification to user %d: tx=%s, amount=%s",
		userID, txInfo.TxHash, txInfo.Amount)
	return nil
}

// SendRechargeSuccessNotification 发送充值成功通知
func (m *mockNotificationService) SendRechargeSuccessNotification(ctx context.Context, userID int64, amount string, orderNo string) error {
	message := fmt.Sprintf("🎉 充值成功通知\n\n💰 充值金额：%s USDT\n📋 订单号：%s\n⏰ 到账时间：现在\n\n您的钱包余额已更新，可以立即使用！",
		amount, orderNo)

	log.Printf("[MOCK NOTIFICATION] SendRechargeSuccessNotification to user %d:\n%s", userID, message)
	return nil
}
