package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"tg-robot-sim/storage/repository"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tg-robot-sim/config"
	"tg-robot-sim/storage/models"
)

// Database 数据库管理器
type Database struct {
	db              *gorm.DB
	config          *config.DatabaseConfig
	userRepo        repository.UserRepository
	sessionRepo     repository.UserSessionRepository
	transactionRepo repository.TransactionRepository
	productRepo     repository.ProductRepository
}

// NewDatabase 创建数据库管理器
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// 配置 GORM 日志
	logLevel := logger.Info
	if cfg.Type == "sqlite" {
		logLevel = logger.Warn // SQLite 减少日志输出
	}

	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if cfg.ConnMaxLife != "" {
		if duration, err := time.ParseDuration(cfg.ConnMaxLife); err == nil {
			sqlDB.SetConnMaxLifetime(duration)
		}
	}

	database := &Database{
		db:     db,
		config: cfg,
	}

	// 初始化仓库
	database.userRepo = repository.NewUserRepository(db)
	database.sessionRepo = repository.NewUserSessionRepository(db)
	database.transactionRepo = repository.NewTransactionRepository(db)
	database.productRepo = repository.NewProductRepository(db)

	return database, nil
}

// GetDB 获取数据库连接
func (d *Database) GetDB() *gorm.DB {
	return d.db
}

// AutoMigrate 自动迁移数据库表结构
func (d *Database) AutoMigrate() error {
	return d.db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.Transaction{},
		&models.Product{},
	)
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 检查数据库连接
func (d *Database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetUserRepository 获取用户仓库
func (d *Database) GetUserRepository() repository.UserRepository {
	return d.userRepo
}

// GetSessionRepository 获取会话仓库
func (d *Database) GetSessionRepository() repository.UserSessionRepository {
	return d.sessionRepo
}

// GetTransactionRepository 获取交易仓库
func (d *Database) GetTransactionRepository() repository.TransactionRepository {
	return d.transactionRepo
}

// GetProductRepository 获取产品仓库
func (d *Database) GetProductRepository() repository.ProductRepository {
	return d.productRepo
}

// Transaction 执行数据库事务
func (d *Database) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn)
}

// HealthCheck 健康检查
func (d *Database) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetStats 获取数据库连接统计信息
func (d *Database) GetStats() (map[string]interface{}, error) {
	sqlDB, err := d.db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}
