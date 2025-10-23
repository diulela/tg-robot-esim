package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// SessionService 会话服务接口
type SessionService interface {
	// GetUserContext 获取用户上下文
	GetUserContext(userID int64) (*UserContext, error)

	// SetUserContext 设置用户上下文
	SetUserContext(userID int64, context *UserContext) error

	// ClearUserContext 清除用户上下文
	ClearUserContext(userID int64) error

	// IsUserActive 检查用户是否活跃
	IsUserActive(userID int64) bool

	// CleanupExpiredSessions 清理过期会话
	CleanupExpiredSessions() error

	// GetActiveSessionCount 获取活跃会话数量
	GetActiveSessionCount() int
}

// sessionService 会话服务实现
type sessionService struct {
	sessionRepo repository.UserSessionRepository
	cache       map[int64]*UserContext
	cacheMutex  sync.RWMutex
	timeout     time.Duration
}

// NewSessionService 创建会话服务
func NewSessionService(sessionRepo repository.UserSessionRepository, timeout time.Duration) SessionService {
	service := &sessionService{
		sessionRepo: sessionRepo,
		cache:       make(map[int64]*UserContext),
		timeout:     timeout,
	}

	// 启动定期清理任务
	go service.startCleanupTask()

	return service
}

// GetUserContext 获取用户上下文
func (s *sessionService) GetUserContext(userID int64) (*UserContext, error) {
	// 首先尝试从缓存获取
	s.cacheMutex.RLock()
	if ctx, exists := s.cache[userID]; exists {
		// 检查是否过期
		if time.Since(ctx.LastActive) <= s.timeout {
			s.cacheMutex.RUnlock()
			// 更新最后活跃时间
			ctx.LastActive = time.Now()
			return ctx, nil
		}
		// 过期了，从缓存中删除
		s.cacheMutex.RUnlock()
		s.cacheMutex.Lock()
		delete(s.cache, userID)
		s.cacheMutex.Unlock()
	} else {
		s.cacheMutex.RUnlock()
	}

	// 从数据库获取
	session, err := s.sessionRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新的用户上下文
			ctx := &UserContext{
				UserID:      userID,
				CurrentMenu: "",
				MenuPath:    []string{},
				Parameters:  make(map[string]interface{}),
				LastActive:  time.Now(),
			}
			return ctx, nil
		}
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	// 检查数据库中的会话是否过期
	if session.IsExpired(s.timeout) {
		// 删除过期会话
		_ = s.sessionRepo.Delete(context.Background(), userID)
		// 返回新的上下文
		ctx := &UserContext{
			UserID:      userID,
			CurrentMenu: "",
			MenuPath:    []string{},
			Parameters:  make(map[string]interface{}),
			LastActive:  time.Now(),
		}
		return ctx, nil
	}

	// 反序列化会话数据
	ctx, err := s.deserializeUserContext(session)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize user context: %w", err)
	}

	// 更新最后活跃时间
	ctx.LastActive = time.Now()

	// 缓存到内存
	s.cacheMutex.Lock()
	s.cache[userID] = ctx
	s.cacheMutex.Unlock()

	return ctx, nil
}

// SetUserContext 设置用户上下文
func (s *sessionService) SetUserContext(userID int64, context *UserContext) error {
	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	context.UserID = userID
	context.LastActive = time.Now()

	// 更新缓存
	s.cacheMutex.Lock()
	s.cache[userID] = context
	s.cacheMutex.Unlock()

	// 序列化并保存到数据库
	session, err := s.serializeUserContext(context)
	if err != nil {
		return fmt.Errorf("failed to serialize user context: %w", err)
	}

	// 尝试更新现有会话
	existingSession, err := s.sessionRepo.GetByUserID(context.Background(), userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing session: %w", err)
	}

	if existingSession != nil {
		// 更新现有会话
		session.ID = existingSession.ID
		session.CreatedAt = existingSession.CreatedAt
		return s.sessionRepo.Update(context.Background(), session)
	} else {
		// 创建新会话
		return s.sessionRepo.Create(context.Background(), session)
	}
}

// ClearUserContext 清除用户上下文
func (s *sessionService) ClearUserContext(userID int64) error {
	// 从缓存中删除
	s.cacheMutex.Lock()
	delete(s.cache, userID)
	s.cacheMutex.Unlock()

	// 从数据库中删除
	return s.sessionRepo.Delete(context.Background(), userID)
}

// IsUserActive 检查用户是否活跃
func (s *sessionService) IsUserActive(userID int64) bool {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	if ctx, exists := s.cache[userID]; exists {
		return time.Since(ctx.LastActive) <= s.timeout
	}

	return false
}

// CleanupExpiredSessions 清理过期会话
func (s *sessionService) CleanupExpiredSessions() error {
	// 清理内存缓存中的过期会话
	s.cacheMutex.Lock()
	for userID, ctx := range s.cache {
		if time.Since(ctx.LastActive) > s.timeout {
			delete(s.cache, userID)
		}
	}
	s.cacheMutex.Unlock()

	// 清理数据库中的过期会话
	return s.sessionRepo.DeleteExpired(context.Background(), s.timeout)
}

// GetActiveSessionCount 获取活跃会话数量
func (s *sessionService) GetActiveSessionCount() int {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	count := 0
	for _, ctx := range s.cache {
		if time.Since(ctx.LastActive) <= s.timeout {
			count++
		}
	}

	return count
}

// serializeUserContext 序列化用户上下文到数据库模型
func (s *sessionService) serializeUserContext(ctx *UserContext) (*models.UserSession, error) {
	parametersJSON, err := json.Marshal(ctx.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	menuPathJSON, err := json.Marshal(ctx.MenuPath)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal menu path: %w", err)
	}

	return &models.UserSession{
		UserID:      ctx.UserID,
		SessionData: string(parametersJSON),
		MenuPath:    string(menuPathJSON),
		CurrentMenu: ctx.CurrentMenu,
		Parameters:  string(parametersJSON),
		LastActive:  ctx.LastActive,
	}, nil
}

// deserializeUserContext 从数据库模型反序列化用户上下文
func (s *sessionService) deserializeUserContext(session *models.UserSession) (*UserContext, error) {
	var parameters map[string]interface{}
	if session.Parameters != "" {
		if err := json.Unmarshal([]byte(session.Parameters), &parameters); err != nil {
			return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
		}
	} else {
		parameters = make(map[string]interface{})
	}

	var menuPath []string
	if session.MenuPath != "" {
		if err := json.Unmarshal([]byte(session.MenuPath), &menuPath); err != nil {
			return nil, fmt.Errorf("failed to unmarshal menu path: %w", err)
		}
	} else {
		menuPath = []string{}
	}

	return &UserContext{
		UserID:      session.UserID,
		CurrentMenu: session.CurrentMenu,
		MenuPath:    menuPath,
		Parameters:  parameters,
		LastActive:  session.LastActive,
	}, nil
}

// startCleanupTask 启动定期清理任务
func (s *sessionService) startCleanupTask() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		if err := s.CleanupExpiredSessions(); err != nil {
			// 这里应该使用日志记录错误，暂时忽略
			continue
		}
	}
}
