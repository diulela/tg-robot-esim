package esim

import "fmt"

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // 待付款 - 订单已创建，等待付款
	OrderStatusPaid       OrderStatus = "paid"       // 已付款 - 订单已付款，处理中
	OrderStatusProcessing OrderStatus = "processing" // 处理中 - 正在为您准备 eSIM
	OrderStatusCompleted  OrderStatus = "completed"  // 已完成 - eSIM 已准备就绪
	OrderStatusCancelled  OrderStatus = "cancelled"  // 已取消 - 订单已取消
	OrderStatusFailed     OrderStatus = "failed"     // 失败 - 订单处理失败
)

// OrderItem 订单项
type OrderItem struct {
	ID          int         `json:"id"`          // 订单项ID
	ProductID   int         `json:"productId"`   // 产品ID
	ProductName string      `json:"productName"` // 产品名称
	ProductType ProductType `json:"productType"` // 产品类型
	Quantity    int         `json:"quantity"`    // 数量
	UnitPrice   float64     `json:"unitPrice"`   // 单价
	Subtotal    float64     `json:"subtotal"`    // 小计
	DataSize    int         `json:"dataSize"`    // 流量大小（MB）
	ValidDays   int         `json:"validDays"`   // 有效天数
	Countries   []Country   `json:"countries"`   // 支持的国家列表
}

// OrderEsim 订单中的eSIM信息
type OrderEsim struct {
	ID                int    `json:"id"`                // eSIM ID
	ICCID             string `json:"iccid"`             // ICCID号码
	Status            string `json:"status"`            // eSIM状态
	HasActivationCode bool   `json:"hasActivationCode"` // 是否有激活码
	HasQrCode         bool   `json:"hasQrCode"`         // 是否有二维码
}

// Order 订单信息
type Order struct {
	ID             int         `json:"id"`             // 订单ID
	OrderNumber    string      `json:"orderNumber"`    // 订单编号
	Status         OrderStatus `json:"status"`         // 订单状态
	TotalAmount    float64     `json:"totalAmount"`    // 订单总金额（美元）
	PayAmount      float64     `json:"payAmount"`      // 实付金额（美元）
	PlatformProfit float64     `json:"platformProfit"` // 平台利润（美元）
	OrderItems     []OrderItem `json:"orderItems"`     // 订单项列表
	Esims          []OrderEsim `json:"esims"`          // eSIM列表
	CreatedAt      string      `json:"createdAt"`      // 创建时间
	UpdatedAt      string      `json:"updatedAt"`      // 更新时间
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ProductID     int    `json:"productId"`               // 产品ID（必填）
	CustomerEmail string `json:"customerEmail"`           // 客户邮箱地址（必填）
	CustomerPhone string `json:"customerPhone,omitempty"` // 客户手机号（可选）
	Quantity      int    `json:"quantity,omitempty"`      // 购买数量，默认为1（可选）
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		OrderID     int         `json:"orderId"`     // 订单ID
		OrderNumber string      `json:"orderNumber"` // 订单编号
		TotalAmount float64     `json:"totalAmount"` // 订单总金额
		PayAmount   float64     `json:"payAmount"`   // 实付金额
		Status      OrderStatus `json:"status"`      // 订单状态
	} `json:"data"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Success   bool             `json:"success"`
	Code      int              `json:"code"`
	Message   OrderListMessage `json:"message"`
	Data      string           `json:"data"`
	Timestamp string           `json:"timestamp"`
}

// OrderListMessage 订单列表消息
type OrderListMessage struct {
	Orders     []Order    `json:"orders"`
	Pagination Pagination `json:"pagination"`
}

// OrderDetailResponse 订单详情响应
type OrderDetailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Order  `json:"data"`
}

// OrderParams 订单查询参数
type OrderParams struct {
	Page      int         `json:"page"`      // 页码，默认为 1
	Limit     int         `json:"limit"`     // 每页数量，默认为 20
	Status    OrderStatus `json:"status"`    // 订单状态筛选
	StartDate string      `json:"startDate"` // 开始日期 (YYYY-MM-DD)
	EndDate   string      `json:"endDate"`   // 结束日期 (YYYY-MM-DD)
}

// CreateOrder 创建订单
func (c *Client) CreateOrder(req CreateOrderRequest) (*CreateOrderResponse, error) {
	var response CreateOrderResponse
	err := c.requestTyped("POST", "/api/v1/orders", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOrders 获取订单列表
func (c *Client) GetOrders(params *OrderParams) (*OrderListResponse, error) {
	queryParams := make(map[string]interface{})
	if params != nil {
		if params.Page > 0 {
			queryParams["page"] = params.Page
		}
		if params.Limit > 0 {
			queryParams["limit"] = params.Limit
		}
		if params.Status != "" {
			queryParams["status"] = params.Status
		}
		if params.StartDate != "" {
			queryParams["startDate"] = params.StartDate
		}
		if params.EndDate != "" {
			queryParams["endDate"] = params.EndDate
		}
	}

	path := "/api/v1/orders" + buildQueryString(queryParams)

	var response OrderListResponse
	err := c.requestTyped("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrder 获取订单详情
func (c *Client) GetOrder(orderID int) (*OrderDetailResponse, error) {
	path := fmt.Sprintf("/api/v1/orders/%d", orderID)

	var response OrderDetailResponse
	err := c.requestTyped("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
