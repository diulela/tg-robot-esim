// 导入所有 API 实例
import { ApiClient, apiClient, type ApiClientConfig } from './client'
import { ProductApi, productApi } from './product'
import { OrderApi, orderApi } from './order'
import { ESIMApi, esimApi } from './esim'
import { WalletApi, walletApi, type USDTRechargeOrder, type RechargeHistoryResponse } from './wallet'
import { RegionApi, regionApi } from './region'
import { UserApi, userApi, type UserProfile, type UserProfileUpdate } from './user'
import { SystemApi, systemApi, type HealthCheckResponse, type SystemConfigResponse } from './system'

// 重新导出所有模块
export { ApiClient, apiClient, type ApiClientConfig }
export { ProductApi, productApi }
export { OrderApi, orderApi }
export { ESIMApi, esimApi }
export { WalletApi, walletApi, type USDTRechargeOrder, type RechargeHistoryResponse }
export { RegionApi, regionApi }
export { UserApi, userApi, type UserProfile, type UserProfileUpdate }
export { SystemApi, systemApi, type HealthCheckResponse, type SystemConfigResponse }

// 统一的 API 对象，保持向后兼容
export default {
  product: productApi,
  order: orderApi,
  esim: esimApi,
  region: regionApi,
  wallet: walletApi,
  user: userApi,
  system: systemApi
}