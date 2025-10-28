/**
 * 产品列表页面
 */

import { productService } from '../services/productService.js';
import { ProductCard } from '../components/ProductCard.js';
import { LoadingSpinner } from '../components/LoadingSpinner.js';
import { showAlert } from '../utils/telegram.js';

export default class ProductListPage {
    constructor(params) {
        this.params = params;
        this.products = [];
        this.currentType = 'all';
        this.searchQuery = '';
        this.loading = false;
    }

    async render(container) {
        container.innerHTML = this.getTemplate();
        await this.loadProducts();
        this.attachEvents(container);
    }

    getTemplate() {
        return `
            <div class="product-list-page">
                <!-- 搜索栏 -->
                <div class="search-bar">
                    <input type="text" 
                           class="search-input" 
                           placeholder="搜索产品..."
                           id="search-input">
                </div>

                <!-- 分类标签 -->
                <div class="category-tabs">
                    <div class="category-tab active" data-type="all">全部</div>
                    <div class="category-tab" data-type="local">本地</div>
                    <div class="category-tab" data-type="regional">区域</div>
                    <div class="category-tab" data-type="global">全球</div>
                </div>

                <!-- 产品列表 -->
                <div class="container">
                    <div id="products-grid" class="grid grid-2"></div>
                </div>
            </div>
        `;
    }

    async loadProducts() {
        if (this.loading) return;

        this.loading = true;
        const grid = document.getElementById('products-grid');
        
        LoadingSpinner.show(grid);

        try {
            const result = await productService.getProducts({
                type: this.currentType,
                search: this.searchQuery,
                limit: 50
            });

            this.products = result.products || [];
            this.renderProducts();
        } catch (error) {
            console.error('Failed to load products:', error);
            await showAlert('加载产品失败，请稍后重试');
            grid.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">😕</div>
                    <div class="empty-state-text">加载失败</div>
                </div>
            `;
        } finally {
            this.loading = false;
        }
    }

    renderProducts() {
        const grid = document.getElementById('products-grid');

        if (this.products.length === 0) {
            grid.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">📦</div>
                    <div class="empty-state-text">暂无产品</div>
                </div>
            `;
            return;
        }

        grid.innerHTML = this.products.map(product => {
            const card = new ProductCard(product, (p) => this.handleProductClick(p));
            return card.render();
        }).join('');

        // 附加事件
        this.products.forEach((product, index) => {
            const cardElement = grid.children[index];
            const card = new ProductCard(product, (p) => this.handleProductClick(p));
            card.attachEvents(cardElement);
        });
    }

    attachEvents(container) {
        // 搜索输入
        const searchInput = container.querySelector('#search-input');
        let searchTimeout;
        searchInput.addEventListener('input', (e) => {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                this.searchQuery = e.target.value;
                this.loadProducts();
            }, 500);
        });

        // 分类标签
        const tabs = container.querySelectorAll('.category-tab');
        tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                // 更新激活状态
                tabs.forEach(t => t.classList.remove('active'));
                tab.classList.add('active');

                // 加载产品
                this.currentType = tab.dataset.type;
                this.loadProducts();
            });
        });
    }

    handleProductClick(product) {
        // 导航到产品详情页
        window.location.hash = `/product/${product.id}`;
    }
}
