package services

import (
	"context"
	"fmt"
	"strings"
	"time"

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

// SendRechargeSuccessNotification 发送充值成功通知
func (n *notificationService) SendRechargeSuccessNotification(ctx context.Context, userID int64, amount string, orderNo string) error {
	message := fmt.Sprintf(
		"🎉 <b>充值成功通知</b>\n\n"+
			"💰 <b>充值金额:</b> %s USDT\n"+
			"📋 <b>订单号:</b> <code>%s</code>\n"+
			"⏰ <b>到账时间:</b> %s\n\n"+
			"您的钱包余额已更新，可以立即使用！",
		amount,
		orderNo,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	// 创建内联键盘
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 查看钱包", "wallet:balance"),
			tgbotapi.NewInlineKeyboardButtonData("📋 充值历史", "wallet:history"),
		),
	)

	msg := tgbotapi.NewMessage(userID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = keyboard

	// 发送消息（带重试机制）
	err := n.sendMessageWithRetry(ctx, msg, 2)
	if err != nil {
		n.logger.Error("发送充值成功通知失败: user_id=%d, error=%v", userID, err)
		return err
	}

	n.logger.Info("充值成功通知已发送: user_id=%d, order_no=%s", userID, orderNo)
	return nil
}

// sendMessageWithRetry 带重试机制的消息发送
func (n *notificationService) sendMessageWithRetry(ctx context.Context, msg tgbotapi.MessageConfig, maxRetries int) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		_, err := n.bot.Send(msg)
		if err == nil {
			return nil
		}

		lastErr = err

		// 如果是用户屏蔽 Bot 的错误，不重试
		if isUserBlockedError(err) {
			n.logger.Warn("用户已屏蔽 Bot: user_id=%d", msg.ChatID)
			return nil // 静默处理，不返回错误
		}

		// 如果不是最后一次重试，等待后重试
		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * time.Second
			n.logger.Warn("发送消息失败，%v 后重试 (第 %d/%d 次): %v", waitTime, i+1, maxRetries, err)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
				// 继续重试
			}
		}
	}

	return fmt.Errorf("重试 %d 次后仍然失败: %w", maxRetries, lastErr)
}

// isUserBlockedError 检查是否是用户屏蔽 Bot 的错误
func isUserBlockedError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// Telegram API 返回的用户屏蔽错误信息
	return strings.Contains(errStr, "blocked by the user") ||
		strings.Contains(errStr, "user is deactivated") ||
		strings.Contains(errStr, "bot was blocked by the user")
}
