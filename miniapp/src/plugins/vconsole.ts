// vConsole 调试工具配置
import VConsole from 'vconsole'

let vConsoleInstance: VConsole | null = null

/**
 * 初始化 vConsole
 * 只在开发环境或需要调试时启用
 */
export function initVConsole() {
  // 检查是否已经初始化
  if (vConsoleInstance) {
    console.log('[vConsole] 已经初始化')
    return vConsoleInstance
  }

  // 检查环境变量
  const isDev = import.meta.env.DEV
  const enableVConsole = import.meta.env.VITE_ENABLE_VCONSOLE === 'true'
  
  // 检查 URL 参数（允许在生产环境通过 URL 参数启用）
  const urlParams = new URLSearchParams(window.location.search)
  const debugParam = urlParams.get('debug') === 'true'

  // 只在开发环境或明确启用时初始化
  if (isDev || enableVConsole || debugParam) {
    console.log('[vConsole] 初始化调试控制台')
    
    vConsoleInstance = new VConsole({
      theme: 'dark', // 主题：dark 或 light
      defaultPlugins: ['system', 'network', 'element', 'storage'], // 默认插件
      maxLogNumber: 1000, // 最大日志数量
      disableLogScrolling: false, // 是否禁用日志滚动
      onReady: () => {
        console.log('[vConsole] 调试控制台已就绪')
      },
      onClearLog: () => {
        console.log('[vConsole] 日志已清除')
      }
    })

    // 添加自定义日志
    console.log('[vConsole] 调试模式已启用')
    console.log('[vConsole] 环境:', isDev ? '开发' : '生产')
    console.log('[vConsole] Telegram WebApp:', window.Telegram?.WebApp ? '已加载' : '未加载')
  }

  return vConsoleInstance
}

/**
 * 销毁 vConsole
 */
export function destroyVConsole() {
  if (vConsoleInstance) {
    console.log('[vConsole] 销毁调试控制台')
    vConsoleInstance.destroy()
    vConsoleInstance = null
  }
}

/**
 * 显示 vConsole
 */
export function showVConsole() {
  if (vConsoleInstance) {
    vConsoleInstance.show()
  }
}

/**
 * 隐藏 vConsole
 */
export function hideVConsole() {
  if (vConsoleInstance) {
    vConsoleInstance.hide()
  }
}

/**
 * 获取 vConsole 实例
 */
export function getVConsole() {
  return vConsoleInstance
}

// 导出默认配置
export default {
  init: initVConsole,
  destroy: destroyVConsole,
  show: showVConsole,
  hide: hideVConsole,
  getInstance: getVConsole
}
