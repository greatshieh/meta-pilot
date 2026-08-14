import type { Plugin, PluginOption } from 'vite'

/**
 * 创建一个Vite插件，用于在HTML的<head>中注入构建时间信息
 * @param buildTime - 构建时间字符串，通常是一个时间戳或格式化日期
 * @returns 返回配置好的Vite插件对象
 */
export function setupHtmlPlugin(buildTime: string): PluginOption {
    // 创建Vite插件配置对象
    const plugin: PluginOption = {
        name: 'html-plugin', // 插件名称
        apply: 'build', // 只在生产构建时应用
        transformIndexHtml(html) {
            // 在<head>标签后插入包含构建时间的meta标签
            return html.replace('<head>', `<head>\n    <meta name="buildTime" content="${buildTime}">`)
        }
    }

    return plugin
}
