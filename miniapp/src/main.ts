// Vue 3.0 + TypeScript 应用入口
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { setupStore } from './stores'
import { setupRouterGuards } from './router/guards'
import vuetify from './plugins/vuetify'
import { initVConsole } from './plugins/vconsole'

// 初始化调试控制台（开发环境或 URL 参数 ?debug=true）
initVConsole()

// 创建 Vue 应用实例
const app = createApp(App)

// 安装插件
app.use(vuetify)
app.use(router)
setupStore(app)

// 设置路由守卫
setupRouterGuards(router)

// 全局错误处理
app.config.errorHandler = (error, _instance, info) => {
  console.error('[Vue] 全局错误:', error, info)
}

// 挂载应用
app.mount('#app')

console.log('[Main] Vue 3.0 应用启动成功')