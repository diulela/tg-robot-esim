/**
 * 交易服务
 */

import { api } from './api.js';

export class TransactionService {
    /**
     * 获取交易历史
     */
    async getTransactionHistory(limit = 20, offset = 0) {
        return api.get('/api/miniapp/transactions', { limit, offset });
    }

    /**
     * 获取交易详情
     */
    async getTransactionByHash(txHash) {
        return api.get(`/api/miniapp/transactions/${txHash}`);
    }

    /**
     * 获取交易统计
     */
    async getTransactionStats() {
        return api.get('/api/miniapp/transactions/stats');
    }
}

// 创建全局实例
export const transactionService = new TransactionService();
