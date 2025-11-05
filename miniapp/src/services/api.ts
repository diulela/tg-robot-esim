// API 服务层 - 重新导出模块化的 API 服务
// 这个文件保持向后兼容，所有新的 API 模块都在 ./api/ 目录中

// 重新导出所有 API 模块，保持向后兼容
export {
  apiClient,
  productApi,
  orderApi,
  regionApi,
  walletApi,
  userApi,
  systemApi,
  default
} from './api/index'