package handlers

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Registry 处理器注册表实现
type Registry struct {
	messageHandlers  []MessageHandler
	callbackHandlers []CallbackHandler
	commandHandlers  map[string]CommandHandler
	middlewares      []Middleware
	mu               sync.RWMutex
}

// NewRegistry 创建新的处理器注册表
func NewRegistry() *Registry {
	return &Registry{
		messageHandlers:  make([]MessageHandler, 0),
		callbackHandlers: make([]CallbackHandler, 0),
		commandHandlers:  make(map[string]CommandHandler),
		middlewares:      make([]Middleware, 0),
	}
}

// RegisterMessageHandler 注册消息处理器
func (r *Registry) RegisterMessageHandler(handler MessageHandler) error {
	if handler == nil {
		return fmt.Errorf("message handler cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.messageHandlers = append(r.messageHandlers, handler)
	return nil
}

// RegisterCallbackHandler 注册回调处理器
func (r *Registry) RegisterCallbackHandler(handler CallbackHandler) error {
	if handler == nil {
		return fmt.Errorf("callback handler cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.callbackHandlers = append(r.callbackHandlers, handler)
	return nil
}

// RegisterCommandHandler 注册命令处理器
func (r *Registry) RegisterCommandHandler(handler CommandHandler) error {
	if handler == nil {
		return fmt.Errorf("command handler cannot be nil")
	}

	command := handler.GetCommand()
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.commandHandlers[command]; exists {
		return fmt.Errorf("command handler for '%s' already registered", command)
	}

	r.commandHandlers[command] = handler
	return nil
}

// RouteMessage 路由消息到合适的处理器
func (r *Registry) RouteMessage(ctx context.Context, message *tgbotapi.Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// 创建处理函数
	handlerFunc := func(ctx context.Context, msg *tgbotapi.Message) error {
		// 首先检查是否是命令
		if msg.IsCommand() {
			command := msg.Command()
			if handler, exists := r.commandHandlers[command]; exists {
				return handler.HandleCommand(ctx, msg)
			}
		}

		// 然后尝试消息处理器
		for _, handler := range r.messageHandlers {
			if handler.CanHandle(msg) {
				return handler.HandleMessage(ctx, msg)
			}
		}

		return fmt.Errorf("no handler found for message")
	}

	// 应用中间件链
	return r.applyMessageMiddlewares(ctx, message, handlerFunc)
}

// RouteCallback 路由回调到合适的处理器
func (r *Registry) RouteCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	if callback == nil {
		return fmt.Errorf("callback cannot be nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// 创建处理函数
	handlerFunc := func(ctx context.Context, cb *tgbotapi.CallbackQuery) error {
		for _, handler := range r.callbackHandlers {
			if handler.CanHandle(cb) {
				return handler.HandleCallback(ctx, cb)
			}
		}
		return fmt.Errorf("no handler found for callback")
	}

	// 应用中间件链
	return r.applyCallbackMiddlewares(ctx, callback, handlerFunc)
}

// GetRegisteredCommands 获取已注册的命令列表
func (r *Registry) GetRegisteredCommands() []tgbotapi.BotCommand {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commands := make([]tgbotapi.BotCommand, 0, len(r.commandHandlers))
	for cmd, handler := range r.commandHandlers {
		commands = append(commands, tgbotapi.BotCommand{
			Command:     cmd,
			Description: handler.GetDescription(),
		})
	}

	return commands
}

// RegisterMiddleware 注册中间件
func (r *Registry) RegisterMiddleware(middleware Middleware) error {
	if middleware == nil {
		return fmt.Errorf("middleware cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.middlewares = append(r.middlewares, middleware)
	return nil
}

// applyMessageMiddlewares 应用消息中间件链
func (r *Registry) applyMessageMiddlewares(ctx context.Context, message *tgbotapi.Message, handler MessageHandlerFunc) error {
	if len(r.middlewares) == 0 {
		return handler(ctx, message)
	}

	// 构建中间件链
	var chainHandler MessageHandlerFunc
	chainHandler = handler

	// 从后往前构建链
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		middleware := r.middlewares[i]
		currentHandler := chainHandler
		chainHandler = func(ctx context.Context, msg *tgbotapi.Message) error {
			return middleware.ProcessMessage(ctx, msg, currentHandler)
		}
	}

	return chainHandler(ctx, message)
}

// applyCallbackMiddlewares 应用回调中间件链
func (r *Registry) applyCallbackMiddlewares(ctx context.Context, callback *tgbotapi.CallbackQuery, handler CallbackHandlerFunc) error {
	if len(r.middlewares) == 0 {
		return handler(ctx, callback)
	}

	// 构建中间件链
	var chainHandler CallbackHandlerFunc
	chainHandler = handler

	// 从后往前构建链
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		middleware := r.middlewares[i]
		currentHandler := chainHandler
		chainHandler = func(ctx context.Context, cb *tgbotapi.CallbackQuery) error {
			return middleware.ProcessCallback(ctx, cb, currentHandler)
		}
	}

	return chainHandler(ctx, callback)
}
