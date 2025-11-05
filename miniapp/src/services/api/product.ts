// 产品相关 API
import { apiClient } from './client'
import type { Product, PaginatedResponse, ProductQueryParams } from '@/types'

export class ProductApi {
  // 转换后端产品数据为前端格式
  transformProduct(backendProduct: any): Product {
    // 解析 countries JSON 字符串
    let coverageAreas: string[] = []
    let countryCode = ''
    let country = ''

    try {
      if (backendProduct.countries) {
        const countriesData = JSON.parse(backendProduct.countries)
        if (Array.isArray(countriesData) && countriesData.length > 0) {
          coverageAreas = countriesData.map((c: any) => c.cn || c.name || '')
          countryCode = countriesData[0].code || ''
          country = countriesData[0].cn || countriesData[0].name || ''
        }
      }
    } catch (e) {
      console.warn('解析国家数据失败:', e)
    }

    // 解析 features JSON 字符串
    let features: string[] = []
    try {
      if (backendProduct.features) {
        features = JSON.parse(backendProduct.features)
      }
    } catch (e) {
      console.warn('解析特性数据失败:', e)
    }

    // 格式化数据量
    const dataAmount = backendProduct.data_size >= 1024
      ? `${(backendProduct.data_size / 1024).toFixed(1)}GB`
      : `${backendProduct.data_size}MB`

    return {
      id: String(backendProduct.id),
      name: backendProduct.name || '',
      description: backendProduct.description || '',
      price: backendProduct.price || 0,
      originalPrice: backendProduct.retail_price || undefined,
      currency: 'USD',
      region: backendProduct.type || 'global',
      country: country,
      countryCode: countryCode,
      dataAmount: dataAmount,
      validDays: backendProduct.valid_days || 0,
      coverage: backendProduct.name || '',
      coverageAreas: coverageAreas,
      features: features,
      icon: backendProduct.image || undefined,
      isActive: backendProduct.status === 'active',
      isPopular: backendProduct.is_hot || false,
      createdAt: backendProduct.created_at || new Date().toISOString(),
      updatedAt: backendProduct.updated_at || new Date().toISOString()
    }
  }

  // 转换后端响应为前端格式
  transformProductResponse(backendData: any): PaginatedResponse<Product> {
    const { products = [], total = 0, limit = 20, offset = 0 } = backendData
    const page = Math.floor(offset / limit) + 1
    const totalPages = Math.ceil(total / limit)

    // 转换每个产品数据
    const transformedProducts = products.map((p: any) => this.transformProduct(p))

    return {
      items: transformedProducts,
      total,
      page,
      pageSize: limit,
      totalPages,
      hasNext: offset + limit < total,
      hasPrev: offset > 0
    }
  }

  // 获取产品列表
  async getProducts(params?: ProductQueryParams): Promise<PaginatedResponse<Product>> {
    // 转换前端参数为后端格式
    const backendParams: any = {}
    if (params?.type) backendParams.type = params.type
    if (params?.country) backendParams.country = params.country
    if (params?.search) backendParams.search = params.search
    if (params?.page) {
      const pageSize = params.pageSize || 20
      backendParams.limit = pageSize
      backendParams.offset = (params.page - 1) * pageSize
    } else {
      backendParams.limit = params?.pageSize || 20
      backendParams.offset = params?.offset || 0
    }

    const data = await apiClient.get('/miniapp/products', backendParams)
    return this.transformProductResponse(data)
  }

  // 获取产品详情
  async getProduct(id: string): Promise<Product> {
    const data = await apiClient.get(`/miniapp/products/${id}`)
    return this.transformProduct(data)
  }

  // 根据类型获取产品
  async getProductsByType(type: string, params?: Omit<ProductQueryParams, 'type'>): Promise<PaginatedResponse<Product>> {
    return this.getProducts({ ...params, type })
  }

  // 根据国家获取产品
  async getProductsByCountry(country: string, params?: Omit<ProductQueryParams, 'country'>): Promise<PaginatedResponse<Product>> {
    return this.getProducts({ ...params, country })
  }

  // 搜索产品
  async searchProducts(query: string, params?: ProductQueryParams): Promise<PaginatedResponse<Product>> {
    return this.getProducts({ ...params, search: query })
  }
}

// 创建产品 API 实例
export const productApi = new ProductApi()