package handlers

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/services"
)

// GeneralMessageHandler 通用消息处理器
type GeneralMessageHandler struct {
	bot           *tgbotapi.BotAPI
	dialogService services.DialogService
}

// NewGeneralMessageHandler 创建通用消息处理器
func NewGeneralMessageHandler(bot *tgbotapi.BotAPI, dialogService services.DialogService) *GeneralMessageHandler {
	return &GeneralMessageHandler{
		bot:           bot,
		dialogService: dialogService,
	}
}

// HandleMessage 处理消息
func (h *GeneralMessageHandler) HandleMessage(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID

	// 使用对话服务处理消息
	response, err := h.dialogService.ProcessMessage(ctx, userID, message.Text)
	if err != nil {
		return fmt.Errorf("failed to process message: %w", err)
	}

	// 发送响应
	return h.sendResponse(message.Chat.ID, response)
}

// CanHandle 判断是否能处理该消息
func (h *GeneralMessageHandler) CanHandle(message *tgbotapi.Message) bool {
	// 处理所有非命令的文本消息
	return message.Text != "" && !message.IsCommand()
}

// GetHandlerName 获取处理器名称
func (h *GeneralMessageHandler) GetHandlerName() string {
	return "general_message"
}

// sendResponse 发送响应
func (h *GeneralMessageHandler) sendResponse(chatID int64, response *services.DialogResponse) error {
	msg := tgbotapi.NewMessage(chatID, response.Message)

	if response.ParseMode != "" {
		msg.ParseMode = response.ParseMode
	}

	// 检查是否有键盘
	if response.Keyboard != nil {
		msg.ReplyMarkup = response.Keyboard
	}

	_, err := h.bot.Send(msg)
	return err
}
