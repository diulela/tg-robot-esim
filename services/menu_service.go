package services

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MenuItem 菜单项结构
type MenuItem struct {
	ID          string     `json:"id"`
	Text        string     `json:"text"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Action      string     `json:"action"`
	SubMenus    []MenuItem `json:"sub_menus,omitempty"`
	Permission  string     `json:"permission,omitempty"`
}

// MenuDefinition 菜单定义
type MenuDefinition struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Items       []MenuItem `json:"items"`
}

// menuService 菜单服务实现
type menuService struct {
	sessionService SessionService
	menus          map[string]*MenuDefinition
	logger         Logger
}

// NewMenuService 创建菜单服务
func NewMenuService(sessionService SessionService, logger Logger) MenuService {
	service := &menuService{
		sessionService: sessionService,
		menus:          make(map[string]*MenuDefinition),
		logger:         logger,
	}

	// 初始化默认菜单
	service.initializeDefaultMenus()
	return service
}

// GetMainMenu 获取主菜单
func (m *menuService) GetMainMenu(userID int64) (*MenuResponse, error) {
	mainMenu := m.menus["main"]
	if mainMenu == nil {
		return nil, fmt.Errorf("main menu not found")
	}

	// 更新用户上下文
	userContext, err := m.sessionService.GetUserContext(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	userContext.CurrentMenu = "main"
	userContext.MenuPath = []string{"main"}

	if err := m.sessionService.SetUserContext(userID, userContext); err != nil {
		m.logger.Warn("Failed to update user context: %v", err)
	}

	return m.buildMenuResponse(mainMenu), nil
}

// HandleMenuAction 处理菜单操作
func (m *menuService) HandleMenuAction(ctx context.Context, userID int64, action string) (*MenuResponse, error) {
	m.logger.Debug("Handling menu action: %s for user %d", action, userID)

	// 获取用户上下文
	userContext, err := m.sessionService.GetUserContext(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	switch action {
	case "main_menu":
		return m.GetMainMenu(userID)
	case "help":
		return m.getHelpMenu(userID)
	case "settings":
		return m.getSettingsMenu(userID)
	case "back":
		return m.NavigateBack(userID)
	default:
		return m.handleCustomAction(ctx, userID, userContext, action)
	}
}

// GetMenuHistory 获取菜单历史
func (m *menuService) GetMenuHistory(userID int64) ([]string, error) {
	userContext, err := m.sessionService.GetUserContext(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	return userContext.MenuPath, nil
}

// NavigateBack 返回上级菜单
func (m *menuService) NavigateBack(userID int64) (*MenuResponse, error) {
	userContext, err := m.sessionService.GetUserContext(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	// 如果已经在主菜单，返回主菜单
	if len(userContext.MenuPath) <= 1 {
		return m.GetMainMenu(userID)
	}

	// 移除当前菜单，返回上级菜单
	userContext.MenuPath = userContext.MenuPath[:len(userContext.MenuPath)-1]
	parentMenuID := userContext.MenuPath[len(userContext.MenuPath)-1]
	userContext.CurrentMenu = parentMenuID

	if err := m.sessionService.SetUserContext(userID, userContext); err != nil {
		m.logger.Warn("Failed to update user context: %v", err)
	}

	// 获取父菜单
	parentMenu := m.menus[parentMenuID]
	if parentMenu == nil {
		return m.GetMainMenu(userID)
	}

	return m.buildMenuResponse(parentMenu), nil
}

// ResetMenu 重置菜单到主菜单
func (m *menuService) ResetMenu(userID int64) (*MenuResponse, error) {
	return m.GetMainMenu(userID)
}

// initializeDefaultMenus 初始化默认菜单
func (m *menuService) initializeDefaultMenus() {
	// 主菜单
	mainMenu := &MenuDefinition{
		ID:          "main",
		Title:       "📱 主菜单",
		Description: "请选择您需要的功能：",
		Items: []MenuItem{
			{
				ID:          "products",
				Text:        "🛍️ 浏览产品",
				Description: "浏览和购买 eSIM 产品",
				Icon:        "🛍️",
				Action:      "products_back",
			},
			{
				ID:          "orders",
				Text:        "📦 我的订单",
				Description: "查看订单和 eSIM 信息",
				Icon:        "📦",
				Action:      "my_orders",
			},
			{
				ID:          "wallet",
				Text:        "💰 钱包管理",
				Description: "查看余额和交易记录",
				Icon:        "💰",
				Action:      "wallet_menu",
			},
			{
				ID:          "settings",
				Text:        "⚙️ 设置",
				Description: "个人设置和偏好",
				Icon:        "⚙️",
				Action:      "settings_menu",
			},
			{
				ID:          "help",
				Text:        "ℹ️ 帮助",
				Description: "使用帮助和常见问题",
				Icon:        "ℹ️",
				Action:      "help_menu",
			},
		},
	}

	// 设置菜单
	settingsMenu := &MenuDefinition{
		ID:          "settings",
		Title:       "⚙️ 设置",
		Description: "个人设置和偏好配置：",
		Items: []MenuItem{
			{
				ID:          "language",
				Text:        "🌐 语言设置",
				Description: "选择界面语言",
				Icon:        "🌐",
				Action:      "language_settings",
			},
			{
				ID:          "notifications",
				Text:        "🔔 通知设置",
				Description: "配置通知偏好",
				Icon:        "🔔",
				Action:      "notification_settings",
			},
			{
				ID:          "back",
				Text:        "🔙 返回主菜单",
				Description: "返回到主菜单",
				Icon:        "🔙",
				Action:      "main_menu",
			},
		},
	}

	m.menus["main"] = mainMenu
	m.menus["settings"] = settingsMenu
}

// buildMenuResponse 构建菜单响应
func (m *menuService) buildMenuResponse(menu *MenuDefinition) *MenuResponse {
	// 创建内联键盘
	var rows [][]tgbotapi.InlineKeyboardButton

	// 每行最多2个按钮
	for i := 0; i < len(menu.Items); i += 2 {
		var row []tgbotapi.InlineKeyboardButton

		// 添加第一个按钮
		item := menu.Items[i]
		button := tgbotapi.NewInlineKeyboardButtonData(
			item.Icon+" "+item.Text,
			item.Action,
		)
		row = append(row, button)

		// 如果还有下一个按钮，添加到同一行
		if i+1 < len(menu.Items) {
			nextItem := menu.Items[i+1]
			nextButton := tgbotapi.NewInlineKeyboardButtonData(
				nextItem.Icon+" "+nextItem.Text,
				nextItem.Action,
			)
			row = append(row, nextButton)
		}

		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &MenuResponse{
		Text:      menu.Title + "\n\n" + menu.Description,
		Keyboard:  keyboard,
		ParseMode: "HTML",
		EditMode:  true,
	}
}

// getHelpMenu 获取帮助菜单
func (m *menuService) getHelpMenu(userID int64) (*MenuResponse, error) {
	helpText := `
ℹ️ <b>帮助信息</b>

<b>主要功能：</b>
• 💰 钱包管理 - 查看余额和交易记录
• 📊 交易监控 - 实时监控区块链交易
• ⚙️ 设置 - 个人偏好配置
• ℹ️ 帮助 - 使用说明和支持

<b>使用提示：</b>
• 点击按钮进行操作
• 使用 /start 重新开始
• 使用 /help 查看帮助

如有问题，请联系管理员。
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回主菜单", "main_menu"),
		),
	)

	return &MenuResponse{
		Text:      helpText,
		Keyboard:  keyboard,
		ParseMode: "HTML",
		EditMode:  true,
	}, nil
}

// getSettingsMenu 获取设置菜单
func (m *menuService) getSettingsMenu(userID int64) (*MenuResponse, error) {
	settingsMenu := m.menus["settings"]
	if settingsMenu == nil {
		return nil, fmt.Errorf("settings menu not found")
	}

	// 更新用户上下文
	userContext, err := m.sessionService.GetUserContext(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	userContext.CurrentMenu = "settings"
	userContext.MenuPath = append(userContext.MenuPath, "settings")

	if err := m.sessionService.SetUserContext(userID, userContext); err != nil {
		m.logger.Warn("Failed to update user context: %v", err)
	}

	return m.buildMenuResponse(settingsMenu), nil
}

// handleCustomAction 处理自定义操作
func (m *menuService) handleCustomAction(ctx context.Context, userID int64, userContext *UserContext, action string) (*MenuResponse, error) {
	switch action {
	case "wallet_menu":
		return m.getWalletMenu(userID)
	case "transactions_menu":
		return m.getTransactionsMenu(userID)
	case "language_settings":
		return m.getLanguageSettings(userID)
	case "notification_settings":
		return m.getNotificationSettings(userID)
	default:
		m.logger.Warn("Unknown menu action: %s", action)
		return m.GetMainMenu(userID)
	}
}

// getWalletMenu 获取钱包菜单（占位符）
func (m *menuService) getWalletMenu(userID int64) (*MenuResponse, error) {
	text := "💰 <b>钱包管理</b>\n\n功能开发中，敬请期待..."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回主菜单", "main_menu"),
		),
	)

	return &MenuResponse{
		Text:      text,
		Keyboard:  keyboard,
		ParseMode: "HTML",
		EditMode:  true,
	}, nil
}

// getTransactionsMenu 获取交易菜单（占位符）
func (m *menuService) getTransactionsMenu(userID int64) (*MenuResponse, error) {
	text := "📊 <b>交易监控</b>\n\n功能开发中，敬请期待..."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回主菜单", "main_menu"),
		),
	)

	return &MenuResponse{
		Text:      text,
		Keyboard:  keyboard,
		ParseMode: "HTML",
		EditMode:  true,
	}, nil
}

// getLanguageSettings 获取语言设置（占位符）
func (m *menuService) getLanguageSettings(userID int64) (*MenuResponse, error) {
	text := "🌐 <b>语言设置</b>\n\n功能开发中，敬请期待..."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回设置", "settings_menu"),
		),
	)

	return &MenuResponse{
		Text:      text,
		Keyboard:  keyboard,
		ParseMode: "HTML",
		EditMode:  true,
	}, nil
}

// getNotificationSettings 获取通知设置（占位符）
func (m *menuService) getNotificationSettings(userID int64) (*MenuResponse, error) {
	text := "🔔 <b>通知设置</b>\n\n功能开发中，敬请期待..."
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 返回设置", "settings_menu"),
		),
	)

	return &MenuResponse{
		Text:      text,
		Keyboard:  keyboard,
		ParseMode: "HTML",
		EditMode:  true,
	}, nil
}
