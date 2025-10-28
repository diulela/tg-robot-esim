/**
 * 状态管理器
 */

export class StateManager {
    constructor() {
        this.state = this.loadState();
        this.listeners = new Map();
    }

    /**
     * 从本地存储加载状态
     */
    loadState() {
        try {
            const saved = localStorage.getItem('app_state');
            return saved ? JSON.parse(saved) : this.getDefaultState();
        } catch (error) {
            console.error('Failed to load state:', error);
            return this.getDefaultState();
        }
    }

    /**
     * 获取默认状态
     */
    getDefaultState() {
        return {
            user: null,
            cart: [],
            selectedProduct: null,
            filters: {
                category: 'all',
                search: ''
            }
        };
    }

    /**
     * 保存状态到本地存储
     */
    saveState() {
        try {
            localStorage.setItem('app_state', JSON.stringify(this.state));
        } catch (error) {
            console.error('Failed to save state:', error);
        }
    }

    /**
     * 获取状态
     */
    get(key) {
        return key ? this.state[key] : this.state;
    }

    /**
     * 设置状态
     */
    set(key, value) {
        this.state[key] = value;
        this.saveState();
        this.notify(key, value);
    }

    /**
     * 更新状态（合并）
     */
    update(key, updates) {
        if (typeof this.state[key] === 'object' && !Array.isArray(this.state[key])) {
            this.state[key] = { ...this.state[key], ...updates };
        } else {
            this.state[key] = updates;
        }
        this.saveState();
        this.notify(key, this.state[key]);
    }

    /**
     * 订阅状态变化
     */
    subscribe(key, callback) {
        if (!this.listeners.has(key)) {
            this.listeners.set(key, []);
        }
        this.listeners.get(key).push(callback);

        // 返回取消订阅函数
        return () => {
            const callbacks = this.listeners.get(key);
            const index = callbacks.indexOf(callback);
            if (index > -1) {
                callbacks.splice(index, 1);
            }
        };
    }

    /**
     * 通知订阅者
     */
    notify(key, value) {
        const callbacks = this.listeners.get(key);
        if (callbacks) {
            callbacks.forEach(callback => callback(value));
        }
    }

    /**
     * 清除状态
     */
    clear() {
        this.state = this.getDefaultState();
        this.saveState();
    }
}

// 创建全局状态管理器实例
export const globalState = new StateManager();
