<template>
    <component :is="activeMenu" />
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import { computed } from 'vue'
import useThemeStore from '@/store/modules/theme/index'
import HorizontalMenu from './modules/HorizontalMenu.vue'
import HorizontalMixMenu from './modules/HorizontalMixMenu.vue'
import VerticalMenu from './modules/VerticalMenu.vue'
import VerticalMixMenu from './modules/VerticalMixMenu.vue'

defineOptions({
    name: 'GlobalMenu'
})

const themeStore = useThemeStore()

const activeMenu = computed(() => {
    const menuMap: Record<UnionKey.ThemeLayoutMode, Component> = {
        vertical: VerticalMenu,
        'vertical-mix': VerticalMixMenu,
        horizontal: HorizontalMenu,
        'horizontal-mix': HorizontalMixMenu
    }

    return menuMap[themeStore.layout.mode]
})
</script>

<style scoped></style>
