<template>
    <el-divider>页面功能</el-divider>
    <div class="flex flex-col items-stretch gap-12px">
        <SettingItem label="滚动模式">
            <ElSelect v-model="themeStore.layout.scrollMode" size="small" class="!w-120px">
                <ElOption v-for="{ label, value } in translateOptions(themeScrollModeOptions)" :key="value" :label="label" :value="value" />
            </ElSelect>
        </SettingItem>
        <SettingItem label="页面切换">
            <ElSwitch v-model="themeStore.page.animate" />
        </SettingItem>
        <SettingItem v-if="themeStore.page.animate" label="页面切换动画">
            <ElSelect v-model="themeStore.page.animateMode" size="small" class="!w-120px">
                <ElOption
                    v-for="{ label, value } in translateOptions(themePageAnimationModeOptions)"
                    :key="value"
                    :label="label"
                    :value="value" />
            </ElSelect>
        </SettingItem>
        <SettingItem v-if="isWrapperScrollMode" label="固定头部和标签栏">
            <ElSwitch v-model="themeStore.fixedHeaderAndTab" />
        </SettingItem>
        <SettingItem label="头部高度">
            <el-input-number v-model="themeStore.header.height" :step="1" size="small" class="!w-120px" />
        </SettingItem>

        <SettingItem label="显示面包屑">
            <el-switch v-model="themeStore.header.breadcrumb.visible" size="small" />
        </SettingItem>

        <SettingItem label="显示标签页">
            <el-switch v-model="themeStore.tab.visible" size="small" />
        </SettingItem>

        <SettingItem label="标签栏风格">
            <ElSelect v-model="themeStore.tab.mode" size="small" class="!w-120px">
                <ElOption v-for="{ label, value } in translateOptions(themeTabModeOptions)" :key="value" :label="label" :value="value" />
            </ElSelect>
        </SettingItem>

        <SettingItem label="标签栏高度">
            <el-input-number v-model="themeStore.tab.height" :step="1" size="small" class="!w-120px" />
        </SettingItem>

        <SettingItem v-if="layoutMode === 'vertical'" key="6-1" label="侧边栏宽度">
            <el-input-number v-model="themeStore.sider.width" size="small" :step="1" class="w-120px" />
        </SettingItem>

        <SettingItem v-if="layoutMode === 'vertical'" key="6-2" label="侧边栏折叠宽度">
            <el-input-number v-model="themeStore.sider.collapsedWidth" size="small" :step="1" class="w-120px" />
        </SettingItem>

        <SettingItem label="显示底部">
            <ElSwitch v-model="themeStore.footer.visible" />
        </SettingItem>

        <SettingItem v-if="themeStore.footer.visible && isWrapperScrollMode" label="固定底部">
            <ElSwitch v-model="themeStore.footer.fixed" />
        </SettingItem>
        <SettingItem v-if="themeStore.footer.visible" label="底部高度">
            <ElInputNumber v-model="themeStore.footer.height" size="small" :step="1" class="w-120px" />
        </SettingItem>
        <SettingItem v-if="themeStore.footer.visible && layoutMode === 'horizontal-mix'" label="底部靠右">
            <ElSwitch v-model="themeStore.footer.right" />
        </SettingItem>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { themePageAnimationModeOptions, themeScrollModeOptions, themeTabModeOptions } from '@/constants/app'
import useThemeStore from '@/store/modules/theme'
import { translateOptions } from '@/utils/common'

defineOptions({ name: 'PageFun' })

const themeStore = useThemeStore()

const layoutMode = computed(() => themeStore.layout.mode)
const isWrapperScrollMode = computed(() => themeStore.layout.scrollMode === 'wrapper')
</script>
