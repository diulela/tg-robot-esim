package esim

import (
	"encoding/json"
	"testing"
)

func TestParseProductDetail(t *testing.T) {
	client := &Client{}

	// 测试数据：产品详情在 message 字段中
	testResponse := `{
		"success": true,
		"code": 200,
		"message": {
			"id": 1439,
			"name": "挪威 无限流量/3天",
			"nameEn": "Lofotel - Unlimited - 3 Days",
			"type": "local",
			"countries": [{"cn": "挪威", "en": "Norway", "code": "NO"}],
			"dataSize": 0,
			"validDays": 3,
			"price": 12.5,
			"costPrice": 5.14,
			"description": "挪威专用套餐",
			"features": ["仅数据流量", "支持漫游"],
			"status": "active"
		},
		"data": "获取产品详情成功",
		"timestamp": "2025-10-27T16:18:50.381Z"
	}`

	var response ProductDetailResponse
	err := json.Unmarshal([]byte(testResponse), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal test response: %v", err)
	}

	// 解析产品详情
	err = client.parseProductDetail(&response)
	if err != nil {
		t.Fatalf("Failed to parse product detail: %v", err)
	}

	// 验证解析结果
	if response.ProductDetail == nil {
		t.Fatal("ProductDetail is nil")
	}

	detail := response.ProductDetail
	if detail.ID != 1439 {
		t.Errorf("Expected ID 1439, got %d", detail.ID)
	}

	if detail.Name != "挪威 无限流量/3天" {
		t.Errorf("Expected name '挪威 无限流量/3天', got '%s'", detail.Name)
	}

	if detail.Price != 12.5 {
		t.Errorf("Expected price 12.5, got %f", detail.Price)
	}

	if len(detail.Countries) != 1 {
		t.Errorf("Expected 1 country, got %d", len(detail.Countries))
	}

	if detail.Countries[0].CN != "挪威" {
		t.Errorf("Expected country '挪威', got '%s'", detail.Countries[0].CN)
	}
}

func TestParseProductDetailFromData(t *testing.T) {
	client := &Client{}

	// 测试数据：产品详情在 data 字段中
	testResponse := `{
		"success": true,
		"message": "获取产品详情成功",
		"data": {
			"id": 1,
			"name": "美国5GB-30天套餐",
			"type": "local",
			"countries": [{"cn": "美国", "en": "United States", "code": "US"}],
			"dataSize": 5120,
			"validDays": 30,
			"price": 50.0,
			"costPrice": 42.5,
			"description": "美国本地5GB流量套餐",
			"features": ["高速4G/5G网络"],
			"status": "active"
		}
	}`

	var response ProductDetailResponse
	err := json.Unmarshal([]byte(testResponse), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal test response: %v", err)
	}

	// 解析产品详情
	err = client.parseProductDetail(&response)
	if err != nil {
		t.Fatalf("Failed to parse product detail: %v", err)
	}

	// 验证解析结果
	if response.ProductDetail == nil {
		t.Fatal("ProductDetail is nil")
	}

	detail := response.ProductDetail
	if detail.ID != 1 {
		t.Errorf("Expected ID 1, got %d", detail.ID)
	}

	if detail.Name != "美国5GB-30天套餐" {
		t.Errorf("Expected name '美国5GB-30天套餐', got '%s'", detail.Name)
	}
}
