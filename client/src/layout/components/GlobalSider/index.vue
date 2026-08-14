<script setup lang="ts">
import { computed } from 'vue'
import { GLOBAL_SIDER_MENU_ID } from '@/constants/app'
import useAppStore from '@/store/modules/app'
import useThemeStore from '@/store/modules/theme'

defineOptions({ name: 'GlobalSider' })

const appStore = useAppStore()
const themeStore = useThemeStore()

const isVerticalMix = computed(() => themeStore.layout.mode === 'vertical-mix')
const isHorizontalMix = computed(() => themeStore.layout.mode === 'horizontal-mix')
// const darkMenu = computed(() => !themeStore.darkMode && !isHorizontalMix.value && themeStore.sider.inverted);
const showLogo = computed(() => !isVerticalMix.value && !isHorizontalMix.value)
const menuWrapperClass = computed(() => (showLogo.value ? 'flex-1 overflow-hidden' : 'h-full'))
</script>

<template>
    <div class="size-full flex-col items-stretch shadow-sider">
        <GlobalLogo v-if="showLogo" :show-title="!appStore.siderCollapse" :style="{ height: themeStore.header.height + 'px' }" />
        <div :id="GLOBAL_SIDER_MENU_ID" :class="menuWrapperClass" />
    </div>
</template>

<style scoped></style>
