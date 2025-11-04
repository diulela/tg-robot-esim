// 产品状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type {
  Product,
  Region,
  Country,
  ProductQueryParams,
  PaginatedResponse,
  HotItem
} from '@/types'
import { productApi, regionApi } from '@/services/api'

// 热门项硬编码数据
export const HOT_ITEMS_DATA = {
  hot: [
    { code: 'CN', name: '中国' },
    { code: 'HK', name: '香港' },
    { code: 'TW', name: '台湾' },
    { code: 'JP', name: '日本' },
    { code: 'VN', name: '越南' },
    { code: 'US', name: '美国' },
    { code: 'MO', name: '澳门' },
    { code: 'TH', name: '泰国' },
    { code: 'KR', name: '韩国' },
    { code: 'SG', name: '新加坡' },
    { code: 'MY', name: '马来西亚' },
    { code: 'AU', name: '澳大利亚' },
    { code: 'GB', name: '英国' }
  ],
  local: [
    { "code": "AD", "name": "安道尔", "en": "Andorra" },
    { "code": "AE", "name": "阿联酋", "en": "United Arab Emirates" },
    { "code": "AF", "name": "阿富汗", "en": "Afghanistan" },
    { "code": "AG", "name": "安提瓜和巴布达", "en": "Antigua and Barbuda" },
    { "code": "AI", "name": "安圭拉", "en": "Anguilla" },
    { "code": "AL", "name": "阿尔巴尼亚", "en": "Albania" },
    { "code": "AM", "name": "亚美尼亚", "en": "Armenia" },
    { "code": "AN", "name": "荷属安的列斯", "en": "Netherlands Antilles" },
    { "code": "AR", "name": "阿根廷", "en": "Argentina" },
    { "code": "AT", "name": "奥地利", "en": "Austria" },
    { "code": "AU", "name": "澳大利亚", "en": "Australia" },
    { "code": "AW", "name": "阿鲁巴", "en": "Aruba" },
    { "code": "AZ", "name": "阿塞拜疆", "en": "Azerbaijan" },
    { "code": "BA", "name": "波黑", "en": "Bosnia and Herzegovina" },
    { "code": "BB", "name": "巴巴多斯", "en": "Barbados" },
    { "code": "BD", "name": "孟加拉国", "en": "Bangladesh" },
    { "code": "BE", "name": "比利时", "en": "Belgium" },
    { "code": "BF", "name": "布基纳法索", "en": "Burkina Faso" },
    { "code": "BG", "name": "保加利亚", "en": "Bulgaria" },
    { "code": "BH", "name": "巴林", "en": "Bahrain" },
    { "code": "BJ", "name": "贝宁", "en": "Benin" },
    { "code": "BL", "name": "圣巴泰勒米", "en": "Saint Barthelemy" },
    { "code": "BM", "name": "百慕大", "en": "Bermuda" },
    { "code": "BN", "name": "文莱", "en": "Brunei" },
    { "code": "BO", "name": "玻利维亚", "en": "Bolivia" },
    { "code": "BQ", "name": "荷兰加勒比区", "en": "Bonaire, Sint Eustatius and Saba" },
    { "code": "BR", "name": "巴西", "en": "Brazil" },
    { "code": "BS", "name": "巴哈马", "en": "Bahamas" },
    { "code": "BT", "name": "不丹", "en": "Bhutan" },
    { "code": "BW", "name": "博茨瓦纳", "en": "Botswana" },
    { "code": "BY", "name": "白俄罗斯", "en": "Belarus" },
    { "code": "BZ", "name": "伯利兹", "en": "Belize" },
    { "code": "CA", "name": "加拿大", "en": "Canada" },
    { "code": "CD", "name": "刚果(金)", "en": "Congo, Democratic Republic" },
    { "code": "CF", "name": "中非共和国", "en": "Central African Republic" },
    { "code": "CG", "name": "刚果(布)", "en": "Congo, Republic" },
    { "code": "CH", "name": "瑞士", "en": "Switzerland" },
    { "code": "CI", "name": "科特迪瓦", "en": "Ivory Coast" },
    { "code": "CL", "name": "智利", "en": "Chile" },
    { "code": "CM", "name": "喀麦隆", "en": "Cameroon" },
    { "code": "CN", "name": "中国", "en": "China" },
    { "code": "CO", "name": "哥伦比亚", "en": "Colombia" },
    { "code": "CR", "name": "哥斯达黎加", "en": "Costa Rica" },
    { "code": "CV", "name": "佛得角", "en": "Cape Verde" },
    { "code": "CW", "name": "库拉索", "en": "Curacao" },
    { "code": "CY", "name": "塞浦路斯", "en": "Cyprus" },
    { "code": "CZ", "name": "捷克", "en": "Czech Republic" },
    { "code": "DE", "name": "德国", "en": "Germany" },
    { "code": "DK", "name": "丹麦", "en": "Denmark" },
    { "code": "DM", "name": "多米尼克", "en": "Dominica" },
    { "code": "DO", "name": "多米尼加", "en": "Dominican Republic" },
    { "code": "DZ", "name": "阿尔及利亚", "en": "Algeria" },
    { "code": "EC", "name": "厄瓜多尔", "en": "Ecuador" },
    { "code": "EE", "name": "爱沙尼亚", "en": "Estonia" },
    { "code": "EG", "name": "埃及", "en": "Egypt" },
    { "code": "ES", "name": "西班牙", "en": "Spain" },
    { "code": "ET", "name": "埃塞俄比亚", "en": "Ethiopia" },
    { "code": "FI", "name": "芬兰", "en": "Finland" },
    { "code": "FJ", "name": "斐济", "en": "Fiji" },
    { "code": "FO", "name": "法罗群岛", "en": "Faroe Islands" },
    { "code": "FR", "name": "法国", "en": "France" },
    { "code": "GA", "name": "加蓬", "en": "Gabon" },
    { "code": "GB", "name": "英国", "en": "United Kingdom" },
    { "code": "GD", "name": "格林纳达", "en": "Grenada" },
    { "code": "GE", "name": "格鲁吉亚", "en": "Georgia" },
    { "code": "GF", "name": "法属圭亚那", "en": "French Guiana" },
    { "code": "GG", "name": "GG", "en": "GG" },
    { "code": "GH", "name": "加纳", "en": "Ghana" },
    { "code": "GI", "name": "直布罗陀", "en": "Gibraltar" },
    { "code": "GL", "name": "格陵兰", "en": "Greenland" },
    { "code": "GM", "name": "冈比亚", "en": "Gambia" },
    { "code": "GN", "name": "几内亚", "en": "Guinea" },
    { "code": "GP", "name": "瓜德罗普", "en": "Guadeloupe" },
    { "code": "GR", "name": "希腊", "en": "Greece" },
    { "code": "GT", "name": "危地马拉", "en": "Guatemala" },
    { "code": "GU", "name": "关岛", "en": "Guam" },
    { "code": "GW", "name": "几内亚比绍", "en": "Guinea-Bissau" },
    { "code": "GY", "name": "圭亚那", "en": "Guyana" },
    { "code": "HK", "name": "香港", "en": "Hong Kong" },
    { "code": "HN", "name": "洪都拉斯", "en": "Honduras" },
    { "code": "HR", "name": "克罗地亚", "en": "Croatia" },
    { "code": "HT", "name": "海地", "en": "Haiti" },
    { "code": "HU", "name": "匈牙利", "en": "Hungary" },
    { "code": "ID", "name": "印度尼西亚", "en": "Indonesia" },
    { "code": "IE", "name": "爱尔兰", "en": "Ireland" },
    { "code": "IL", "name": "以色列", "en": "Israel" },
    { "code": "IM", "name": "马恩岛", "en": "Isle of Man" },
    { "code": "IN", "name": "印度", "en": "India" },
    { "code": "IQ", "name": "伊拉克", "en": "Iraq" },
    { "code": "IR", "name": "伊朗", "en": "Iran" },
    { "code": "IS", "name": "冰岛", "en": "Iceland" },
    { "code": "IT", "name": "意大利", "en": "Italy" },
    { "code": "JE", "name": "泽西岛", "en": "Jersey" },
    { "code": "JM", "name": "牙买加", "en": "Jamaica" },
    { "code": "JO", "name": "约旦", "en": "Jordan" },
    { "code": "JP", "name": "日本", "en": "Japan" },
    { "code": "KE", "name": "肯尼亚", "en": "Kenya" },
    { "code": "KG", "name": "吉尔吉斯斯坦", "en": "Kyrgyzstan" },
    { "code": "KH", "name": "柬埔寨", "en": "Cambodia" },
    { "code": "KN", "name": "圣基茨和尼维斯", "en": "Saint Kitts and Nevis" },
    { "code": "KR", "name": "韩国", "en": "South Korea" },
    { "code": "KW", "name": "科威特", "en": "Kuwait" },
    { "code": "KY", "name": "开曼群岛", "en": "Cayman Islands" },
    { "code": "KZ", "name": "哈萨克斯坦", "en": "Kazakhstan" },
    { "code": "LA", "name": "老挝", "en": "Laos" },
    { "code": "LB", "name": "黎巴嫩", "en": "Lebanon" },
    { "code": "LC", "name": "圣卢西亚", "en": "Saint Lucia" },
    { "code": "LI", "name": "列支敦士登", "en": "Liechtenstein" },
    { "code": "LK", "name": "斯里兰卡", "en": "Sri Lanka" },
    { "code": "LR", "name": "利比里亚", "en": "Liberia" },
    { "code": "LS", "name": "莱索托", "en": "Lesotho" },
    { "code": "LT", "name": "立陶宛", "en": "Lithuania" },
    { "code": "LU", "name": "卢森堡", "en": "Luxembourg" },
    { "code": "LV", "name": "拉脱维亚", "en": "Latvia" },
    { "code": "MA", "name": "摩洛哥", "en": "Morocco" },
    { "code": "MC", "name": "摩纳哥", "en": "Monaco" },
    { "code": "MD", "name": "摩尔多瓦", "en": "Moldova" },
    { "code": "ME", "name": "黑山", "en": "Montenegro" },
    { "code": "MF", "name": "法属圣马丁", "en": "Saint Martin" },
    { "code": "MG", "name": "马达加斯加", "en": "Madagascar" },
    { "code": "MK", "name": "北马其顿", "en": "North Macedonia" },
    { "code": "ML", "name": "马里", "en": "Mali" },
    { "code": "MN", "name": "蒙古", "en": "Mongolia" },
    { "code": "MO", "name": "澳门", "en": "Macau" },
    { "code": "MQ", "name": "马提尼克", "en": "Martinique" },
    { "code": "MR", "name": "毛里塔尼亚", "en": "Mauritania" },
    { "code": "MS", "name": "蒙特塞拉特", "en": "Montserrat" },
    { "code": "MT", "name": "马耳他", "en": "Malta" },
    { "code": "MU", "name": "毛里求斯", "en": "Mauritius" },
    { "code": "MW", "name": "马拉维", "en": "Malawi" },
    { "code": "MX", "name": "墨西哥", "en": "Mexico" },
    { "code": "MY", "name": "马来西亚", "en": "Malaysia" },
    { "code": "MZ", "name": "莫桑比克", "en": "Mozambique" },
    { "code": "NA", "name": "纳米比亚", "en": "Namibia" },
    { "code": "NE", "name": "尼日尔", "en": "Niger" },
    { "code": "NG", "name": "尼日利亚", "en": "Nigeria" },
    { "code": "NI", "name": "尼加拉瓜", "en": "Nicaragua" },
    { "code": "NL", "name": "荷兰", "en": "Netherlands" },
    { "code": "NO", "name": "挪威", "en": "Norway" },
    { "code": "NP", "name": "尼泊尔", "en": "Nepal" },
    { "code": "NR", "name": "瑙鲁", "en": "Nauru" },
    { "code": "NZ", "name": "新西兰", "en": "New Zealand" },
    { "code": "OM", "name": "阿曼", "en": "Oman" },
    { "code": "PA", "name": "巴拿马", "en": "Panama" },
    { "code": "PE", "name": "秘鲁", "en": "Peru" },
    { "code": "PF", "name": "法属波利尼西亚", "en": "French Polynesia" },
    { "code": "PG", "name": "巴布亚新几内亚", "en": "Papua New Guinea" },
    { "code": "PH", "name": "菲律宾", "en": "Philippines" },
    { "code": "PK", "name": "巴基斯坦", "en": "Pakistan" },
    { "code": "PL", "name": "波兰", "en": "Poland" },
    { "code": "PR", "name": "波多黎各", "en": "Puerto Rico" },
    { "code": "PS", "name": "巴勒斯坦", "en": "Palestine" },
    { "code": "PT", "name": "葡萄牙", "en": "Portugal" },
    { "code": "PY", "name": "巴拉圭", "en": "Paraguay" },
    { "code": "QA", "name": "卡塔尔", "en": "Qatar" },
    { "code": "RE", "name": "留尼汪", "en": "Reunion" },
    { "code": "RO", "name": "罗马尼亚", "en": "Romania" },
    { "code": "RS", "name": "塞尔维亚", "en": "Serbia" },
    { "code": "RU", "name": "俄罗斯", "en": "Russia" },
    { "code": "RW", "name": "卢旺达", "en": "Rwanda" },
    { "code": "SA", "name": "沙特阿拉伯", "en": "Saudi Arabia" },
    { "code": "SC", "name": "塞舌尔", "en": "Seychelles" },
    { "code": "SD", "name": "苏丹", "en": "Sudan" },
    { "code": "SE", "name": "瑞典", "en": "Sweden" },
    { "code": "SG", "name": "新加坡", "en": "Singapore" },
    { "code": "SI", "name": "斯洛文尼亚", "en": "Slovenia" },
    { "code": "SK", "name": "斯洛伐克", "en": "Slovakia" },
    { "code": "SL", "name": "塞拉利昂", "en": "Sierra Leone" },
    { "code": "SN", "name": "塞内加尔", "en": "Senegal" },
    { "code": "SR", "name": "苏里南", "en": "Suriname" },
    { "code": "SV", "name": "萨尔瓦多", "en": "El Salvador" },
    { "code": "SZ", "name": "斯威士兰", "en": "Swaziland" },
    { "code": "TC", "name": "特克斯和凯科斯群岛", "en": "Turks and Caicos Islands" },
    { "code": "TD", "name": "乍得", "en": "Chad" },
    { "code": "TG", "name": "多哥", "en": "Togo" },
    { "code": "TH", "name": "泰国", "en": "Thailand" },
    { "code": "TJ", "name": "塔吉克斯坦", "en": "Tajikistan" },
    { "code": "TL", "name": "东帝汶", "en": "Timor-Leste" },
    { "code": "TN", "name": "突尼斯", "en": "Tunisia" },
    { "code": "TO", "name": "汤加", "en": "Tonga" },
    { "code": "TR", "name": "土耳其", "en": "Turkey" },
    { "code": "TT", "name": "特立尼达和多巴哥", "en": "Trinidad and Tobago" },
    { "code": "TW", "name": "台湾", "en": "Taiwan" },
    { "code": "TZ", "name": "坦桑尼亚", "en": "Tanzania" },
    { "code": "UA", "name": "乌克兰", "en": "Ukraine" },
    { "code": "UG", "name": "乌干达", "en": "Uganda" },
    { "code": "US", "name": "美国", "en": "United States" },
    { "code": "UY", "name": "乌拉圭", "en": "Uruguay" },
    { "code": "UZ", "name": "乌兹别克斯坦", "en": "Uzbekistan" },
    { "code": "VA", "name": "梵蒂冈", "en": "Vatican City" },
    { "code": "VC", "name": "圣文森特和格林纳丁斯", "en": "Saint Vincent and the Grenadines" },
    { "code": "VE", "name": "委内瑞拉", "en": "Venezuela" },
    { "code": "VG", "name": "英属维尔京群岛", "en": "British Virgin Islands" },
    { "code": "VI", "name": "美属维尔京群岛", "en": "U.S. Virgin Islands" },
    { "code": "VN", "name": "越南", "en": "Vietnam" },
    { "code": "VU", "name": "瓦努阿图", "en": "Vanuatu" },
    { "code": "WS", "name": "萨摩亚", "en": "Samoa" },
    { "code": "XK", "name": "科索沃", "en": "Kosovo" },
    { "code": "YT", "name": "马约特", "en": "Mayotte" },
    { "code": "ZA", "name": "南非", "en": "South Africa" },
    { "code": "ZM", "name": "赞比亚", "en": "Zambia" },
    { "code": "ZW", "name": "津巴布韦", "en": "Zimbabwe" }
  ],
  region: [
    { code: 'JP', name: '日本' },
    { code: 'KR', name: '韩国' },
    { code: 'TH', name: '泰国' },
    { code: 'VN', name: '越南' },
    { code: 'SG', name: '新加坡' },
    { code: 'MY', name: '马来西亚' },
    { code: 'PH', name: '菲律宾' },
    { code: 'ID', name: '印度尼西亚' }
  ],
  global: [
    { code: 'US', name: '美国' },
    { code: 'CA', name: '加拿大' },
    { code: 'GB', name: '英国' },
    { code: 'FR', name: '法国' },
    { code: 'DE', name: '德国' },
    { code: 'IT', name: '意大利' },
    { code: 'ES', name: '西班牙' },
    { code: 'AU', name: '澳大利亚' },
    { code: 'NZ', name: '新西兰' }
  ]
}

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

  // 热门项相关状态
  const currentHotItem = ref<HotItem | null>(null)
  const hotItemProducts = ref<Record<string, Product[]>>({})
  const hotItemPagination = ref<Record<string, {
    page: number
    pageSize: number
    total: number
    totalPages: number
    hasNext: boolean
    hasPrev: boolean
  }>>({})

  // 全球商品相关状态
  const globalProducts = ref<Product[]>([])
  const globalProductsPagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0,
    hasNext: false,
    hasPrev: false
  })
  const currentCategory = ref<'hot' | 'local' | 'region' | 'global' | null>(null)

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
      const region = product.region || 'unknown'
      if (!grouped[region]) {
        grouped[region] = []
      }
      grouped[region].push(product)
    })

    return grouped
  })

  const productsByCountry = computed(() => {
    const grouped: Record<string, Product[]> = {}

    products.value.forEach(product => {
      const countryCode = product.countryCode || 'unknown'
      if (!grouped[countryCode]) {
        grouped[countryCode] = []
      }
      grouped[countryCode].push(product)
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

  // 热门项相关方法
  const getHotItemsByCategory = (category: 'hot' | 'local' | 'region' | 'global'): HotItem[] => {
    return HOT_ITEMS_DATA[category] || []
  }

  const fetchProductsByHotItem = async (hotItemCode: string): Promise<void> => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const queryParams = {
        country: hotItemCode,
        page: 1,
        pageSize: 20
      }

      const response: PaginatedResponse<Product> = await productApi.getProductsByCountry(hotItemCode, queryParams)

      // 存储商品数据
      hotItemProducts.value[hotItemCode] = response.items

      // 存储分页信息
      hotItemPagination.value[hotItemCode] = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Products] 热门项商品获取成功:', {
        hotItemCode,
        count: response.items.length,
        total: response.total
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取热门商品失败'
      error.value = errorMessage
      console.error('[Products] 获取热门商品失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const setCurrentHotItem = (hotItem: HotItem | null): void => {
    currentHotItem.value = hotItem
  }

  const loadMoreHotItemProducts = async (hotItemCode: string): Promise<void> => {
    const currentPagination = hotItemPagination.value[hotItemCode]
    if (!currentPagination || !currentPagination.hasNext || isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const queryParams = {
        country: hotItemCode,
        page: currentPagination.page + 1,
        pageSize: currentPagination.pageSize
      }

      const response: PaginatedResponse<Product> = await productApi.getProductsByCountry(hotItemCode, queryParams)

      // 追加商品数据
      const existingProducts = hotItemProducts.value[hotItemCode] || []
      hotItemProducts.value[hotItemCode] = [...existingProducts, ...response.items]

      // 更新分页信息
      hotItemPagination.value[hotItemCode] = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Products] 加载更多热门商品成功:', {
        hotItemCode,
        page: response.page,
        count: response.items.length
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '加载更多商品失败'
      error.value = errorMessage
      console.error('[Products] 加载更多商品失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const clearHotItemProducts = (hotItemCode: string): void => {
    delete hotItemProducts.value[hotItemCode]
    delete hotItemPagination.value[hotItemCode]
  }

  const getHotItemProducts = (hotItemCode: string): Product[] => {
    return hotItemProducts.value[hotItemCode] || []
  }

  const getHotItemPagination = (hotItemCode: string) => {
    return hotItemPagination.value[hotItemCode] || {
      page: 1,
      pageSize: 20,
      total: 0,
      totalPages: 0,
      hasNext: false,
      hasPrev: false
    }
  }

  // 全球商品相关方法
  const fetchGlobalProducts = async (params?: ProductQueryParams): Promise<void> => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const queryParams = {
        ...params,
        type: 'global',
        page: 1,
        pageSize: 20
      }

      const response: PaginatedResponse<Product> = await productApi.getProductsByType('global', queryParams)

      globalProducts.value = response.items
      globalProductsPagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Products] 全球商品获取成功:', {
        count: response.items.length,
        total: response.total
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取全球商品失败'
      error.value = errorMessage
      console.error('[Products] 获取全球商品失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const loadMoreGlobalProducts = async (): Promise<void> => {
    if (!globalProductsPagination.value.hasNext || isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const queryParams = {
        type: 'global',
        page: globalProductsPagination.value.page + 1,
        pageSize: globalProductsPagination.value.pageSize
      }

      const response: PaginatedResponse<Product> = await productApi.getProductsByType('global', queryParams)

      // 追加商品数据
      globalProducts.value = [...globalProducts.value, ...response.items]

      // 更新分页信息
      globalProductsPagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: response.total,
        totalPages: response.totalPages,
        hasNext: response.hasNext,
        hasPrev: response.hasPrev
      }

      console.log('[Products] 加载更多全球商品成功:', {
        page: response.page,
        count: response.items.length
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '加载更多全球商品失败'
      error.value = errorMessage
      console.error('[Products] 加载更多全球商品失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const getGlobalProducts = (): Product[] => {
    return globalProducts.value
  }

  const getGlobalProductsPagination = () => {
    return globalProductsPagination.value
  }

  const setCurrentCategory = (category: 'hot' | 'local' | 'region' | 'global' | null): void => {
    currentCategory.value = category
  }

  const clearGlobalProducts = (): void => {
    globalProducts.value = []
    globalProductsPagination.value = {
      page: 1,
      pageSize: 20,
      total: 0,
      totalPages: 0,
      hasNext: false,
      hasPrev: false
    }
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
    currentHotItem: readonly(currentHotItem),
    hotItemProducts: readonly(hotItemProducts),
    hotItemPagination: readonly(hotItemPagination),

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
    getProductsByValidDays,

    // 热门项方法
    getHotItemsByCategory,
    fetchProductsByHotItem,
    setCurrentHotItem,
    loadMoreHotItemProducts,
    clearHotItemProducts,
    getHotItemProducts,
    getHotItemPagination,

    // 全球商品方法
    fetchGlobalProducts,
    loadMoreGlobalProducts,
    getGlobalProducts,
    getGlobalProductsPagination,
    setCurrentCategory,
    clearGlobalProducts
  }
})