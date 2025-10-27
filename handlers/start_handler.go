package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/services"
	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// StartHandler 处理 /start 命令
type StartHandler struct {
	bot             *tgbotapi.BotAPI
	userRepo        repository.UserRepository
	dialogService   services.DialogService
	productsHandler *ProductsHandler // 添加 ProductsHandler 引用
}

// NewStartHandler 创建 Start 命令处理器
func NewStartHandler(bot *tgbotapi.BotAPI, userRepo repository.UserRepository, dialogService services.DialogService, productsHandler *ProductsHandler) *StartHandler {
	return &StartHandler{
		bot:             bot,
		userRepo:        userRepo,
		dialogService:   dialogService,
		productsHandler: productsHandler,
	}
}

// HandleCommand 处理命令
func (h *StartHandler) HandleCommand(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID

	// 确保用户存在于数据库中
	if err := h.ensureUserExists(ctx, message.From); err != nil {
		return fmt.Errorf("failed to ensure user exists: %w", err)
	}

	// 检查是否有深度链接参数
	args := message.CommandArguments()
	if args != "" {
		return h.handleDeepLink(ctx, message.Chat.ID, args)
	}

	// 使用对话服务处理 start 命令
	response, err := h.dialogService.ProcessMessage(ctx, userID, "/start")
	if err != nil {
		return fmt.Errorf("failed to process start command: %w", err)
	}

	// 发送响应
	return h.sendResponse(message.Chat.ID, response)
}

// handleDeepLink 处理深度链接
func (h *StartHandler) handleDeepLink(ctx context.Context, chatID int64, args string) error {
	switch {
	case args == "inline_products":
		return h.handleInlineProductsDeepLink(ctx, chatID)
	case strings.HasPrefix(args, "product_detail_"):
		// 提取产品ID并显示产品详情
		productIDStr := strings.TrimPrefix(args, "product_detail_")
		return h.handleProductDetailDeepLink(ctx, chatID, productIDStr)
	case strings.HasPrefix(args, "product_buy_"):
		// 提取产品ID并开始购买流程
		productIDStr := strings.TrimPrefix(args, "product_buy_")
		return h.handleProductBuyDeepLink(ctx, chatID, productIDStr)
	default:
		// 未知参数，显示默认欢迎消息
		return h.handleInlineProductsDeepLink(ctx, chatID)
	}
}

// handleInlineProductsDeepLink 处理从 Inline Mode 切换过来的用户
func (h *StartHandler) handleInlineProductsDeepLink(ctx context.Context, chatID int64) error {
	text := "<b>🎉 欢迎使用 eSIM 机器人！</b>\n\n"
	text += "您刚才在 Inline Mode 中浏览产品，现在可以在这里进行更多操作：\n\n"
	text += "• 📱 查看完整产品列表\n"
	text += "• 🛒 购买 eSIM 产品\n"
	text += "• 💰 管理钱包和订单\n"
	text += "• 🔍 搜索特定产品\n\n"
	text += "<i>点击下方按钮开始浏览产品！</i>"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛍️ 浏览产品", "products_back"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ 帮助", "help"),
			tgbotapi.NewInlineKeyboardButtonData("📞 联系客服", "contact"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	return err
}

// handleProductDetailDeepLink 处理产品详情深度链接
func (h *StartHandler) handleProductDetailDeepLink(ctx context.Context, chatID int64, productIDStr string) error {
	if h.productsHandler == nil {
		// 产品服务未配置，显示友好提示
		return h.sendServiceUnavailableMessage(ctx, chatID, "产品详情服务暂时不可用")
	}

	// 解析产品ID
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %s", productIDStr)
	}

	// 直接调用 ProductsHandler 的方法显示产品详情
	return h.productsHandler.ShowProductDetailToUser(ctx, chatID, productID)
}

// handleProductBuyDeepLink 处理产品购买深度链接
func (h *StartHandler) handleProductBuyDeepLink(ctx context.Context, chatID int64, productIDStr string) error {
	if h.productsHandler == nil {
		// 产品服务未配置，显示友好提示
		return h.sendServiceUnavailableMessage(ctx, chatID, "购买服务暂时不可用")
	}

	// 解析产品ID
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %s", productIDStr)
	}

	// 直接调用 ProductsHandler 的方法开始购买流程
	return h.productsHandler.StartPurchaseToUser(ctx, chatID, productID)
}

// sendServiceUnavailableMessage 发送服务不可用消息
func (h *StartHandler) sendServiceUnavailableMessage(ctx context.Context, chatID int64, message string) error {
	text := fmt.Sprintf("<b>⚠️ %s</b>\n\n", message)
	text += "请稍后重试或联系客服获取帮助。\n\n"
	text += "<i>您可以使用其他功能：</i>"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ 帮助", "help"),
			tgbotapi.NewInlineKeyboardButtonData("📞 联系客服", "contact"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	return err
}

// GetCommand 获取处理的命令名称
func (h *StartHandler) GetCommand() string {
	return "start"
}

// GetDescription 获取命令描述
func (h *StartHandler) GetDescription() string {
	return "开始使用机器人"
}

// ensureUserExists 确保用户存在于数据库中
func (h *StartHandler) ensureUserExists(ctx context.Context, from *tgbotapi.User) error {
	// 检查用户是否已存在
	_, err := h.userRepo.GetByTelegramID(ctx, from.ID)
	if err == nil {
		return nil // 用户已存在
	}

	// 创建新用户
	user := &models.User{
		TelegramID: from.ID,
		Username:   from.UserName,
		FirstName:  from.FirstName,
		LastName:   from.LastName,
		Language:   from.LanguageCode,
		IsActive:   true,
	}

	return h.userRepo.Create(ctx, user)
}

// sendResponse 发送响应
func (h *StartHandler) sendResponse(chatID int64, response *services.DialogResponse) error {
	msg := tgbotapi.NewMessage(chatID, response.Message)

	if response.ParseMode != "" {
		msg.ParseMode = response.ParseMode
	}

	// 始终显示主菜单按钮
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛍️ 浏览产品", "products_back"),
			tgbotapi.NewInlineKeyboardButtonData("📦 我的订单", "my_orders"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 钱包管理", "wallet_menu"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ 设置", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ 帮助", "help"),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	return err
}
