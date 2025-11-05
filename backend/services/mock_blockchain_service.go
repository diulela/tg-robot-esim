package services

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// mockBlockchainService 模拟区块链服务实现（用于测试）
type mockBlockchainService struct{}

// NewMockBlockchainService 创建模拟区块链服务
func NewMockBlockchainService() BlockchainService {
	return &mockBlockchainService{}
}

// StartMonitoring 开始监控区块链交易
func (m *mockBlockchainService) StartMonitoring(ctx context.Context) error {
	return nil
}

// StopMonitoring 停止监控
func (m *mockBlockchainService) StopMonitoring() error {
	return nil
}

// ValidateTransaction 验证交易
func (m *mockBlockchainService) ValidateTransaction(txHash string) (*TransactionInfo, error) {
	return &TransactionInfo{
		TxHash:        txHash,
		FromAddress:   "TFromAddressExample",
		ToAddress:     "TToAddressExample",
		Amount:        "100.1234",
		Confirmations: 20,
		BlockNumber:   12345678,
		Timestamp:     time.Now(),
		Status:        string(TransactionStatusConfirmed),
	}, nil
}

// GetTransactionStatus 获取交易状态
func (m *mockBlockchainService) GetTransactionStatus(txHash string) (TransactionStatus, error) {
	return TransactionStatusConfirmed, nil
}

// MonitorAddress 监控指定地址的交易
func (m *mockBlockchainService) MonitorAddress(address string) error {
	return nil
}

// GetAddressTransactions 获取地址的交易记录
func (m *mockBlockchainService) GetAddressTransactions(address string, limit int) ([]*TransactionInfo, error) {
	return []*TransactionInfo{}, nil
}

// IsTransactionConfirmed 检查交易是否已确认
func (m *mockBlockchainService) IsTransactionConfirmed(txHash string, requiredConfirmations int) (bool, error) {
	return true, nil
}

// GetAddressIncomingTransactions 获取地址的入账交易
func (m *mockBlockchainService) GetAddressIncomingTransactions(ctx context.Context, address string, minAmount string) ([]*TransactionInfo, error) {
	// 模拟返回一个匹配的交易
	return []*TransactionInfo{
		{
			TxHash:        fmt.Sprintf("mock_tx_%d", time.Now().Unix()),
			FromAddress:   "TFromAddressExample",
			ToAddress:     address,
			Amount:        minAmount, // 返回匹配的金额
			Confirmations: 20,        // 足够的确认数
			BlockNumber:   12345678,
			Timestamp:     time.Now(),
			Status:        string(TransactionStatusConfirmed),
		},
	}, nil
}

// GetTransactionByHash 根据哈希获取交易详情
func (m *mockBlockchainService) GetTransactionByHash(ctx context.Context, txHash string) (*TransactionInfo, error) {
	return m.ValidateTransaction(txHash)
}

// MatchTransactionAmount 匹配交易金额
func (m *mockBlockchainService) MatchTransactionAmount(txAmount string, targetAmount string) bool {
	// 简单的字符串比较
	txFloat, err1 := strconv.ParseFloat(txAmount, 64)
	targetFloat, err2 := strconv.ParseFloat(targetAmount, 64)

	if err1 != nil || err2 != nil {
		return txAmount == targetAmount
	}

	// 允许小的浮点数误差
	diff := txFloat - targetFloat
	if diff < 0 {
		diff = -diff
	}

	return diff < 0.0001 // 允许 0.0001 的误差
}
