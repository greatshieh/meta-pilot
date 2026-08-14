// @unocss-include

import { DARK_CLASS } from '@/constants/app'
import { SetupStoreId } from '@/store/enum'
import { toggleHtmlClass } from '@/utils/common'

/**
 * 设置应用加载动画
 * 功能：在应用初始化时显示一个带有品牌Logo和动画效果的加载界面
 */
export function setupLoading() {
    // 从本地存储获取主题配置
    const themeStorage = localStorage.getItem(SetupStoreId.Theme)
    const themeStore = themeStorage ? JSON.parse(themeStorage) : {}

    // 如果当前是暗色模式，添加暗色类名
    if (themeStore?.isDark) {
        toggleHtmlClass(DARK_CLASS).add()
    }

    // 定义四个动画点的位置和动画延迟
    const loadingClasses = [
        'left-0 top-0',
        'left-0 bottom-0 animate-delay-500',
        'right-0 top-0 animate-delay-1000',
        'right-0 bottom-0 animate-delay-1500'
    ]

    // 创建四个动画点元素
    const dot = loadingClasses
        .map(item => {
            return `<div class="absolute w-16px h-16px bg-primary rounded-8px animate-pulse ${item}"></div>`
        })
        .join('\n')

    // 构建完整的加载动画HTML结构
    const loading = `
  <div class="fixed-center flex-col bg-page">
    <img src="/src/assets/logo-transparent-circle.png" alt="logo" class="size-128px" />
    <div class="w-56px h-56px my-36px">
      <div class="relative h-full animate-spin">
        ${dot}
      </div>
    </div>
    <h2 class="text-28px font-500 text-primary">${import.meta.env.VITE_APP_TITLE}</h2>
  </div>`

    // 获取应用根元素并注入加载动画
    const app = document.getElementById('app')

    if (app) {
        app.innerHTML = loading
    }
}
