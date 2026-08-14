import path from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import { setupVitePlugins } from './build/plugins'
import { createViteProxy, getBuildTime } from './build/plugins/config'

// 加载环境变量，根据当前NODE_ENV获取对应环境配置
const env = loadEnv(process.env.NODE_ENV as string, process.cwd()) as Env.ImportMeta

// 导出Vite配置
export default defineConfig(configEnv => {
    // 根据当前构建模式加载环境变量
    const viteEnv = loadEnv(configEnv.mode, process.cwd()) as Env.ImportMeta

    // 获取构建时间戳
    const buildTime = getBuildTime()

    // 判断是否启用代理（仅在开发服务器且非预览模式下启用）
    const enableProxy = configEnv.command === 'serve' && !configEnv.isPreview

    return {
        // 基础公共路径
        base: viteEnv.VITE_BASE_URL,
        // 解析配置
        resolve: {
            // 路径别名
            alias: {
                '@': path.resolve(__dirname, 'src') // 将@映射到src目录
            }
        },
        // CSS预处理器配置
        css: {
            preprocessorOptions: {
                scss: {
                    api: 'modern-compiler', // 使用现代编译器API
                    additionalData: `@use "@/styles/element-variables.scss" as *;`
                }
            }
        },
        // 插件配置
        plugins: setupVitePlugins(viteEnv, buildTime),
        // 开发服务器配置
        server: {
            open: true, // 自动打开浏览器
            port: env.VITE_CLI_PORT || 9527, // 服务器端口
            proxy:{
                '/api/v1': 'http://example.com/'
            },
            // host: 'http://example.com', // 服务器主机
            // proxy: createViteProxy(viteEnv, enableProxy) // 代理配置
        },
        // 依赖优化配置
        optimizeDeps: {
            exclude: ['@wiris/mathtype-html-integration-devkit'] // 排除特定依赖的优化
        }
    }
})
