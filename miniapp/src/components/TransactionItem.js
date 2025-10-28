/**
 * 交易项组件
 */

import { formatAmount, formatRelativeTime } from '../utils/formatter.js';

export class TransactionItem {
    constructor(transaction) {
        this.transaction = transaction;
    }

    render() {
        const { amount, status, timestamp, tx_hash } = this.transaction;
        
        const isPositive = parseFloat(amount) > 0;
        const amountClass = isPositive ? 'positive' : 'negative';
        const amountPrefix = isPositive ? '+' : '';
        
        const statusText = this.getStatusText(status);

        return `
            <div class="transaction-item">
                <div class="transaction-info">
                    <div class="transaction-title">${statusText}</div>
                    <div class="transaction-time">${formatRelativeTime(timestamp)}</div>
                </div>
                <div class="transaction-amount ${amountClass}">
                    ${amountPrefix}${formatAmount(amount)}
                </div>
            </div>
        `;
    }

    getStatusText(status) {
        const statusMap = {
            'pending': '待确认',
            'confirmed': '已确认',
            'failed': '失败'
        };
        return statusMap[status] || status;
    }
}
