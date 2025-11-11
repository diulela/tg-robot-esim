// 钱包相关 API
import { apiClient } from './client'
import type { Wallet, WalletTransaction, PaginatedResponse, WalletRechargeRequest } from '@/types'

// USDT 充值订单接口
export interface USDTRechargeOrder {
  order_no: string
  amount: string
  exact_amount: string
  wallet_address: string
  status: string
  tx_hash: string
  confirmations: number
  expires_at: string
  confirmed_at?: string
  created_at: string
}

// 充值历史响应接口
export interface RechargeHistoryResponse {
  orders: Array<{
    order_no: string
    amount: string
    status: string
    tx_hash: string
    created_at: string
    confirmed_at?: string
  }>
  total: number
  limit: number
  offset: number
}

export class WalletApi {
  // 获取钱包余额
  async getWallet(): Promise<Wallet> {
    return apiClient.get('/miniapp/wallet/balance')
  }

  // 获取钱包历史记录
  async getTransactions(params?: { 
    limit?: number
    offset?: number
    type?: string
    status?: string
    start_date?: string
    end_date?: string
  }): Promise<PaginatedResponse<WalletTransaction>> {
    return apiClient.get('/miniapp/wallet/history', params)
  }

  // 获取钱包历史统计
  async getWalletHistoryStats(): Promise<{
    total_records: number
    total_income: string
    total_expense: string
    pending_amount: string
    completed_amount: string
  }> {
    return apiClient.get('/miniapp/wallet/history/stats')
  }

  // 获取单条历史记录详情
  async getHistoryRecord(recordId: number): Promise<WalletTransaction> {
    return apiClient.get(`/miniapp/wallet/history/${recordId}`)
  }

  // 钱包充值（旧版本，保持兼容性）
  async recharge(data: WalletRechargeRequest): Promise<{ 
    paymentUrl: string
    transactionId: string 
  }> {
    return apiClient.post('/miniapp/wallet/recharge', data)
  }

  // 获取充值状态（旧版本）
  async getRechargeStatus(transactionId: string): Promise<{ 
    status: string
    amount?: number 
  }> {
    return apiClient.get(`/miniapp/wallet/recharge/${transactionId}/status`)
  }

  // 创建 USDT 充值订单
  async createRechargeOrder(data: { amount: string }): Promise<{
    order_no: string
    amount: string
    exact_amount: string
    wallet_address: string
    status: string
    expires_at: string
    created_at: string
  }> {
    return apiClient.post('/miniapp/wallet/recharge', data)
  }

  // 获取充值订单详情
  async getRechargeOrder(orderNo: string): Promise<USDTRechargeOrder> {
    return apiClient.get(`/miniapp/wallet/recharge/${orderNo}`)
  }

  // 手动检查充值状态
  async checkRechargeStatus(orderNo: string): Promise<{
    order_no: string
    status: string
    tx_hash: string
    confirmations: number
    confirmed_at?: string
  }> {
    return apiClient.post(`/miniapp/wallet/recharge/${orderNo}/check`)
  }

  // 获取充值历史
  async getRechargeHistory(params?: { 
    limit?: number
    offset?: number 
  }): Promise<RechargeHistoryResponse> {
    return apiClient.get('/miniapp/wallet/recharge/history', params)
  }
}

// 创建钱包 API 实例
export const walletApi = new WalletApi()