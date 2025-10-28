/**
 * 性能优化工具
 */

/**
 * 图片懒加载
 */
export function lazyLoadImages() {
    const images = document.querySelectorAll('img[data-src]');
    
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.removeAttribute('data-src');
                observer.unobserve(img);
            }
        });
    });
    
    images.forEach(img => imageObserver.observe(img));
}

/**
 * 缓存管理器
 */
export class CacheManager {
    constructor(prefix = 'app_cache_') {
        this.prefix = prefix;
        this.maxAge = 5 * 60 * 1000; // 5 分钟
    }

    /**
     * 设置缓存
     */
    set(key, value, maxAge = this.maxAge) {
        const item = {
            value,
            timestamp: Date.now(),
            maxAge
        };
        
        try {
            localStorage.setItem(this.prefix + key, JSON.stringify(item));
        } catch (error) {
            console.error('Failed to set cache:', error);
        }
    }

    /**
     * 获取缓存
     */
    get(key) {
        try {
            const itemStr = localStorage.getItem(this.prefix + key);
            if (!itemStr) return null;
            
            const item = JSON.parse(itemStr);
            const now = Date.now();
            
            // 检查是否过期
            if (now - item.timestamp > item.maxAge) {
                this.remove(key);
                return null;
            }
            
            return item.value;
        } catch (error) {
            console.error('Failed to get cache:', error);
            return null;
        }
    }

    /**
     * 移除缓存
     */
    remove(key) {
        try {
            localStorage.removeItem(this.prefix + key);
        } catch (error) {
            console.error('Failed to remove cache:', error);
        }
    }

    /**
     * 清除所有缓存
     */
    clear() {
        try {
            const keys = Object.keys(localStorage);
            keys.forEach(key => {
                if (key.startsWith(this.prefix)) {
                    localStorage.removeItem(key);
                }
            });
        } catch (error) {
            console.error('Failed to clear cache:', error);
        }
    }
}

/**
 * 请求去重（防止重复请求）
 */
export class RequestDeduplicator {
    constructor() {
        this.pendingRequests = new Map();
    }

    /**
     * 执行请求（自动去重）
     */
    async request(key, requestFn) {
        // 如果已有相同请求在进行中，返回该请求的 Promise
        if (this.pendingRequests.has(key)) {
            return this.pendingRequests.get(key);
        }

        // 创建新请求
        const promise = requestFn()
            .finally(() => {
                // 请求完成后移除
                this.pendingRequests.delete(key);
            });

        this.pendingRequests.set(key, promise);
        return promise;
    }
}

/**
 * 虚拟滚动（用于长列表优化）
 */
export class VirtualScroller {
    constructor(container, itemHeight, renderItem) {
        this.container = container;
        this.itemHeight = itemHeight;
        this.renderItem = renderItem;
        this.items = [];
        this.visibleStart = 0;
        this.visibleEnd = 0;
    }

    /**
     * 设置数据
     */
    setItems(items) {
        this.items = items;
        this.render();
    }

    /**
     * 渲染可见项
     */
    render() {
        const scrollTop = this.container.scrollTop;
        const containerHeight = this.container.clientHeight;
        
        this.visibleStart = Math.floor(scrollTop / this.itemHeight);
        this.visibleEnd = Math.ceil((scrollTop + containerHeight) / this.itemHeight);
        
        // 添加缓冲区
        const bufferSize = 5;
        const start = Math.max(0, this.visibleStart - bufferSize);
        const end = Math.min(this.items.length, this.visibleEnd + bufferSize);
        
        // 渲染可见项
        const fragment = document.createDocumentFragment();
        for (let i = start; i < end; i++) {
            const item = this.renderItem(this.items[i], i);
            fragment.appendChild(item);
        }
        
        this.container.innerHTML = '';
        this.container.appendChild(fragment);
    }

    /**
     * 初始化滚动监听
     */
    init() {
        this.container.addEventListener('scroll', () => {
            requestAnimationFrame(() => this.render());
        });
        
        this.render();
    }
}

/**
 * 批量 DOM 操作
 */
export function batchDOMUpdates(updates) {
    requestAnimationFrame(() => {
        updates.forEach(update => update());
    });
}

/**
 * 预加载资源
 */
export function preloadImage(src) {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = reject;
        img.src = src;
    });
}

/**
 * 预加载多个图片
 */
export async function preloadImages(srcs) {
    return Promise.all(srcs.map(src => preloadImage(src)));
}

// 创建全局缓存管理器实例
export const cacheManager = new CacheManager();

// 创建全局请求去重器实例
export const requestDeduplicator = new RequestDeduplicator();
