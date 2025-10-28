package services

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// notificationService 通知服务实现
type notificationService struct {
	bot    *tgbotapi.BotAPI
	logger Logger
}

// NewNotificationService 创建通知服务
func NewNotificationService(bot *tgbotapi.BotAPI, logger Logger) NotificationService {
	return &notificationService{
		bot:    bot,
		logger: logger,
	}
}

// SendMessage 发送消息
func (n *notificationService) SendMessage(ctx context.Context, userID int64, message string) error {
	msg := tgbotapi.NewMessage(userID, message)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := n.bot.Send(msg)
	if err != nil {
		n.logger.Error("Failed to send message to user %d: %v", userID, err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	n.logger.Debug("Message sent to user %d", userID)
	return nil
}

// SendMenuMessage 发送菜单消息
func (n *notificationService) SendMenuMessage(ctx context.Context, userID int64, response *MenuResponse) error {
	msg := tgbotapi.NewMessage(userID, response.Text)

	if response.ParseMode != "" {
		msg.ParseMode = response.ParseMode
	}

	if response.Keyboard != nil {
		msg.ReplyMarkup = response.Keyboard
	}

	_, err := n.bot.Send(msg)
	if err != nil {
		n.logger.Error("Failed to send menu message to user %d: %v", userID, err)
		return fmt.Errorf("failed to send menu message: %w", err)
	}

	n.logger.Debug("Menu message sent to user %d", userID)
	return nil
}

// EditMessage 编辑消息
func (n *notificationService) EditMessage(ctx context.Context, userID int64, messageID int, newText string) error {
	editMsg := tgbotapi.NewEditMessageText(userID, messageID, newText)
	editMsg.ParseMode = tgbotapi.ModeHTML

	_, err := n.bot.Send(editMsg)
	if err != nil {
		n.logger.Error("Failed to edit message for user %d: %v", userID, err)
		return fmt.Errorf("failed to edit message: %w", err)
	}

	n.logger.Debug("Message edited for user %d", userID)
	return nil
}

// SendTransactionNotification 发送交易通知
func (n *notificationService) SendTransactionNotification(ctx context.Context, userID int64, txInfo *TransactionInfo) error {
	var statusIcon string
	var statusText string

	switch txInfo.Status {
	case string(TransactionStatusConfirmed):
		statusIcon = "✅"
		statusText = "已确认"
	case string(TransactionStatusFailed):
		statusIcon = "❌"
		statusText = "失败"
	default:
		statusIcon = "⏳"
		statusText = "待确认"
	}

	message := fmt.Sprintf(
		"%s <b>交易通知</b>\n\n"+
			"<b>状态:</b> %s %s\n"+
			"<b>交易哈希:</b> <code>%s</code>\n"+
			"<b>金额:</b> %s USDT\n"+
			"<b>从:</b> <code>%s</code>\n"+
			"<b>到:</b> <code>%s</code>\n"+
			"<b>确认数:</b> %d\n"+
			"<b>时间:</b> %s",
		statusIcon,
		statusIcon, statusText,
		txInfo.TxHash,
		txInfo.Amount,
		txInfo.FromAddress,
		txInfo.ToAddress,
		txInfo.Confirmations,
		txInfo.Timestamp.Format("2006-01-02 15:04:05"),
	)

	return n.SendMessage(ctx, userID, message)
}
