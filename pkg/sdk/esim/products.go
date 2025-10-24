package esim

import "fmt"

// ProductType 产品类型
type ProductType string

const (
	ProductTypeLocal    ProductType = "local"    // 本地
	ProductTypeRegional ProductType = "regional" // 区域
	ProductTypeGlobal   ProductType = "global"   // 全球
)

// Country 国家信息
type Country struct {
	CN   string `json:"cn"`   // 中文名
	EN   string `json:"en"`   // 英文名
	Code string `json:"code"` // 国家代码
}

// Product 产品信息
type Product struct {
	ID             int         `json:"id"`             // 产品唯一标识
	Name           string      `json:"name"`           // 产品中文名称
	NameEn         string      `json:"nameEn"`         // 产品英文名称
	Description    string      `json:"description"`    // 产品详细描述
	DescriptionEn  string      `json:"descriptionEn"`  // 产品英文描述
	Type           ProductType `json:"type"`           // 产品类型
	Countries      []Country   `json:"countries"`      // 支持的国家列表
	DataSize       int         `json:"dataSize"`       // 流量大小（MB）
	ValidDays      int         `json:"validDays"`      // 有效天数
	Features       []string    `json:"features"`       // 产品特性列表
	Image          string      `json:"image"`          // 产品图片URL
	Gallery        interface{} `json:"gallery"`        // 图片库
	Price          float64     `json:"price"`          // 价格（美元）- API 返回字段
	RetailPrice    float64     `json:"retailPrice"`    // 零售价格（美元）
	AgentPrice     float64     `json:"agentPrice"`     // 代理商价格（美元）
	CostPrice      float64     `json:"costPrice"`      // 成本价格（美元）
	PlatformProfit float64     `json:"platformProfit"` // 平台利润（美元）
	ThirdPartyID   string      `json:"thirdPartyId"`   // 第三方ID
	ThirdPartyData interface{} `json:"thirdPartyData"` // 第三方数据
	IsHot          bool        `json:"isHot"`          // 是否为热门产品
	IsRecommend    bool        `json:"isRecommend"`    // 是否为推荐产品
	SortOrder      int         `json:"sortOrder"`      // 排序顺序
	Status         string      `json:"status"`         // 产品状态
	CreatedAt      string      `json:"createdAt"`      // 创建时间
}

// Pagination 分页信息
type Pagination struct {
	Page       int `json:"page"`       // 当前页码
	Limit      int `json:"limit"`      // 每页数量
	Total      int `json:"total"`      // 总记录数
	TotalPages int `json:"totalPages"` // 总页数
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	Success   bool               `json:"success"`
	Code      int                `json:"code"`
	Message   ProductListMessage `json:"message"`
	Data      string             `json:"data"`
	Timestamp string             `json:"timestamp"`
}

// ProductListMessage 产品列表消息
type ProductListMessage struct {
	Products   []Product  `json:"products"`
	Pagination Pagination `json:"pagination"`
}

// ProductDetailResponse 产品详情响应
type ProductDetailResponse struct {
	Success   bool    `json:"success"`
	Code      int     `json:"code"`
	Message   Product `json:"message"` // 注意：这里 message 字段包含产品数据
	Data      string  `json:"data"`
	Timestamp string  `json:"timestamp"`
}

// ProductParams 产品查询参数
type ProductParams struct {
	Page    int         `json:"page"`    // 页码，默认为 1
	Limit   int         `json:"limit"`   // 每页数量，默认为 20，最大 100
	Country string      `json:"country"` // 国家代码，如 CN、US、JP
	Type    ProductType `json:"type"`    // 产品类型
}

// GetProducts 获取产品列表
func (c *Client) GetProducts(params *ProductParams) (*ProductListResponse, error) {
	queryParams := make(map[string]interface{})
	if params != nil {
		if params.Page > 0 {
			queryParams["page"] = params.Page
		}
		if params.Limit > 0 {
			queryParams["limit"] = params.Limit
		}
		if params.Country != "" {
			queryParams["country"] = params.Country
		}
		if params.Type != "" {
			queryParams["type"] = params.Type
		}
	}

	path := "/api/v1/products" + buildQueryString(queryParams)

	var response ProductListResponse
	err := c.requestTyped("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetProduct 获取产品详情
func (c *Client) GetProduct(productID int) (*ProductDetailResponse, error) {
	path := fmt.Sprintf("/api/v1/products/%d", productID)

	var response ProductDetailResponse
	err := c.requestTyped("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCountries 获取支持的国家列表
func (c *Client) GetCountries() (map[string]interface{}, error) {
	return c.request("GET", "/api/v1/countries", nil)
}
