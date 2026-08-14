<template>
    <div class="bg-page h-full flex items-center px-12px shadow-header transition-300">
        <GlobalLogo v-if="showLogo" class="h-full" :style="{ width: themeStore.sider.width + 'px' }" />
        <Hamburger v-if="showMenuToggler" id="hamburger-container" :collapsed="appStore.siderCollapse" @click="appStore.toggleSidebar" />
        <div v-if="showMenu" :id="GLOBAL_HEADER_MENU_ID" class="h-full flex items-center flex-1 overflow-hidden"></div>
        <div v-else class="h-full flex items-center flex-1 overflow-hidden">
            <GlobalBreadcrumb v-if="!appStore.isMobile" id="breadcrumb-container" class="ml-12px" />
        </div>
        <div class="h-full flex items-center justify-end">
            <ElTooltip content="屏幕模式">
                <Fullscreen v-if="appStore.device === 'desktop'" class="text-18px cursor-pointer min-w-40px text-center" />
            </ElTooltip>

            <ElTooltip :content="themeStore.isDark ? '暗黑模式' : '明亮模式'">
                <DarkToggle id="dartheme" class="text-18px cursor-pointer min-w-40px text-center" />
            </ElTooltip>

            <ElTooltip content="主题配置">
                <ThemeButton class="text-18px cursor-pointer min-w-40px text-center" />
            </ElTooltip>
            <UserProfile />
        </div>
    </div>
</template>

<script setup lang="ts">
import { GLOBAL_HEADER_MENU_ID } from '@/constants/app'
import useAppStore from '@/store/modules/app'
import useThemeStore from '@/store/modules/theme'

defineOptions({ name: 'GlobalHeader' })

interface Props {
    /** Whether to show the logo */
    showLogo?: App.Global.HeaderProps['showLogo']
    /** Whether to show the menu toggler */
    showMenuToggler?: App.Global.HeaderProps['showMenuToggler']
    /** Whether to show the menu */
    showMenu?: App.Global.HeaderProps['showMenu']
}

defineProps<Props>()

const appStore = useAppStore()
const themeStore = useThemeStore()
// const { isFullscreen, toggle } = useFullscreen();
</script>

<style scoped></style>
