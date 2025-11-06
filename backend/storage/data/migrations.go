package data

import (
	"fmt"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	// 迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductDetail{},
		&models.Transaction{},
		&models.UserSession{},
		&models.Wallet{},
		&models.Order{},
		&models.OrderDetail{},
		&models.RechargeOrder{},
		&models.WalletHistory{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// MigrateWallets 为现有用户创建钱包
func MigrateWallets(db *gorm.DB) error {
	// 查找所有没有钱包的用户
	var users []models.User
	err := db.Where("id NOT IN (SELECT user_id FROM wallets)").Find(&users).Error
	if err != nil {
		return fmt.Errorf("failed to find users without wallets: %w", err)
	}

	// 为每个用户创建钱包
	for _, user := range users {
		wallet := models.Wallet{
			UserID:        user.TelegramID,
			Balance:       "0",
			FrozenBalance: "0",
			TotalIncome:   "0",
			TotalExpense:  "0",
		}

		if err := db.Create(&wallet).Error; err != nil {
			return fmt.Errorf("failed to create wallet for user %d: %w", user.ID, err)
		}
	}

	return nil
}
