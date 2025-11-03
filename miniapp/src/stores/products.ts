// 产品状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { 
  Product, 
  Region, 
  Country, 
  ProductQueryParams,
  PaginatedResponse 
} from '@/types'
import { productApi, regionApi } from '@/services/api'

export const useProductsStore = defineStore('products', () => {
  // 状态
  const products = ref<Product[]>([])
  const regions = ref<Region[]>([])
  const countries = ref<Country[]>([])
  const currentProduct = ref<Product | null>(null)
  const selectedRegion = ref<Region | null>(null)
  const selectedCountry = ref<Country | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0,
    hasNext: false,
    hasPrev: false
  })
  const filters = ref<ProductQueryParams>({
    sortBy: 'price',
    sortOrder: 'asc'
  })

  // 计算属性
  const popularRegions = computed(() => 
    regions.value.filter(region => region.isPopular)
  )

  const popularCountries = computed(() => 
    countries.value.filter(country => country.isPopular)
  )

  const filteredProducts = computed(() => {
    let filtered = products.value

    // 按搜索查询过滤
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter(product => 
        product.name.toLowerCase().includes(query) ||
        product.description.toLowerCase().includes(query) ||
        product.country.toLowerCase().includes(query) ||
        product.region.toLowerCase().includes(query)
      )
    }

    // 按区域过滤
    if (selectedRegion.value) {
      filtered = filtered.filter(product => 
        product.region === selectedRegion.value!.code
      )
    }

    // 按国家过滤
    if (selectedCountry.value) {
      filtered = filtered.filter(product => 
        product.countryCode === selectedCountry.value!.code
      )
    }

    return filtered
  })

  const productsByRegion = computed(() => {
    const grouped: Record<string, Product[]> = {}
    
    products.value.forEach(product => {
      if (!grouped[product.region]) {
        grouped[product.region] = []
      }
      grouped[product.region].push(product)
    })

    return grouped
  })

  const productsByCountry = computed(() => {
    const grouped: Record<string, Product[]> = {}
    
    products.value.forEach(product => {
      if (!grouped[product.countryCode]) {
        grouped[product.countryCode] = []
      }
      grouped[product.countryCode].push(product)
    })

    return grouped
  })

  const priceRange = computed(() => {
    if (products.value.length === 0) return { min: 0, max: 0 }
    
    const prices = products.value.map(p => p.price)
    return {
      min: Math.min(...prices),
      max: Math.max(...prices)
    }
  })

  const availableDataAmounts = computed(() => {
    const amounts = new Set(products.value.map(p => p.dataAmount))
    return Array.from(amounts).sort()
  })

  const availableValidDays = computed(() => {
    const days = new Set(products.value.map(p => p.validDays))
    return Array.from(days).sort((a, b) => a - b)
  })

  const hasProducts = computed(() => products.value.length > 0)

  const canLoadMore = computed(() => pagination.value.hasNext)

  // 操作方法
  const fetchProducts = async (params?: ProductQueryParams, append = false): Promise<void> => {
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

      const response: PaginatedResponse<Product> = await productApi.getProducts(queryParams)

      if (append) {
        products.value = [...products.value, ...response.items]
      } else {
        products.value = response.items
      }

      pagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Products] 产品列表获取成功:', {
        count: response.items.length,
        total: response.total,
        page: response.page
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取产品列表失败'
      error.value = errorMessage
      console.error('[Products] 获取产品列表失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const fetchProductById = async (id: string): Promise<Product> => {
    isLoading.value = true
    error.value = null

    try {
      const product = await productApi.getProduct(id)
      
      // 设置为当前产品
      currentProduct.value = product

      // 更新产品列表中的对应项
      const index = products.value.findIndex(p => p.id === id)
      if (index !== -1) {
        products.value[index] = product
      }

      console.log('[Products] 产品详情获取成功:', product.name)
      return product
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取产品详情失败'
      error.value = errorMessage
      console.error('[Products] 获取产品详情失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const fetchRegions = async (): Promise<void> => {
    if (regions.value.length > 0) return // 避免重复加载

    isLoading.value = true
    error.value = null

    try {
      const regionList = await regionApi.getRegions()
      regions.value = regionList

      console.log('[Products] 区域列表获取成功:', regionList.length)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取区域列表失败'
      error.value = errorMessage
      console.error('[Products] 获取区域列表失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const fetchCountries = async (regionCode?: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      let countryList: Country[]
      
      if (regionCode) {
        countryList = await regionApi.getCountriesByRegion(regionCode)
      } else {
        countryList = await regionApi.getCountries()
      }
      
      countries.value = countryList

      console.log('[Products] 国家列表获取成功:', countryList.length)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取国家列表失败'
      error.value = errorMessage
      console.error('[Products] 获取国家列表失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const fetchProductsByRegion = async (regionCode: string, params?: Omit<ProductQueryParams, 'region'>): Promise<void> => {
    const queryParams = { ...params, region: regionCode }
    await fetchProducts(queryParams)
  }

  const fetchProductsByCountry = async (countryCode: string, params?: Omit<ProductQueryParams, 'country'>): Promise<void> => {
    const queryParams = { ...params, country: countryCode }
    await fetchProducts(queryParams)
  }

  const searchProducts = async (query: string, params?: ProductQueryParams): Promise<void> => {
    searchQuery.value = query
    
    if (!query.trim()) {
      await fetchProducts(params)
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const response: PaginatedResponse<Product> = await productApi.searchProducts(query, params)
      
      products.value = response.items
      pagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Products] 产品搜索成功:', {
        query,
        count: response.items.length,
        total: response.total
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '搜索产品失败'
      error.value = errorMessage
      console.error('[Products] 搜索产品失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const searchCountries = async (query: string): Promise<Country[]> => {
    if (!query.trim()) return countries.value

    try {
      const results = await regionApi.searchCountries(query)
      console.log('[Products] 国家搜索成功:', { query, count: results.length })
      return results
    } catch (err) {
      console.error('[Products] 搜索国家失败:', err)
      return []
    }
  }

  const loadMore = async (): Promise<void> => {
    if (!canLoadMore.value || isLoading.value) return
    await fetchProducts(filters.value, true)
  }

  const refresh = async (): Promise<void> => {
    pagination.value.page = 1
    await fetchProducts(filters.value, false)
  }

  const setFilters = async (newFilters: ProductQueryParams): Promise<void> => {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.page = 1
    await fetchProducts(filters.value, false)
  }

  const clearFilters = async (): Promise<void> => {
    filters.value = {
      sortBy: 'price',
      sortOrder: 'asc'
    }
    searchQuery.value = ''
    selectedRegion.value = null
    selectedCountry.value = null
    pagination.value.page = 1
    await fetchProducts(filters.value, false)
  }

  const setSelectedRegion = async (region: Region | null): Promise<void> => {
    selectedRegion.value = region
    selectedCountry.value = null // 清除国家选择
    
    if (region) {
      await fetchCountries(region.code)
      await fetchProductsByRegion(region.code)
    } else {
      await fetchCountries()
      await fetchProducts()
    }
  }

  const setSelectedCountry = async (country: Country | null): Promise<void> => {
    selectedCountry.value = country
    
    if (country) {
      await fetchProductsByCountry(country.code)
    } else if (selectedRegion.value) {
      await fetchProductsByRegion(selectedRegion.value.code)
    } else {
      await fetchProducts()
    }
  }

  const setCurrentProduct = (product: Product | null): void => {
    currentProduct.value = product
  }

  const clearError = (): void => {
    error.value = null
  }

  const clearProducts = (): void => {
    products.value = []
    currentProduct.value = null
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
  const findProductById = (id: string): Product | undefined => {
    return products.value.find(product => product.id === id)
  }

  const findRegionByCode = (code: string): Region | undefined => {
    return regions.value.find(region => region.code === code)
  }

  const findCountryByCode = (code: string): Country | undefined => {
    return countries.value.find(country => country.code === code)
  }

  const formatProductPrice = (product: Product): string => {
    return `$${product.price.toFixed(2)}`
  }

  const formatProductDataAmount = (dataAmount: string): string => {
    // 格式化数据量显示，如 "1GB", "500MB"
    return dataAmount.toUpperCase()
  }

  const formatProductValidDays = (validDays: number): string => {
    if (validDays === 1) return '1天'
    if (validDays === 7) return '1周'
    if (validDays === 30) return '1个月'
    if (validDays === 365) return '1年'
    return `${validDays}天`
  }

  const getProductSummary = (product: Product): string => {
    return `${product.region} ${formatProductDataAmount(product.dataAmount)}/${formatProductValidDays(product.validDays)}`
  }

  const isProductAvailable = (product: Product): boolean => {
    return product.isActive
  }

  const getProductsByPriceRange = (minPrice: number, maxPrice: number): Product[] => {
    return products.value.filter(product => 
      product.price >= minPrice && product.price <= maxPrice
    )
  }

  const getProductsByDataAmount = (dataAmount: string): Product[] => {
    return products.value.filter(product => 
      product.dataAmount === dataAmount
    )
  }

  const getProductsByValidDays = (validDays: number): Product[] => {
    return products.value.filter(product => 
      product.validDays === validDays
    )
  }

  // 返回状态和方法
  return {
    // 只读状态
    products: readonly(products),
    regions: readonly(regions),
    countries: readonly(countries),
    currentProduct: readonly(currentProduct),
    selectedRegion: readonly(selectedRegion),
    selectedCountry: readonly(selectedCountry),
    isLoading: readonly(isLoading),
    error: readonly(error),
    searchQuery: readonly(searchQuery),
    pagination: readonly(pagination),
    filters: readonly(filters),

    // 计算属性
    popularRegions,
    popularCountries,
    filteredProducts,
    productsByRegion,
    productsByCountry,
    priceRange,
    availableDataAmounts,
    availableValidDays,
    hasProducts,
    canLoadMore,

    // 操作方法
    fetchProducts,
    fetchProductById,
    fetchRegions,
    fetchCountries,
    fetchProductsByRegion,
    fetchProductsByCountry,
    searchProducts,
    searchCountries,
    loadMore,
    refresh,
    setFilters,
    clearFilters,
    setSelectedRegion,
    setSelectedCountry,
    setCurrentProduct,
    clearError,
    clearProducts,

    // 工具方法
    findProductById,
    findRegionByCode,
    findCountryByCode,
    formatProductPrice,
    formatProductDataAmount,
    formatProductValidDays,
    getProductSummary,
    isProductAvailable,
    getProductsByPriceRange,
    getProductsByDataAmount,
    getProductsByValidDays
  }
})