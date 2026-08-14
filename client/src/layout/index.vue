<script setup lang="ts" name="Layout">
import { computed } from 'vue'
import useAppStore from '@/store/modules/app'
import useThemeStore from '@/store/modules/theme'
import { LAYOUT_SCROLL_EL_ID } from './components/AdminLayout/shared'
import { setupMixMenuContext } from './context'
import type { LayoutMode } from './types'

const appStore = useAppStore()
const themeStore = useThemeStore()

const { childLevelMenus, isActiveFirstLevelMenuHasChildren } = setupMixMenuContext()

const layoutMode = computed(() => {
    const vertical: LayoutMode = 'vertical'
    const horizontal: LayoutMode = 'horizontal'
    return themeStore.layout.mode.includes(vertical) ? vertical : horizontal
})

const siderVisible = computed(() => themeStore.layout.mode !== 'horizontal')

const isVerticalMix = computed(() => themeStore.layout.mode === 'vertical-mix')

const isHorizontalMix = computed(() => themeStore.layout.mode === 'horizontal-mix')

const siderWidth = computed(() => getSiderWidth())

const siderCollapsedWidth = computed(() => getSiderCollapsedWidth())

const headerProps = computed(() => {
    const { mode } = themeStore.layout

    const headerPropsConfig: Record<UnionKey.ThemeLayoutMode, App.Global.HeaderProps> = {
        vertical: {
            showLogo: false,
            showMenu: false,
            showMenuToggler: true
        },
        'vertical-mix': {
            showLogo: false,
            showMenu: false,
            showMenuToggler: false
        },
        horizontal: {
            showLogo: true,
            showMenu: true,
            showMenuToggler: false
        },
        'horizontal-mix': {
            showLogo: true,
            showMenu: true,
            showMenuToggler: true
        }
    }

    return headerPropsConfig[mode]
})

function getSiderWidth() {
    const { reverseHorizontalMix } = themeStore.layout
    const { width, mixWidth, mixChildMenuWidth } = themeStore.sider

    if (isHorizontalMix.value && reverseHorizontalMix) {
        return isActiveFirstLevelMenuHasChildren.value ? width : 0
    }

    let w = isVerticalMix.value || isHorizontalMix.value ? mixWidth : width

    if (isVerticalMix.value && appStore.mixSiderFixed && childLevelMenus.value.length) {
        w += mixChildMenuWidth
    }

    return w
}

function getSiderCollapsedWidth() {
    const { reverseHorizontalMix } = themeStore.layout
    const { collapsedWidth, mixCollapsedWidth, mixChildMenuWidth } = themeStore.sider

    if (isHorizontalMix.value && reverseHorizontalMix) {
        return isActiveFirstLevelMenuHasChildren.value ? collapsedWidth : 0
    }

    let w = isVerticalMix.value || isHorizontalMix.value ? mixCollapsedWidth : collapsedWidth

    if (isVerticalMix.value && appStore.mixSiderFixed && childLevelMenus.value.length) {
        w += mixChildMenuWidth
    }

    return w
}
</script>

<template>
    <AdminLayout
        v-model:sider-collapse="appStore.siderCollapse"
        :mode="layoutMode"
        :scroll-el-id="LAYOUT_SCROLL_EL_ID"
        :scroll-mode="themeStore.layout.scrollMode"
        :is-mobile="appStore.device !== 'desktop'"
        :full-content="appStore.fullContent"
        :fixed-top="themeStore.fixedHeaderAndTab"
        :header-height="themeStore.header.height"
        :tab-visible="themeStore.tab.visible"
        :tab-height="themeStore.tab.height"
        :content-class="appStore.contentXScrollable ? 'overflow-x-hidden' : ''"
        :sider-visible="siderVisible"
        :sider-width="siderWidth"
        :sider-collapsed-width="siderCollapsedWidth"
        :footer-visible="themeStore.footer.visible"
        :footer-height="themeStore.footer.height"
        :fixed-footer="themeStore.footer.fixed"
        :right-footer="themeStore.footer.right">
        <template #header>
            <GlobalHeader v-bind="headerProps" />
        </template>
        <template #tab>
            <GlobalTab />
        </template>
        <template #sider>
            <GlobalSider />
        </template>
        <GlobalMenu />
        <GlobalContent />
        <ThemeDrawer />
        <template #footer>
            <GlobalFooter />
        </template>
    </AdminLayout>
</template>

<style lang="scss" scoped></style>
