package bot

import (
	"context"
	"fmt"
	"runtime"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/config"
	"tg-robot-sim/handlers"
	"tg-robot-sim/pkg/logger"
)

// Bot Telegram 机器人管理器
type Bot struct {
	api          *tgbotapi.BotAPI
	config       *config.TelegramConfig
	logger       *logger.Logger
	registry     *handlers.Registry
	updates      tgbotapi.UpdatesChannel
	stopChan     chan struct{}
	errorHandler *ErrorHandler
}

// NewBot 创建新的机器人实例
func NewBot(cfg *config.TelegramConfig, log *logger.Logger) (*Bot, error) {
	// 创建 Telegram Bot API 实例
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = cfg.Debug

	// 创建处理器注册表
	registry := handlers.NewRegistry()

	// 创建错误处理器
	errorHandler := NewErrorHandler(log)

	bot := &Bot{
		api:          api,
		config:       cfg,
		logger:       log,
		registry:     registry,
		stopChan:     make(chan struct{}),
		errorHandler: errorHandler,
	}

	log.Info("Bot created successfully: @%s", api.Self.UserName)
	return bot, nil
}

// GetAPI 获取 Telegram Bot API 实例
func (b *Bot) GetAPI() *tgbotapi.BotAPI {
	return b.api
}

// GetRegistry 获取处理器注册表
func (b *Bot) GetRegistry() *handlers.Registry {
	return b.registry
}

// Start 启动机器人
func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("Starting bot...")

	// 设置 Webhook 或长轮询
	if b.config.WebhookURL != "" {
		return b.startWebhook(ctx)
	} else {
		return b.startPolling(ctx)
	}
}

// startPolling 启动长轮询模式
func (b *Bot) startPolling(ctx context.Context) error {
	b.logger.Info("Starting polling mode...")

	// 删除现有的 webhook
	_, err := b.api.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		b.logger.Warn("Failed to delete webhook: %v", err)
	}

	// 配置更新
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(b.config.Timeout.ToDuration().Seconds())

	// 获取更新通道
	b.updates = b.api.GetUpdatesChan(updateConfig)

	// 设置机器人命令
	if err := b.setupBotCommands(); err != nil {
		b.logger.Warn("Failed to setup bot commands: %v", err)
	}

	// 处理更新
	go b.handleUpdates(ctx)

	b.logger.Info("Bot started in polling mode")
	return nil
}

// startWebhook 启动 Webhook 模式
func (b *Bot) startWebhook(ctx context.Context) error {
	b.logger.Info("Starting webhook mode...")

	// 设置 webhook
	webhookConfig, _ := tgbotapi.NewWebhook(b.config.WebhookURL)

	_, err := b.api.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	// 设置机器人命令
	if err := b.setupBotCommands(); err != nil {
		b.logger.Warn("Failed to setup bot commands: %v", err)
	}

	b.logger.Info("Bot started in webhook mode")
	return nil
}

// Stop 停止机器人
func (b *Bot) Stop() {
	b.logger.Info("Stopping bot...")

	close(b.stopChan)

	if b.updates != nil {
		b.api.StopReceivingUpdates()
	}

	b.logger.Info("Bot stopped")
}

// setupBotCommands 设置机器人命令
func (b *Bot) setupBotCommands() error {
	commands := b.registry.GetRegisteredCommands()
	if len(commands) == 0 {
		return nil
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := b.api.Request(config)
	if err != nil {
		return fmt.Errorf("failed to set bot commands: %w", err)
	}

	b.logger.Info("Set %d bot commands", len(commands))
	return nil
}

// handleUpdates 处理更新
func (b *Bot) handleUpdates(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			b.logger.Info("Context cancelled, stopping update handler")
			return
		case <-b.stopChan:
			b.logger.Info("Stop signal received, stopping update handler")
			return
		case update := <-b.updates:
			go b.processUpdate(ctx, update)
		}
	}
}

// processUpdate 处理单个更新
func (b *Bot) processUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			// 获取详细的堆栈跟踪信息
			stack := make([]byte, 4096)
			length := runtime.Stack(stack, false)
			b.logger.Error("Panic in update processing: %v\nStack trace:\n%s", r, string(stack[:length]))

			// 打印更新内容以便调试
			b.logger.Error("Update that caused panic: %+v", update)
		}
	}()

	// 处理消息
	if update.Message != nil {
		if err := b.registry.RouteMessage(ctx, update.Message); err != nil {
			b.logger.Error("Failed to route message: %v", err)
			b.sendErrorMessage(update.Message.Chat.ID, "处理消息时发生错误，请稍后重试")
		}
		return
	}

	// 处理回调查询
	if update.CallbackQuery != nil {
		if err := b.registry.RouteCallback(ctx, update.CallbackQuery); err != nil {
			b.logger.Error("Failed to route callback: %v", err)
			b.answerCallbackQuery(update.CallbackQuery.ID, "处理请求时发生错误")
		}
		return
	}

	// 处理 Inline 查询
	if update.InlineQuery != nil {
		if err := b.registry.RouteInlineQuery(ctx, update.InlineQuery); err != nil {
			b.logger.Error("Failed to route inline query: %v", err)
		}
		return
	}

	// 其他类型的更新暂时忽略
	b.logger.Debug("Received unhandled update type")
}

// sendErrorMessage 发送错误消息
func (b *Bot) sendErrorMessage(chatID int64, message string) {
	if err := b.SendMessage(chatID, message); err != nil {
		b.logger.Error("Failed to send error message: %v", err)
	}
}

// answerCallbackQuery 回答回调查询
func (b *Bot) answerCallbackQuery(callbackQueryID, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := b.errorHandler.HandleAPIRequest(ctx, func() (tgbotapi.APIResponse, error) {
		callback := tgbotapi.NewCallback(callbackQueryID, text)
		resp, err := b.api.Request(callback)
		if err != nil {
			return tgbotapi.APIResponse{}, err
		}
		return *resp, nil
	})

	if err != nil {
		b.logger.Error("Failed to answer callback query: %v", err)
	}
}

// SendMessage 发送消息的便捷方法
func (b *Bot) SendMessage(chatID int64, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := b.errorHandler.HandleAPICall(ctx, func() (tgbotapi.Message, error) {
		msg := tgbotapi.NewMessage(chatID, text)
		return b.api.Send(msg)
	})
	return err
}

// SendMessageWithKeyboard 发送带键盘的消息
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := b.errorHandler.HandleAPICall(ctx, func() (tgbotapi.Message, error) {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = keyboard
		return b.api.Send(msg)
	})
	return err
}
