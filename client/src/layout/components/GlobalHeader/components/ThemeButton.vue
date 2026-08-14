<template>
    <div @click="toggle">
        <SvgIcon icon="bxs:color" />
    </div>
</template>

<script setup lang="ts">
import useAppStore from '@/store/modules/app'

defineOptions({ name: 'ThemeButton' })

const appStore = useAppStore()

function toggle() {
    appStore.themeDrawerVisible = !appStore.themeDrawerVisible
}
</script>

<style lang="scss" scoped></style>
