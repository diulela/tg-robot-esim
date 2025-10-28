/**
 * 输入验证工具函数
 */

/**
 * 验证必填字段
 */
export function validateRequired(value, fieldName = '此字段') {
    if (!value || (typeof value === 'string' && value.trim() === '')) {
        return `${fieldName}不能为空`;
    }
    return null;
}

/**
 * 验证金额
 */
export function validateAmount(value, min = 0, max = Infinity) {
    if (!value) {
        return '请输入金额';
    }
    
    const amount = parseFloat(value);
    
    if (isNaN(amount)) {
        return '请输入有效的金额';
    }
    
    if (amount <= min) {
        return `金额必须大于 ${min}`;
    }
    
    if (amount > max) {
        return `金额不能超过 ${max}`;
    }
    
    return null;
}

/**
 * 验证整数
 */
export function validateInteger(value, min = 0, max = Infinity) {
    if (!value && value !== 0) {
        return '请输入数字';
    }
    
    const num = parseInt(value);
    
    if (isNaN(num) || num !== parseFloat(value)) {
        return '请输入有效的整数';
    }
    
    if (num < min) {
        return `数字不能小于 ${min}`;
    }
    
    if (num > max) {
        return `数字不能大于 ${max}`;
    }
    
    return null;
}

/**
 * 验证邮箱
 */
export function validateEmail(email) {
    if (!email) {
        return '请输入邮箱地址';
    }
    
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        return '请输入有效的邮箱地址';
    }
    
    return null;
}

/**
 * 验证手机号（中国）
 */
export function validatePhone(phone) {
    if (!phone) {
        return '请输入手机号';
    }
    
    const phoneRegex = /^1[3-9]\d{9}$/;
    if (!phoneRegex.test(phone)) {
        return '请输入有效的手机号';
    }
    
    return null;
}

/**
 * 验证长度
 */
export function validateLength(value, min = 0, max = Infinity, fieldName = '此字段') {
    if (!value) {
        return `${fieldName}不能为空`;
    }
    
    const length = value.length;
    
    if (length < min) {
        return `${fieldName}长度不能少于 ${min} 个字符`;
    }
    
    if (length > max) {
        return `${fieldName}长度不能超过 ${max} 个字符`;
    }
    
    return null;
}

/**
 * 验证钱包地址（TRON）
 */
export function validateTronAddress(address) {
    if (!address) {
        return '请输入钱包地址';
    }
    
    // TRON 地址以 T 开头，长度为 34
    if (!address.startsWith('T') || address.length !== 34) {
        return '请输入有效的 TRON 钱包地址';
    }
    
    return null;
}

/**
 * 验证交易哈希
 */
export function validateTxHash(hash) {
    if (!hash) {
        return '请输入交易哈希';
    }
    
    // 交易哈希通常是 64 位十六进制字符串
    const hashRegex = /^[a-fA-F0-9]{64}$/;
    if (!hashRegex.test(hash)) {
        return '请输入有效的交易哈希';
    }
    
    return null;
}

/**
 * 批量验证
 */
export function validateAll(validations) {
    const errors = {};
    
    for (const [field, validator] of Object.entries(validations)) {
        const error = validator();
        if (error) {
            errors[field] = error;
        }
    }
    
    return Object.keys(errors).length > 0 ? errors : null;
}

/**
 * 显示验证错误
 */
export function showValidationError(error) {
    if (typeof error === 'string') {
        return error;
    }
    
    if (typeof error === 'object') {
        return Object.values(error).join('\n');
    }
    
    return '验证失败';
}
