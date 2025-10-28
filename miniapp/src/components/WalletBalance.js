/**
 * 钱包余额组件
 */

import { formatCurrency } from '../utils/formatter.js';

export class WalletBalance {
    constructor(balance, onRecharge) {
        this.balance = balance;
        this.onRecharge = onRecharge;
    }

    render() {
        const { balance, frozen_balance } = this.balance;

        return `
            <div class="wallet-balance-card">
                <div class="wallet-balance-label">可用余额</div>
                <div class="wallet-balance-amount">${formatCurrency(balance)}</div>
                ${frozen_balance && parseFloat(frozen_balance) > 0 ? `
                    <div class="text-muted" style="font-size: 14px;">
                        冻结余额: ${formatCurrency(frozen_balance)}
                    </div>
                ` : ''}
                <div class="wallet-balance-actions">
                    <button class="btn btn-primary" id="recharge-btn">充值</button>
                </div>
            </div>
        `;
    }

    attachEvents(element) {
        const rechargeBtn = element.querySelector('#recharge-btn');
        if (rechargeBtn && this.onRecharge) {
            rechargeBtn.addEventListener('click', this.onRecharge);
        }
    }
}
