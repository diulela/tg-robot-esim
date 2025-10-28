/**
 * 产品卡片组件
 */

import { formatAmount, formatDataSize, formatValidDays } from '../utils/formatter.js';

export class ProductCard {
    constructor(product, onClick) {
        this.product = product;
        this.onClick = onClick;
    }

    render() {
        const { id, name, price, image, data_size, valid_days, is_hot, is_recommend } = this.product;

        const badges = [];
        if (is_hot) badges.push('<span class="product-card-badge badge-hot">热门</span>');
        if (is_recommend) badges.push('<span class="product-card-badge badge-recommend">推荐</span>');

        return `
            <div class="product-card" data-product-id="${id}">
                <img src="${image || 'https://via.placeholder.com/300x150'}" 
                     alt="${name}" 
                     class="product-card-image"
                     onerror="this.src='https://via.placeholder.com/300x150'">
                <div class="product-card-body">
                    <div class="product-card-title">${name}</div>
                    <div class="text-muted" style="font-size: 12px; margin-top: 4px;">
                        ${formatDataSize(data_size)} · ${formatValidDays(valid_days)}
                    </div>
                    <div class="product-card-price">${formatAmount(price)}</div>
                    ${badges.length > 0 ? badges.join(' ') : ''}
                </div>
            </div>
        `;
    }

    attachEvents(element) {
        element.addEventListener('click', () => {
            if (this.onClick) {
                this.onClick(this.product);
            }
        });
    }
}
