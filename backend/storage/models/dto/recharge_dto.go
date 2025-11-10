package dto

import (
	"tg-robot-sim/storage/models"
)

// RechargeOrderWithUser 充值订单及用户信息
// 用于返回充值订单详情时包含用户信息
type RechargeOrderWithUser struct {
	models.RechargeOrder
	User *models.User `json:"user,omitempty"`
}
