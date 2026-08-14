import { useElementSize } from '@vueuse/core'
import type { HeatmapSeriesOption } from 'echarts'
import type {
    BarSeriesOption,
    GaugeSeriesOption,
    LineSeriesOption,
    PictorialBarSeriesOption,
    PieSeriesOption,
    RadarSeriesOption,
    ScatterSeriesOption
} from 'echarts/charts'
import { BarChart, GaugeChart, LineChart, PictorialBarChart, PieChart, RadarChart, ScatterChart } from 'echarts/charts'
import type {
    DatasetComponentOption,
    GridComponentOption,
    LegendComponentOption,
    TitleComponentOption,
    ToolboxComponentOption,
    TooltipComponentOption
} from 'echarts/components'
import {
    DatasetComponent,
    GridComponent,
    LegendComponent,
    TitleComponent,
    ToolboxComponent,
    TooltipComponent,
    TransformComponent
} from 'echarts/components'
import * as echarts from 'echarts/core'
import { LabelLayout, UniversalTransition } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'
import { computed, effectScope, nextTick, onScopeDispose, ref, watch } from 'vue'
import useThemeStore from '@/store/modules/theme'
import { darkTheme, lightTheme } from '@/styles/echartsTheme'

export type ECOption = echarts.ComposeOption<
    | BarSeriesOption
    | LineSeriesOption
    | PieSeriesOption
    | ScatterSeriesOption
    | PictorialBarSeriesOption
    | RadarSeriesOption
    | GaugeSeriesOption
    | TitleComponentOption
    | LegendComponentOption
    | TooltipComponentOption
    | GridComponentOption
    | ToolboxComponentOption
    | DatasetComponentOption
    | HeatmapSeriesOption
>

echarts.use([
    TitleComponent,
    LegendComponent,
    TooltipComponent,
    GridComponent,
    DatasetComponent,
    TransformComponent,
    ToolboxComponent,
    BarChart,
    LineChart,
    PieChart,
    ScatterChart,
    PictorialBarChart,
    RadarChart,
    GaugeChart,
    LabelLayout,
    UniversalTransition,
    CanvasRenderer
])

interface ChartHooks {
    onRender?: (chart: echarts.ECharts) => void | Promise<void>
    onUpdated?: (chart: echarts.ECharts) => void | Promise<void>
    onDestroy?: (chart: echarts.ECharts) => void | Promise<void>
}

echarts.registerTheme('light', lightTheme)
echarts.registerTheme('dark', darkTheme)

/**
 * 自定义ECharts钩子函数
 * @param optionsFactory - 生成ECharts配置项的函数
 * @param hooks - 包含生命周期回调函数的对象
 * @returns 返回包含DOM引用和操作方法的对象
 */
export function useEcharts<T extends ECOption>(optionsFactory: () => T, hooks: ChartHooks = {}) {
    // 创建effect作用域，用于管理副作用
    const scope = effectScope()

    // 获取主题状态
    const themeStore = useThemeStore()
    const darkMode = computed(() => themeStore.isDark)

    // DOM引用和尺寸状态
    const domRef = ref<HTMLElement | null>(null)
    const initialSize = { width: 0, height: 0 }
    const { width, height } = useElementSize(domRef, initialSize)

    // ECharts实例和配置项
    let chart: echarts.ECharts | null = null
    const chartOptions: T = optionsFactory()

    // 默认生命周期钩子
    const {
        // 渲染时的默认回调：显示加载状态
        onRender = instance => {
            const textColor = darkMode.value ? 'rgb(224, 224, 224)' : 'rgb(31, 31, 31)'
            const maskColor = darkMode.value ? 'rgba(0, 0, 0, 0.4)' : 'rgba(255, 255, 255, 0.8)'

            instance.showLoading({
                color: themeStore.themeColor,
                textColor,
                fontSize: 14,
                maskColor
            })
        },
        // 更新后的默认回调：隐藏加载状态
        onUpdated = instance => {
            instance.hideLoading()
        },
        // 销毁时的回调
        onDestroy
    } = hooks

    /**
     * 检查是否可以渲染图表
     * @returns 返回布尔值，表示是否满足渲染条件
     */
    function canRender() {
        return domRef.value && initialSize.width > 0 && initialSize.height > 0
    }

    /**
     * 检查图表是否已渲染
     * @returns 返回布尔值，表示图表是否已渲染
     */
    function isRendered() {
        return Boolean(domRef.value && chart)
    }

    /**
     * 更新图表配置项
     * @param callback - 用于生成新配置项的回调函数
     */
    async function updateOptions(callback: (opts: T, optsFactory: () => T) => ECOption = () => chartOptions) {
        if (!isRendered()) return

        const updatedOpts = callback(chartOptions, optionsFactory)
        Object.assign(chartOptions, updatedOpts)

        if (isRendered()) {
            chart?.clear()
        }

        chart?.setOption({ ...updatedOpts, backgroundColor: 'transparent' })
        if (chart) await onUpdated?.(chart)
    }

    /**
     * 直接设置图表配置项
     * @param options - 新的配置项
     */
    function setOptions(options: T) {
        chart?.setOption(options)
    }

    /**
     * 渲染图表
     */
    async function render() {
        if (!isRendered()) {
            const chartTheme = darkMode.value ? 'dark' : 'light'
            await nextTick()
            chart = echarts.init(domRef.value, chartTheme)
            chart.setOption({
                ...chartOptions,
                backgroundColor: 'transparent'
            })
            await onRender?.(chart)
        }
    }

    /**
     * 调整图表尺寸
     */
    function resize() {
        chart?.resize()
    }

    /**
     * 销毁图表实例
     */
    async function destroy() {
        if (!chart) return
        await onDestroy?.(chart)
        chart?.dispose()
        chart = null
    }

    /**
     * 切换图表主题
     */
    async function changeTheme() {
        await destroy()
        await render()
        if (chart) await onUpdated?.(chart)
    }

    /**
     * 根据尺寸渲染图表
     * @param w - 宽度
     * @param h - 高度
     */
    async function renderChartBySize(w: number, h: number) {
        initialSize.width = w
        initialSize.height = h

        if (!canRender()) {
            await destroy()
            return
        }

        if (isRendered()) {
            resize()
        }

        await render()
    }

    // 设置响应式监听
    scope.run(() => {
        // 监听尺寸变化
        watch([width, height], ([newWidth, newHeight]) => {
            renderChartBySize(newWidth, newHeight)
        })

        // 监听暗黑模式变化
        watch(darkMode, () => {
            changeTheme()
        })
    })

    // 组件卸载时清理
    onScopeDispose(() => {
        destroy()
        scope.stop()
    })

    // 返回公共方法和引用
    return {
        domRef,
        updateOptions,
        setOptions
    }
}
