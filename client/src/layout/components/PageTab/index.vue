<script setup lang="ts">
import type { Component } from 'vue'
import { computed } from 'vue'
import type { PageTabMode, PageTabProps } from '../../types'
import ButtonTab from './button-tab.vue'
import ChromeTab from './chrome-tab.vue'
import SvgClose from './SvgClose.vue'
import SvgPin from './SvgPin.vue'
import { ACTIVE_COLOR, createTabCssVars } from './shared'

defineOptions({
    name: 'PageTab'
})

const props = withDefaults(defineProps<PageTabProps>(), {
    mode: 'chrome',
    commonClass: 'transition-all-300',
    activeColor: ACTIVE_COLOR,
    closable: true
})

type Emits = (e: 'close') => void

const emit = defineEmits<Emits>()

const activeTabComponent = computed(() => {
    const { mode, chromeClass, buttonClass } = props

    const tabComponentMap = {
        chrome: {
            component: ChromeTab,
            class: chromeClass
        },
        button: {
            component: ButtonTab,
            class: buttonClass
        }
    } satisfies Record<PageTabMode, { component: Component; class?: string }>

    return tabComponentMap[mode]
})

const cssVars = computed(() => createTabCssVars())

const bindProps = computed(() => {
    const { chromeClass: chromeCls, buttonClass: btnCls, ...rest } = props

    return rest
})

function handleClose() {
    emit('close')
}
</script>

<template>
    <component :is="activeTabComponent.component" :class="activeTabComponent.class" :style="cssVars" v-bind="bindProps">
        <template #prefix>
            <slot name="prefix"></slot>
        </template>
        <slot></slot>
        <template #suffix>
            <slot name="suffix">
                <SvgClose v-if="closable" class="svg-close" @click.stop="handleClose" />
                <SvgPin v-else />
            </slot>
        </template>
    </component>
</template>

<style lang="scss" scoped>
/* 关闭按钮样式 */
.svg-close:hover {
    font-size: 12px;
}
</style>
