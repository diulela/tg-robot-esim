package handlers

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/services"
)

// HelpHandler 处理 /help 命令
type HelpHandler struct {
	bot           *tgbotapi.BotAPI
	dialogService services.DialogService
}

// NewHelpHandler 创建 Help 命令处理器
func NewHelpHandler(bot *tgbotapi.BotAPI, dialogService services.DialogService) *HelpHandler {
	return &HelpHandler{
		bot:           bot,
		dialogService: dialogService,
	}
}

// HandleCommand 处理命令
func (h *HelpHandler) HandleCommand(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID

	// 使用对话服务处理 help 命令
	response, err := h.dialogService.ProcessMessage(ctx, userID, "/help")
	if err != nil {
		return fmt.Errorf("failed to process help command: %w", err)
	}

	// 发送响应
	msg := tgbotapi.NewMessage(message.Chat.ID, response.Message)
	if response.ParseMode != "" {
		msg.ParseMode = response.ParseMode
	}

	_, err = h.bot.Send(msg)
	return err
}

// GetCommand 获取处理的命令名称
func (h *HelpHandler) GetCommand() string {
	return "help"
}

// GetDescription 获取命令描述
func (h *HelpHandler) GetDescription() string {
	return "显示帮助信息"
}
