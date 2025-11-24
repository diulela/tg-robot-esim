package dto

import (
	"tg-robot-sim/storage/models"
)

// OrderWithRelations 订单及其关联数据
// 用于返回订单详情时包含用户、产品和订单详情信息
type OrderWithRelations struct {
	models.Order
	User    *models.User    `json:"user,omitempty"`
	Product *models.Product `json:"product,omitempty"`
}

// OrderListItem 订单列表项（包含产品信息）
// 用于返回订单列表时包含产品信息（不包含用户信息以减少数据量）
type OrderListItem struct {
	models.Order
	Product *models.Product `json:"product,omitempty"`
}
