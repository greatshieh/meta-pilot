import ElementPlus from 'element-plus'
import type { App } from 'vue'

/**
 * 设置Element Plus UI组件库
 * 功能：全局注册Element Plus组件到Vue应用
 * @param app - Vue应用实例
 */
export function setupUI(app: App) {
    // 使用Element Plus插件
    app.use(ElementPlus)
}
