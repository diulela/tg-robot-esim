import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'

// 浅色主题配置
const lightTheme = {
  dark: false,
  colors: {
    primary: '#6366F1',      // 紫色主色调 (与原型图一致)
    secondary: '#EC4899',    // 粉色辅助色
    accent: '#10B981',       // 绿色强调色
    error: '#EF4444',        // 红色错误色
    warning: '#F59E0B',      // 黄色警告色
    info: '#3B82F6',         // 蓝色信息色
    success: '#10B981',      // 绿色成功色
    surface: '#FFFFFF',      // 表面色
    background: '#F8FAFC',   // 背景色
    'on-primary': '#FFFFFF',
    'on-secondary': '#FFFFFF',
    'on-surface': '#1F2937',
    'on-background': '#1F2937',
  }
}

// 深色主题配置
const darkTheme = {
  dark: true,
  colors: {
    primary: '#818CF8',      // 浅紫色
    secondary: '#F472B6',    // 浅粉色
    accent: '#34D399',       // 浅绿色
    error: '#F87171',        // 浅红色
    warning: '#FBBF24',      // 浅黄色
    info: '#60A5FA',         // 浅蓝色
    success: '#34D399',      // 浅绿色
    surface: '#1E293B',      // 深色表面
    background: '#0F172A',   // 深色背景
    'on-primary': '#1E1B4B',
    'on-secondary': '#831843',
    'on-surface': '#F1F5F9',
    'on-background': '#F1F5F9',
  }
}

export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: lightTheme,
      dark: darkTheme,
    },
    variations: {
      colors: ['primary', 'secondary', 'accent'],
      lighten: 4,
      darken: 4,
    },
  },
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi },
  },
  display: {
    mobileBreakpoint: 'sm',
    thresholds: {
      xs: 0,
      sm: 600,
      md: 960,
      lg: 1280,
      xl: 1920,
    },
  },
  defaults: {
    VBtn: {
      style: 'text-transform: none;', // 禁用按钮文字大写
      rounded: 'lg',
    },
    VCard: {
      rounded: 'lg',
      elevation: 2,
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VAutocomplete: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VTextarea: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VChip: {
      rounded: 'lg',
    },
    VAppBar: {
      elevation: 0,
    },
  },
})