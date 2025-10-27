package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/storage/repository"
)

// ProductSelectionHandler 产品选择消息处理器
type ProductSelectionHandler struct {
	bot         *tgbotapi.BotAPI
	productRepo repository.ProductRepository
	logger      Logger
}

// NewProductSelectionHandler 创建产品选择消息处理器
func NewProductSelectionHandler(bot *tgbotapi.BotAPI, productRepo repository.ProductRepository, logger Logger) *ProductSelectionHandler {
	return &ProductSelectionHandler{
		bot:         bot,
		productRepo: productRepo,
		logger:      logger,
	}
}

// HandleMessage 处理消息
func (h *ProductSelectionHandler) HandleMessage(ctx context.Context, message *tgbotapi.Message) error {
	// 解析用户输入的数字
	text := strings.TrimSpace(message.Text)
	productIndex, err := strconv.Atoi(text)
	if err != nil || productIndex < 1 {
		return h.sendError(message.Chat.ID, "请输入有效的产品编号（1-5）")
	}

	// 获取当前页的产品列表（默认第1页）
	// TODO: 这里应该从会话中获取当前页码，暂时使用第1页
	page := 1
	products, _, err := h.getAsiaProducts(ctx, page, 5)
	if err != nil {
		h.logger.Error("Failed to get products: %v", err)
		return h.sendError(message.Chat.ID, "获取产品列表失败")
	}

	// 检查产品编号是否有效
	if productIndex > len(products) {
		return h.sendError(message.Chat.ID, "产品编号无效，请输入1-%d之间的数字", len(products))
	}

	// 获取对应的产品
	product := products[productIndex-1]

	// 显示产品详情
	return h.showProductDetail(ctx, message, product.ID)
}

// CanHandle 判断是否能处理该消息
func (h *ProductSelectionHandler) CanHandle(message *tgbotapi.Message) bool {
	// 只处理纯数字消息（1-5）
	text := strings.TrimSpace(message.Text)
	if num, err := strconv.Atoi(text); err == nil && num >= 1 && num <= 5 {
		return true
	}
	return false
}

// GetHandlerName 获取处理器名称
func (h *ProductSelectionHandler) GetHandlerName() string {
	return "product_selection"
}

// getAsiaProducts 获取亚洲产品列表
func (h *ProductSelectionHandler) getAsiaProducts(ctx context.Context, page, limit int) ([]*repository.ProductModel, int64, error) {
	params := repository.ListParams{
		Type:      "regional",
		Status:    "active",
		NameLike:  "亚洲",
		Page:      page,
		Limit:     limit,
		OrderBy:   "sort_order",
		OrderDesc: false,
	}
	return h.productRepo.List(ctx, params)
}

// showProductDetail 显示产品详情
func (h *ProductSelectionHandler) showProductDetail(ctx context.Context, message *tgbotapi.Message, productID int) error {
	// 这里复用 ProductsHandler 的逻辑
	// 为了简化，直接发送一个简单的消息
	text := "正在加载产品详情..."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 查看详情", "product_detail:"+strconv.Itoa(productID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回产品列表", "products_back"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	return err
}

// sendError 发送错误消息
func (h *ProductSelectionHandler) sendError(chatID int64, format string, args ...interface{}) error {
	text := "❌ "
	if len(args) > 0 {
		text += fmt.Sprintf(format, args...)
	} else {
		text += format
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.bot.Send(msg)
	return err
}
