<script setup lang="ts">
import { autoUpdate, computePosition, flip, offset, shift } from '@floating-ui/dom'
import { useElementBounding } from '@vueuse/core'
import { computed, nextTick, ref, useTemplateRef, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type BetterScroll from '@/components/BetterScroll/index.vue'
import useAppStore from '@/store/modules/app'
import useThemeStore from '@/store/modules/theme'
import type ContextMenu from './context-menu.vue'
import useTabStore from '@/store/modules/tab'

defineOptions({ name: 'GlobalTab' })

const $route = useRoute()
const $router = useRouter()
const appStore = useAppStore()
const themeStore = useThemeStore()
const tabStore = useTabStore()

/* 响应式定义 */
// 标签栏对象
const bsWrapper = useTemplateRef<HTMLElement>('bsWrapper')
// 标签栏对象的尺寸
const { width: bsWrapperWidth, left: bsWrapperLeft } = useElementBounding(bsWrapper)
// 滚动条对象
const bsScroll = useTemplateRef<InstanceType<typeof BetterScroll>>('bsScroll')

const tabRef = useTemplateRef<HTMLElement>('tabRef')
// 上下文对象
const contextMenu = useTemplateRef<InstanceType<typeof ContextMenu>>('contextMenu')
// 是否pc
const isPCFlag = computed(() => appStore.device === 'desktop')

const TAB_DATA_ID = 'data-tab-id'

type TabNamedNodeMap = NamedNodeMap & {
    [TAB_DATA_ID]: Attr
}

async function scrollToActiveTab() {
    await nextTick()
    if (!tabRef.value) return

    const { children } = tabRef.value

    for (let i = 0; i < children.length; i += 1) {
        const child = children[i]

        const { value: tabId } = (child.attributes as TabNamedNodeMap)[TAB_DATA_ID]

        if (tabId === tabStore.activeTabName) {
            const { left, width } = child.getBoundingClientRect()
            const clientX = left + width / 2

            setTimeout(() => {
                scrollByClientX(clientX)
            }, 50)

            break
        }
    }
}

function scrollByClientX(clientX: number) {
    const currentX = clientX - bsWrapperLeft.value
    const deltaX = currentX - bsWrapperWidth.value / 2

    if (bsScroll.value?.instance) {
        const { maxScrollX, x: leftX, scrollBy } = bsScroll.value.instance

        const rightX = maxScrollX - leftX
        const update = deltaX > 0 ? Math.max(-deltaX, rightX) : Math.min(-deltaX, -leftX)

        scrollBy(update, 0, 300)
    }
}

function handleCloseTab(tab: App.Global.Tab) {
    tabStore.removeTab(tab.name)
}

interface DropdownConfig {
    visible: boolean
    x: number
    y: number
    tabName: string
}

const dropdown = ref<DropdownConfig>({
    visible: false,
    x: 0,
    y: 0,
    tabName: ''
})

function setDropdown(config: Partial<DropdownConfig>) {
    Object.assign(dropdown.value, config)
}

function init() {
    tabStore.initTabStore($route)
}

// 按指定方向切换标签
function switchTag(dir: string) {
    const visitedViews = tabStore.tabs
    for (let i = 0; i < visitedViews.length; i += 1) {
        const tag = visitedViews[i]
        if (tag.path === $route.path) {
            if (dir === 'next') {
                if (i < visitedViews.length - 1) {
                    $router.push({
                        path: visitedViews[i + 1].path,
                        query: visitedViews[i + 1].query
                    })
                } else {
                    $router.push({
                        path: visitedViews[0].path,
                        query: visitedViews[0].query
                    })
                }
            } else {
                if (i > 0) {
                    $router.push({
                        path: visitedViews[i - 1].path,
                        query: visitedViews[i - 1].query
                    })
                } else {
                    $router.push({
                        path: visitedViews[visitedViews.length - 1].path,
                        query: visitedViews[visitedViews.length - 1].query
                    })
                }
            }
        }
    }
}

// 重新加载选择的标签
async function refreshSelectedTag() {
    appStore.reloadPage(500)
}

function updateContextMenuPos(target: HTMLElement) {
    if (!contextMenu.value) return
    computePosition(target, contextMenu.value.$el, {
        placement: 'right-end',
        middleware: [
            offset({ mainAxis: -40 }),
            flip({ crossAxis: true }), // 自动翻转
            shift({ crossAxis: true }) // 防止超出屏幕
        ]
    }).then(({ x, y }) => {
        setDropdown({ x, y, visible: true })
    })
}

// 清理floating元素的自动更新
let clearup: (() => void) | undefined = () => {}

// 打开上下文菜单
async function handleContextMenu(e: MouseEvent, tag: App.Global.Tab) {
    const target = e.currentTarget as HTMLElement
    if (!target || !contextMenu.value) return

    dropdown.value.tabName = tag.name

    clearup = autoUpdate(target, contextMenu.value.$el, () => updateContextMenuPos(target))
}

function handleDropdownVisible(visible: boolean | undefined) {
    closeContextMenu()
}

// 关闭上下文菜单
function closeContextMenu() {
    dropdown.value = {
        visible: false,
        x: 0,
        y: 0,
        tabName: ''
    }
    if (clearup) clearup()
}

function removeFocus() {
    ;(document.activeElement as HTMLElement)?.blur()
}

// 回到顶部
function scrollToTop() {
    const target = document.getElementById('__SCROLL_EL_ID__')

    target?.scrollTo(0, 0)
}

watch(
    () => $route.fullPath,
    () => {
        tabStore.addTab($route)
    }
)

watch(
    () => tabStore.activeTabName,
    () => {
        scrollToActiveTab()
    }
)

init()
</script>

<template>
    <DarkModeContainer class="size-full flex items-center px-16px shadow-tab">
        <div ref="bsWrapper" class="h-full flex-1 overflow-hidden">
            <BetterScroll ref="bsScroll" :options="{ scrollX: true, scrollY: false, click: !isPCFlag }" @click="removeFocus">
                <div
                    ref="tabRef"
                    class="h-full flex pr-18px"
                    :class="[themeStore.tab.mode === 'chrome' ? 'items-end' : 'items-center gap-12px']">
                    <PageTab
                        v-for="tab in tabStore.tabs"
                        :key="tab.name"
                        :[TAB_DATA_ID]="tab.name"
                        :mode="themeStore.tab.mode"
                        :dark-mode="themeStore.isDark"
                        :active="tab.name === tabStore.activeTabName"
                        :closable="!tabStore.isTabRetain(tab.name)"
                        @click="tabStore.switchRouteByTab(tab)"
                        @close="handleCloseTab(tab)"
                        @contextmenu.prevent.stop="handleContextMenu($event, tab)">
                        <div class="max-w-240px ellipsis-text">{{ tab.title }}</div>
                    </PageTab>
                </div>
            </BetterScroll>
        </div>

        <div>
            <ElTooltip content="上一个页面" placement="bottom">
                <ElButton :disabled="tabStore.tabs.length === 1" size="small" circle @click="switchTag('pre')">
                    <SvgIcon icon="ant-design:double-left-outlined" />
                </ElButton>
            </ElTooltip>
            <ElTooltip content="下一个页面" placement="bottom">
                <ElButton :disabled="tabStore.tabs.length === 1" size="small" circle @click="switchTag('next')">
                    <SvgIcon icon="ant-design:double-right-outlined" />
                </ElButton>
            </ElTooltip>
            <ElTooltip content="重新加载" placement="bottom">
                <ElButton size="small" circle @click="refreshSelectedTag">
                    <SvgIcon icon="ant-design:reload-outlined" :class="{ 'animate-spin animate-duration-750': !appStore.reloadFlag }" />
                </ElButton>
            </ElTooltip>
            <ElTooltip content="返回顶部" placement="bottom">
                <ElButton size="small" circle @click="scrollToTop">
                    <SvgIcon icon="ant-design:to-top-outlined" />
                </ElButton>
            </ElTooltip>
        </div>
    </DarkModeContainer>
    <ContextMenu
        ref="contextMenu"
        v-model="dropdown.visible"
        :left="dropdown.x"
        :top="dropdown.y"
        :tab-name="dropdown.tabName"
        :is-affix="tabStore.isTabRetain(dropdown.tabName)"
        @update:v-model="handleDropdownVisible" />
</template>

<style scoped></style>
