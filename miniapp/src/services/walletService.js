/**
 * 钱包服务
 */

import { api } from './api.js';

export class WalletService {
    /**
     * 获取钱包余额
     */
    async getBalance() {
        return api.get('/api/miniapp/wallet/balance');
    }

    /**
     * 创建充值订单
     */
    async createRechargeOrder(amount) {
        return api.post('/api/miniapp/wallet/recharge', { amount });
    }

    /**
     * 获取充值订单详情
     */
    async getRechargeOrder(orderNo) {
        return api.get(`/api/miniapp/wallet/recharge/${orderNo}`);
    }
}

// 创建全局实例
export const walletService = new WalletService();
