package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/services"
)

// MenuHandler 处理 /menu 命令
type MenuHandler struct {
	bot         *tgbotapi.BotAPI
	menuService services.MenuService
}

// NewMenuHandler 创建菜单命令处理器
func NewMenuHandler(bot *tgbotapi.BotAPI, menuService services.MenuService) *MenuHandler {
	return &MenuHandler{
		bot:         bot,
		menuService: menuService,
	}
}

// HandleCommand 处理命令
func (h *MenuHandler) HandleCommand(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID

	// 获取主菜单
	menuResponse, err := h.menuService.GetMainMenu(userID)
	if err != nil {
		return h.sendError(message.Chat.ID, "获取菜单失败")
	}

	// 发送菜单
	msg := tgbotapi.NewMessage(message.Chat.ID, menuResponse.Text)
	if menuResponse.ParseMode != "" {
		msg.ParseMode = menuResponse.ParseMode
	}
	if menuResponse.Keyboard != nil {
		msg.ReplyMarkup = menuResponse.Keyboard
	}

	_, err = h.bot.Send(msg)
	return err
}

// GetCommand 获取处理的命令名称
func (h *MenuHandler) GetCommand() string {
	return "menu"
}

// GetDescription 获取命令描述
func (h *MenuHandler) GetDescription() string {
	return "显示主菜单"
}

// sendError 发送错误消息
func (h *MenuHandler) sendError(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, "❌ "+message)
	_, err := h.bot.Send(msg)
	return err
}
