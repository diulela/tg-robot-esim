package dto

import (
	"tg-robot-sim/storage/models"
)

// WalletHistoryWithUser 钱包历史及用户信息
// 用于返回钱包历史记录时包含用户信息
type WalletHistoryWithUser struct {
	models.WalletHistory
	User *models.User `json:"user,omitempty"`
}
