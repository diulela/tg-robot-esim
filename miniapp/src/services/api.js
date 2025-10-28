/**
 * API 基础配置和请求封装
 */

import { getInitData } from '../utils/telegram.js';
import { cacheManager, requestDeduplicator } from '../utils/performance.js';

// API 基础 URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * API 错误类
 */
export class ApiError extends Error {
    constructor(status, message, details = '') {
        super(message);
        this.name = 'ApiError';
        this.status = status;
        this.details = details;
    }
}

/**
 * API 服务类
 */
export class ApiService {
    constructor() {
        this.baseURL = API_BASE_URL;
    }

    /**
     * 发送 HTTP 请求
     */
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        // 添加 Telegram 初始化数据
        const initData = getInitData();
        if (initData) {
            headers['X-Telegram-Init-Data'] = initData;
        }

        // 添加时间戳（防止重放攻击）
        headers['X-Request-Time'] = Date.now().toString();

        try {
            const response = await fetch(url, {
                ...options,
                headers
            });

            // 解析响应
            const data = await response.json();

            if (!response.ok) {
                throw new ApiError(
                    response.status,
                    data.message || 'Request failed',
                    data.details || ''
                );
            }

            return data.data || data;
        } catch (error) {
            if (error instanceof ApiError) {
                throw error;
            }

            // 网络错误
            if (error.name === 'TypeError' && error.message.includes('fetch')) {
                throw new ApiError(0, '网络连接失败，请检查网络设置');
            }

            throw new ApiError(500, '发生未知错误', error.message);
        }
    }

    /**
     * GET 请求（带缓存）
     */
    async get(endpoint, params = {}, useCache = true) {
        const queryString = new URLSearchParams(params).toString();
        const url = queryString ? `${endpoint}?${queryString}` : endpoint;
        
        // 尝试从缓存获取
        if (useCache) {
            const cacheKey = `get_${url}`;
            const cached = cacheManager.get(cacheKey);
            if (cached) {
                return cached;
            }
            
            // 使用请求去重
            return requestDeduplicator.request(cacheKey, async () => {
                const data = await this.request(url, { method: 'GET' });
                cacheManager.set(cacheKey, data);
                return data;
            });
        }
        
        return this.request(url, {
            method: 'GET'
        });
    }

    /**
     * POST 请求
     */
    async post(endpoint, data = {}) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    /**
     * PUT 请求
     */
    async put(endpoint, data = {}) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    /**
     * DELETE 请求
     */
    async delete(endpoint) {
        return this.request(endpoint, {
            method: 'DELETE'
        });
    }
}

// 创建全局 API 实例
export const api = new ApiService();
