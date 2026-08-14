<template>
    <el-divider>主题配置</el-divider>
    <div class="flex flex-col items-stretch gap-16px">
        <div class="inline-flex items-center justify-center">
            <ElTabs v-model="themeStore.themeScheme" type="border-card" class="segment" @tab-change="handleSegmentChange">
                <ElTabPane v-for="(_, key) in themeSchemaRecord" :key="key" :name="key">
                    <template #label>
                        <SvgIcon :icon="icons[key]" class="h-23px text-icon-small" />
                    </template>
                </ElTabPane>
            </ElTabs>
        </div>
    </div>
</template>

<script setup lang="ts">
import { themeSchemaRecord } from '@/constants/app'
import useThemeStore from '@/store/modules/theme'

defineOptions({ name: 'DarkMode' })

const themeStore = useThemeStore()

const icons: Record<UnionKey.ThemeScheme, string> = {
    light: 'material-symbols:sunny',
    dark: 'material-symbols:nightlight-rounded'
}

function handleSegmentChange(value: string | number) {
    themeStore.setThemeScheme(value as UnionKey.ThemeScheme)
}
</script>
