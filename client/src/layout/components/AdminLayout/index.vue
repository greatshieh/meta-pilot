<template>
    <div :class="layoutKls" :style="cssVars">
        <div :id="isWrapperScroll ? scrollElId : undefined" :class="wrapperKls">
            <!-- Header -->
            <template v-if="showHeader">
                <header v-show="!fullContent" :class="[...headerKls]">
                    <slot name="header"></slot>
                </header>
                <div v-show="!fullContent && fixedHeaderAndTab" :class="headerPlacementKls"></div>
            </template>

            <!-- Tab -->
            <template v-if="showTab">
                <div :class="tabKls">
                    <slot name="tab"></slot>
                </div>
                <div v-show="fullContent || fixedHeaderAndTab" :class="tabPlacementKls"></div>
            </template>

            <!-- Sider -->
            <template v-if="showSider">
                <aside v-show="!fullContent" :class="siderKls">
                    <slot name="sider"></slot>
                </aside>
            </template>

            <!-- Mobile Sider -->
            <template v-if="showMobileSider">
                <aside :class="mobileSiderKls">
                    <slot name="sider"></slot>
                </aside>
                <div v-show="!siderCollapse" :class="mobileSiderMaskKls" @click="handleClickMask"></div>
            </template>

            <!-- Main Content -->
            <main :id="isContentScroll ? scrollElId : undefined" :class="contentKls">
                <slot></slot>
            </main>

            <!-- Footer -->
            <template v-if="showFooter">
                <footer v-show="!fullContent" :class="footerKls">
                    <slot name="footer"></slot>
                </footer>
                <div v-show="!fullContent && fixedFooter" :class="footerPlacementKls"></div>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { AdminLayoutProps } from '../../types'
import style from './index.module.css'
import { createLayoutCssVars, LAYOUT_MAX_Z_INDEX, LAYOUT_SCROLL_EL_ID } from './shared'
import { useNamespace } from '@/hooks/use-namespace'

defineOptions({
    name: 'AdminLayout'
})

const props = withDefaults(defineProps<AdminLayoutProps>(), {
    mode: 'vertical',
    scrollMode: 'content',
    scrollElId: LAYOUT_SCROLL_EL_ID,
    commonClass: 'transition-all-300',
    fixedTop: true,
    maxZIndex: LAYOUT_MAX_Z_INDEX,
    headerVisible: true,
    headerHeight: 56,
    tabVisible: true,
    tabHeight: 48,
    siderVisible: true,
    siderCollapse: false,
    siderWidth: 220,
    siderCollapsedWidth: 64,
    footerVisible: true,
    footerHeight: 48,
    rightFooter: false
})

interface Emits {
    /** Update siderCollapse */
    (e: 'update:siderCollapse', collapse: boolean): void
}

const emit = defineEmits<Emits>()

const ns = useNamespace('layout', ref(''))

type SlotFn = (props?: Record<string, unknown>) => unknown

type Slots = {
    /** Main */
    default?: SlotFn
    /** Header */
    header?: SlotFn
    /** Tab */
    tab?: SlotFn
    /** Sider */
    sider?: SlotFn
    /** Footer */
    footer?: SlotFn
}

const slots = defineSlots<Slots>()

const cssVars = computed(() => createLayoutCssVars(props))

// config visible
const showHeader = computed(() => Boolean(slots.header) && props.headerVisible)
const showTab = computed(() => Boolean(slots.tab) && props.tabVisible)
const showSider = computed(() => !props.isMobile && Boolean(slots.sider) && props.siderVisible)
const showMobileSider = computed(() => props.isMobile && Boolean(slots.sider) && props.siderVisible)
const showFooter = computed(() => Boolean(slots.footer) && props.footerVisible)

// scroll mode
const isWrapperScroll = computed(() => props.scrollMode === 'wrapper')
const isContentScroll = computed(() => props.scrollMode === 'content')

// layout direction
const isVertical = computed(() => props.mode === 'vertical')
const isHorizontal = computed(() => props.mode === 'horizontal')

const fixedHeaderAndTab = computed(() => props.fixedTop || (isHorizontal.value && isWrapperScroll.value))

// css
const leftGapClass = computed(() => {
    if (!props.fullContent && showSider.value) {
        return props.siderCollapse ? 'pl-[var(--layout-sider-collapsed-width)]' : 'pl-[var(--layout-sider-width)]'
    }

    return ''
})

const headerLeftGapClass = computed(() => (isVertical.value ? leftGapClass.value : ''))

const footerLeftGapClass = computed(() => {
    const condition1 = isVertical.value
    const condition2 = isHorizontal.value && isWrapperScroll.value && !props.fixedFooter
    const condition3 = Boolean(isHorizontal.value && props.rightFooter)

    if (condition1 || condition2 || condition3) {
        return leftGapClass.value
    }

    return ''
})

const layoutKls = computed(() => {
    return [ns.b(''), props.commonClass, ns.is('fixed', props.fixedTop || (isHorizontal.value && isWrapperScroll.value))]
})

const wrapperKls = computed(() => {
    return [ns.e('wrapper'), props.commonClass, props.scrollWrapperClass, ns.is('scroll', isWrapperScroll.value)]
})

const headerKls = computed(() => {
    return [ns.e('header'), props.commonClass, props.headerClass, headerLeftGapClass.value]
})

const headerPlacementKls = computed(() => {
    return [ns.e('header-placement')]
})

const tabKls = computed(() => {
    return [ns.e('tab'), props.commonClass, props.tabClass, leftGapClass.value, ns.is('show', !props.fullContent && showHeader.value)]
})

const tabPlacementKls = computed(() => {
    return [ns.e('tab-placement')]
})

const siderKls = computed(() => {
    const kls = [ns.e('sider'), props.commonClass, props.siderClass, ns.is('collapsed', props.siderCollapse)]

    if (showHeader.value && !headerLeftGapClass.value) {
        kls.push(ns.em('sider', 'top'))
    }
    if (showFooter.value && !footerLeftGapClass.value) {
        kls.push(ns.em('sider', 'bottom'))
    }

    return kls
})

const mobileSiderKls = computed(() => {
    return [
        ns.em('sider', 'mobile'),
        props.commonClass,
        props.mobileSiderClass,
        props.siderCollapse ? 'overflow-hidden' : 'w-[var(--layout-sider-width)]'
    ]
})

const mobileSiderMaskKls = computed(() => {
    return [ns.em('sider-mask', 'mobile')]
})

const contentKls = computed(() => {
    return [ns.e('content'), props.commonClass, props.contentClass, leftGapClass.value, ns.is('scroll', isContentScroll.value)]
})

const footerKls = computed(() => {
    return [ns.e('footer'), props.commonClass, props.footerClass, footerLeftGapClass.value, ns.is('fixed', props.fixedFooter)]
})

const footerPlacementKls = computed(() => {
    return [ns.e('footer-placement')]
})

function handleClickMask() {
    emit('update:siderCollapse', true)
}
</script>

<style lang="scss">
@use '@/styles/scrollbar.scss' as *;

#_SCROLL_EL_ID__ {
    @include scrollbar();
}
</style>

<style lang="scss" scoped>
.layout {
    position: relative;
    height: 100%;

    .layout__wrapper {
        display: flex;
        flex-direction: column;
        height: 100%;

        &.is-scroll {
            overflow-y: auto;
        }
    }

    .layout__header {
        z-index: var(--layout-header-z-index);
        height: var(--layout-header-height);
        flex-shrink: 0;
    }

    .layout__header-placement {
        height: var(--layout-header-height);
        flex-shrink: 0;
        overflow: hidden;
    }

    .layout__tab {
        top: 0;
        height: var(--layout-tab-height);
        z-index: var(--layout-tab-z-index);
        flex-shrink: 0;

        &.is-show {
            top: var(--layout-header-height);
        }
    }

    .layout__tab-placement {
        height: var(--layout-tab-height);
        flex-shrink: 0;
        overflow: hidden;
    }

    .layout__sider {
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        background-color: var(--el-bg-color-page);
        width: var(--layout-sider-width);
        z-index: var(--layout-sider-z-index);
        &--top {
            padding-top: var(--layout-header-height);
        }
        &--bottom {
            padding-bottom: var(--layout-footer-height);
        }
        &.is-collapsed {
            width: var(--layout-sider-collapsed-width);
        }
    }

    .layout__sider--mobile {
        position: absolute;
        top: 0;
        left: 0;
        width: 0;
        height: 100%;
        background-color: #fff;
        z-index: var(--layout-sider-z-index);
    }

    .layout__sider-mask--mobile {
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        width: 100%;
        background-color: rgba(0, 0, 0, 0.2);
        z-index: var(--layout-mobile-sider-z-index);
    }

    .layout__content {
        display: flex;
        flex-direction: column;
        flex-grow: 1;
        background-color: var(--el-bg-color-page);

        &.is-scroll {
            overflow-y: auto;
        }
    }

    .layout__footer {
        height: var(--layout-footer-height);
        z-index: var(--layout-footer-z-index);
        flex-shrink: 0;

        &.is-fixed {
            position: absolute;
            left: 0;
            bottom: 0;
            width: 100%;
        }
    }

    .layout__footer-placement {
        flex-shrink: 0;
        overflow: hidden;
        height: var(--layout-footer-height);
    }

    &.is-fixed {
        .layout__header {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
        }

        .layout__tab {
            position: absolute;
            left: 0;
            width: 100%;
        }
    }
}
</style>
