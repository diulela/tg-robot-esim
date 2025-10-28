package services

import (
	"context"
	"fmt"
	"strings"

	"tg-robot-sim/storage/repository"
)

// dialogService 对话服务实现
type dialogService struct {
	sessionService SessionService
	userRepo       repository.UserRepository
	menuService    MenuService
	logger         Logger
}

// Logger 日志接口
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// NewDialogService 创建对话服务
func NewDialogService(sessionService SessionService, userRepo repository.UserRepository, menuService MenuService, logger Logger) DialogService {
	return &dialogService{
		sessionService: sessionService,
		userRepo:       userRepo,
		menuService:    menuService,
		logger:         logger,
	}
}

// ProcessMessage 处理用户消息
func (d *dialogService) ProcessMessage(ctx context.Context, userID int64, message string) (*DialogResponse, error) {
	d.logger.Debug("Processing message from user %d: %s", userID, message)

	// 获取用户上下文
	userContext, err := d.sessionService.GetUserContext(userID)
	if err != nil {
		d.logger.Error("Failed to get user context: %v", err)
		return nil, fmt.Errorf("failed to get user context: %w", err)
	}

	// 解析消息
	command, args := d.parseMessage(message)
	d.logger.Debug("Parsed command: %s, args: %v", command, args)

	// 根据当前状态和命令生成响应
	response, err := d.generateResponse(ctx, userID, userContext, command, args)
	if err != nil {
		d.logger.Error("Failed to generate response: %v", err)
		return d.createErrorResponse("处理消息时发生错误，请稍后重试"), nil
	}

	// 更新用户上下文
	if err := d.sessionService.SetUserContext(userID, userContext); err != nil {
		d.logger.Error("Failed to update user context: %v", err)
		// 不返回错误，因为响应已经生成
	}

	return response, nil
}

// GetUserContext 获取用户上下文
func (d *dialogService) GetUserContext(userID int64) (*UserContext, error) {
	return d.sessionService.GetUserContext(userID)
}

// SetUserContext 设置用户上下文
func (d *dialogService) SetUserContext(userID int64, context *UserContext) error {
	return d.sessionService.SetUserContext(userID, context)
}

// ClearUserContext 清除用户上下文
func (d *dialogService) ClearUserContext(userID int64) error {
	return d.sessionService.ClearUserContext(userID)
}

// IsUserActive 检查用户是否活跃
func (d *dialogService) IsUserActive(userID int64) bool {
	return d.sessionService.IsUserActive(userID)
}

// parseMessage 解析消息
func (d *dialogService) parseMessage(message string) (command string, args []string) {
	message = strings.TrimSpace(message)

	// 检查是否是命令
	if strings.HasPrefix(message, "/") {
		parts := strings.Fields(message)
		if len(parts) > 0 {
			command = strings.TrimPrefix(parts[0], "/")
			if len(parts) > 1 {
				args = parts[1:]
			}
		}
		return command, args
	}

	// 普通文本消息
	return "text", []string{message}
}

// generateResponse 生成响应
func (d *dialogService) generateResponse(ctx context.Context, userID int64, userContext *UserContext, command string, args []string) (*DialogResponse, error) {
	switch command {
	case "start":
		return d.handleStartCommand(ctx, userID, userContext)
	case "help":
		return d.handleHelpCommand(ctx, userID, userContext)
	case "menu":
		return d.handleMenuCommand(ctx, userID, userContext)
	case "text":
		return d.handleTextMessage(ctx, userID, userContext, strings.Join(args, " "))
	default:
		return d.handleUnknownCommand(ctx, userID, userContext, command)
	}
}

// handleStartCommand 处理 start 命令
func (d *dialogService) handleStartCommand(ctx context.Context, userID int64, userContext *UserContext) (*DialogResponse, error) {
	// 获取主菜单
	menuResponse, err := d.menuService.GetMainMenu(userID)
	if err != nil {
		d.logger.Error("Failed to get main menu: %v", err)
		return d.createErrorResponse("获取主菜单失败"), nil
	}

	// 转换为对话响应
	response := &DialogResponse{
		Message:   "🤖 欢迎使用 Telegram 机器人！\n\n" + menuResponse.Text,
		ParseMode: menuResponse.ParseMode,
		Keyboard:  menuResponse.Keyboard,
	}

	return response, nil
}

// handleHelpCommand 处理 help 命令
func (d *dialogService) handleHelpCommand(ctx context.Context, userID int64, userContext *UserContext) (*DialogResponse, error) {
	helpText := `
🤖 <b>机器人帮助</b>

<b>可用命令：</b>
/start - 开始使用机器人
/help - 显示帮助信息
/menu - 显示主菜单

<b>功能介绍：</b>
• 💬 智能对话处理
• 📱 交互式菜单系统
• 💰 区块链交易监控
• 🔔 实时通知服务

如有问题，请联系管理员。
`

	response := &DialogResponse{
		Message:   helpText,
		ParseMode: "HTML",
	}

	return response, nil
}

// handleMenuCommand 处理 menu 命令
func (d *dialogService) handleMenuCommand(ctx context.Context, userID int64, userContext *UserContext) (*DialogResponse, error) {
	// 获取主菜单
	menuResponse, err := d.menuService.GetMainMenu(userID)
	if err != nil {
		d.logger.Error("Failed to get main menu: %v", err)
		return d.createErrorResponse("获取主菜单失败"), nil
	}

	// 转换为对话响应
	response := &DialogResponse{
		Message:   menuResponse.Text,
		ParseMode: menuResponse.ParseMode,
		Keyboard:  menuResponse.Keyboard,
	}

	return response, nil
}

// handleTextMessage 处理文本消息
func (d *dialogService) handleTextMessage(ctx context.Context, userID int64, userContext *UserContext, text string) (*DialogResponse, error) {
	// 根据当前菜单状态处理文本消息
	switch userContext.CurrentMenu {
	case "":
		// 没有活跃菜单，显示帮助
		return d.handleHelpCommand(ctx, userID, userContext)
	default:
		// 在菜单中，提示使用按钮
		response := &DialogResponse{
			Message:   "请使用下方的按钮进行操作，或发送 /help 查看帮助。",
			ParseMode: "HTML",
		}
		return response, nil
	}
}

// handleUnknownCommand 处理未知命令
func (d *dialogService) handleUnknownCommand(ctx context.Context, userID int64, userContext *UserContext, command string) (*DialogResponse, error) {
	response := &DialogResponse{
		Message: fmt.Sprintf(
			"❓ 未知命令: /%s\n\n发送 /help 查看可用命令。",
			command,
		),
		ParseMode: "HTML",
	}

	return response, nil
}

// createErrorResponse 创建错误响应
func (d *dialogService) createErrorResponse(message string) *DialogResponse {
	return &DialogResponse{
		Message:   "❌ " + message,
		ParseMode: "HTML",
	}
}
