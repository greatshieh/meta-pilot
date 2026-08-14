import { defu } from 'defu'
import { DARK_CLASS } from '@/constants/app'
import { overrideThemeSettings, themeSettings } from '@/theme/settings'
import { toggleHtmlClass } from '@/utils/common'
import { localStg } from '@/utils/storage'

/**
 * 初始化主题设置
 * 功能：根据环境变量决定使用本地缓存设置还是默认设置
 * 开发模式：直接使用默认主题设置
 * 生产模式：优先使用本地存储的设置，并应用覆盖设置
 */
export function initThemeSettings() {
    const isProd = import.meta.env.NODE_ENV !== 'development'

    // 开发环境直接返回默认设置
    if (!isProd) return themeSettings

    // 获取本地存储的设置
    const localSettings = localStg.get('themeSettings')
    // 合并本地和默认设置
    let settings = defu(localSettings, themeSettings)
    // 应用覆盖设置
    settings = defu(overrideThemeSettings, settings)
    return settings
}

/**
 * 切换CSS暗黑模式
 * 功能：根据参数动态添加或移除暗黑模式类名
 * @param darkMode - 布尔值，true表示启用暗黑模式，false表示禁用
 */
export function toggleCssDarkMode(darkMode = false) {
    // 获取添加/移除类名的方法
    const { add, remove } = toggleHtmlClass(DARK_CLASS)
    // 根据参数调用对应方法
    if (darkMode) {
        add()
    } else {
        remove()
    }
}
