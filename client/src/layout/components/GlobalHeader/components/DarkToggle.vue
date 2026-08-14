<template>
    <div @click="toggle">
        <SvgIcon :icon="themeStore.isDark ? 'material-symbols:nightlight-rounded' : 'material-symbols:sunny'" />
    </div>
</template>

<script setup lang="ts">
import useThemeStore from '@/store/modules/theme'

defineOptions({ name: 'DarkToggle' })

const themeStore = useThemeStore()
function toggle() {
    themeStore.setThemeScheme(themeStore.isDark ? 'light' : 'dark')
}
</script>

<style lang="scss" scoped></style>
