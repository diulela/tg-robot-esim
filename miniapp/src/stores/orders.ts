// 订单状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { 
  Order, 
  OrderStatus, 
  OrderQueryParams, 
  CreateOrderRequest,
  PaginatedResponse 
} from '@/types'
import { orderApi } from '@/services/api'

export const useOrdersStore = defineStore('orders', () => {
  // 状态
  const orders = ref<Order[]>([])
  const currentOrder = ref<Order | null>(null)
  const isLoading = ref(false)
  const isCreating = ref(false)
  const error = ref<string | null>(null)
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0,
    hasNext: false,
    hasPrev: false
  })
  const filters = ref<OrderQueryParams>({
    sortBy: 'createdAt',
    sortOrder: 'desc'
  })

  // 计算属性
  const pendingOrders = computed(() => 
    orders.value.filter(order => order.status === OrderStatus.PENDING)
  )

  const paidOrders = computed(() => 
    orders.value.filter(order => order.status === OrderStatus.PAID)
  )

  const completedOrders = computed(() => 
    orders.value.filter(order => order.status === OrderStatus.COMPLETED)
  )

  const cancelledOrders = computed(() => 
    orders.value.filter(order => order.status === OrderStatus.CANCELLED)
  )

  const totalOrders = computed(() => orders.value.length)

  const totalSpent = computed(() => 
    completedOrders.value.reduce((sum, order) => sum + order.amount, 0)
  )

  const recentOrders = computed(() => 
    orders.value
      .slice()
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, 5)
  )

  const ordersByStatus = computed(() => {
    const grouped: Record<OrderStatus, Order[]> = {
      [OrderStatus.PENDING]: [],
      [OrderStatus.PAID]: [],
      [OrderStatus.PROCESSING]: [],
      [OrderStatus.COMPLETED]: [],
      [OrderStatus.CANCELLED]: [],
      [OrderStatus.REFUNDED]: []
    }

    orders.value.forEach(order => {
      grouped[order.status].push(order)
    })

    return grouped
  })

  const hasOrders = computed(() => orders.value.length > 0)

  const canLoadMore = computed(() => pagination.value.hasNext)

  // 操作方法
  const fetchOrders = async (params?: OrderQueryParams, append = false): Promise<void> => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const queryParams = {
        ...filters.value,
        ...params,
        page: append ? pagination.value.page + 1 : 1,
        pageSize: pagination.value.pageSize
      }

      const response: PaginatedResponse<Order> = await orderApi.getOrders(queryParams)

      if (append) {
        orders.value = [...orders.value, ...response.items]
      } else {
        orders.value = response.items
      }

      pagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Orders] 订单列表获取成功:', {
        count: response.items.length,
        total: response.total,
        page: response.page
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

  const fetchOrderById = async (id: string): Promise<Order> => {
    isLoading.value = true
    error.value = null

    try {
      const order = await orderApi.getOrder(id)
      
      // 更新当前订单
      currentOrder.value = order

      // 更新订单列表中的对应项
      const index = orders.value.findIndex(o => o.id === id)
      if (index !== -1) {
        orders.value[index] = order
      }

      console.log('[Orders] 订单详情获取成功:', order.orderNumber)
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

  const createOrder = async (orderData: CreateOrderRequest): Promise<Order> => {
    if (isCreating.value) {
      throw new Error('正在创建订单，请稍候')
    }

    isCreating.value = true
    error.value = null

    try {
      const newOrder = await orderApi.createOrder(orderData)
      
      // 将新订单添加到列表开头
      orders.value.unshift(newOrder)
      
      // 设置为当前订单
      currentOrder.value = newOrder

      // 更新统计
      pagination.value.total += 1

      console.log('[Orders] 订单创建成功:', newOrder.orderNumber)
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

  const cancelOrder = async (id: string, reason?: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      const updatedOrder = await orderApi.cancelOrder(id, reason)
      
      // 更新订单状态
      updateOrderInList(updatedOrder)

      console.log('[Orders] 订单取消成功:', updatedOrder.orderNumber)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '取消订单失败'
      error.value = errorMessage
      console.error('[Orders] 取消订单失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const retryPayment = async (id: string): Promise<string> => {
    isLoading.value = true
    error.value = null

    try {
      const result = await orderApi.retryPayment(id)
      console.log('[Orders] 重新支付链接获取成功')
      return result.paymentUrl
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取支付链接失败'
      error.value = errorMessage
      console.error('[Orders] 重新支付失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const refreshOrderStatus = async (id: string): Promise<void> => {
    try {
      const statusResult = await orderApi.getOrderStatus(id)
      
      // 如果状态有变化，重新获取完整订单信息
      const currentOrderInList = orders.value.find(o => o.id === id)
      if (currentOrderInList && currentOrderInList.status !== statusResult.status) {
        await fetchOrderById(id)
      }
    } catch (err) {
      console.warn('[Orders] 刷新订单状态失败:', err)
    }
  }

  const loadMore = async (): Promise<void> => {
    if (!canLoadMore.value || isLoading.value) return
    await fetchOrders(filters.value, true)
  }

  const refresh = async (): Promise<void> => {
    pagination.value.page = 1
    await fetchOrders(filters.value, false)
  }

  const setFilters = async (newFilters: OrderQueryParams): Promise<void> => {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.page = 1
    await fetchOrders(filters.value, false)
  }

  const clearFilters = async (): Promise<void> => {
    filters.value = {
      sortBy: 'createdAt',
      sortOrder: 'desc'
    }
    pagination.value.page = 1
    await fetchOrders(filters.value, false)
  }

  const setCurrentOrder = (order: Order | null): void => {
    currentOrder.value = order
  }

  const clearError = (): void => {
    error.value = null
  }

  const clearOrders = (): void => {
    orders.value = []
    currentOrder.value = null
    pagination.value = {
      page: 1,
      pageSize: 20,
      total: 0,
      totalPages: 0,
      hasNext: false,
      hasPrev: false
    }
  }

  // 工具方法
  const updateOrderInList = (updatedOrder: Order): void => {
    const index = orders.value.findIndex(o => o.id === updatedOrder.id)
    if (index !== -1) {
      orders.value[index] = updatedOrder
    }

    // 如果是当前订单，也更新
    if (currentOrder.value?.id === updatedOrder.id) {
      currentOrder.value = updatedOrder
    }
  }

  const findOrderById = (id: string): Order | undefined => {
    return orders.value.find(order => order.id === id)
  }

  const findOrderByNumber = (orderNumber: string): Order | undefined => {
    return orders.value.find(order => order.orderNumber === orderNumber)
  }

  const getOrderStatusText = (status: OrderStatus): string => {
    const statusMap: Record<OrderStatus, string> = {
      [OrderStatus.PENDING]: '待支付',
      [OrderStatus.PAID]: '已支付',
      [OrderStatus.PROCESSING]: '处理中',
      [OrderStatus.COMPLETED]: '已完成',
      [OrderStatus.CANCELLED]: '已取消',
      [OrderStatus.REFUNDED]: '已退款'
    }
    return statusMap[status] || '未知状态'
  }

  const getOrderStatusColor = (status: OrderStatus): string => {
    const colorMap: Record<OrderStatus, string> = {
      [OrderStatus.PENDING]: 'warning',
      [OrderStatus.PAID]: 'info',
      [OrderStatus.PROCESSING]: 'primary',
      [OrderStatus.COMPLETED]: 'success',
      [OrderStatus.CANCELLED]: 'error',
      [OrderStatus.REFUNDED]: 'secondary'
    }
    return colorMap[status] || 'default'
  }

  const formatOrderAmount = (order: Order): string => {
    return `$${order.amount.toFixed(2)}`
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
    cancelOrder,
    retryPayment,
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