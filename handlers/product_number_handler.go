package handlers

import (
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ProductNumberHandler 产品编号消息处理器
type ProductNumberHandler struct {
	productsHandler *ProductsHandler
}

// NewProductNumberHandler 创建产品编号消息处理器
func NewProductNumberHandler(productsHandler *ProductsHandler) *ProductNumberHandler {
	return &ProductNumberHandler{
		productsHandler: productsHandler,
	}
}

// HandleMessage 处理消息
func (h *ProductNumberHandler) HandleMessage(ctx context.Context, message *tgbotapi.Message) error {
	// 解析用户输入的数字
	text := strings.TrimSpace(message.Text)
	productIndex, err := strconv.Atoi(text)
	if err != nil || productIndex < 1 {
		return h.productsHandler.sendError(message.Chat.ID, "请输入有效的产品编号（1-5）")
	}

	// 调用 ProductsHandler 的内部方法处理产品选择
	return h.productsHandler.handleProductSelection(ctx, message, productIndex)
}

// CanHandle 判断是否能处理该消息
func (h *ProductNumberHandler) CanHandle(message *tgbotapi.Message) bool {
	// 只处理纯数字消息（1-5）
	text := strings.TrimSpace(message.Text)
	if num, err := strconv.Atoi(text); err == nil && num >= 1 && num <= 5 {
		return true
	}
	return false
}

// GetHandlerName 获取处理器名称
func (h *ProductNumberHandler) GetHandlerName() string {
	return "product_number"
}
