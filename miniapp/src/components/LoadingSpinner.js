/**
 * 加载指示器组件
 */

export class LoadingSpinner {
    constructor(message = '加载中...') {
        this.message = message;
    }

    render() {
        return `
            <div class="loading-container">
                <div class="spinner"></div>
                <p>${this.message}</p>
            </div>
        `;
    }

    static show(container, message = '加载中...') {
        const spinner = new LoadingSpinner(message);
        container.innerHTML = spinner.render();
    }

    static hide(container) {
        const loadingEl = container.querySelector('.loading-container');
        if (loadingEl) {
            loadingEl.remove();
        }
    }
}
