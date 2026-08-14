import path from 'node:path'
import process from 'node:process'
import AutoImport from 'unplugin-auto-import/vite'
import { FileSystemIconLoader } from 'unplugin-icons/loaders'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import type { PluginOption } from 'vite'

/**
 * 配置Vite的Unplugin插件集合
 * @param viteEnv - Vite环境变量配置对象，包含图标相关配置
 * @returns 返回配置好的插件数组
 */
export function setupUnplugin(viteEnv: Env.ImportMeta): PluginOption[] {
    // 从环境变量中解构图标前缀配置
    const { VITE_ICON_PREFIX, VITE_ICON_LOCAL_PREFIX, VITE_ICON_SVG_ASSET } = viteEnv

    // 获取本地图标存放路径
    const localIconPath = path.join(process.cwd(), VITE_ICON_SVG_ASSET)

    /** 本地图标集合的名称，通过移除前缀得到 */
    const collectionName = VITE_ICON_LOCAL_PREFIX.replace(`${VITE_ICON_PREFIX}-`, '')

    // 配置插件数组
    const plugins: PluginOption[] = [
        // 配置unplugin-icons插件
        Icons({
            compiler: 'vue3', // 指定Vue3编译器
            customCollections: {
                // 自定义图标集合
                [collectionName]: FileSystemIconLoader(
                    localIconPath, // 本地图标路径
                    svg => svg.replace(/^<svg\s/, '<svg width="1em" height="1em" ') // 处理SVG属性
                )
            },
            scale: 1, // 图标缩放比例
            defaultClass: 'inline-block' // 默认CSS类
        }),
        // 配置自动导入组件插件
        Components({
            dts: 'src/types/components.d.ts', // 类型声明文件路径
            dirs: ['src/components', 'src/**/components'],
            resolvers: [
                // 配置Element Plus组件自动导入解析器
                ElementPlusResolver({
                    importStyle: 'sass' // 不自动导入样式
                }),
                // 配置图标组件解析器
                IconsResolver({
                    customCollections: [collectionName], // 使用的自定义图标集合
                    componentPrefix: VITE_ICON_PREFIX // 组件前缀
                })
            ]
        }),
        // 配置自动导入插件
        AutoImport({
            // 配置解析器数组
            resolvers: [
                // Element Plus组件自动导入解析器，不自动导入样式
                ElementPlusResolver({ importStyle: 'sass' }),
                // 图标组件解析器配置
                IconsResolver({
                    // 指定使用的自定义图标集合名称
                    customCollections: [collectionName],
                    // 设置图标组件前缀为'icon'
                    componentPrefix: 'i'
                })
            ],
            // 指定自动导入类型声明文件的输出路径
            dts: path.resolve(__dirname, 'src/types/auto-imports.d.ts')
        })
    ]

    return plugins
}
