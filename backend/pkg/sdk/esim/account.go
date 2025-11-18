package esim

import "encoding/json"

// AccountStatus 账户状态
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"   // 活跃
	AccountStatusInactive AccountStatus = "inactive" // 未激活
)

// AccountInfo 账户基本信息
type AccountInfo struct {
	ID          int           `json:"id"`          // 代理商ID
	Name        string        `json:"name"`        // 代理商名称
	Email       string        `json:"email"`       // 联系邮箱
	Phone       string        `json:"phone"`       // 联系电话
	Level       int           `json:"level"`       // 代理商等级（1-3）
	Status      AccountStatus `json:"status"`      // 账户状态
	Discount    float64       `json:"discount"`    // 折扣率（0-1之间）
	MemberSince string        `json:"memberSince"` // 注册时间（ISO格式）
}

// AccountInfoResponse 账户信息响应（兼容不同格式）
type AccountInfoResponse struct {
	Success   bool            `json:"success"`
	Code      int             `json:"code"`
	Message   json.RawMessage `json:"message"` // 可能是字符串或对象
	Data      json.RawMessage `json:"data"`    // 可能是字符串或对象
	Timestamp string          `json:"timestamp"`

	// 解析后的账户信息
	AccountInfo *AccountInfo `json:"-"`
}

// BalanceInfo 账户余额信息
type BalanceInfo struct {
	Balance         float64 `json:"balance"`         // 当前账户余额
	TotalCommission float64 `json:"totalCommission"` // 累计佣金收入
	TotalWithdraw   float64 `json:"totalWithdraw"`   // 累计提现金额
	MonthSales      float64 `json:"monthSales"`      // 本月销售额
	TotalSales      float64 `json:"totalSales"`      // 累计销售额
	Level           int     `json:"level"`           // 代理商等级
	Discount        float64 `json:"discount"`        // 折扣率（0-1之间）
	UpdatedAt       string  `json:"updatedAt"`       // 更新时间
}

// BalanceResponse 账户余额响应（兼容不同格式）
type BalanceResponse struct {
	Success   bool            `json:"success"`
	Code      int             `json:"code"`
	Message   json.RawMessage `json:"message"` // 可能是字符串或对象
	Data      json.RawMessage `json:"data"`    // 可能是字符串或对象
	Timestamp string          `json:"timestamp"`

	// 解析后的余额信息
	BalanceInfo *BalanceInfo `json:"-"`
}

// GetAccount 获取账户基本信息
// 获取代理商账户基本信息，包括名称、等级、联系方式等
func (c *Client) GetAccount() (*AccountInfoResponse, error) {
	var response AccountInfoResponse
	err := c.requestTyped("GET", "/api/v1/account", nil, &response)
	if err != nil {
		return nil, err
	}

	// 解析账户信息（兼容处理）
	if err := c.parseAccountInfo(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// parseAccountInfo 解析账户信息（兼容处理）
func (c *Client) parseAccountInfo(response *AccountInfoResponse) error {
	var accountInfo AccountInfo
	var err error

	// 首先尝试从 Message 字段解析
	if len(response.Message) > 0 {
		if response.Message[0] == '{' {
			err = json.Unmarshal(response.Message, &accountInfo)
			if err == nil {
				response.AccountInfo = &accountInfo
				return nil
			}
		}
	}

	// 如果 Message 不是对象，尝试从 Data 字段解析
	if len(response.Data) > 0 {
		if response.Data[0] == '{' {
			err = json.Unmarshal(response.Data, &accountInfo)
			if err == nil {
				response.AccountInfo = &accountInfo
				return nil
			}
		}
	}

	return nil // 账户信息可能为空，不返回错误
}

// GetBalance 获取账户余额
// 获取账户余额、佣金、销售额等财务信息
func (c *Client) GetBalance() (*BalanceResponse, error) {
	var response BalanceResponse
	err := c.requestTyped("GET", "/api/v1/account/balance", nil, &response)
	if err != nil {
		return nil, err
	}

	// 解析余额信息（兼容处理）
	if err := c.parseBalanceInfo(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// parseBalanceInfo 解析余额信息（兼容处理）
func (c *Client) parseBalanceInfo(response *BalanceResponse) error {
	var balanceInfo BalanceInfo
	var err error

	// 首先尝试从 Message 字段解析
	if len(response.Message) > 0 {
		if response.Message[0] == '{' {
			err = json.Unmarshal(response.Message, &balanceInfo)
			if err == nil {
				response.BalanceInfo = &balanceInfo
				return nil
			}
		}
	}

	// 如果 Message 不是对象，尝试从 Data 字段解析
	if len(response.Data) > 0 {
		if response.Data[0] == '{' {
			err = json.Unmarshal(response.Data, &balanceInfo)
			if err == nil {
				response.BalanceInfo = &balanceInfo
				return nil
			}
		}
	}

	return nil // 余额信息可能为空，不返回错误
}
