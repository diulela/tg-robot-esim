// 订单状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type {
  EsimOrder,
  EsimOrderDetail,
  CreateEsimOrderRequest,
  OrderQueryParams,
  EsimOrderPaginatedResponse
} from '@/types/esim-order'
import { EsimOrderStatus } from '@/types/esim-order'
import { orderApi } from '@/services/api'

export const useOrdersStore = defineStore('orders', () => {
  // 状态
  const orders = ref<EsimOrder[]>([])
  const currentOrder = ref<EsimOrderDetail | null>(null)
  const isLoading = ref(false)
  const isCreating = ref(false)
  const error = ref<string | null>(null)
  const pagination = ref({
    limit: 20,
    offset: 0,
    total: 0
  })
  const filters = ref<OrderQueryParams>({})

  // 计算属性
  const pendingOrders = computed(() =>
    orders.value.filter(order => order.status === EsimOrderStatus.PENDING)
  )

  const paidOrders = computed(() =>
    orders.value.filter(order => order.status === EsimOrderStatus.PAID)
  )

  const completedOrders = computed(() =>
    orders.value.filter(order => order.status === EsimOrderStatus.COMPLETED)
  )

  const cancelledOrders = computed(() =>
    orders.value.filter(order => order.status === EsimOrderStatus.CANCELLED)
  )

  const totalOrders = computed(() => pagination.value.total)

  const totalSpent = computed(() =>
    completedOrders.value.reduce((sum, order) => sum + parseFloat(order.totalAmount), 0)
  )

  const recentOrders = computed(() =>
    orders.value
      .slice()
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, 5)
  )

  const ordersByStatus = computed(() => {
    const grouped: Record<EsimOrderStatus, EsimOrder[]> = {
      [EsimOrderStatus.PENDING]: [],
      [EsimOrderStatus.PAID]: [],
      [EsimOrderStatus.PROCESSING]: [],
      [EsimOrderStatus.COMPLETED]: [],
      [EsimOrderStatus.FAILED]: [],
      [EsimOrderStatus.CANCELLED]: []
    }

    orders.value.forEach(order => {
      grouped[order.status].push(order)
    })

    return grouped
  })

  const hasOrders = computed(() => orders.value.length > 0)

  const canLoadMore = computed(() =>
    pagination.value.offset + orders.value.length < pagination.value.total
  )

  // 操作方法
  const fetchOrders = async (params?: OrderQueryParams, append = false): Promise<void> => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const queryParams: OrderQueryParams = {
        ...filters.value,
        ...params,
        limit: pagination.value.limit,
        offset: append ? pagination.value.offset + pagination.value.limit : 0
      }

      const response: EsimOrderPaginatedResponse<EsimOrder> = await orderApi.getEsimOrders(queryParams)

      if (append) {
        orders.value = [...orders.value, ...response.items]
      } else {
        orders.value = response.items
      }

      pagination.value = {
        limit: response.limit,
        offset: response.offset,
        total: response.total
      }

      console.log('[Orders] 订单列表获取成功:', {
        count: response.items.length,
        total: response.total,
        offset: response.offset
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取订单列表失败'
      error.value = errorMessage
      console.error('[Orders] 获取订单列表失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const fetchOrderById = async (id: string): Promise<EsimOrderDetail> => {
    isLoading.value = true
    error.value = null

    try {
      const order = await orderApi.getEsimOrderDetail(id)

      // 更新当前订单
      currentOrder.value = order

      // 更新订单列表中的对应项（使用基本信息）
      const index = orders.value.findIndex(o => o.id === id)
      if (index !== -1) {
        const basicOrder: EsimOrder = {
          id: order.id,
          orderNo: order.orderNo,
          productId: order.productId,
          productName: order.productName,
          quantity: order.quantity,
          unitPrice: order.unitPrice,
          totalAmount: order.totalAmount,
          status: order.status,
          createdAt: order.createdAt,
          updatedAt: order.updatedAt
        }
        if (order.providerOrderId) basicOrder.providerOrderId = order.providerOrderId
        if (order.completedAt) basicOrder.completedAt = order.completedAt
        orders.value[index] = basicOrder
      }

      console.log('[Orders] 订单详情获取成功:', order.orderNo)
      return order
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取订单详情失败'
      error.value = errorMessage
      console.error('[Orders] 获取订单详情失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const createOrder = async (orderData: CreateEsimOrderRequest): Promise<EsimOrder> => {
    if (isCreating.value) {
      throw new Error('正在创建订单，请稍候')
    }

    isCreating.value = true
    error.value = null

    try {
      const newOrder = await orderApi.createEsimOrder(orderData)

      // 将新订单添加到列表开头
      orders.value.unshift(newOrder)

      // 更新统计
      pagination.value.total += 1

      console.log('[Orders] 订单创建成功:', newOrder.orderNo)
      return newOrder
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '创建订单失败'
      error.value = errorMessage
      console.error('[Orders] 创建订单失败:', err)
      throw new Error(errorMessage)
    } finally {
      isCreating.value = false
    }
  }

  const refreshOrderStatus = async (id: string): Promise<void> => {
    try {
      // 重新获取订单详情
      await fetchOrderById(id)
    } catch (err) {
      console.warn('[Orders] 刷新订单状态失败:', err)
    }
  }

  const loadMore = async (): Promise<void> => {
    if (!canLoadMore.value || isLoading.value) return
    await fetchOrders(filters.value, true)
  }

  const refresh = async (): Promise<void> => {
    pagination.value.offset = 0
    await fetchOrders(filters.value, false)
  }

  const setFilters = async (newFilters: OrderQueryParams): Promise<void> => {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.offset = 0
    await fetchOrders(filters.value, false)
  }

  const clearFilters = async (): Promise<void> => {
    filters.value = {}
    pagination.value.offset = 0
    await fetchOrders(filters.value, false)
  }

  const setCurrentOrder = (order: EsimOrderDetail | null): void => {
    currentOrder.value = order
  }

  const clearError = (): void => {
    error.value = null
  }

  const clearOrders = (): void => {
    orders.value = []
    currentOrder.value = null
    pagination.value = {
      limit: 20,
      offset: 0,
      total: 0
    }
  }

  // 工具方法
  const findOrderById = (id: string): EsimOrder | undefined => {
    return orders.value.find(order => order.id === id)
  }

  const findOrderByNumber = (orderNumber: string): EsimOrder | undefined => {
    return orders.value.find(order => order.orderNo === orderNumber)
  }

  const getOrderStatusText = (status: EsimOrderStatus): string => {
    const statusMap: Record<EsimOrderStatus, string> = {
      [EsimOrderStatus.PENDING]: '待支付',
      [EsimOrderStatus.PAID]: '已支付',
      [EsimOrderStatus.PROCESSING]: '处理中',
      [EsimOrderStatus.COMPLETED]: '已完成',
      [EsimOrderStatus.CANCELLED]: '已取消',
      [EsimOrderStatus.FAILED]: '失败'
    }
    return statusMap[status] || '未知状态'
  }

  const getOrderStatusColor = (status: EsimOrderStatus): string => {
    const colorMap: Record<EsimOrderStatus, string> = {
      [EsimOrderStatus.PENDING]: 'warning',
      [EsimOrderStatus.PAID]: 'info',
      [EsimOrderStatus.PROCESSING]: 'primary',
      [EsimOrderStatus.COMPLETED]: 'success',
      [EsimOrderStatus.CANCELLED]: 'error',
      [EsimOrderStatus.FAILED]: 'error'
    }
    return colorMap[status] || 'default'
  }

  const formatOrderAmount = (order: EsimOrder): string => {
    return `$${order.totalAmount}`
  }

  const formatOrderDate = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getOrderSummary = () => {
    return {
      total: totalOrders.value,
      pending: pendingOrders.value.length,
      paid: paidOrders.value.length,
      completed: completedOrders.value.length,
      cancelled: cancelledOrders.value.length,
      totalSpent: totalSpent.value,
      recent: recentOrders.value
    }
  }

  // 返回状态和方法
  return {
    // 只读状态
    orders: readonly(orders),
    currentOrder: readonly(currentOrder),
    isLoading: readonly(isLoading),
    isCreating: readonly(isCreating),
    error: readonly(error),
    pagination: readonly(pagination),
    filters: readonly(filters),

    // 计算属性
    pendingOrders,
    paidOrders,
    completedOrders,
    cancelledOrders,
    totalOrders,
    totalSpent,
    recentOrders,
    ordersByStatus,
    hasOrders,
    canLoadMore,

    // 操作方法
    fetchOrders,
    fetchOrderById,
    createOrder,
    refreshOrderStatus,
    loadMore,
    refresh,
    setFilters,
    clearFilters,
    setCurrentOrder,
    clearError,
    clearOrders,

    // 工具方法
    findOrderById,
    findOrderByNumber,
    getOrderStatusText,
    getOrderStatusColor,
    formatOrderAmount,
    formatOrderDate,
    getOrderSummary
  }
})