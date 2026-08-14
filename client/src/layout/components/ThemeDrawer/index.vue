<template>
    <el-drawer v-model="appStore.themeDrawerVisible" title="主题配置" :size="360">
        <DarkMode />
        <LayoutMode />
        <PageFun />
    </el-drawer>
</template>

<script setup lang="ts">
import useAppStore from '@/store/modules/app'

const appStore = useAppStore()
</script>

<style scoped></style>
