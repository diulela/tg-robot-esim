package esim

import (
	"encoding/json"
	"fmt"
)

// EsimStatus eSIM状态
type EsimStatus string

const (
	EsimStatusPending    EsimStatus = "pending"    // 待激活 - eSIM 已分配，等待用户激活
	EsimStatusActive     EsimStatus = "active"     // 已激活 - eSIM 已激活，可正常使用
	EsimStatusExpired    EsimStatus = "expired"    // 已过期 - eSIM 套餐已过期
	EsimStatusSuspended  EsimStatus = "suspended"  // 已暂停 - eSIM 服务已暂停
	EsimStatusTerminated EsimStatus = "terminated" // 已终止 - eSIM 服务已终止
)

// EsimUsageInfo eSIM使用信息
type EsimUsageInfo struct {
	ICCID           string     `json:"iccid"`           // ICCID号码
	Status          EsimStatus `json:"status"`          // eSIM状态
	ActivationTime  string     `json:"activationTime"`  // 激活时间
	ExpireTime      string     `json:"expireTime"`      // 过期时间
	DataUsed        int        `json:"dataUsed"`        // 已使用流量（MB）
	DataTotal       int        `json:"dataTotal"`       // 总流量（MB）
	DataRemaining   int        `json:"dataRemaining"`   // 剩余流量（MB）
	UsagePercentage string     `json:"usagePercentage"` // 使用百分比
}

// EsimUsageData eSIM使用数据
type EsimUsageData struct {
	OrderID int           `json:"orderId"` // 订单ID
	Esim    EsimUsageInfo `json:"esim"`    // eSIM使用信息
}

// EsimUsageResponse eSIM使用统计响应（兼容不同格式）
type EsimUsageResponse struct {
	Success   bool            `json:"success"`
	Code      int             `json:"code"`
	Message   json.RawMessage `json:"message"` // 可能是字符串或对象
	Data      json.RawMessage `json:"data"`    // 可能是字符串或对象
	Timestamp string          `json:"timestamp"`

	// 解析后的使用数据
	UsageData *EsimUsageData `json:"-"`
}

// TopupPackage 充值套餐
type TopupPackage struct {
	ID          string  `json:"id"`          // 套餐ID
	Title       string  `json:"title"`       // 套餐标题
	Data        string  `json:"data"`        // 流量大小
	Price       float64 `json:"price"`       // 价格（美元）
	Validity    int     `json:"validity"`    // 有效期（天）
	Description string  `json:"description"` // 套餐描述
}

// TopupPackagesData 充值套餐数据
type TopupPackagesData struct {
	OrderID  int            `json:"orderId"`  // 订单ID
	Packages []TopupPackage `json:"packages"` // 充值套餐列表
}

// TopupPackagesResponse 充值套餐响应（兼容不同格式）
type TopupPackagesResponse struct {
	Success   bool            `json:"success"`
	Code      int             `json:"code"`
	Message   json.RawMessage `json:"message"` // 可能是字符串或对象
	Data      json.RawMessage `json:"data"`    // 可能是字符串或对象
	Timestamp string          `json:"timestamp"`

	// 解析后的套餐数据
	PackagesData *TopupPackagesData `json:"-"`
}

// TopupRequest 充值请求
type TopupRequest struct {
	PackageID   string `json:"packageId"`             // 充值套餐ID（必填）
	Description string `json:"description,omitempty"` // 充值说明，默认为"代理商API充值"（可选）
}

// TopupData 充值数据
type TopupData struct {
	TopupOrderID int     `json:"topupOrderId"` // 充值订单ID
	OrderID      int     `json:"orderId"`      // 原订单ID
	PackageID    string  `json:"packageId"`    // 套餐ID
	Amount       float64 `json:"amount"`       // 充值金额
	Status       string  `json:"status"`       // 充值状态
}

// TopupResponse 充值响应（兼容不同格式）
type TopupResponse struct {
	Success   bool            `json:"success"`
	Code      int             `json:"code"`
	Message   json.RawMessage `json:"message"` // 可能是字符串或对象
	Data      json.RawMessage `json:"data"`    // 可能是字符串或对象
	Timestamp string          `json:"timestamp"`

	// 解析后的充值数据
	TopupData *TopupData `json:"-"`
}

// GetEsimUsage 获取eSIM使用统计
func (c *Client) GetEsimUsage(orderID int) (*EsimUsageResponse, error) {
	path := fmt.Sprintf("/api/v1/esims/%d/usage", orderID)

	var response EsimUsageResponse
	err := c.requestTyped("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	// 解析使用数据（兼容处理）
	if err := c.parseEsimUsageData(&response); err != nil {
		return nil, fmt.Errorf("parse esim usage data: %w", err)
	}

	return &response, nil
}

// parseEsimUsageData 解析eSIM使用数据（兼容处理）
func (c *Client) parseEsimUsageData(response *EsimUsageResponse) error {
	var usageData EsimUsageData
	var err error

	// 首先尝试从 Data 字段解析
	if len(response.Data) > 0 {
		if response.Data[0] == '{' {
			err = json.Unmarshal(response.Data, &usageData)
			if err == nil {
				response.UsageData = &usageData
				return nil
			}
		}
	}

	// 如果 Data 不是对象，尝试从 Message 字段解析
	if len(response.Message) > 0 {
		if response.Message[0] == '{' {
			err = json.Unmarshal(response.Message, &usageData)
			if err == nil {
				response.UsageData = &usageData
				return nil
			}
		}
	}

	return fmt.Errorf("no valid esim usage data found in response")
}

// GetTopupPackages 获取充值套餐
func (c *Client) GetTopupPackages(orderID int) (*TopupPackagesResponse, error) {
	path := fmt.Sprintf("/api/v1/esims/%d/topup-packages", orderID)

	var response TopupPackagesResponse
	err := c.requestTyped("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	// 解析套餐数据（兼容处理）
	if err := c.parseTopupPackagesData(&response); err != nil {
		return nil, fmt.Errorf("parse topup packages data: %w", err)
	}

	return &response, nil
}

// parseTopupPackagesData 解析充值套餐数据（兼容处理）
func (c *Client) parseTopupPackagesData(response *TopupPackagesResponse) error {
	var packagesData TopupPackagesData
	var err error

	// 首先尝试从 Data 字段解析
	if len(response.Data) > 0 {
		if response.Data[0] == '{' {
			err = json.Unmarshal(response.Data, &packagesData)
			if err == nil {
				response.PackagesData = &packagesData
				return nil
			}
		}
	}

	// 如果 Data 不是对象，尝试从 Message 字段解析
	if len(response.Message) > 0 {
		if response.Message[0] == '{' {
			err = json.Unmarshal(response.Message, &packagesData)
			if err == nil {
				response.PackagesData = &packagesData
				return nil
			}
		}
	}

	return fmt.Errorf("no valid topup packages data found in response")
}

// TopupEsim eSIM充值
func (c *Client) TopupEsim(orderID int, req TopupRequest) (*TopupResponse, error) {
	path := fmt.Sprintf("/api/v1/esims/%d/topup", orderID)

	var response TopupResponse
	err := c.requestTyped("POST", path, req, &response)
	if err != nil {
		return nil, err
	}

	// 解析充值数据（兼容处理）
	if err := c.parseTopupData(&response); err != nil {
		return nil, fmt.Errorf("parse topup data: %w", err)
	}

	return &response, nil
}

// parseTopupData 解析充值数据（兼容处理）
func (c *Client) parseTopupData(response *TopupResponse) error {
	var topupData TopupData
	var err error

	// 首先尝试从 Data 字段解析
	if len(response.Data) > 0 {
		if response.Data[0] == '{' {
			err = json.Unmarshal(response.Data, &topupData)
			if err == nil {
				response.TopupData = &topupData
				return nil
			}
		}
	}

	// 如果 Data 不是对象，尝试从 Message 字段解析
	if len(response.Message) > 0 {
		if response.Message[0] == '{' {
			err = json.Unmarshal(response.Message, &topupData)
			if err == nil {
				response.TopupData = &topupData
				return nil
			}
		}
	}

	return fmt.Errorf("no valid topup data found in response")
}
