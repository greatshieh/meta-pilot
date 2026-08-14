import path from 'node:path'
import type { Preset } from '@unocss/core'
import type { Theme } from '@unocss/preset-uno'
import { presetIcons } from 'unocss'
import { loadEnv } from 'vite'
import { FileSystemIconLoader } from 'unplugin-icons/loaders'

/**
 * 创建Unocss预设配置
 * 功能：定义项目中常用的CSS工具类快捷方式
 * @returns 返回配置好的Unocss预设对象
 */
export function presetUniAdmin(): Preset<Theme> {
    const preset: Preset<Theme> = {
        name: 'preset-uni-admin', // 预设名称
        shortcuts: [
            // 第一组：Flex布局相关快捷方式
            {
                'flex-center': 'flex justify-center items-center', // 水平垂直居中
                'flex-x-center': 'flex justify-center', // 水平居中
                'flex-y-center': 'flex items-center', // 垂直居中
                'flex-col': 'flex flex-col', // 垂直排列
                'flex-col-center': 'flex-center flex-col', // 垂直排列+居中
                'flex-col-stretch': 'flex-col items-stretch', // 垂直排列+拉伸子项
                'i-flex-center': 'inline-flex justify-center items-center', // 行内元素水平垂直居中
                'i-flex-x-center': 'inline-flex justify-center', // 行内元素水平居中
                'i-flex-y-center': 'inline-flex items-center', // 行内元素垂直居中
                'i-flex-col': 'flex-col inline-flex', // 行内元素垂直排列
                'i-flex-col-center': 'flex-col i-flex-center', // 行内元素垂直排列+居中
                'i-flex-col-stretch': 'i-flex-col items-stretch', // 行内元素垂直排列+拉伸
                'flex-1-hidden': 'flex-1 overflow-hidden' // 弹性填充+隐藏溢出
            },
            // 第二组：定位相关快捷方式
            {
                'absolute-lt': 'absolute left-0 top-0', // 绝对定位左上
                'absolute-lb': 'absolute left-0 bottom-0', // 绝对定位左下
                'absolute-rt': 'absolute right-0 top-0', // 绝对定位右上
                'absolute-rb': 'absolute right-0 bottom-0', // 绝对定位右下
                'absolute-tl': 'absolute-lt', // 同左上(别名)
                'absolute-tr': 'absolute-rt', // 同右上(别名)
                'absolute-bl': 'absolute-lb', // 同左下(别名)
                'absolute-br': 'absolute-rb', // 同右下(别名)
                'absolute-center': 'absolute-lt flex-center size-full', // 绝对定位居中
                'fixed-lt': 'fixed left-0 top-0', // 固定定位左上
                'fixed-lb': 'fixed left-0 bottom-0', // 固定定位左下
                'fixed-rt': 'fixed right-0 top-0', // 固定定位右上
                'fixed-rb': 'fixed right-0 bottom-0', // 固定定位右下
                'fixed-tl': 'fixed-lt', // 同左上(别名)
                'fixed-tr': 'fixed-rt', // 同右上(别名)
                'fixed-bl': 'fixed-lb', // 同左下(别名)
                'fixed-br': 'fixed-rb', // 同右下(别名)
                'fixed-center': 'fixed-lt flex-center size-full' // 固定定位居中
            },
            // 第三组：文本溢出处理快捷方式
            {
                'nowrap-hidden': 'overflow-hidden whitespace-nowrap', // 不换行+隐藏溢出
                'ellipsis-text': 'nowrap-hidden text-ellipsis' // 单行文本溢出显示省略号
            },
            {
                'card-wrapper': 'rd-10px shadow-sm'
            }
        ]
    }

    return preset
}

export function presetIcon(): Preset<Theme> {
    const env = loadEnv(process.env.NODE_ENV as string, process.cwd()) as Env.ImportMeta
    // 从环境变量中解构图标前缀配置
    const { VITE_ICON_PREFIX, VITE_ICON_LOCAL_PREFIX, VITE_ICON_SVG_ASSET } = env
    // 获取本地图标存放路径
    const localIconPath = path.join(process.cwd(), VITE_ICON_SVG_ASSET)

    /** 本地图标集合的名称，通过移除前缀得到 */
    const collectionName = VITE_ICON_LOCAL_PREFIX.replace(`${VITE_ICON_PREFIX}-`, '')

    return presetIcons({
        // prefix: `${VITE_ICON_PREFIX}-`, // 图标类名前缀
        scale: 1, // 图标缩放比例
        extraProperties: {
            // 额外的CSS属性
            display: 'inline-block'
        },
        collections: {
            // 自定义图标集合配置
            [collectionName]: FileSystemIconLoader(localIconPath, svg => {
                return svg.replace(/^<svg\s/, '<svg width="1em" height="1em" ') // 处理SVG属性
            })
        },
        warn: true // 启用警告信息
    })
}
