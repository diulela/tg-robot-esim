/**
 * 购买页面
 */

import { productService } from '../services/productService.js';
import { walletService } from '../services/walletService.js';
import { orderService } from '../services/orderService.js';
import { LoadingSpinner } from '../components/LoadingSpinner.js';
import { formatAmount } from '../utils/formatter.js';
import { showAlert, showConfirm, showBackButton, hideBackButton } from '../utils/telegram.js';

export default class PurchasePage {
    constructor(params) {
        this.params = params;
        this.product = null;
        this.balance = null;
        this.processing = false;
    }

    async render(container) {
        showBackButton(() => window.history.back());

        const productId = new URLSearchParams(window.location.hash.split('?')[1]).get('product_id');
        if (!productId) {
            await showAlert('产品ID无效');
            window.history.back();
            return;
        }

        container.innerHTML = '<div id="purchase-container"></div>';
        const purchaseContainer = container.querySelector('#purchase-container');
        
        LoadingSpinner.show(purchaseContainer);

        try {
            // 加载产品和余额信息
            [this.product, this.balance] = await Promise.all([
                productService.getProductById(productId),
                walletService.getBalance()
            ]);

            purchaseContainer.innerHTML = this.getTemplate();
            this.attachEvents(purchaseContainer);
        } catch (error) {
            console.error('Failed to load purchase page:', error);
            await showAlert('加载失败，请稍后重试');
            window.history.back();
        }
    }

    getTemplate() {
        const { name, price } = this.product;
        const { balance } = this.balance;
        const balanceNum = parseFloat(balance);
        const priceNum = parseFloat(price);
        const insufficient = balanceNum < priceNum;

        return `
            <div class="purchase-page">
                <div class="container">
                    <h2 style="margin-bottom: 24px;">确认购买</h2>

                    <!-- 产品信息 -->
                    <div class="card">
                        <h3 style="margin-bottom: 12px;">产品信息</h3>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
                            <span class="text-muted">产品名称</span>
                            <span style="font-weight: 500;">${name}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between;">
                            <span class="text-muted">产品价格</span>
                            <span style="font-weight: 700; color: var(--primary-color);">
                                ${formatAmount(price)}
                            </span>
                        </div>
                    </div>

                    <!-- 余额信息 -->
                    <div class="card">
                        <h3 style="margin-bottom: 12px;">账户余额</h3>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 8px;">
                            <span class="text-muted">当前余额</span>
                            <span style="font-weight: 500;">${formatAmount(balance)}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between;">
                            <span class="text-muted">支付后余额</span>
                            <span style="font-weight: 500; color: ${insufficient ? 'var(--error-color)' : 'var(--success-color)'};">
                                ${formatAmount(balanceNum - priceNum)}
                            </span>
                        </div>
                    </div>

                    ${insufficient ? `
                        <div class="card" style="background: #fff3cd; border: 1px solid #ffc107;">
                            <p style="margin: 0; color: #856404;">
                                ⚠️ 余额不足，请先充值
                            </p>
                        </div>
                    ` : ''}

                    <!-- 操作按钮 -->
                    <div style="display: flex; gap: 12px; margin-top: 24px;">
                        ${insufficient ? `
                            <button class="btn btn-primary btn-block" id="recharge-btn">
                                去充值
                            </button>
                        ` : `
                            <button class="btn btn-secondary" style="flex: 1;" id="cancel-btn">
                                取消
                            </button>
                            <button class="btn btn-primary" style="flex: 2;" id="confirm-btn">
                                确认支付
                            </button>
                        `}
                    </div>
                </div>
            </div>
        `;
    }

    attachEvents(container) {
        const confirmBtn = container.querySelector('#confirm-btn');
        const cancelBtn = container.querySelector('#cancel-btn');
        const rechargeBtn = container.querySelector('#recharge-btn');

        if (confirmBtn) {
            confirmBtn.addEventListener('click', () => this.handleConfirm());
        }

        if (cancelBtn) {
            cancelBtn.addEventListener('click', () => window.history.back());
        }

        if (rechargeBtn) {
            rechargeBtn.addEventListener('click', () => {
                window.location.hash = '/recharge';
            });
        }
    }

    async handleConfirm() {
        if (this.processing) return;

        const confirmed = await showConfirm('确认购买此产品？');
        if (!confirmed) return;

        this.processing = true;
        const confirmBtn = document.querySelector('#confirm-btn');
        if (confirmBtn) {
            confirmBtn.disabled = true;
            confirmBtn.textContent = '处理中...';
        }

        try {
            const order = await orderService.createPurchaseOrder(this.product.id);
            await showAlert('购买成功！');
            window.location.hash = '/wallet';
        } catch (error) {
            console.error('Purchase failed:', error);
            await showAlert(error.message || '购买失败，请稍后重试');
            if (confirmBtn) {
                confirmBtn.disabled = false;
                confirmBtn.textContent = '确认支付';
            }
        } finally {
            this.processing = false;
        }
    }

    destroy() {
        hideBackButton();
    }
}
