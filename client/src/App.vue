<template>
    <ElConfigProvider :locale="zhCn">
        <div class="h-full">
            <RouterView />
        </div>
    </ElConfigProvider>
</template>

<script setup lang="ts" name="App">
import { zhCn } from 'element-plus/es/locale/index.mjs'
import useThemeStore from '@/store/modules/theme'

useThemeStore()
</script>
