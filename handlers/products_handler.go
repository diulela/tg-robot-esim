package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/services"
	"tg-robot-sim/storage/repository"
)

// ProductsHandler 商品列表处理器
type ProductsHandler struct {
	bot         *tgbotapi.BotAPI
	esimService services.EsimService
	productRepo repository.ProductRepository
	logger      Logger
}

// NewProductsHandler 创建商品处理器
func NewProductsHandler(bot *tgbotapi.BotAPI, esimService services.EsimService, productRepo repository.ProductRepository, logger Logger) *ProductsHandler {
	return &ProductsHandler{
		bot:         bot,
		esimService: esimService,
		productRepo: productRepo,
		logger:      logger,
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

// showAsiaProducts 显示亚洲产品列表（编辑消息）
func (h *ProductsHandler) showAsiaProducts(ctx context.Context, message *tgbotapi.Message, page int) error {
	products, total, err := h.getAsiaProducts(ctx, page, 5)
	if err != nil {
		h.logger.Error("Failed to get Asia products: %v", err)
		return h.sendError(message.Chat.ID, "获取产品列表失败")
	}

	if len(products) == 0 {
		return h.sendError(message.Chat.ID, "暂无产品")
	}

	// 构建消息文本
	text := h.buildAsiaProductListText(products, page, total, 5)

	// 构建键盘
	keyboard := h.buildAsiaProductKeyboard(products, page, total, 5)

	// 编辑消息
	return h.editOrSendMessage(message, text, keyboard)
}

// showAsiaProductsNew 显示亚洲产品列表（新消息）
func (h *ProductsHandler) showAsiaProductsNew(ctx context.Context, chatID int64, page int) error {
	products, total, err := h.getAsiaProducts(ctx, page, 5)
	if err != nil {
		h.logger.Error("Failed to get Asia products: %v", err)
		return h.sendError(chatID, "获取产品列表失败")
	}

	if len(products) == 0 {
		return h.sendError(chatID, "暂无产品")
	}

	// 构建消息文本
	text := h.buildAsiaProductListText(products, page, total, 5)

	// 构建键盘
	keyboard := h.buildAsiaProductKeyboard(products, page, total, 5)

	// 发送新消息
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	return err
}

// showProducts 显示产品列表
func (h *ProductsHandler) showProducts(ctx context.Context, message *tgbotapi.Message, productType esim.ProductType, page int) error {
	// 获取产品列表
	resp, err := h.esimService.GetProducts(ctx, &esim.ProductParams{
		Type:  productType,
		Page:  page,
		Limit: 5,
	})

	if err != nil {
		h.logger.Error("Failed to get products: %v", err)
		return h.sendError(message.Chat.ID, "获取产品列表失败")
	}

	if !resp.Success || len(resp.Message.Products) == 0 {
		return h.sendError(message.Chat.ID, "暂无产品")
	}

	// 调试：打印第一个产品的详细信息
	if len(resp.Message.Products) > 0 {
		p := resp.Message.Products[0]
		h.logger.Debug("First product details: ID=%d, Name=%s, DataSize=%d, ValidDays=%d, RetailPrice=%.2f, AgentPrice=%.2f, Countries=%d",
			p.ID, p.Name, p.DataSize, p.ValidDays, p.RetailPrice, p.AgentPrice, len(p.Countries))
	}

	// 构建消息文本
	text := h.buildProductListText(productType, resp.Message.Products, resp.Message.Pagination)

	// 构建键盘
	keyboard := h.buildProductListKeyboard(resp.Message.Products, productType, resp.Message.Pagination)

	// 编辑或发送消息
	return h.editOrSendMessage(message, text, keyboard)
}

// buildProductListText 构建产品列表文本
func (h *ProductsHandler) buildProductListText(productType esim.ProductType, products []esim.Product, pagination esim.Pagination) string {
	typeText := map[esim.ProductType]string{
		esim.ProductTypeLocal:    "🏠 本地产品",
		esim.ProductTypeRegional: "🌏 区域产品",
		esim.ProductTypeGlobal:   "🌍 全球产品",
	}

	text := fmt.Sprintf("*%s*\n\n", typeText[productType])
	text += fmt.Sprintf("📄 第 %d/%d 页 (共 %d 个产品)\n",
		pagination.Page, pagination.TotalPages, pagination.Total)
	text += "━━━━━━━━━━━━━━━━━━\n\n"

	for i, product := range products {
		// 产品标题
		text += fmt.Sprintf("*%d\\. %s*\n", i+1, escapeMarkdown(product.Name))

		// 国家信息（简化显示）
		if len(product.Countries) > 0 {
			if len(product.Countries) == 1 {
				text += fmt.Sprintf("🗺️ %s\n", product.Countries[0].CN)
			} else if len(product.Countries) <= 3 {
				countryNames := make([]string, len(product.Countries))
				for j, country := range product.Countries {
					countryNames[j] = country.CN
				}
				text += fmt.Sprintf("🗺️ %s\n", strings.Join(countryNames, "、"))
			} else {
				text += fmt.Sprintf("🗺️ %s、%s 等%d国\n",
					product.Countries[0].CN, product.Countries[1].CN, len(product.Countries))
			}
		}

		// 流量和有效期
		text += fmt.Sprintf("📊 %s  ⏰ %d天\n",
			formatDataSize(product.DataSize), product.ValidDays)

		// 价格
		text += fmt.Sprintf("💵 代理价: *$%.2f*  💰 零售价: $%.2f\n",
			product.AgentPrice, product.RetailPrice)

		text += "\n"
	}

	return text
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

// showProductDetail 显示产品详情
func (h *ProductsHandler) showProductDetail(ctx context.Context, message *tgbotapi.Message, productID int) error {
	resp, err := h.esimService.GetProduct(ctx, productID)
	if err != nil {
		h.logger.Error("Failed to get product detail: %v", err)
		return h.sendError(message.Chat.ID, "获取产品详情失败")
	}

	if !resp.Success {
		return h.sendError(message.Chat.ID, "产品不存在")
	}

	// 注意：产品数据在 Message 字段中
	product := resp.Message
	text := services.FormatProductMessage(&product)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛒 立即购买", fmt.Sprintf("product_buy:%d", productID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回列表", fmt.Sprintf("products_%s", product.Type)),
		),
	)

	return h.editOrSendMessage(message, text, keyboard)
}

// searchProductsByCountry 按国家代码搜索产品
func (h *ProductsHandler) searchProductsByCountry(ctx context.Context, chatID int64, countryCode string) error {
	resp, err := h.esimService.GetProducts(ctx, &esim.ProductParams{
		Country: countryCode,
		Limit:   10,
	})

	if err != nil {
		h.logger.Error("Failed to search products: %v", err)
		return h.sendError(chatID, "搜索产品失败")
	}

	if !resp.Success || len(resp.Message.Products) == 0 {
		text := fmt.Sprintf("未找到国家代码 *%s* 的产品\n\n", countryCode)
		text += "请检查国家代码是否正确，或使用 /products 浏览所有产品"

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		_, err := h.bot.Send(msg)
		return err
	}

	text := fmt.Sprintf("🔍 *搜索结果: %s*\n\n", countryCode)
	text += fmt.Sprintf("找到 %d 个产品\n", len(resp.Message.Products))
	text += "━━━━━━━━━━━━━━━━━━\n\n"

	for i, product := range resp.Message.Products {
		text += fmt.Sprintf("*%d\\. %s*\n", i+1, escapeMarkdown(product.Name))
		text += fmt.Sprintf("📊 %s  ⏰ %d天\n",
			formatDataSize(product.DataSize), product.ValidDays)
		text += fmt.Sprintf("💵 代理价: *$%.2f*  💰 零售价: $%.2f\n\n",
			product.AgentPrice, product.RetailPrice)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	// 每行2个按钮
	for i := 0; i < len(resp.Message.Products); i += 2 {
		var row []tgbotapi.InlineKeyboardButton

		btn1 := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. 详情", i+1),
			fmt.Sprintf("product_detail:%d", resp.Message.Products[i].ID),
		)
		row = append(row, btn1)

		if i+1 < len(resp.Message.Products) {
			btn2 := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d. 详情", i+2),
				fmt.Sprintf("product_detail:%d", resp.Message.Products[i+1].ID),
			)
			row = append(row, btn2)
		}

		rows = append(rows, row)
	}

	// 添加返回按钮
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 返回", "products_back"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	return err
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
	editMsg.ParseMode = "Markdown"
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

// buildAsiaProductListText 构建亚洲产品列表文本
func (h *ProductsHandler) buildAsiaProductListText(products []*repository.ProductModel, page int, total int64, limit int) string {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	text := "*🌏 亚洲区域产品*\n\n"
	text += fmt.Sprintf("📄 第 %d/%d 页 (共 %d 个产品)\n",
		page, totalPages, total)
	text += "━━━━━━━━━━━━━━━━━━\n\n"

	for i, product := range products {
		// 产品标题
		text += fmt.Sprintf("*%d\\. %s*\n", i+1, escapeMarkdown(product.Name))

		// 流量和有效期
		text += fmt.Sprintf("📊 %s  ⏰ %d天\n",
			formatDataSize(product.DataSize), product.ValidDays)

		// 价格（只显示零售价，单位改为 USDT）
		text += fmt.Sprintf("💰 价格: *%.2f USDT*\n", product.Price)

		text += "\n"
	}

	return text
}

// buildAsiaProductKeyboard 构建亚洲产品键盘
func (h *ProductsHandler) buildAsiaProductKeyboard(products []*repository.ProductModel, page int, total int64, limit int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	totalPages := int((total + int64(limit) - 1) / int64(limit))

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
