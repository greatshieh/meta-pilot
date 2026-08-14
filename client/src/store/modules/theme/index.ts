import { defineStore } from 'pinia'
import { computed, effectScope, onScopeDispose, type Ref, ref, toRefs, watch } from 'vue'
import { SetupStoreId } from '../../enum'
import { initThemeSettings, toggleCssDarkMode } from './shared'

const useThemeStore = defineStore(SetupStoreId.Theme, () => {
    // 创建effect作用域，用于管理副作用
    const scope = effectScope()

    // 初始化主题设置，使用ref创建响应式引用
    const settings: Ref<App.Theme.ThemeSetting> = ref(initThemeSettings())

    /**
     * 暗黑模式计算属性
     * 功能：根据当前主题方案判断是否为暗黑模式
     * 返回：布尔值，true表示暗黑模式，false表示亮色模式
     */
    const isDark = computed(() => {
        return settings.value.themeScheme === 'dark'
    })

    /**
     * 设置水平混合布局的反向模式
     * 功能：控制水平混合布局的排列方向（正向/反向）
     * @param reverse - 是否启用反向模式
     */
    function setLayoutReverseHorizontalMix(reverse: boolean) {
        settings.value.layout.reverseHorizontalMix = reverse
    }

    /**
     * 设置主题配色方案
     * 功能：切换应用的主题模式（亮色/暗色）
     * @param themeScheme - 主题模式，可选值为 'light' 或 'dark'
     */
    function setThemeScheme(themeScheme: UnionKey.ThemeScheme) {
        settings.value.themeScheme = themeScheme
    }

    // 在effect作用域内设置状态监听
    scope.run(() => {
        // 监听暗黑模式状态变化
        watch(
            isDark, // 监听的暗黑模式计算属性
            val => {
                toggleCssDarkMode(val) // 切换CSS暗黑模式类名
            },
            { immediate: true } // 立即执行一次回调
        )
    })

    /**
     * 作用域清理回调
     * 功能：在组件/存储卸载时自动停止所有监听和副作用
     */
    onScopeDispose(() => {
        scope.stop() // 停止effectScope内的所有响应式效果
    })

    return {
        ...toRefs(settings.value),
        isDark,
        setThemeScheme,
        setLayoutReverseHorizontalMix
    }
})

export default useThemeStore
