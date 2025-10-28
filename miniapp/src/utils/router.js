/**
 * 简单的基于 Hash 的路由系统
 */

export class Router {
    constructor() {
        this.routes = new Map();
        this.guards = [];
        this.currentRoute = null;
        this.currentPage = null;
        this.contentElement = document.getElementById('content');
    }

    /**
     * 添加路由
     */
    addRoute(path, handler) {
        this.routes.set(path, handler);
    }

    /**
     * 添加路由守卫
     */
    addGuard(guard) {
        this.guards.push(guard);
    }

    /**
     * 启动路由
     */
    start() {
        window.addEventListener('hashchange', () => this.handleRoute());
        this.handleRoute();
    }

    /**
     * 处理路由变化
     */
    async handleRoute() {
        const hash = window.location.hash.slice(1) || '/';
        const { route, params } = this.matchRoute(hash);

        if (!route) {
            console.error('Route not found:', hash);
            this.navigate('/');
            return;
        }

        try {
            // 执行路由守卫
            for (const guard of this.guards) {
                const canActivate = await guard(hash, params);
                if (!canActivate) {
                    console.log('Route guard blocked navigation to:', hash);
                    return;
                }
            }

            // 清理上一个页面
            if (this.currentPage && typeof this.currentPage.destroy === 'function') {
                this.currentPage.destroy();
            }

            this.currentRoute = { path: hash, params };
            
            // 动态导入页面模块
            const module = await route();
            const PageClass = module.default;
            
            // 创建页面实例
            this.currentPage = new PageClass(params);
            
            // 渲染页面
            await this.currentPage.render(this.contentElement);
        } catch (error) {
            console.error('Failed to load route:', error);
            this.showError('页面加载失败');
        }
    }

    /**
     * 匹配路由
     */
    matchRoute(path) {
        // 精确匹配
        if (this.routes.has(path)) {
            return { route: this.routes.get(path), params: {} };
        }

        // 参数匹配
        for (const [routePath, handler] of this.routes) {
            const params = this.extractParams(routePath, path);
            if (params) {
                return { route: handler, params };
            }
        }

        return { route: null, params: {} };
    }

    /**
     * 提取路由参数
     */
    extractParams(routePath, actualPath) {
        const routeParts = routePath.split('/');
        const actualParts = actualPath.split('/');

        if (routeParts.length !== actualParts.length) {
            return null;
        }

        const params = {};
        for (let i = 0; i < routeParts.length; i++) {
            if (routeParts[i].startsWith(':')) {
                const paramName = routeParts[i].slice(1);
                params[paramName] = actualParts[i];
            } else if (routeParts[i] !== actualParts[i]) {
                return null;
            }
        }

        return params;
    }

    /**
     * 导航到指定路径
     */
    navigate(path, params = {}) {
        let url = path;
        
        // 替换路径参数
        Object.keys(params).forEach(key => {
            url = url.replace(`:${key}`, params[key]);
        });
        
        window.location.hash = url;
    }

    /**
     * 返回上一页
     */
    back() {
        window.history.back();
    }

    /**
     * 显示错误
     */
    showError(message) {
        this.contentElement.innerHTML = `
            <div class="container">
                <div class="empty-state">
                    <div class="empty-state-icon">⚠️</div>
                    <div class="empty-state-text">${message}</div>
                    <button class="btn btn-primary mt-md" onclick="location.reload()">重新加载</button>
                </div>
            </div>
        `;
    }
}
