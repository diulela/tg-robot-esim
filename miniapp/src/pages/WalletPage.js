/**
 * 钱包管理页面
 */

import { walletService } from '../services/walletService.js';
import { transactionService } from '../services/transactionService.js';
import { WalletBalance } from '../components/WalletBalance.js';
import { TransactionItem } from '../components/TransactionItem.js';
import { LoadingSpinner } from '../components/LoadingSpinner.js';
import { showAlert } from '../utils/telegram.js';

export default class WalletPage {
    constructor(params) {
        this.params = params;
        this.balance = null;
        this.transactions = [];
    }

    async render(container) {
        container.innerHTML = this.getTemplate();
        await this.loadData();
    }

    getTemplate() {
        return `
            <div class="wallet-page">
                <!-- 余额卡片 -->
                <div id="balance-container"></div>

                <!-- 交易历史 -->
                <div class="container">
                    <h3 style="margin: 24px 0 16px 0;">交易历史</h3>
                    <div class="card" style="padding: 0;">
                        <div id="transactions-container"></div>
                    </div>
                </div>
            </div>
        `;
    }

    async loadData() {
        const balanceContainer = document.getElementById('balance-container');
        const transactionsContainer = document.getElementById('transactions-container');

        LoadingSpinner.show(balanceContainer);
        LoadingSpinner.show(transactionsContainer);

        try {
            // 加载余额和交易历史
            [this.balance, this.transactions] = await Promise.all([
                walletService.getBalance(),
                transactionService.getTransactionHistory(20, 0)
            ]);

            this.renderBalance();
            this.renderTransactions();
        } catch (error) {
            console.error('Failed to load wallet data:', error);
            await showAlert('加载失败，请稍后重试');
        }
    }

    renderBalance() {
        const container = document.getElementById('balance-container');
        const balanceComponent = new WalletBalance(this.balance, () => {
            window.location.hash = '/recharge';
        });
        
        container.innerHTML = balanceComponent.render();
        balanceComponent.attachEvents(container);
    }

    renderTransactions() {
        const container = document.getElementById('transactions-container');
        const transactions = this.transactions.transactions || [];

        if (transactions.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">📝</div>
                    <div class="empty-state-text">暂无交易记录</div>
                </div>
            `;
            return;
        }

        container.innerHTML = transactions.map(tx => {
            const item = new TransactionItem(tx);
            return item.render();
        }).join('');
    }
}
