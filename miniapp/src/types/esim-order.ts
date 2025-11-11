// eSIM 订单相关类型定义

// eSIM 订单状态枚举
export enum EsimOrderStatus {
  PENDING = 'pending',           // 待支付
  PROCESSING = 'processing',     // 处理中
  PAID = 'paid',                 // 已支付
  COMPLETED = 'completed',       // 已完成
  FAILED = 'failed',             // 失败
  CANCELLED = 'cancelled'        // 已取消
}

// 创建 eSIM 订单请求
export interface CreateEsimOrderRequest {
  productId: number
  quantity: number
  totalAmount: string
  remark?: string
}

// eSIM 订单基本信息
export interface EsimOrder {
  id: string
  orderNo: string
  productId: number
  productName: string
  quantity: number
  unitPrice: string
  totalAmount: string
  status: EsimOrderStatus
  providerOrderId?: string
  createdAt: string
  updatedAt: string
  completedAt?: string
}

// 订单项详情
export interface OrderItemDetail {
  iccid: string
  dataLimit: string
  validityDays: number
}

// eSIM 详情
export interface EsimDetail {
  iccid: string
  activationCode: string
  qrCode: string
  status: string
}

// eSIM 订单详情
export interface EsimOrderDetail extends EsimOrder {
  userId: number
  orderItems: OrderItemDetail[]
  esims: EsimDetail[]
}

// 订单查询参数
export interface OrderQueryParams {
  status?: EsimOrderStatus
  limit?: number
  offset?: number
}

// eSIM 订单分页响应（简化版本）
export interface EsimOrderPaginatedResponse<T> {
  items: T[]
  total: number
  limit: number
  offset: number
}
