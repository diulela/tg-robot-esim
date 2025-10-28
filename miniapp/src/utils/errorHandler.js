/**
 * 全局错误处理和监控
 */

import { showAlert } from './telegram.js';

/**
 * 错误日志类
 */
class ErrorLogger {
    constructor() {
        this.errors = [];
        this.maxErrors = 100;
    }

    /**
     * 记录错误
     */
    log(error, context = {}) {
        const errorLog = {
            message: error.message || String(error),
            stack: error.stack,
            context,
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent,
            url: window.location.href
        };

        this.errors.push(errorLog);

        // 限制错误日志数量
        if (this.errors.length > this.maxErrors) {
            this.errors.shift();
        }

        // 输出到控制台
        console.error('Error logged:', errorLog);

        // 可以在这里发送到服务器
        // this.sendToServer(errorLog);
    }

    /**
     * 获取所有错误日志
     */
    getErrors() {
        return this.errors;
    }

    /**
     * 清除错误日志
     */
    clear() {
        this.errors = [];
    }

    /**
     * 发送错误到服务器（可选）
     */
    async sendToServer(errorLog) {
        try {
            // await fetch('/api/errors', {
            //     method: 'POST',
            //     headers: { 'Content-Type': 'application/json' },
            //     body: JSON.stringify(errorLog)
            // });
        } catch (error) {
            console.error('Failed to send error to server:', error);
        }
    }
}

// 创建全局错误日志实例
export const errorLogger = new ErrorLogger();

/**
 * 全局错误处理器
 */
export function setupGlobalErrorHandler() {
    // 捕获未处理的错误
    window.addEventListener('error', (event) => {
        errorLogger.log(event.error || new Error(event.message), {
            type: 'uncaught_error',
            filename: event.filename,
            lineno: event.lineno,
            colno: event.colno
        });
    });

    // 捕获未处理的 Promise 拒绝
    window.addEventListener('unhandledrejection', (event) => {
        errorLogger.log(event.reason || new Error('Unhandled Promise Rejection'), {
            type: 'unhandled_rejection',
            promise: event.promise
        });
    });
}

/**
 * 错误边界包装器
 */
export function withErrorBoundary(fn, fallback) {
    return async function(...args) {
        try {
            return await fn.apply(this, args);
        } catch (error) {
            errorLogger.log(error, {
                type: 'error_boundary',
                function: fn.name
            });

            if (fallback) {
                return fallback(error);
            }

            await showAlert('操作失败，请稍后重试');
            throw error;
        }
    };
}

/**
 * 用户操作日志
 */
class ActionLogger {
    constructor() {
        this.actions = [];
        this.maxActions = 50;
    }

    /**
     * 记录用户操作
     */
    log(action, data = {}) {
        const actionLog = {
            action,
            data,
            timestamp: new Date().toISOString(),
            url: window.location.href
        };

        this.actions.push(actionLog);

        // 限制日志数量
        if (this.actions.length > this.maxActions) {
            this.actions.shift();
        }

        console.log('Action logged:', actionLog);
    }

    /**
     * 获取所有操作日志
     */
    getActions() {
        return this.actions;
    }

    /**
     * 清除操作日志
     */
    clear() {
        this.actions = [];
    }
}

// 创建全局操作日志实例
export const actionLogger = new ActionLogger();

/**
 * 性能监控
 */
class PerformanceMonitor {
    constructor() {
        this.metrics = {};
    }

    /**
     * 开始计时
     */
    start(label) {
        this.metrics[label] = {
            startTime: performance.now()
        };
    }

    /**
     * 结束计时
     */
    end(label) {
        if (!this.metrics[label]) {
            console.warn(`No start time found for: ${label}`);
            return;
        }

        const endTime = performance.now();
        const duration = endTime - this.metrics[label].startTime;

        this.metrics[label].endTime = endTime;
        this.metrics[label].duration = duration;

        console.log(`Performance [${label}]: ${duration.toFixed(2)}ms`);

        return duration;
    }

    /**
     * 获取所有指标
     */
    getMetrics() {
        return this.metrics;
    }

    /**
     * 清除指标
     */
    clear() {
        this.metrics = {};
    }
}

// 创建全局性能监控实例
export const performanceMonitor = new PerformanceMonitor();

/**
 * 网络状态监控
 */
export function setupNetworkMonitor() {
    window.addEventListener('online', () => {
        console.log('Network: Online');
        actionLogger.log('network_online');
    });

    window.addEventListener('offline', () => {
        console.log('Network: Offline');
        actionLogger.log('network_offline');
        showAlert('网络连接已断开');
    });
}

/**
 * 初始化所有监控
 */
export function initializeMonitoring() {
    setupGlobalErrorHandler();
    setupNetworkMonitor();
    console.log('Monitoring initialized');
}
