<template>
    <template v-if="isLocalIcon">
        <svg aria-hidden="true" width="1em" height="1em" v-bind="bindAttrs">
            <use :xlink:href="symbolId" fill="currentColor" />
        </svg>
    </template>
    <template v-else>
        <Icon v-if="icon" :icon="icon" v-bind="bindAttrs" />
    </template>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { computed, useAttrs } from 'vue'

defineOptions({ name: 'SvgIcon', inheritAttrs: false })

interface Props {
    icon: string
}

const props = defineProps<Props>()

const attrs = useAttrs()

const bindAttrs = computed<{ class: string; style: string }>(() => ({
    class: (attrs.class as string) || '',
    style: (attrs.style as string) || ''
}))

const symbolId = computed(() => {
    const { VITE_ICON_LOCAL_PREFIX: prefix } = import.meta.env

    return `#${prefix}-${props.icon}`
})

/**
 * 计算属性，判断是否应该渲染本地SVG图标
 * 判断逻辑：
 * 1. 如果icon属性中不包含':'字符（即不是Iconify图标格式），则认为是本地图标，返回true
 * 2. 如果icon属性包含':'字符（符合Iconify图标格式，如"mdi:home"），则返回false
 * 3. 如果icon属性为空或undefined，使用可选链操作符避免错误，返回false
 *
 * Iconify图标的格式通常为"集合名:图标名"，例如"mdi:home"、"fa:check"等
 * 本地SVG图标通常直接使用图标名称，不包含':'字符
 */
const isLocalIcon = computed(() => {
    return !props.icon?.includes(':') || false
})
</script>

<style scoped></style>
