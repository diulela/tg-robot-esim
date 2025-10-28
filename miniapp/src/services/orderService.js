/**
 * 订单服务
 */

import { api } from './api.js';

export class OrderService {
    /**
     * 创建购买订单
     */
    async createPurchaseOrder(productId) {
        return api.post('/api/miniapp/purchase', { product_id: productId });
    }

    /**
     * 获取订单列表
     */
    async getOrders(limit = 20, offset = 0) {
        return api.get('/api/miniapp/orders', { limit, offset });
    }

    /**
     * 获取订单详情
     */
    async getOrderByNo(orderNo) {
        return api.get(`/api/miniapp/orders/${orderNo}`);
    }

    /**
     * 取消订单
     */
    async cancelOrder(orderNo) {
        return api.post(`/api/miniapp/orders/${orderNo}/cancel`);
    }
}

// 创建全局实例
export const orderService = new OrderService();
