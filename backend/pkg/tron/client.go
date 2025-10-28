package tron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client TRON API 客户端
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     Logger
}

// Logger 日志接口
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// NewClient 创建 TRON API 客户端
func NewClient(baseURL, apiKey string, logger Logger) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// TransactionInfo TRON 交易信息
type TransactionInfo struct {
	TxID          string `json:"txID"`
	BlockNumber   int64  `json:"blockNumber"`
	BlockHash     string `json:"blockHash"`
	Timestamp     int64  `json:"timestamp"`
	From          string `json:"from"`
	To            string `json:"to"`
	Amount        string `json:"amount"`
	TokenSymbol   string `json:"tokenSymbol"`
	Confirmations int    `json:"confirmations"`
	Status        string `json:"status"`
	GasUsed       int64  `json:"gasUsed"`
	GasPrice      string `json:"gasPrice"`
}

// BlockInfo TRON 区块信息
type BlockInfo struct {
	BlockNumber int64  `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
	Timestamp   int64  `json:"timestamp"`
	TxCount     int    `json:"txCount"`
}

// GetTransaction 获取交易信息
func (c *Client) GetTransaction(ctx context.Context, txHash string) (*TransactionInfo, error) {
	url := fmt.Sprintf("%s/v1/transactions/%s", c.baseURL, txHash)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var txInfo TransactionInfo
	if err := json.Unmarshal(body, &txInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &txInfo, nil
}

// GetAddressTransactions 获取地址的交易记录
func (c *Client) GetAddressTransactions(ctx context.Context, address string, limit int) ([]*TransactionInfo, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s/transactions", c.baseURL, address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	// 添加查询参数
	q := req.URL.Query()
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Data []*TransactionInfo `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Data, nil
}

// GetLatestBlock 获取最新区块信息
func (c *Client) GetLatestBlock(ctx context.Context) (*BlockInfo, error) {
	url := fmt.Sprintf("%s/wallet/getnowblock", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var blockInfo BlockInfo
	if err := json.Unmarshal(body, &blockInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &blockInfo, nil
}

// ValidateAddress 验证 TRON 地址格式
func (c *Client) ValidateAddress(ctx context.Context, address string) (bool, error) {
	url := fmt.Sprintf("%s/wallet/validateaddress", c.baseURL)

	requestBody := map[string]string{
		"address": address,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Result bool `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Result, nil
}
