package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/services"
)

// ProductsHandler 商品列表处理器
type ProductsHandler struct {
	bot         *tgbotapi.BotAPI
	esimService services.EsimService
	logger      Logger
}

// NewProductsHandler 创建商品处理器
func NewProductsHandler(bot *tgbotapi.BotAPI, esimService services.EsimService, logger Logger) *ProductsHandler {
	return &ProductsHandler{
		bot:         bot,
		esimService: esimService,
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
	case "products_local":
		return h.showProducts(ctx, callback.Message, esim.ProductTypeLocal, 1)
	case "products_regional":
		return h.showProducts(ctx, callback.Message, esim.ProductTypeRegional, 1)
	case "products_global":
		return h.showProducts(ctx, callback.Message, esim.ProductTypeGlobal, 1)
	case "products_page":
		if len(parts) >= 3 {
			productType := esim.ProductType(parts[1])
			page, _ := strconv.Atoi(parts[2])
			return h.showProducts(ctx, callback.Message, productType, page)
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
	case "products_back":
		return h.showMainMenu(callback.Message)
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
	// 解析命令参数（国家代码）
	args := strings.TrimSpace(strings.TrimPrefix(message.Text, "/products"))

	if args == "" {
		// 没有参数，显示产品类型选择菜单
		return h.showProductTypeMenu(message.Chat.ID)
	}

	// 有参数，按国家代码搜索
	countryCode := strings.ToUpper(args)
	return h.searchProductsByCountry(ctx, message.Chat.ID, countryCode)
}

// GetCommand 获取处理的命令名称
func (h *ProductsHandler) GetCommand() string {
	return "products"
}

// GetDescription 获取命令描述
func (h *ProductsHandler) GetDescription() string {
	return "浏览 eSIM 产品"
}

// showProductTypeMenu 显示产品类型选择菜单
func (h *ProductsHandler) showProductTypeMenu(chatID int64) error {
	text := "📱 *eSIM 产品商城*\n\n"
	text += "请选择产品类型：\n\n"
	text += "🏠 *本地* - 单个国家使用\n"
	text += "🌏 *区域* - 多个国家使用\n"
	text += "🌍 *全球* - 全球通用\n\n"
	text += "💡 提示：您也可以使用 `/products 国家代码` 搜索特定国家的产品"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 本地产品", "products_local"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌏 区域产品", "products_regional"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌍 全球产品", "products_global"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回主菜单", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
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
	text += fmt.Sprintf("📄 第 %d/%d 页 (共 %d 个产品)\n\n",
		pagination.Page, pagination.TotalPages, pagination.Total)

	for i, product := range products {
		text += fmt.Sprintf("*%d. %s*\n", i+1, product.Name)

		// 国家信息
		if len(product.Countries) > 0 {
			countries := make([]string, 0, 3)
			for j, country := range product.Countries {
				if j >= 3 {
					countries = append(countries, fmt.Sprintf("等%d国", len(product.Countries)))
					break
				}
				countries = append(countries, country.CN)
			}
			text += fmt.Sprintf("   🗺️ %s\n", strings.Join(countries, ", "))
		}

		text += fmt.Sprintf("   📊 %s | ⏰ %d天\n",
			formatDataSize(product.DataSize), product.ValidDays)
		text += fmt.Sprintf("   💰 $%.2f | 💵 $%.2f\n\n",
			product.RetailPrice, product.AgentPrice)
	}

	return text
}

// buildProductListKeyboard 构建产品列表键盘
func (h *ProductsHandler) buildProductListKeyboard(products []esim.Product, productType esim.ProductType, pagination esim.Pagination) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// 产品按钮
	for i, product := range products {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. 查看详情", i+1),
			fmt.Sprintf("product_detail:%d", product.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
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
			fmt.Sprintf("%d/%d", pagination.Page, pagination.TotalPages),
			"noop",
		))

		if pagination.Page < pagination.TotalPages {
			pageRow = append(pageRow, tgbotapi.NewInlineKeyboardButtonData(
				"➡️ 下一页",
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

	product := resp.Data
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
	text += fmt.Sprintf("找到 %d 个产品\n\n", len(resp.Message.Products))

	for i, product := range resp.Message.Products {
		text += fmt.Sprintf("*%d. %s*\n", i+1, product.Name)
		text += fmt.Sprintf("   📊 %s | ⏰ %d天\n",
			formatDataSize(product.DataSize), product.ValidDays)
		text += fmt.Sprintf("   💰 $%.2f | 💵 $%.2f\n\n",
			product.RetailPrice, product.AgentPrice)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, product := range resp.Message.Products {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. 查看详情", i+1),
			fmt.Sprintf("product_detail:%d", product.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

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

	// 忽略 "message is not modified" 错误
	if err != nil && strings.Contains(err.Error(), "message is not modified") {
		h.logger.Debug("Message content unchanged, skipping edit")
		return nil
	}

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
