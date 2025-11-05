// 用户相关 API
import { apiClient } from './client'
import type { DashboardStats } from '@/types'

// 用户信息接口
export interface UserProfile {
  id: string
  telegramId: number
  firstName: string
  lastName?: string
  username?: string
}

// 用户信息更新接口
export interface UserProfileUpdate {
  firstName?: string
  lastName?: string
}

export class UserApi {
  // 获取用户信息
  async getProfile(): Promise<UserProfile> {
    return apiClient.get('/user/profile')
  }

  // 更新用户信息
  async updateProfile(data: UserProfileUpdate): Promise<void> {
    return apiClient.put('/user/profile', data)
  }

  // 获取用户统计
  async getStats(): Promise<DashboardStats> {
    return apiClient.get('/user/stats')
  }
}

// 创建用户 API 实例
export const userApi = new UserApi()