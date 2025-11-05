// 区域和国家相关 API
import { apiClient } from './client'
import type { Region, Country } from '@/types'

export class RegionApi {
  // 获取所有区域
  async getRegions(): Promise<Region[]> {
    return apiClient.get('/regions')
  }

  // 获取热门区域
  async getPopularRegions(): Promise<Region[]> {
    return apiClient.get('/regions/popular')
  }

  // 获取区域详情
  async getRegion(id: string): Promise<Region> {
    return apiClient.get(`/regions/${id}`)
  }

  // 获取所有国家
  async getCountries(): Promise<Country[]> {
    return apiClient.get('/countries')
  }

  // 根据区域获取国家
  async getCountriesByRegion(region: string): Promise<Country[]> {
    return apiClient.get(`/countries/region/${region}`)
  }

  // 获取热门国家
  async getPopularCountries(): Promise<Country[]> {
    return apiClient.get('/countries/popular')
  }

  // 搜索国家
  async searchCountries(query: string): Promise<Country[]> {
    return apiClient.get('/countries/search', { q: query })
  }
}

// 创建区域 API 实例
export const regionApi = new RegionApi()