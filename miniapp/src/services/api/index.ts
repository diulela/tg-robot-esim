// API 服务统一导出
export { ApiClient, apiClient, type ApiClientConfig } from './client'
export { ProductApi, productApi } from './product'
export { OrderApi, orderApi } from './order'
export { WalletApi, walletApi, type USDTRechargeOrder, type RechargeHistoryResponse } from './wallet'
export { RegionApi, regionApi } from './region'
export { UserApi, userApi, type UserProfile, type UserProfileUpdate } from './user'
export { SystemApi, systemApi, type HealthCheckResponse, type SystemConfigResponse } from './system'

// 统一的 API 对象，保持向后兼容
export default {
  product: productApi,
  order: orderApi,
  region: regionApi,
  wallet: walletApi,
  user: userApi,
  system: systemApi
}