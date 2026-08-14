<script setup lang="ts">
import { computed, ref } from 'vue'
import type { PageTabProps } from '../../types'
import { useNamespace } from '@/hooks/use-namespace'

defineOptions({
    name: 'ButtonTab'
})

const props = defineProps<PageTabProps>()

type SlotFn = (props?: Record<string, unknown>) => unknown

type Slots = {
    /**
     * Slot
     *
     * The center content of the tab
     */
    default?: SlotFn
    /**
     * Slot
     *
     * The left content of the tab
     */
    prefix?: SlotFn
    /**
     * Slot
     *
     * The right content of the tab
     */
    suffix?: SlotFn
}

defineSlots<Slots>()

const ns = useNamespace('tab', ref('button'))

const tabKls = computed(() => [ns.b(), ns.is('active', props.active)])
</script>

<template>
    <div
        class="relative inline-flex cursor-pointer items-center justify-center gap-12px whitespace-nowrap border-(1px solid) rounded-4px px-12px py-4px"
        :class="[tabKls]">
        <slot name="prefix"></slot>
        <slot></slot>
        <slot name="suffix"></slot>
    </div>
</template>

<style lang="scss" scoped>
.button-tab {
    border-color: var(--el-border-color);

    &.is-active {
        color: var(--tab-active-text-color);
        border-color: var(--tab-border-color);
        background-color: var(--tab-active-bg-color);

        .svg-close:hover {
            color: #ffffff;
            background-color: var(--tab-active-text-color);
        }
    }

    &:hover {
        color: var(--tab-active-text-color);
        border-color: var(--tab-border-color);
    }
}
</style>
