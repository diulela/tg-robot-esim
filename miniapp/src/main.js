import { initTelegramWebApp } from './utils/telegram.js';
import { Router } from './utils/router.js';
import { StateManager } from './utils/state.js';
import { initializeMonitoring } from './utils/errorHandler.js';

// 导入样式
import '../styles/components.css';

// 初始化监控
initializeMonitoring();

// 初始化应用
class App {
    constructor() {
        this.router = null;
        this.state = null;
    }

    async init() {
        try {
            // 初始化 Telegram Web App
            const tgApp = initTelegramWebApp();
            console.log('Telegram Web App initialized:', tgApp);

            // 初始化状态管理
            this.state = new StateManager();
            
            // 初始化路由
            this.router = new Router();
            this.setupRoutes();
            
            // 隐藏加载指示器
            document.getElementById('loading').style.display = 'none';
            document.getElementById('content').style.display = 'block';
            
            // 启动路由
            this.router.start();
            
            console.log('App initialized successfully');
        } catch (error) {
            console.error('Failed to initialize app:', error);
            this.showError('应用初始化失败，请刷新重试');
        }
    }

    setupRoutes() {
        // 路由将在后续任务中配置
        this.router.addRoute('/', () => import('./pages/ProductListPage.js'));
        this.router.addRoute('/product/:id', () => import('./pages/ProductDetailPage.js'));
        this.router.addRoute('/purchase', () => import('./pages/PurchasePage.js'));
        this.router.addRoute('/wallet', () => import('./pages/WalletPage.js'));
        this.router.addRoute('/recharge', () => import('./pages/WalletRechargePage.js'));
    }

    showError(message) {
        const loading = document.getElementById('loading');
        loading.innerHTML = `
            <div style="text-align: center; padding: 20px;">
                <p style="color: var(--error-color); margin-bottom: 16px;">${message}</p>
                <button class="btn btn-primary" onclick="location.reload()">重新加载</button>
            </div>
        `;
    }
}

// 启动应用
const app = new App();
app.init();
