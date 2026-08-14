declare namespace Env {
    interface ImportMeta extends ImportMetaEnv {
        /** 当前环境模式：development/production */
        readonly NODE_ENV?: string
        /** 基础URL地址 */
        readonly VITE_BASE_URL?: string
        /** 应用标题 */
        readonly VITE_APP_TITLE?: string
        /** 后端服务器基础URL地址 */
        readonly VITE_SERVICE_BASE_URL: string
        /** API前缀 */
        readonly VITE_API_PREFIX?: string
        /** 是否启用HTTP代理 (Y/N) */
        readonly VITE_HTTP_PROXY?: string
        /** 客户端开发服务器端口 */
        readonly VITE_CLI_PORT?: number
        /** 服务端端口 */
        readonly VITE_SERVER_PORT?: number
        /** 服务请求成功状态码 */
        readonly VITE_SERVICE_SUCCESS_CODE: number
        /** 本地图标前缀 */
        readonly VITE_ICON_LOCAL_PREFIX: string
        /** 图标组件前缀 */
        readonly VITE_ICON_PREFIX: string
        /** SVG图标资源路径 */
        readonly VITE_ICON_SVG_ASSET: string
    }
}

/** 扩展ImportMeta接口，添加env属性 */
interface ImportMeta {
    readonly env: Env.ImportMeta
}
