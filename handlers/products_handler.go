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

// showAsiaProducts 显示亚洲产品列表（编辑消息）
func (h *ProductsHandler) showAsiaProducts(ctx context.Context, message *tgbotapi.Message, page int) error {
	products, total, err := h.getAsiaProducts(ctx, page, 100)
	if err != nil {
		h.logger.Error("Failed to get Asia products: %v", err)
		return h.sendError(message.Chat.ID, "获取产品列表失败")
	}

	if len(products) == 0 {
		return h.sendError(message.Chat.ID, "暂无产品")
	}

	// 构建消息文本
	text := h.buildAsiaProductListText(products, page, total, 100)

	// 构建键盘
	keyboard := h.buildAsiaProductKeyboard(products, page, total, 100)

	// 编辑消息
	return h.editOrSendMessage(message, text, keyboard)
}

// showAsiaProductsNew 显示亚洲产品列表（新消息）
func (h *ProductsHandler) showAsiaProductsNew(ctx context.Context, chatID int64, page int) error {
	products, total, err := h.getAsiaProducts(ctx, page, 100)
	if err != nil {
		h.logger.Error("Failed to get Asia products: %v", err)
		return h.sendError(chatID, "获取产品列表失败")
	}

	if len(products) == 0 {
		return h.sendError(chatID, "暂无产品")
	}

	// 构建消息文本
	text := h.buildAsiaProductListText(products, page, total, 100)

	// 构建键盘
	keyboard := h.buildAsiaProductKeyboard(products, page, total, 100)

	// 发送新消息
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
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

// formatProductDetailFromAPI 格式化API返回的产品详情
func (h *ProductsHandler) formatProductDetailFromAPI(detail *esim.ProductDetail) string {
	text := fmt.Sprintf("📱 *%s*\n\n", escapeMarkdown(detail.Name))

	// 产品类型
	typeText := map[string]string{
		"local":    "🏠 本地",
		"regional": "🌏 区域",
		"global":   "🌍 全球",
	}
	if t, ok := typeText[detail.Type]; ok {
		text += fmt.Sprintf("类型: %s\n", t)
	}

	// 国家列表
	if len(detail.Countries) > 0 {
		text += "🗺️ 支持国家: "
		if len(detail.Countries) <= 5 {
			countryNames := make([]string, len(detail.Countries))
			for i, c := range detail.Countries {
				countryNames[i] = c.CN
			}
			text += strings.Join(countryNames, "、")
		} else {
			countryNames := make([]string, 5)
			for i := 0; i < 5; i++ {
				countryNames[i] = detail.Countries[i].CN
			}
			text += strings.Join(countryNames, "、")
			text += fmt.Sprintf(" 等%d个国家", len(detail.Countries))
		}
		text += "\n"
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
	text += fmt.Sprintf("📊 流量: %s\n", dataSize)
	text += fmt.Sprintf("⏰ 有效期: %d天\n", detail.ValidDays)

	// 价格（只显示零售价，单位 USDT）
	text += fmt.Sprintf("\n💰 价格: *%.2f USDT*\n", detail.Price)

	// 产品描述
	if detail.Description != "" {
		text += fmt.Sprintf("\n📝 描述:\n%s\n", detail.Description)
	}

	// 产品特性
	if len(detail.Features) > 0 {
		text += "\n✨ 特性:\n"
		for _, feature := range detail.Features {
			text += fmt.Sprintf("  • %s\n", feature)
		}
	}

	// 添加数据来源标识
	text += "\n_数据来源: 实时API_"

	return text
}

// formatProductDetailFromDetailDB 格式化产品详情消息（从产品详情表）
func (h *ProductsHandler) formatProductDetailFromDetailDB(detail *models.ProductDetail) string {
	text := fmt.Sprintf("📱 *%s*\n\n", escapeMarkdown(detail.Name))

	// 产品类型
	typeText := map[string]string{
		"local":    "🏠 本地",
		"regional": "🌏 区域",
		"global":   "🌍 全球",
	}
	if t, ok := typeText[detail.Type]; ok {
		text += fmt.Sprintf("类型: %s\n", t)
	}

	// 解析国家列表
	var countries []string
	if err := json.Unmarshal([]byte(detail.Countries), &countries); err == nil && len(countries) > 0 {
		text += "🗺️ 支持国家: "
		if len(countries) <= 5 {
			text += strings.Join(countries, "、")
		} else {
			text += strings.Join(countries[:5], "、")
			text += fmt.Sprintf(" 等%d个国家", len(countries))
		}
		text += "\n"
	}

	// 流量和有效期
	text += fmt.Sprintf("📊 流量: %s\n", detail.DataSize)
	text += fmt.Sprintf("⏰ 有效期: %d天\n", detail.ValidDays)

	// 价格（只显示零售价，单位 USDT）
	text += fmt.Sprintf("\n💰 价格: *%.2f USDT*\n", detail.Price)

	// 产品描述
	if detail.Description != "" {
		text += fmt.Sprintf("\n📝 描述:\n%s\n", detail.Description)
	}

	// 解析特性列表
	var features []string
	if err := json.Unmarshal([]byte(detail.Features), &features); err == nil && len(features) > 0 {
		text += "\n✨ 特性:\n"
		for _, feature := range features {
			text += fmt.Sprintf("  • %s\n", feature)
		}
	}

	// 添加数据来源标识和同步时间
	text += fmt.Sprintf("\n_数据来源: 本地缓存 (同步时间: %s)_",
		detail.SyncedAt.Format("01-02 15:04"))

	return text
}

// promptProductSelection 提示用户输入产品编号
func (h *ProductsHandler) promptProductSelection(ctx context.Context, message *tgbotapi.Message) error {
	text := "<b>🔍 选择产品</b>\n\n"
	text += "请回复您想查看的产品编号\n"
	text += "例如：回复 <code>1</code> 查看产品1的详情\n\n"
	text += "<i>💡 提示：直接输入数字即可</i>"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回产品列表", "products_back"),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, text)
	editMsg.ParseMode = "HTML"
	editMsg.ReplyMarkup = &keyboard

	_, err := h.bot.Send(editMsg)
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

	text := "*🌏 亚洲区域产品*\n\n"

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

	// 返回按钮
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 返回主菜单", "main_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
