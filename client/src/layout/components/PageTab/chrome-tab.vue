<script setup lang="ts">
import { computed, ref } from 'vue'
import type { PageTabProps } from '../../types'
import ChromeTabBg from './chrome-tab-bg.vue'
import { useNamespace } from '@/hooks/use-namespace'

defineOptions({
    name: 'ChromeTab'
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

const ns = useNamespace('tab', ref('chrome'))

const tabKls = computed(() => [ns.b(), ns.is('active', props.active)])
</script>

<template>
    <div
        class="relative inline-flex cursor-pointer items-center justify-center gap-16px whitespace-nowrap px-24px py-6px -mr-18px"
        :class="tabKls">
        <div class="pointer-events-none absolute left-0 top-0 h-full w-full -z-1" :class="[ns.e('bg')]">
            <ChromeTabBg />
        </div>
        <slot name="prefix"></slot>
        <slot></slot>
        <slot name="suffix"></slot>
        <div class="absolute right-7px h-16px w-1px bg-#1f2225" :class="[ns.e('divider')]"></div>
    </div>
</template>

<style lang="scss" scoped>
.chrome-tab {
    &:hover {
        z-index: 9;

        .chrome-tab__bg {
            color: var(--tab-hover-bg-color);
        }

        .chrome-tab__divider {
            opacity: 0;
        }
    }

    &.is-active {
        z-index: 10;
        color: var(--tab-active-text-color);

        .chrome-tab__bg {
            color: var(--tab-active-bg-color);
        }

        .chrome-tab__divider {
            opacity: 0;
        }

        .svg-close:hover {
            color: #ffffff;
            background-color: var(--tab-active-text-color);
        }

        // &:hover .chrome-tab__bg {
        //     color: var(--tab-hover-bg-color);
        // }
    }

    .chrome-tab__bg {
        color: transparent;
    }
}
</style>
