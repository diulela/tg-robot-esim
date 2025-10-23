package handlers

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/services"
)

// CallbackQueryHandler 回调查询处理器
type CallbackQueryHandler struct {
	bot         *tgbotapi.BotAPI
	menuService services.MenuService
	logger      Logger
}

// NewCallbackQueryHandler 创建回调查询处理器
func NewCallbackQueryHandler(bot *tgbotapi.BotAPI, menuService services.MenuService, logger Logger) *CallbackQueryHandler {
	return &CallbackQueryHandler{
		bot:         bot,
		menuService: menuService,
		logger:      logger,
	}
}

// HandleCallback 处理回调查询
func (h *CallbackQueryHandler) HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	userID := callback.From.ID
	data := callback.Data

	h.logger.Debug("Handling callback: %s from user %d", data, userID)

	// 首先回答回调查询
	if err := h.answerCallback(callback.ID, ""); err != nil {
		h.logger.Error("Failed to answer callback: %v", err)
	}

	// 解析回调数据
	action, params := h.parseCallbackData(data)

	// 处理菜单操作
	response, err := h.menuService.HandleMenuAction(ctx, userID, action)
	if err != nil {
		h.logger.Error("Failed to handle menu action: %v", err)
		return h.sendErrorMessage(callback.Message.Chat.ID, "处理请求时发生错误")
	}

	// 发送或编辑消息
	return h.sendMenuResponse(callback.Message, response, params)
}

// CanHandle 判断是否能处理该回调
func (h *CallbackQueryHandler) CanHandle(callback *tgbotapi.CallbackQuery) bool {
	// 处理所有回调查询
	return callback.Data != ""
}

// GetHandlerName 获取处理器名称
func (h *CallbackQueryHandler) GetHandlerName() string {
	return "callback_query"
}

// parseCallbackData 解析回调数据
func (h *CallbackQueryHandler) parseCallbackData(data string) (action string, params map[string]string) {
	params = make(map[string]string)

	// 简单的格式：action:param1=value1:param2=value2
	parts := strings.Split(data, ":")
	if len(parts) > 0 {
		action = parts[0]
	}

	// 解析参数
	for i := 1; i < len(parts); i++ {
		if paramParts := strings.Split(parts[i], "="); len(paramParts) == 2 {
			params[paramParts[0]] = paramParts[1]
		}
	}

	return action, params
}

// sendMenuResponse 发送菜单响应
func (h *CallbackQueryHandler) sendMenuResponse(originalMessage *tgbotapi.Message, response *services.MenuResponse, params map[string]string) error {
	if response.EditMode && originalMessage != nil {
		// 编辑现有消息
		return h.editMessage(originalMessage, response)
	} else {
		// 发送新消息
		return h.sendNewMessage(originalMessage.Chat.ID, response)
	}
}

// editMessage 编辑消息
func (h *CallbackQueryHandler) editMessage(originalMessage *tgbotapi.Message, response *services.MenuResponse) error {
	editMsg := tgbotapi.NewEditMessageText(
		originalMessage.Chat.ID,
		originalMessage.MessageID,
		response.Text,
	)

	if response.ParseMode != "" {
		editMsg.ParseMode = response.ParseMode
	}

	if response.Keyboard != nil {
		if keyboard, ok := response.Keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
			editMsg.ReplyMarkup = &keyboard
		}
	}

	_, err := h.bot.Send(editMsg)
	return err
}

// sendNewMessage 发送新消息
func (h *CallbackQueryHandler) sendNewMessage(chatID int64, response *services.MenuResponse) error {
	msg := tgbotapi.NewMessage(chatID, response.Text)

	if response.ParseMode != "" {
		msg.ParseMode = response.ParseMode
	}

	if response.Keyboard != nil {
		msg.ReplyMarkup = response.Keyboard
	}

	_, err := h.bot.Send(msg)
	return err
}

// answerCallback 回答回调查询
func (h *CallbackQueryHandler) answerCallback(callbackQueryID, text string) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	_, err := h.bot.Request(callback)
	return err
}

// sendErrorMessage 发送错误消息
func (h *CallbackQueryHandler) sendErrorMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, "❌ "+message)
	_, err := h.bot.Send(msg)
	return err
}
