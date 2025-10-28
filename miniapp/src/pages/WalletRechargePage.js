/**
 * 钱包充值页面
 */

import { walletService } from '../services/walletService.js';
import { LoadingSpinner } from '../components/LoadingSpinner.js';
import { validateAmount } from '../utils/validator.js';
import { showAlert, showBackButton, hideBackButton } from '../utils/telegram.js';
import QRCode from 'qrcode';

export default class WalletRechargePage {
    constructor(params) {
        this.params = params;
        this.rechargeOrder = null;
        this.selectedAmount = '';
    }

    async render(container) {
        showBackButton(() => window.history.back());

        container.innerHTML = this.getTemplate();
        this.attachEvents(container);
    }

    getTemplate() {
        const presetAmounts = [10, 20, 50, 100, 200, 500];

        return `
            <div class="recharge-page">
                <div class="container">
                    <h2 style="margin-bottom: 24px;">钱包充值</h2>

                    <!-- 充值金额选择 -->
                    <div class="card">
                        <h3 style="margin-bottom: 16px;">选择充值金额</h3>
                        <div class="grid grid-2" style="gap: 12px; margin-bottom: 16px;">
                            ${presetAmounts.map(amount => `
                                <button class="btn btn-secondary amount-btn" data-amount="${amount}">
                                    ${amount} USDT
                                </button>
                            `).join('')}
                        </div>
                        
                        <div style="margin-top: 16px;">
                            <label style="display: block; margin-bottom: 8px; font-weight: 500;">
                                自定义金额
                            </label>
                            <input type="number" 
                                   class="input" 
                                   id="custom-amount"
                                   placeholder="请输入充值金额"
                                   min="1"
                                   step="0.01">
                        </div>

                        <button class="btn btn-primary btn-block mt-md" id="create-order-btn">
                            生成充值订单
                        </button>
                    </div>

                    <!-- 充值信息（初始隐藏） -->
                    <div id="recharge-info" class="hidden">
                        <div class="card">
                            <h3 style="margin-bottom: 16px;">充值信息</h3>
                            
                            <div style="margin-bottom: 16px;">
                                <div class="text-muted" style="margin-bottom: 4px;">充值金额</div>
                                <div style="font-size: 24px; font-weight: 700; color: var(--primary-color);" id="order-amount">
                                </div>
                            </div>

                            <div style="margin-bottom: 16px;">
                                <div class="text-muted" style="margin-bottom: 4px;">充值地址</div>
                                <div style="font-family: monospace; word-break: break-all; font-size: 14px;" id="wallet-address">
                                </div>
                            </div>

                            <div style="text-align: center; margin: 24px 0;">
                                <canvas id="qrcode" style="max-width: 200px;"></canvas>
                            </div>

                            <div style="background: #e3f2fd; padding: 12px; border-radius: 8px; font-size: 14px;">
                                <p style="margin: 0 0 8px 0; font-weight: 500;">充值说明：</p>
                                <ul style="margin: 0; padding-left: 20px; line-height: 1.6;">
                                    <li>请向上述地址转账 USDT (TRC20)</li>
                                    <li>充值将在区块链确认后到账</li>
                                    <li>订单有效期 30 分钟</li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    attachEvents(container) {
        // 预设金额按钮
        const amountBtns = container.querySelectorAll('.amount-btn');
        amountBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                amountBtns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                this.selectedAmount = btn.dataset.amount;
                document.getElementById('custom-amount').value = '';
            });
        });

        // 自定义金额输入
        const customAmountInput = container.querySelector('#custom-amount');
        customAmountInput.addEventListener('input', () => {
            amountBtns.forEach(b => b.classList.remove('active'));
            this.selectedAmount = customAmountInput.value;
        });

        // 创建订单按钮
        const createOrderBtn = container.querySelector('#create-order-btn');
        createOrderBtn.addEventListener('click', () => this.handleCreateOrder());
    }

    async handleCreateOrder() {
        // 验证金额
        const error = validateAmount(this.selectedAmount, 1, 10000);
        if (error) {
            await showAlert(error);
            return;
        }

        const createOrderBtn = document.getElementById('create-order-btn');
        createOrderBtn.disabled = true;
        createOrderBtn.textContent = '生成中...';

        try {
            this.rechargeOrder = await walletService.createRechargeOrder(this.selectedAmount);
            await this.showRechargeInfo();
        } catch (error) {
            console.error('Failed to create recharge order:', error);
            await showAlert(error.message || '创建充值订单失败');
        } finally {
            createOrderBtn.disabled = false;
            createOrderBtn.textContent = '生成充值订单';
        }
    }

    async showRechargeInfo() {
        const rechargeInfo = document.getElementById('recharge-info');
        rechargeInfo.classList.remove('hidden');

        // 显示金额
        document.getElementById('order-amount').textContent = `${this.rechargeOrder.amount} USDT`;

        // 显示地址
        document.getElementById('wallet-address').textContent = this.rechargeOrder.wallet_address;

        // 生成二维码
        const canvas = document.getElementById('qrcode');
        try {
            await QRCode.toCanvas(canvas, this.rechargeOrder.wallet_address, {
                width: 200,
                margin: 1
            });
        } catch (error) {
            console.error('Failed to generate QR code:', error);
        }

        // 滚动到充值信息
        rechargeInfo.scrollIntoView({ behavior: 'smooth' });
    }

    destroy() {
        hideBackButton();
    }
}
