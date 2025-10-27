package handlers

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/storage/repository"
)

// InlineHandler 内联查询处理器
type InlineHandler struct {
	bot         *tgbotapi.BotAPI
	productRepo repository.ProductRepository
	logger      Logger
}

// NewInlineHandler 创建内联查询处理器
func NewInlineHandler(bot *tgbotapi.BotAPI, productRepo repository.ProductRepository, logger Logger) *InlineHandler {
	return &InlineHandler{
		bot:         bot,
		productRepo: productRepo,
		logger:      logger,
	}
}

// HandleInlineQuery 处理内联查询
func (h *InlineHandler) HandleInlineQuery(ctx context.Context, query *tgbotapi.InlineQuery) error {
	h.logger.Debug("Processing inline query: '%s' from user %d", query.Query, query.From.ID)

	// 解析查询
	searchTerm := strings.TrimSpace(query.Query)
	if searchTerm == "" {
		// 空查询，显示热门产品
		return h.showHotProducts(ctx, query)
	}

	// 搜索产品
	return h.searchProducts(ctx, query, searchTerm)
}

// showHotProducts 显示热门产品
func (h *InlineHandler) showHotProducts(ctx context.Context, query *tgbotapi.InlineQuery) error {
	// 获取热门产品
	isHot := true
	params := repository.ListParams{
		Status:  "active",
		IsHot:   &isHot,
		Limit:   10,
		OrderBy: "sort_order",
	}

	products, _, err := h.productRepo.List(ctx, params)
	if err != nil {
		h.logger.Error("Failed to get hot products: %v", err)
		return h.answerEmptyQuery(query, "获取产品失败")
	}

	if len(products) == 0 {
		return h.answerEmptyQuery(query, "暂无热门产品")
	}

	// 构建结果
	results := make([]interface{}, 0, len(products))
	for i, product := range products {
		if i >= 10 { // 限制结果数量
			break
		}

		result := h.createProductInlineResult(product, fmt.Sprintf("hot_%d", i))
		results = append(results, result)
	}

	return h.answerInlineQuery(query, results)
}

// searchProducts 搜索产品
func (h *InlineHandler) searchProducts(ctx context.Context, query *tgbotapi.InlineQuery, searchTerm string) error {
	// 搜索产品（按名称）
	params := repository.ListParams{
		Status:   "active",
		NameLike: searchTerm,
		Limit:    20,
		OrderBy:  "sort_order",
	}

	products, _, err := h.productRepo.List(ctx, params)
	if err != nil {
		h.logger.Error("Failed to search products: %v", err)
		return h.answerEmptyQuery(query, "搜索失败")
	}

	if len(products) == 0 {
		return h.answerEmptyQuery(query, fmt.Sprintf("未找到包含 '%s' 的产品", searchTerm))
	}

	// 构建结果
	results := make([]interface{}, 0, len(products))
	for i, product := range products {
		if i >= 20 { // 限制结果数量
			break
		}

		result := h.createProductInlineResult(product, fmt.Sprintf("search_%d", i))
		results = append(results, result)
	}

	return h.answerInlineQuery(query, results)
}

// createProductInlineResult 创建产品内联结果
func (h *InlineHandler) createProductInlineResult(product *repository.ProductModel, resultID string) tgbotapi.InlineQueryResultArticle {
	// 构建标题和描述
	title := product.Name
	description := fmt.Sprintf("📊 %s | ⏰ %d天 | 💰 %.2f USDT",
		formatDataSize(product.DataSize), product.ValidDays, product.Price)

	// 构建消息内容
	messageText := h.buildInlineProductMessage(product)

	// 构建内联键盘
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 查看详情", fmt.Sprintf("product_detail:%d", product.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", product.ID)),
		),
	)

	return tgbotapi.NewInlineQueryResultArticle(resultID, title, messageText).
		SetDescription(description).
		SetReplyMarkup(&keyboard).
		SetParseMode("Markdown")
}

// buildInlineProductMessage 构建内联产品消息
func (h *InlineHandler) buildInlineProductMessage(product *repository.ProductModel) string {
	text := fmt.Sprintf("📱 *%s*\n\n", escapeMarkdown(product.Name))

	// 产品类型
	typeText := map[string]string{
		"local":    "🏠 本地",
		"regional": "🌏 区域",
		"global":   "🌍 全球",
	}
	if t, ok := typeText[product.Type]; ok {
		text += fmt.Sprintf("类型: %s\n", t)
	}

	// 流量和有效期
	text += fmt.Sprintf("📊 流量: %s\n", formatDataSize(product.DataSize))
	text += fmt.Sprintf("⏰ 有效期: %d天\n", product.ValidDays)

	// 价格
	text += fmt.Sprintf("\n💰 价格: *%.2f USDT*\n", product.Price)

	// 简短描述
	if product.Description != "" && len(product.Description) > 0 {
		desc := product.Description
		if len(desc) > 100 {
			desc = desc[:100] + "..."
		}
		text += fmt.Sprintf("\n📝 %s\n", desc)
	}

	text += "\n_点击按钮查看详情或购买_"

	return text
}

// answerInlineQuery 回答内联查询
func (h *InlineHandler) answerInlineQuery(query *tgbotapi.InlineQuery, results []interface{}) error {
	config := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		Results:       results,
		CacheTime:     300, // 缓存5分钟
		IsPersonal:    true,
	}

	_, err := h.bot.Request(config)
	return err
}

// answerEmptyQuery 回答空查询
func (h *InlineHandler) answerEmptyQuery(query *tgbotapi.InlineQuery, message string) error {
	result := tgbotapi.NewInlineQueryResultArticle("empty", "暂无结果", message).
		SetDescription(message)

	config := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		Results:       []interface{}{result},
		CacheTime:     60,
		IsPersonal:    true,
	}

	_, err := h.bot.Request(config)
	return err
}

// CanHandle 判断是否能处理该内联查询
func (h *InlineHandler) CanHandle(query *tgbotapi.InlineQuery) bool {
	return true // 处理所有内联查询
}

// GetHandlerName 获取处理器名称
func (h *InlineHandler) GetHandlerName() string {
	return "inline_query"
}
