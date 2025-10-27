package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/services"
	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// ProductsHandler 商品列表处理器
type ProductsHandler struct {
	bot               *tgbotapi.BotAPI
	esimService       services.EsimService
	productRepo       repository.ProductRepository
	productDetailRepo repository.ProductDetailRepository
	logger            Logger
}

// NewProductsHandler 创建商品处理器
func NewProductsHandler(bot *tgbotapi.BotAPI, esimService services.EsimService, productRepo repository.ProductRepository, productDetailRepo repository.ProductDetailRepository, logger Logger) *ProductsHandler {
	return &ProductsHandler{
		bot:               bot,
		esimService:       esimService,
		productRepo:       productRepo,
		productDetailRepo: productDetailRepo,
		logger:            logger,
	}
}

// HandleCallback 处理回调查询
func (h *ProductsHandler) HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	data := callback.Data
	userID := callback.From.ID

	h.logger.Debug("Products handler processing callback: %s", data)

	// 回答回调
	if err := h.answerCallback(callback.ID); err != nil {
		h.logger.Error("Failed to answer callback: %v", err)
	}

	// 解析回调数据
	parts := strings.Split(data, ":")
	action := parts[0]

	switch action {
	case "products_back":
		// 直接显示亚洲产品列表
		return h.showAsiaProducts(ctx, callback.Message, 1)
	case "products_page":
		if len(parts) >= 2 {
			page, _ := strconv.Atoi(parts[1])
			return h.showAsiaProducts(ctx, callback.Message, page)
		}
	case "product_detail":
		if len(parts) >= 2 {
			productID, _ := strconv.Atoi(parts[1])
			return h.showProductDetail(ctx, callback.Message, productID)
		}
	case "product_buy":
		if len(parts) >= 2 {
			productID, _ := strconv.Atoi(parts[1])
			return h.startPurchase(ctx, callback.Message, userID, productID)
		}
	}

	return nil
}

// CanHandle 判断是否能处理该回调
func (h *ProductsHandler) CanHandle(callback *tgbotapi.CallbackQuery) bool {
	return strings.HasPrefix(callback.Data, "products_") ||
		strings.HasPrefix(callback.Data, "product_")
}

// GetHandlerName 获取处理器名称
func (h *ProductsHandler) GetHandlerName() string {
	return "products"
}

// HandleCommand 处理命令
func (h *ProductsHandler) HandleCommand(ctx context.Context, message *tgbotapi.Message) error {
	// 直接显示亚洲产品列表
	return h.showAsiaProductsNew(ctx, message.Chat.ID, 1)
}

// GetCommand 获取处理的命令名称
func (h *ProductsHandler) GetCommand() string {
	return "products"
}

// GetDescription 获取命令描述
func (h *ProductsHandler) GetDescription() string {
	return "浏览 eSIM 产品"
}

// showAsiaProducts 显示亚洲产品列表（编辑消息）- 多消息卡片模式
func (h *ProductsHandler) showAsiaProducts(ctx context.Context, message *tgbotapi.Message, page int) error {
	// 删除旧消息
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	h.bot.Send(deleteMsg)

	// 使用新消息方式显示
	return h.showAsiaProductsNew(ctx, message.Chat.ID, page)
}

// showAsiaProductsNew 显示亚洲产品列表（新消息）- 每个产品一条消息
func (h *ProductsHandler) showAsiaProductsNew(ctx context.Context, chatID int64, page int) error {
	products, total, err := h.getAsiaProducts(ctx, page, 5)
	if err != nil {
		h.logger.Error("Failed to get Asia products: %v", err)
		return h.sendError(chatID, "获取产品列表失败")
	}

	if len(products) == 0 {
		return h.sendError(chatID, "暂无产品")
	}

	totalPages := int((total + 4) / 5)

	// 1. 发送标题消息
	headerText := h.buildProductListHeader(page, totalPages, int(total))
	headerMsg := tgbotapi.NewMessage(chatID, headerText)
	headerMsg.ParseMode = "HTML"
	if _, err := h.bot.Send(headerMsg); err != nil {
		h.logger.Error("Failed to send header: %v", err)
		return err
	}

	// 2. 为每个产品发送独立的卡片消息（带按钮）
	for i, product := range products {
		cardText := h.buildSingleProductCard(product, i+1)
		cardKeyboard := h.buildSingleProductKeyboard(product.ID)

		cardMsg := tgbotapi.NewMessage(chatID, cardText)
		cardMsg.ParseMode = "HTML"
		cardMsg.ReplyMarkup = cardKeyboard

		if _, err := h.bot.Send(cardMsg); err != nil {
			h.logger.Error("Failed to send product card %d: %v", i+1, err)
		}
	}

	// 3. 发送分页导航消息
	navText := "<i>💡 点击产品卡片上的按钮查看详情或购买</i>"
	navKeyboard := h.buildProductListNavigation(page, totalPages)

	navMsg := tgbotapi.NewMessage(chatID, navText)
	navMsg.ParseMode = "HTML"
	navMsg.ReplyMarkup = navKeyboard

	_, err = h.bot.Send(navMsg)
	return err
}

// escapeMarkdown 转义 Markdown 特殊字符
func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

// buildProductListKeyboard 构建产品列表键盘
func (h *ProductsHandler) buildProductListKeyboard(products []esim.Product, productType esim.ProductType, pagination esim.Pagination) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// 产品按钮 - 每行2个
	for i := 0; i < len(products); i += 2 {
		var row []tgbotapi.InlineKeyboardButton

		// 第一个按钮
		btn1 := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. 详情", i+1),
			fmt.Sprintf("product_detail:%d", products[i].ID),
		)
		row = append(row, btn1)

		// 第二个按钮（如果存在）
		if i+1 < len(products) {
			btn2 := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d. 详情", i+2),
				fmt.Sprintf("product_detail:%d", products[i+1].ID),
			)
			row = append(row, btn2)
		}

		rows = append(rows, row)
	}

	// 分页按钮
	if pagination.TotalPages > 1 {
		var pageRow []tgbotapi.InlineKeyboardButton

		if pagination.Page > 1 {
			pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
				"⬅️ 上一页",
				fmt.Sprintf("products_page:%s:%d", productType, pagination.Page-1),
			))
		}

		pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("📄 %d/%d", pagination.Page, pagination.TotalPages),
			"noop",
		))

		if pagination.Page < pagination.TotalPages {
			pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
				"下一页 ➡️",
				fmt.Sprintf("products_page:%s:%d", productType, pagination.Page+1),
			))
		}

		rows = append(rows, pageRow)
	}

	// 返回按钮
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 返回", "products_back"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// showProductDetail 显示产品详情（优先从数据库获取，降级到API）
func (h *ProductsHandler) showProductDetail(ctx context.Context, message *tgbotapi.Message, productID int) error {
	var text string
	var err error

	// 首先尝试从产品详情表获取
	productDetail, err := h.productDetailRepo.GetByProductID(ctx, productID)
	if err == nil && productDetail != nil {
		h.logger.Debug("Got product detail from database for product %d", productID)
		text = h.formatProductDetailFromDetailDB(productDetail)
	} else {
		h.logger.Debug("Product detail not found in database for product %d, trying API", productID)

		// 从数据库获取失败，尝试从API获取
		text, err = h.getProductDetailFromAPI(ctx, productID)
		if err != nil {
			h.logger.Error("Failed to get product detail from API: %v", err)
			return h.sendError(message.Chat.ID, "产品详情不存在")
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", productID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回列表", "products_back"),
		),
	)

	return h.editOrSendMessage(message, text, keyboard)
}

// getProductDetailFromAPI 从API获取产品详情
func (h *ProductsHandler) getProductDetailFromAPI(ctx context.Context, productID int) (string, error) {
	// 首先从产品表获取基本信息，以获取第三方ID
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return "", fmt.Errorf("product not found: %w", err)
	}

	// 提取第三方ID
	thirdPartyID := extractThirdPartyIDFromString(product.ThirdPartyID)
	if thirdPartyID == 0 {
		return "", fmt.Errorf("invalid third party ID: %s", product.ThirdPartyID)
	}

	// 调用API获取详情
	resp, err := h.esimService.GetProduct(ctx, thirdPartyID)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	if !resp.Success || resp.ProductDetail == nil {
		return "", fmt.Errorf("API returned no product detail")
	}

	// 格式化API返回的详情
	return h.formatProductDetailFromAPI(resp.ProductDetail), nil
}

// extractThirdPartyIDFromString 从字符串中提取第三方ID
func extractThirdPartyIDFromString(thirdPartyID string) int {
	// 如果是 "product-123" 格式，提取数字
	if strings.HasPrefix(thirdPartyID, "product-") {
		idStr := strings.TrimPrefix(thirdPartyID, "product-")
		if id, err := strconv.Atoi(idStr); err == nil {
			return id
		}
	}
	// 尝试直接转换
	if id, err := strconv.Atoi(thirdPartyID); err == nil {
		return id
	}
	return 0
}

// formatProductDetailFromAPI 格式化API返回的产品详情（使用 HTML 格式）
func (h *ProductsHandler) formatProductDetailFromAPI(detail *esim.ProductDetail) string {
	// 产品标题 - 使用大标题样式
	text := fmt.Sprintf("<b>📱 %s</b>\n\n", escapeHTML(detail.Name))

	// 产品类型标签
	typeText := map[string]string{
		"local":    "� 本地",
		"regional": " 区域",
		"global":   "🌍 全球",
	}
	if t, ok := typeText[detail.Type]; ok {
		text += fmt.Sprintf("<blockquote><b>类型：</b>%s</blockquote>\n", t)
	}

	// 国家列表
	if len(detail.Countries) > 0 {
		text += "<blockquote><b>🗺️ 支持国家：</b>\n"
		if len(detail.Countries) <= 5 {
			countryNames := make([]string, len(detail.Countries))
			for i, c := range detail.Countries {
				countryNames[i] = c.CN
			}
			text += strings.Join(countryNames, " • ")
		} else {
			countryNames := make([]string, 5)
			for i := 0; i < 5; i++ {
				countryNames[i] = detail.Countries[i].CN
			}
			text += strings.Join(countryNames, " • ")
			text += fmt.Sprintf(" <i>等%d个国家</i>", len(detail.Countries))
		}
		text += "</blockquote>\n"
	}

	// 流量和有效期
	dataSize := "无限流量"
	if detail.DataSize > 0 {
		if detail.DataSize >= 1024 {
			dataSize = fmt.Sprintf("%.1fGB", float64(detail.DataSize)/1024)
		} else {
			dataSize = fmt.Sprintf("%dMB", detail.DataSize)
		}
	}

	text += "<blockquote>"
	text += fmt.Sprintf("<b>📊 流量：</b><code>%s</code>\n", dataSize)
	text += fmt.Sprintf("<b>⏰ 有效期：</b><code>%d天</code>", detail.ValidDays)
	text += "</blockquote>\n"

	// 价格 - 突出显示
	text += fmt.Sprintf("\n<blockquote><b>💰 价格：</b><u><b>%.2f USDT</b></u></blockquote>\n", detail.Price)

	// 产品描述
	if detail.Description != "" {
		text += fmt.Sprintf("\n<blockquote expandable><b>📝 产品描述</b>\n\n%s</blockquote>\n", escapeHTML(detail.Description))
	}

	// 产品特性
	if len(detail.Features) > 0 {
		text += "\n<blockquote><b>✨ 产品特性</b>\n"
		for _, feature := range detail.Features {
			text += fmt.Sprintf("  • %s\n", escapeHTML(feature))
		}
		text += "</blockquote>\n"
	}

	// 添加数据来源标识
	text += "\n<i>🔄 数据来源：实时API</i>"

	return text
}

// formatProductDetailFromDetailDB 格式化产品详情消息（从产品详情表，使用 HTML 格式）
func (h *ProductsHandler) formatProductDetailFromDetailDB(detail *models.ProductDetail) string {
	// 产品标题 - 使用大标题样式
	text := fmt.Sprintf("<b>📱 %s</b>\n\n", escapeHTML(detail.Name))

	// 产品类型标签
	typeText := map[string]string{
		"local":    " 本地",
		"regional": " 区域",
		"global":   "🌍 全球",
	}
	if t, ok := typeText[detail.Type]; ok {
		text += fmt.Sprintf("<blockquote><b>类型：</b>%s</blockquote>\n", t)
	}

	// 解析国家列表
	var countries []string
	if err := json.Unmarshal([]byte(detail.Countries), &countries); err == nil && len(countries) > 0 {
		text += "<blockquote><b>🗺️ 支持国家：</b>\n"
		if len(countries) <= 5 {
			text += strings.Join(countries, " • ")
		} else {
			text += strings.Join(countries[:5], " • ")
			text += fmt.Sprintf(" <i>等%d个国家</i>", len(countries))
		}
		text += "</blockquote>\n"
	}

	// 产品规格 - 使用表格式布局
	text += "<blockquote>"
	text += fmt.Sprintf("<b>📊 流量：</b><code>%s</code>\n", detail.DataSize)
	text += fmt.Sprintf("<b>⏰ 有效期：</b><code>%d天</code>", detail.ValidDays)
	text += "</blockquote>\n"

	// 价格 - 突出显示
	text += fmt.Sprintf("\n<blockquote><b>💰 价格：</b><u><b>%.2f USDT</b></u></blockquote>\n", detail.Price)

	// 产品描述
	if detail.Description != "" {
		text += fmt.Sprintf("\n<blockquote expandable><b>📝 产品描述</b>\n\n%s</blockquote>\n", escapeHTML(detail.Description))
	}

	// 解析特性列表
	var features []string
	if err := json.Unmarshal([]byte(detail.Features), &features); err == nil && len(features) > 0 {
		text += "\n<blockquote><b>✨ 产品特性</b>\n"
		for _, feature := range features {
			text += fmt.Sprintf("  • %s\n", escapeHTML(feature))
		}
		text += "</blockquote>\n"
	}

	return text
}

// startPurchase 开始购买流程
func (h *ProductsHandler) startPurchase(ctx context.Context, message *tgbotapi.Message, userID int64, productID int) error {
	text := "🛒 *开始购买流程*\n\n"
	text += "请提供以下信息：\n"
	text += "1. 客户邮箱地址（必填）\n"
	text += "2. 客户手机号（可选）\n"
	text += "3. 购买数量（默认1）\n\n"
	text += "请按以下格式发送：\n"
	text += "`customer@example.com`\n"
	text += "或\n"
	text += "`customer@example.com,+86 138 0000 0000,2`"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ 取消", fmt.Sprintf("product_detail:%d", productID)),
		),
	)

	return h.editOrSendMessage(message, text, keyboard)
}

// Helper functions

func (h *ProductsHandler) editOrSendMessage(message *tgbotapi.Message, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, text)
	editMsg.ParseMode = "HTML"
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	return err
}

func (h *ProductsHandler) sendError(chatID int64, errorMsg string) error {
	msg := tgbotapi.NewMessage(chatID, "❌ "+errorMsg)
	_, err := h.bot.Send(msg)
	return err
}

func (h *ProductsHandler) answerCallback(callbackID string) error {
	callback := tgbotapi.NewCallback(callbackID, "")
	_, err := h.bot.Request(callback)
	return err
}

func (h *ProductsHandler) showMainMenu(message *tgbotapi.Message) error {
	text := "📱 <b>主菜单</b>\n\n请选择您需要的功能："

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

	editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, text)
	editMsg.ParseMode = "HTML"
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
	return err
}

func formatDataSize(sizeMB int) string {
	if sizeMB >= 1024 {
		return fmt.Sprintf("%.1fGB", float64(sizeMB)/1024)
	}
	return fmt.Sprintf("%dMB", sizeMB)
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

// getAsiaProducts 获取亚洲产品列表
func (h *ProductsHandler) getAsiaProducts(ctx context.Context, page, limit int) ([]*repository.ProductModel, int64, error) {
	// 从数据库获取 type=regional 且 name 包含"亚洲"的产品
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

// buildAsiaProductListText 构建亚洲产品列表文本（使用 HTML 格式）
func (h *ProductsHandler) buildAsiaProductListText(products []*repository.ProductModel, page int, total int64, limit int) string {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	text := "<b>🌏 亚洲区域产品</b>\n\n"
	text += fmt.Sprintf("📄 第 <b>%d</b>/<b>%d</b> 页 (共 <b>%d</b> 个产品)\n",
		page, totalPages, total)
	text += "━━━━━━━━━━━━━━━━━━\n\n"

	for i, product := range products {
		// 产品卡片开始 - 使用引用块样式
		text += "<blockquote>"

		// 产品标题 - 加粗并使用 emoji
		text += fmt.Sprintf("<b>%d. 📱 %s</b>\n\n", i+1, escapeHTML(product.Name))

		// 产品信息 - 使用表格式布局
		text += fmt.Sprintf("📊 <b>流量：</b><code>%s</code>  ", formatDataSize(product.DataSize))
		text += fmt.Sprintf("⏰ <b>有效期：</b><code>%d天</code>\n", product.ValidDays)

		// 价格 - 突出显示
		text += fmt.Sprintf("\n💰 <b>价格：</b><u>%.2f USDT</u>", product.Price)

		text += "</blockquote>\n\n"
	}

	text += "<i>💡 点击下方按钮查看产品详情</i>"

	return text
}

// buildProductListHeader 构建产品列表标题
func (h *ProductsHandler) buildProductListHeader(page, totalPages, total int) string {
	text := "<b>🌏 亚洲区域产品</b>\n\n"
	text += fmt.Sprintf("📄 第 <b>%d</b>/<b>%d</b> 页 (共 <b>%d</b> 个产品)\n",
		page, totalPages, total)
	text += "━━━━━━━━━━━━━━━━━━"
	return text
}

// buildSingleProductCard 构建单个产品卡片（独立消息）
func (h *ProductsHandler) buildSingleProductCard(product *repository.ProductModel, index int) string {
	text := "<blockquote expandable>"

	// 产品标题
	text += fmt.Sprintf("<b>📱 %s</b>\n", escapeHTML(product.Name))
	text += "━━━━━━━━━━━━━━━━\n\n"

	// 产品规格
	text += fmt.Sprintf("📊 <b>流量：</b><code>%s</code>\n", formatDataSize(product.DataSize))
	text += fmt.Sprintf("⏰ <b>有效期：</b><code>%d天</code>\n\n", product.ValidDays)

	// 价格 - 突出显示
	text += fmt.Sprintf("💰 <b>价格：</b><u><b>%.2f USDT</b></u>", product.Price)

	text += "</blockquote>"

	return text
}

// buildSingleProductKeyboard 构建单个产品的按钮
func (h *ProductsHandler) buildSingleProductKeyboard(productID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 查看详情", fmt.Sprintf("product_detail:%d", productID)),
			tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", productID)),
		),
	)
}

// buildProductListNavigation 构建产品列表导航按钮
func (h *ProductsHandler) buildProductListNavigation(page, totalPages int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// 分页按钮
	if totalPages > 1 {
		var pageRow []tgbotapi.InlineKeyboardButton

		if page > 1 {
			pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
				"⬅️ 上一页",
				fmt.Sprintf("products_page:%d", page-1),
			))
		}

		pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("📄 %d/%d", page, totalPages),
			"noop",
		))

		if page < totalPages {
			pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
				"下一页 ➡️",
				fmt.Sprintf("products_page:%d", page+1),
			))
		}

		rows = append(rows, pageRow)
	}

	// 返回按钮
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 返回主菜单", "main_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
