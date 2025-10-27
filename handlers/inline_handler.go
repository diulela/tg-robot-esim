package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/storage/repository"
)

// InlineHandler Inline 查询处理器
type InlineHandler struct {
	bot         *tgbotapi.BotAPI
	productRepo repository.ProductRepository
	logger      Logger
}

// NewInlineHandler 创建 Inline 查询处理器
func NewInlineHandler(bot *tgbotapi.BotAPI, productRepo repository.ProductRepository, logger Logger) *InlineHandler {
	return &InlineHandler{
		bot:         bot,
		productRepo: productRepo,
		logger:      logger,
	}
}

// HandleInlineQuery 处理 Inline 查询
func (h *InlineHandler) HandleInlineQuery(ctx context.Context, query *tgbotapi.InlineQuery) error {
	h.logger.Debug("Processing inline query: %s", query.Query)

	queryText := strings.TrimSpace(query.Query)

	// 根据查询内容决定返回什么结果
	var results []interface{}
	var err error

	if queryText == "" || strings.Contains(strings.ToLower(queryText), "产品") || strings.Contains(strings.ToLower(queryText), "product") {
		// 显示产品列表
		results, err = h.buildProductListResults(ctx)
	} else if strings.Contains(strings.ToLower(queryText), "详情") || strings.Contains(strings.ToLower(queryText), "detail") {
		// 显示产品详情（如果查询包含产品ID）
		results, err = h.buildProductDetailResults(ctx, queryText)
	} else {
		// 搜索产品
		results, err = h.searchProducts(ctx, queryText)
	}

	if err != nil {
		h.logger.Error("Failed to build inline results: %v", err)
		// 返回错误结果
		results = []interface{}{
			tgbotapi.NewInlineQueryResultArticle(
				"error",
				"❌ 查询失败",
				"抱歉，查询产品时出现错误，请稍后重试。",
			),
		}
	}

	// 发送结果
	config := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		Results:       results,
		CacheTime:     300, // 缓存5分钟
	}

	_, err = h.bot.Request(config)
	return err
}

// buildProductListResults 构建产品列表结果
func (h *InlineHandler) buildProductListResults(ctx context.Context) ([]interface{}, error) {
	// 获取亚洲产品列表
	products, _, err := h.getAsiaProducts(ctx, 1, 10) // 获取前10个产品
	if err != nil {
		return nil, err
	}

	var results []interface{}

	// 添加产品列表标题
	listResult := tgbotapi.NewInlineQueryResultArticle(
		"product_list",
		"🌏 亚洲区域产品列表",
		h.buildProductListSummary(products),
	)
	listResult.Description = fmt.Sprintf("共 %d 个产品可选", len(products))

	// 设置消息内容
	messageText := h.buildInlineProductListText(products)
	listResult.InputMessageContent = tgbotapi.InputTextMessageContent{
		Text:      messageText,
		ParseMode: "HTML",
	}

	// 添加 Inline Keyboard
	keyboard := h.buildInlineProductListKeyboard(products)
	listResult.ReplyMarkup = &keyboard

	results = append(results, listResult)

	// 为每个产品添加单独的结果项
	for i, product := range products {
		productResult := tgbotapi.NewInlineQueryResultArticle(
			fmt.Sprintf("product_%d", product.ID),
			fmt.Sprintf("%d. %s", i+1, product.Name),
			h.buildProductSummary(product),
		)

		productResult.Description = fmt.Sprintf("%.2f USDT | %s | %d天",
			product.Price, formatDataSize(product.DataSize), product.ValidDays)

		// 设置点击后发送的消息
		productText := h.buildSingleProductInlineText(product, i+1)
		productResult.InputMessageContent = tgbotapi.InputTextMessageContent{
			Text:      productText,
			ParseMode: "HTML",
		}

		// 添加产品操作按钮
		productKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📖 查看详情", fmt.Sprintf("product_detail:%d", product.ID)),
				tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", product.ID)),
			),
		)
		productResult.ReplyMarkup = &productKeyboard

		results = append(results, productResult)
	}

	return results, nil
}

// buildProductDetailResults 构建产品详情结果
func (h *InlineHandler) buildProductDetailResults(ctx context.Context, query string) ([]interface{}, error) {
	// 尝试从查询中提取产品ID
	// 例如: "详情 1" 或 "detail 1"
	parts := strings.Fields(query)
	if len(parts) < 2 {
		return h.buildProductListResults(ctx) // 如果没有指定产品，返回产品列表
	}

	productIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return h.buildProductListResults(ctx)
	}

	// 获取产品列表
	products, _, err := h.getAsiaProducts(ctx, 1, 10)
	if err != nil || productIndex < 1 || productIndex > len(products) {
		return h.buildProductListResults(ctx)
	}

	product := products[productIndex-1]

	var results []interface{}

	detailResult := tgbotapi.NewInlineQueryResultArticle(
		fmt.Sprintf("detail_%d", product.ID),
		fmt.Sprintf("📱 %s - 详情", product.Name),
		h.buildProductDetailSummary(product),
	)

	detailResult.Description = fmt.Sprintf("完整产品信息 | %.2f USDT", product.Price)

	// 这里应该调用详细的产品信息格式化
	detailText := h.buildProductDetailInlineText(product)
	detailResult.InputMessageContent = tgbotapi.InputTextMessageContent{
		Text:      detailText,
		ParseMode: "HTML",
	}

	// 添加操作按钮
	detailKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", product.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回列表", "products_back"),
		),
	)
	detailResult.ReplyMarkup = &detailKeyboard

	results = append(results, detailResult)

	return results, nil
}

// searchProducts 搜索产品
func (h *InlineHandler) searchProducts(ctx context.Context, query string) ([]interface{}, error) {
	// 简单的搜索实现，可以根据需要扩展
	products, _, err := h.getAsiaProducts(ctx, 1, 10)
	if err != nil {
		return nil, err
	}

	var filteredProducts []*repository.ProductModel
	queryLower := strings.ToLower(query)

	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Name), queryLower) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	if len(filteredProducts) == 0 {
		// 没有找到匹配的产品，返回所有产品
		filteredProducts = products
	}

	var results []interface{}

	for i, product := range filteredProducts {
		result := tgbotapi.NewInlineQueryResultArticle(
			fmt.Sprintf("search_%d", product.ID),
			fmt.Sprintf("🔍 %s", product.Name),
			h.buildProductSummary(product),
		)

		result.Description = fmt.Sprintf("搜索结果 | %.2f USDT", product.Price)

		productText := h.buildSingleProductInlineText(product, i+1)
		result.InputMessageContent = tgbotapi.InputTextMessageContent{
			Text:      productText,
			ParseMode: "HTML",
		}

		productKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📖 查看详情", fmt.Sprintf("product_detail:%d", product.ID)),
				tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", product.ID)),
			),
		)
		result.ReplyMarkup = &productKeyboard

		results = append(results, result)
	}

	return results, nil
}

// 辅助方法

func (h *InlineHandler) getAsiaProducts(ctx context.Context, page, limit int) ([]*repository.ProductModel, int64, error) {
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

func (h *InlineHandler) buildProductListSummary(products []*repository.ProductModel) string {
	if len(products) == 0 {
		return "暂无产品"
	}
	return fmt.Sprintf("共 %d 个亚洲区域 eSIM 产品", len(products))
}

func (h *InlineHandler) buildProductSummary(product *repository.ProductModel) string {
	return fmt.Sprintf("📊 %s | ⏰ %d天 | 💰 %.2f USDT",
		formatDataSize(product.DataSize), product.ValidDays, product.Price)
}

func (h *InlineHandler) buildProductDetailSummary(product *repository.ProductModel) string {
	return fmt.Sprintf("完整产品信息：%s | %s | %d天 | %.2f USDT",
		product.Name, formatDataSize(product.DataSize), product.ValidDays, product.Price)
}

func (h *InlineHandler) buildInlineProductListText(products []*repository.ProductModel) string {
	text := "<b>🌏 亚洲区域产品</b>\n\n"

	for i, product := range products {
		text += fmt.Sprintf("<b>%d.</b> %s\n", i+1, escapeHTML(product.Name))
		text += fmt.Sprintf("   📊 %s  ⏰ %d天  \n💰 <b>%.2f USDT</b>\n\n",
			formatDataSize(product.DataSize), product.ValidDays, product.Price)
	}

	text += "<i>💡 点击下方按钮查看详情或购买</i>"
	return text
}

func (h *InlineHandler) buildSingleProductInlineText(product *repository.ProductModel, index int) string {
	text := fmt.Sprintf("<b>产品 %d</b>\n", index)
	text += fmt.Sprintf("<b>📱 %s</b>\n\n", escapeHTML(product.Name))
	text += fmt.Sprintf("📊 <b>流量：</b>%s\n", formatDataSize(product.DataSize))
	text += fmt.Sprintf("⏰ <b>有效期：</b>%d天\n", product.ValidDays)
	text += fmt.Sprintf("💰 <b>价格：</b><u>%.2f USDT</u>", product.Price)
	return text
}

func (h *InlineHandler) buildProductDetailInlineText(product *repository.ProductModel) string {
	text := fmt.Sprintf("<b>📱 %s</b>\n\n", escapeHTML(product.Name))
	text += "<blockquote>"
	text += fmt.Sprintf("<b>📊 流量：</b><code>%s</code>\n", formatDataSize(product.DataSize))
	text += fmt.Sprintf("<b>⏰ 有效期：</b><code>%d天</code>\n\n", product.ValidDays)
	text += fmt.Sprintf("<b>💰 价格：</b><u><b>%.2f USDT</b></u>", product.Price)
	text += "</blockquote>\n\n"
	text += "<i>🔄 通过 Inline Mode 查看</i>"
	return text
}

func (h *InlineHandler) buildInlineProductListKeyboard(products []*repository.ProductModel) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// 添加快速操作按钮
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔍 选择产品", "product_select"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// escapeHTML 转义 HTML 特殊字符
func escapeHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
	)
	return replacer.Replace(text)
}

// GetHandlerName 获取处理器名称
func (h *InlineHandler) GetHandlerName() string {
	return "inline_query"
}
