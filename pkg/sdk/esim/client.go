package esim

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Config SDK配置
type Config struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// Client eSIM API客户端
type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建新的SDK客户端
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.your-domain.com"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Client{
		apiKey:     config.APIKey,
		apiSecret:  config.APISecret,
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
	}
}

// generateSignature 生成API签名
func (c *Client) generateSignature(method, path, body, timestamp, nonce string) string {
	signString := method + path + body + timestamp + nonce
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(signString))
	return hex.EncodeToString(h.Sum(nil))
}

// generateNonce 生成随机字符串
func generateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// request 发送API请求
func (c *Client) request(method, path string, data interface{}) (map[string]interface{}, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := generateNonce(16)

	var bodyStr string
	var bodyReader io.Reader

	if data != nil {
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshal request data: %w", err)
		}
		bodyStr = string(bodyBytes)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	signature := c.generateSignature(method, path, bodyStr, timestamp, nonce)

	reqURL := c.baseURL + path
	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "EsimSDK-Go/1.0.0")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-nonce", nonce)
	req.Header.Set("x-signature", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if msg, ok := errResp["message"].(string); ok {
				return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, msg)
			}
		}
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

// requestTyped 发送API请求并解析到指定类型
func (c *Client) requestTyped(method, path string, data interface{}, result interface{}) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := generateNonce(16)

	var bodyStr string
	var bodyReader io.Reader

	if data != nil {
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshal request data: %w", err)
		}
		bodyStr = string(bodyBytes)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	signature := c.generateSignature(method, path, bodyStr, timestamp, nonce)

	reqURL := c.baseURL + path
	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "EsimSDK-Go/1.0.0")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-nonce", nonce)
	req.Header.Set("x-signature", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if msg, ok := errResp["message"].(string); ok {
				return fmt.Errorf("API error %d: %s", resp.StatusCode, msg)
			}
		}
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

// buildQueryString 构建查询字符串
func buildQueryString(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for k, v := range params {
		if v != nil {
			values.Add(k, fmt.Sprintf("%v", v))
		}
	}

	query := values.Encode()
	if query != "" {
		return "?" + query
	}
	return ""
}
