/**
 * 产品服务
 */

import { api } from './api.js';

export class ProductService {
    /**
     * 获取产品列表
     */
    async getProducts(filters = {}) {
        const params = {
            type: filters.type || 'all',
            search: filters.search || '',
            limit: filters.limit || 20,
            offset: filters.offset || 0
        };

        return api.get('/api/miniapp/products', params);
    }

    /**
     * 获取产品详情
     */
    async getProductById(id) {
        return api.get(`/api/miniapp/products/${id}`);
    }

    /**
     * 搜索产品
     */
    async searchProducts(query, limit = 20) {
        return api.get('/api/miniapp/products', {
            search: query,
            limit
        });
    }

    /**
     * 获取热门产品
     */
    async getHotProducts(limit = 10) {
        return api.get('/api/miniapp/products', {
            type: 'hot',
            limit
        });
    }

    /**
     * 获取推荐产品
     */
    async getRecommendedProducts(limit = 10) {
        return api.get('/api/miniapp/products', {
            type: 'recommend',
            limit
        });
    }
}

// 创建全局实例
export const productService = new ProductService();
