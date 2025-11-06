package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"tg-robot-sim/config"
	"tg-robot-sim/pkg/tron"
	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// blockchainService 区块链服务实现
type blockchainService struct {
	tronClient   *tron.Client
	txRepo       repository.TransactionRepository
	config       *config.BlockchainConfig
	logger       Logger
	isMonitoring bool
	stopChan     chan struct{}
	monitorMutex sync.RWMutex
	watchedAddrs map[string]bool
	addrMutex    sync.RWMutex
}

// NewBlockchainService 创建区块链服务
func NewBlockchainService(
	tronClient *tron.Client,
	config *config.BlockchainConfig,
	logger Logger,
) BlockchainService {
	return &blockchainService{
		tronClient:   tronClient,
		config:       config,
		logger:       logger,
		stopChan:     make(chan struct{}),
		watchedAddrs: make(map[string]bool),
	}
}

// ValidateTransaction 验证交易
func (b *blockchainService) ValidateTransaction(txHash string) (*TransactionInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 从 TRON API 获取交易信息
	tronTx, err := b.tronClient.GetTransaction(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction from TRON API: %w", err)
	}

	// 转换为内部格式
	txInfo := &TransactionInfo{
		TxHash:        tronTx.TxID,
		FromAddress:   tronTx.From,
		ToAddress:     tronTx.To,
		Amount:        tronTx.Amount,
		Confirmations: tronTx.Confirmations,
		BlockNumber:   tronTx.BlockNumber,
		Timestamp:     time.Unix(tronTx.Timestamp/1000, 0),
		Status:        b.mapTronStatus(tronTx.Status),
	}

	return txInfo, nil
}

// GetTransactionStatus 获取交易状态
func (b *blockchainService) GetTransactionStatus(txHash string) (TransactionStatus, error) {
	// 首先从数据库查询
	tx, err := b.txRepo.GetByTxHash(context.Background(), txHash)
	if err == nil {
		return TransactionStatus(tx.Status), nil
	}

	// 如果数据库中没有，从区块链查询
	txInfo, err := b.ValidateTransaction(txHash)
	if err != nil {
		return TransactionStatusFailed, fmt.Errorf("failed to validate transaction: %w", err)
	}

	return TransactionStatus(txInfo.Status), nil
}

// MonitorAddress 监控指定地址的交易
func (b *blockchainService) MonitorAddress(address string) error {
	// 验证地址格式
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isValid, err := b.tronClient.ValidateAddress(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to validate address: %w", err)
	}

	if !isValid {
		return fmt.Errorf("invalid TRON address: %s", address)
	}

	// 添加到监控列表
	b.addrMutex.Lock()
	b.watchedAddrs[address] = true
	b.addrMutex.Unlock()

	b.logger.Info("Added address to monitoring: %s", address)
	return nil
}

// GetAddressTransactions 获取地址的交易记录
func (b *blockchainService) GetAddressTransactions(address string, limit int) ([]*TransactionInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tronTxs, err := b.tronClient.GetAddressTransactions(ctx, address, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get address transactions: %w", err)
	}

	var transactions []*TransactionInfo
	for _, tronTx := range tronTxs {
		txInfo := &TransactionInfo{
			TxHash:        tronTx.TxID,
			FromAddress:   tronTx.From,
			ToAddress:     tronTx.To,
			Amount:        tronTx.Amount,
			Confirmations: tronTx.Confirmations,
			BlockNumber:   tronTx.BlockNumber,
			Timestamp:     time.Unix(tronTx.Timestamp/1000, 0),
			Status:        b.mapTronStatus(tronTx.Status),
		}
		transactions = append(transactions, txInfo)
	}

	return transactions, nil
}

// IsTransactionConfirmed 检查交易是否已确认
func (b *blockchainService) IsTransactionConfirmed(txHash string, requiredConfirmations int) (bool, error) {
	txInfo, err := b.ValidateTransaction(txHash)
	if err != nil {
		return false, fmt.Errorf("failed to validate transaction: %w", err)
	}

	return txInfo.Confirmations >= requiredConfirmations, nil
}

// checkPendingTransactions 检查待确认的交易
func (b *blockchainService) checkPendingTransactions() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 获取所有待确认的交易
	pendingTxs, err := b.txRepo.GetPendingTransactions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending transactions: %w", err)
	}

	b.logger.Debug("Checking %d pending transactions", len(pendingTxs))

	for _, tx := range pendingTxs {
		if err := b.updateTransactionStatus(ctx, tx); err != nil {
			b.logger.Error("Failed to update transaction %s: %v", tx.TxHash, err)
			continue
		}
	}

	return nil
}

// updateTransactionStatus 更新交易状态
func (b *blockchainService) updateTransactionStatus(ctx context.Context, tx *models.Transaction) error {
	// 从区块链获取最新状态
	tronTx, err := b.tronClient.GetTransaction(ctx, tx.TxHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction from blockchain: %w", err)
	}

	// 检查状态是否有变化
	newStatus := b.mapTronStatus(tronTx.Status)
	if string(tx.Status) == newStatus && tx.Confirmations == tronTx.Confirmations {
		return nil // 没有变化
	}

	// 更新数据库中的交易状态
	err = b.txRepo.UpdateStatus(ctx, tx.TxHash, models.TransactionStatus(newStatus), tronTx.Confirmations)
	if err != nil {
		return fmt.Errorf("failed to update transaction status in database: %w", err)
	}

	b.logger.Info("Updated transaction %s: status=%s, confirmations=%d",
		tx.TxHash, newStatus, tronTx.Confirmations)

	return nil
}

// GetAddressIncomingTransactions 获取地址的入账交易
func (b *blockchainService) GetAddressIncomingTransactions(ctx context.Context, address string, minAmount string) ([]*TransactionInfo, error) {
	// 获取地址的所有交易
	allTxs, err := b.GetAddressTransactions(address, 50) // 获取最近50笔交易
	if err != nil {
		return nil, fmt.Errorf("获取地址交易失败: %w", err)
	}

	var incomingTxs []*TransactionInfo
	for _, tx := range allTxs {
		// 只返回入账交易（ToAddress 是目标地址）
		if tx.ToAddress == address {
			incomingTxs = append(incomingTxs, tx)
		}
	}

	return incomingTxs, nil
}

// GetTransactionByHash 根据哈希获取交易详情
func (b *blockchainService) GetTransactionByHash(ctx context.Context, txHash string) (*TransactionInfo, error) {
	return b.ValidateTransaction(txHash)
}

// MatchTransactionAmount 匹配交易金额
func (b *blockchainService) MatchTransactionAmount(txAmount string, targetAmount string) bool {
	// 简单的字符串比较，实际应用中可能需要更精确的数值比较
	return txAmount == targetAmount
}

// mapTronStatus 映射 TRON 状态到内部状态
func (b *blockchainService) mapTronStatus(tronStatus string) string {
	switch tronStatus {
	case "SUCCESS":
		return string(TransactionStatusConfirmed)
	case "FAILED":
		return string(TransactionStatusFailed)
	default:
		return string(TransactionStatusPending)
	}
}
