/**
 * Telegram Web App 工具函数
 */

let tg = null;

/**
 * 初始化 Telegram Web App
 */
export function initTelegramWebApp() {
    if (typeof window.Telegram === 'undefined' || !window.Telegram.WebApp) {
        console.warn('Telegram WebApp not available, using mock data');
        return createMockTelegramApp();
    }

    tg = window.Telegram.WebApp;
    
    // 展开应用到全屏
    tg.expand();
    
    // 启用关闭确认
    tg.enableClosingConfirmation();
    
    // 设置主题颜色
    applyTheme();
    
    // 准备就绪
    tg.ready();
    
    return tg;
}

/**
 * 获取 Telegram Web App 实例
 */
export function getTelegramApp() {
    return tg;
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
    if (!tg || !tg.initDataUnsafe || !tg.initDataUnsafe.user) {
        return {
            id: 123456789,
            first_name: 'Test',
            last_name: 'User',
            username: 'testuser',
            language_code: 'zh'
        };
    }
    
    return tg.initDataUnsafe.user;
}

/**
 * 获取初始化数据
 */
export function getInitData() {
    if (!tg) return '';
    return tg.initData;
}

/**
 * 显示主按钮
 */
export function showMainButton(text, onClick) {
    if (!tg || !tg.MainButton) return;
    
    tg.MainButton.setText(text);
    tg.MainButton.show();
    tg.MainButton.onClick(onClick);
}

/**
 * 隐藏主按钮
 */
export function hideMainButton() {
    if (!tg || !tg.MainButton) return;
    tg.MainButton.hide();
}

/**
 * 显示返回按钮
 */
export function showBackButton(onClick) {
    if (!tg || !tg.BackButton) return;
    
    tg.BackButton.show();
    tg.BackButton.onClick(onClick);
}

/**
 * 隐藏返回按钮
 */
export function hideBackButton() {
    if (!tg || !tg.BackButton) return;
    tg.BackButton.hide();
}

/**
 * 显示弹窗
 */
export function showAlert(message) {
    if (!tg || !tg.showAlert) {
        alert(message);
        return Promise.resolve();
    }
    
    return new Promise((resolve) => {
        tg.showAlert(message, resolve);
    });
}

/**
 * 显示确认对话框
 */
export function showConfirm(message) {
    if (!tg || !tg.showConfirm) {
        return Promise.resolve(confirm(message));
    }
    
    return new Promise((resolve) => {
        tg.showConfirm(message, resolve);
    });
}

/**
 * 显示弹出窗口
 */
export function showPopup(params) {
    if (!tg || !tg.showPopup) {
        alert(params.message);
        return Promise.resolve();
    }
    
    return new Promise((resolve) => {
        tg.showPopup(params, resolve);
    });
}

/**
 * 关闭 Mini App
 */
export function close() {
    if (!tg || !tg.close) {
        window.close();
        return;
    }
    
    tg.close();
}

/**
 * 应用主题
 */
function applyTheme() {
    if (!tg || !tg.themeParams) return;
    
    const root = document.documentElement;
    const theme = tg.themeParams;
    
    if (theme.bg_color) root.style.setProperty('--tg-theme-bg-color', theme.bg_color);
    if (theme.text_color) root.style.setProperty('--tg-theme-text-color', theme.text_color);
    if (theme.hint_color) root.style.setProperty('--tg-theme-hint-color', theme.hint_color);
    if (theme.link_color) root.style.setProperty('--tg-theme-link-color', theme.link_color);
    if (theme.button_color) root.style.setProperty('--tg-theme-button-color', theme.button_color);
    if (theme.button_text_color) root.style.setProperty('--tg-theme-button-text-color', theme.button_text_color);
    if (theme.secondary_bg_color) root.style.setProperty('--tg-theme-secondary-bg-color', theme.secondary_bg_color);
}

/**
 * 创建模拟的 Telegram App（用于开发测试）
 */
function createMockTelegramApp() {
    return {
        initData: '',
        initDataUnsafe: {
            user: {
                id: 123456789,
                first_name: 'Test',
                last_name: 'User',
                username: 'testuser',
                language_code: 'zh'
            }
        },
        version: '6.0',
        platform: 'web',
        colorScheme: 'light',
        themeParams: {},
        isExpanded: true,
        viewportHeight: window.innerHeight,
        viewportStableHeight: window.innerHeight,
        expand: () => console.log('Mock: expand'),
        close: () => console.log('Mock: close'),
        ready: () => console.log('Mock: ready'),
        enableClosingConfirmation: () => console.log('Mock: enableClosingConfirmation'),
        MainButton: {
            text: '',
            color: '#2481cc',
            textColor: '#ffffff',
            isVisible: false,
            isActive: true,
            setText: (text) => console.log('Mock MainButton setText:', text),
            show: () => console.log('Mock MainButton show'),
            hide: () => console.log('Mock MainButton hide'),
            onClick: (callback) => console.log('Mock MainButton onClick')
        },
        BackButton: {
            isVisible: false,
            show: () => console.log('Mock BackButton show'),
            hide: () => console.log('Mock BackButton hide'),
            onClick: (callback) => console.log('Mock BackButton onClick')
        },
        showAlert: (message, callback) => {
            alert(message);
            if (callback) callback();
        },
        showConfirm: (message, callback) => {
            const result = confirm(message);
            if (callback) callback(result);
        }
    };
}
