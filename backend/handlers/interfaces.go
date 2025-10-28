package handlers

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageHandler 定义消息处理器接口
// 负责处理来自用户的文本消息和命令
type MessageHandler interface {
	// HandleMessage 处理消息
	HandleMessage(ctx context.Context, message *tgbotapi.Message) error

	// CanHandle 判断是否能处理该消息
	CanHandle(message *tgbotapi.Message) bool

	// GetHandlerName 获取处理器名称
	GetHandlerName() string
}

// CallbackHandler 定义回调处理器接口
// 负责处理内联键盘按钮点击事件
type CallbackHandler interface {
	// HandleCallback 处理回调查询
	HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error

	// CanHandle 判断是否能处理该回调
	CanHandle(callback *tgbotapi.CallbackQuery) bool

	// GetHandlerName 获取处理器名称
	GetHandlerName() string
}

// CommandHandler 定义命令处理器接口
// 负责处理特定的 Telegram 命令
type CommandHandler interface {
	// HandleCommand 处理命令
	HandleCommand(ctx context.Context, message *tgbotapi.Message) error

	// GetCommand 获取处理的命令名称
	GetCommand() string

	// GetDescription 获取命令描述
	GetDescription() string
}

// InlineQueryHandler 定义 Inline 查询处理器接口
// 负责处理 Inline 查询
type InlineQueryHandler interface {
	// HandleInlineQuery 处理 Inline 查询
	HandleInlineQuery(ctx context.Context, query *tgbotapi.InlineQuery) error

	// GetHandlerName 获取处理器名称
	GetHandlerName() string
}

// HandlerRegistry 定义处理器注册表接口
// 负责管理所有处理器的注册和路由
type HandlerRegistry interface {
	// RegisterMessageHandler 注册消息处理器
	RegisterMessageHandler(handler MessageHandler) error

	// RegisterCallbackHandler 注册回调处理器
	RegisterCallbackHandler(handler CallbackHandler) error

	// RegisterCommandHandler 注册命令处理器
	RegisterCommandHandler(handler CommandHandler) error

	// RegisterInlineHandler 注册 Inline 查询处理器
	RegisterInlineHandler(handler InlineQueryHandler) error

	// RouteMessage 路由消息到合适的处理器
	RouteMessage(ctx context.Context, message *tgbotapi.Message) error

	// RouteCallback 路由回调到合适的处理器
	RouteCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error

	// RouteInlineQuery 路由 Inline 查询到合适的处理器
	RouteInlineQuery(ctx context.Context, query *tgbotapi.InlineQuery) error
}
