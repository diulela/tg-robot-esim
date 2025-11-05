// 订单相关 API
import { apiClient } from './client'
import type { Order, PaginatedResponse, OrderQueryParams, CreateOrderRequest } from '@/types'
import { OrderStatus } from '@/types'

export class OrderApi {
  // 转换后端订单数据为前端格式
  transformOrder(backendOrder: any): Order {
    const order: Order = {
      id: String(backendOrder.id),
      orderNumber: backendOrder.order_number || '',
      productId: String(backendOrder.product_id),
      productName: backendOrder.product_name || '',
      amount: backendOrder.amount || 0,
      currency: backendOrder.currency || 'USD',
      status: backendOrder.status || OrderStatus.PENDING,
      createdAt: backendOrder.created_at || new Date().toISOString()
    }

    // 可选字段
    if (backendOrder.payment_method) order.paymentMethod = backendOrder.payment_method
    if (backendOrder.paid_at) order.paidAt = backendOrder.paid_at
    if (backendOrder.completed_at) order.completedAt = backendOrder.completed_at
    if (backendOrder.transaction_hash) order.transactionHash = backendOrder.transaction_hash
    if (backendOrder.refund_reason) order.refundReason = backendOrder.refund_reason

    // eSIM 信息
    if (backendOrder.esim_info) {
      order.esimInfo = {
        iccid: backendOrder.esim_info.iccid || '',
        activationCode: backendOrder.esim_info.activation_code || '',
        qrCode: backendOrder.esim_info.qr_code || '',
        apnType: backendOrder.esim_info.apn_type || 'automatic',
        isRoaming: backendOrder.esim_info.is_roaming || false,
        activatedAt: backendOrder.esim_info.activated_at,
        expiresAt: backendOrder.esim_info.expires_at,
        dataUsed: backendOrder.esim_info.data_used,
        dataRemaining: backendOrder.esim_info.data_remaining
      }
    }

    return order
  }

  // 转换后端订单响应为前端格式
  transformOrderResponse(backendData: any): PaginatedResponse<Order> {
    const { orders = [], stats, limit = 20, offset = 0 } = backendData
    const total = stats?.total_orders || 0
    const page = Math.floor(offset / limit) + 1
    const totalPages = Math.ceil(total / limit)

    // 转换每个订单数据
    const transformedOrders = orders.map((o: any) => this.transformOrder(o))

    return {
      items: transformedOrders,
      total,
      page,
      pageSize: limit,
      totalPages,
      hasNext: offset + limit < total,
      hasPrev: offset > 0
    }
  }

  // 获取订单列表
  async getOrders(params?: OrderQueryParams): Promise<PaginatedResponse<Order>> {
    // 转换前端参数为后端格式
    const backendParams: any = {}
    if (params?.status) backendParams.status = params.status
    if (params?.startDate) backendParams.start_date = params.startDate
    if (params?.endDate) backendParams.end_date = params.endDate
    if (params?.page) {
      const pageSize = params.pageSize || 20
      backendParams.limit = pageSize
      backendParams.offset = (params.page - 1) * pageSize
    } else {
      backendParams.limit = params?.pageSize || 20
      backendParams.offset = 0
    }

    const data = await apiClient.get('/miniapp/orders', backendParams)
    return this.transformOrderResponse(data)
  }

  // 获取订单详情
  async getOrder(id: string): Promise<Order> {
    const data = await apiClient.get(`/miniapp/orders/${id}`)
    return this.transformOrder(data)
  }

  // 创建订单 (购买产品)
  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const result = await apiClient.post('/miniapp/purchase', data)
    return this.transformOrder(result)
  }

  // 取消订单
  async cancelOrder(id: string, reason?: string): Promise<Order> {
    const result = await apiClient.patch(`/miniapp/orders/${id}/cancel`, { reason })
    return this.transformOrder(result)
  }

  // 获取订单状态
  async getOrderStatus(id: string): Promise<{ status: string; updatedAt: string }> {
    return apiClient.get(`/miniapp/orders/${id}/status`)
  }

  // 重新支付订单
  async retryPayment(id: string): Promise<{ paymentUrl: string }> {
    return apiClient.post(`/miniapp/orders/${id}/retry-payment`)
  }
}

// 创建订单 API 实例
export const orderApi = new OrderApi()