import { addAPIProvider } from '@iconify/vue'

/**
 * 设置Iconify离线模式
 * 用于配置Iconify图标库的离线资源URL
 */
export function setupIconifyOffline() {
    // 从环境变量中获取Iconify资源URL
    const { VITE_ICONIFY_URL } = import.meta.env

    // 如果配置了Iconify资源URL，则添加API提供者
    if (VITE_ICONIFY_URL) {
        addAPIProvider('', { resources: [VITE_ICONIFY_URL] })
    }
}
