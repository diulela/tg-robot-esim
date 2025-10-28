/**
 * 安全工具函数
 */

/**
 * 转义 HTML 特殊字符，防止 XSS 攻击
 */
export function escapeHtml(text) {
    if (!text) return '';
    
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    
    return text.replace(/[&<>"']/g, (m) => map[m]);
}

/**
 * 清理用户输入
 */
export function sanitizeInput(input) {
    if (typeof input !== 'string') return input;
    
    // 移除潜在的脚本标签
    return input
        .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
        .replace(/javascript:/gi, '')
        .replace(/on\w+\s*=/gi, '');
}

/**
 * 验证 URL 是否安全
 */
export function isSafeUrl(url) {
    if (!url) return false;
    
    try {
        const parsed = new URL(url);
        // 只允许 http 和 https 协议
        return ['http:', 'https:'].includes(parsed.protocol);
    } catch {
        return false;
    }
}

/**
 * 生成随机字符串（用于 CSRF token 等）
 */
export function generateRandomString(length = 32) {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    const randomValues = new Uint8Array(length);
    crypto.getRandomValues(randomValues);
    
    for (let i = 0; i < length; i++) {
        result += chars[randomValues[i] % chars.length];
    }
    
    return result;
}

/**
 * 验证数据完整性
 */
export function validateDataIntegrity(data, expectedFields) {
    if (!data || typeof data !== 'object') {
        return false;
    }
    
    for (const field of expectedFields) {
        if (!(field in data)) {
            return false;
        }
    }
    
    return true;
}

/**
 * 限流函数（防止频繁请求）
 */
export function createRateLimiter(maxRequests, timeWindow) {
    const requests = [];
    
    return function() {
        const now = Date.now();
        
        // 移除过期的请求记录
        while (requests.length > 0 && requests[0] < now - timeWindow) {
            requests.shift();
        }
        
        // 检查是否超过限制
        if (requests.length >= maxRequests) {
            return false;
        }
        
        requests.push(now);
        return true;
    };
}

/**
 * 防抖函数
 */
export function debounce(func, wait) {
    let timeout;
    
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

/**
 * 节流函数
 */
export function throttle(func, limit) {
    let inThrottle;
    
    return function(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}
