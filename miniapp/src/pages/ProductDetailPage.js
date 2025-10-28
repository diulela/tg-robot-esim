/**
 * 产品详情页面
 */

import { productService } from '../services/productService.js';
import { LoadingSpinner } from '../components/LoadingSpinner.js';
import { formatAmount, formatDataSize, formatValidDays, formatCountries, safeParseJSON } from '../utils/formatter.js';
import { showAlert, showBackButton, hideBackButton, showMainButton, hideMainButton } from '../utils/telegram.js';

export default class ProductDetailPage {
    constructor(params) {
        this.params = params;
        this.product = null;
    }

    async render(container) {
        // 显示返回按钮
        showBackButton(() => {
            window.history.back();
        });

        container.innerHTML = '<div id="product-detail-container"></div>';
        const detailContainer = container.querySelector('#product-detail-container');
        
        LoadingSpinner.show(detailContainer);

        try {
            this.product = await productService.getProductById(this.params.id);
            detailContainer.innerHTML = this.getTemplate();
            this.setupMainButton();
        } catch (error) {
            console.error('Failed to load product:', error);
            await showAlert('加载产品详情失败');
            detailContainer.innerHTML = `
                <div class="container">
                    <div class="empty-state">
                        <div class="empty-state-icon">😕</div>
                        <div class="empty-state-text">加载失败</div>
                    </div>
                </div>
            `;
        }
    }

    getTemplate() {
        const { name, name_en, description, description_en, price, image, 
                data_size, valid_days, countries, features, type } = this.product;

        const featuresList = safeParseJSON(features, []);
        const countriesText = formatCountries(countries);

        return `
            <div class="product-detail-page">
                <!-- 产品图片 -->
                <div style="width: 100%; background: var(--tg-theme-secondary-bg-color);">
                    <img src="${image || 'https://via.placeholder.com/400x300'}" 
                         alt="${name}"
                         style="width: 100%; height: 250px; object-fit: cover;"
                         onerror="this.src='https://via.placeholder.com/400x300'">
                </div>

                <div class="container">
                    <!-- 产品标题和价格 -->
                    <div class="card">
                        <h2 style="margin-bottom: 8px;">${name}</h2>
                        ${name_en ? `<div class="text-muted" style="margin-bottom: 16px;">${name_en}</div>` : ''}
                        <div style="font-size: 28px; font-weight: 700; color: var(--primary-color);">
                            ${formatAmount(price)}
                        </div>
                    </div>

                    <!-- 产品规格 -->
                    <div class="card">
                        <h3 style="margin-bottom: 12px;">产品规格</h3>
                        <div style="display: flex; flex-direction: column; gap: 8px;">
                            <div style="display: flex; justify-content: space-between;">
                                <span class="text-muted">数据流量</span>
                                <span style="font-weight: 500;">${formatDataSize(data_size)}</span>
                            </div>
                            <div style="display: flex; justify-content: space-between;">
                                <span class="text-muted">有效期</span>
                                <span style="font-weight: 500;">${formatValidDays(valid_days)}</span>
                            </div>
                            <div style="display: flex; justify-content: space-between;">
                                <span class="text-muted">类型</span>
                                <span style="font-weight: 500;">${this.getTypeText(type)}</span>
                            </div>
                            <div style="display: flex; justify-content: space-between;">
                                <span class="text-muted">覆盖地区</span>
                                <span style="font-weight: 500;">${countriesText}</span>
                            </div>
                        </div>
                    </div>

                    <!-- 产品描述 -->
                    ${description ? `
                        <div class="card">
                            <h3 style="margin-bottom: 12px;">产品描述</h3>
                            <p style="line-height: 1.6; color: var(--tg-theme-text-color);">
                                ${description}
                            </p>
                        </div>
                    ` : ''}

                    <!-- 产品特性 -->
                    ${featuresList.length > 0 ? `
                        <div class="card">
                            <h3 style="margin-bottom: 12px;">产品特性</h3>
                            <ul style="padding-left: 20px; line-height: 1.8;">
                                ${featuresList.map(f => `<li>${f}</li>`).join('')}
                            </ul>
                        </div>
                    ` : ''}

                    <div style="height: 80px;"></div>
                </div>
            </div>
        `;
    }

    getTypeText(type) {
        const typeMap = {
            'local': '本地',
            'regional': '区域',
            'global': '全球'
        };
        return typeMap[type] || type;
    }

    setupMainButton() {
        showMainButton('立即购买', () => {
            this.handlePurchase();
        });
    }

    async handlePurchase() {
        // 导航到购买页面
        window.location.hash = `/purchase?product_id=${this.product.id}`;
    }

    destroy() {
        hideBackButton();
        hideMainButton();
    }
}
