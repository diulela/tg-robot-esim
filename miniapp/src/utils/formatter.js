/**
 * 数据格式化工具函数
 */

/**
 * 格式化金额
 */
export function formatAmount(amount, currency = 'USDT') {
    const num = parseFloat(amount);
    if (isNaN(num)) return '0.00';
    
    return `${num.toFixed(2)} ${currency}`;
}

/**
 * 格式化大数字金额（带千分位）
 */
export function formatCurrency(amount, currency = 'USDT') {
    const num = parseFloat(amount);
    if (isNaN(num)) return '0.00';
    
    const formatted = num.toLocaleString('zh-CN', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    });
    
    return `${formatted} ${currency}`;
}

/**
 * 格式化日期时间
 */
export function formatDateTime(date) {
    if (!date) return '';
    
    const d = new Date(date);
    if (isNaN(d.getTime())) return '';
    
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}`;
}

/**
 * 格式化相对时间
 */
export function formatRelativeTime(date) {
    if (!date) return '';
    
    const d = new Date(date);
    if (isNaN(d.getTime())) return '';
    
    const now = new Date();
    const diff = now - d;
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    
    if (seconds < 60) return '刚刚';
    if (minutes < 60) return `${minutes}分钟前`;
    if (hours < 24) return `${hours}小时前`;
    if (days < 7) return `${days}天前`;
    
    return formatDateTime(date);
}

/**
 * 格式化数据大小
 */
export function formatDataSize(mb) {
    const num = parseFloat(mb);
    if (isNaN(num)) return '0 MB';
    
    if (num >= 1024) {
        return `${(num / 1024).toFixed(1)} GB`;
    }
    
    return `${num} MB`;
}

/**
 * 格式化有效期
 */
export function formatValidDays(days) {
    const num = parseInt(days);
    if (isNaN(num)) return '0天';
    
    if (num >= 365) {
        const years = Math.floor(num / 365);
        return `${years}年`;
    }
    
    if (num >= 30) {
        const months = Math.floor(num / 30);
        return `${months}个月`;
    }
    
    return `${num}天`;
}

/**
 * 截断文本
 */
export function truncateText(text, maxLength = 50) {
    if (!text) return '';
    if (text.length <= maxLength) return text;
    
    return text.substring(0, maxLength) + '...';
}

/**
 * 格式化交易哈希（显示前后部分）
 */
export function formatTxHash(hash, startChars = 6, endChars = 4) {
    if (!hash) return '';
    if (hash.length <= startChars + endChars) return hash;
    
    return `${hash.substring(0, startChars)}...${hash.substring(hash.length - endChars)}`;
}

/**
 * 格式化地址
 */
export function formatAddress(address, startChars = 6, endChars = 4) {
    return formatTxHash(address, startChars, endChars);
}

/**
 * 解析 JSON 字符串（安全）
 */
export function safeParseJSON(jsonString, defaultValue = null) {
    try {
        return JSON.parse(jsonString);
    } catch (error) {
        console.error('Failed to parse JSON:', error);
        return defaultValue;
    }
}

/**
 * 格式化国家列表
 */
export function formatCountries(countriesJson) {
    const countries = safeParseJSON(countriesJson, []);
    if (!Array.isArray(countries) || countries.length === 0) {
        return '全球通用';
    }
    
    if (countries.length <= 3) {
        return countries.join(', ');
    }
    
    return `${countries.slice(0, 3).join(', ')} 等${countries.length}个国家`;
}
