<template>
    <Teleport defer :to="`#${GLOBAL_SIDER_MENU_ID}`">
        <div class="h-full flex" @mouseleave="handleResetActiveMenu">
            <FirstLevelMenu
                :menus="allMenus"
                :active-menu-key="activeFirstLevelMenuKey"
                :inverted="inverted"
                :sider-collapse="appStore.siderCollapse"
                :dark-mode="themeStore.isDark"
                @select="handleSelectMixMenu"
                @toggle-sider-collapse="appStore.toggleSidebar">
                <GlobalLogo :show-title="false" :style="{ height: themeStore.header.height + 'px' }" />
            </FirstLevelMenu>
            <div
                class="relative h-full transition-width-300"
                :style="{
                    width: appStore.mixSiderFixed && hasChildMenus ? themeStore.sider.mixChildMenuWidth + 'px' : '0px'
                }">
                <DarkModeContainer
                    class="absolute left-0 top-0 h-full flex flex-col items-stretch overflow-hidden whitespace-nowrap shadow-sm transition-all-300"
                    :inverted="inverted"
                    :style="{
                        width: showDrawer ? themeStore.sider.mixChildMenuWidth + 'px' : '0px'
                    }">
                    <header class="flex justify-between items-center px-12px" :style="{ height: themeStore.header.height + 'px' }">
                        <h2 class="text-16px text-primary font-bold">成绩智能管理系统</h2>
                        <PinToggler
                            :pin="appStore.mixSiderFixed"
                            :class="{ 'text-white:88 !hover:text-white': inverted }"
                            @click="appStore.toggleMixSiderFixed" />
                    </header>
                    <SimpleScrollbar>
                        <ElMenu
                            mode="vertical"
                            :default-active="selectedKey"
                            router
                            unique-opened
                            @select="val => routerPushByKeyWithMetaQuery(val)">
                            <MenuItem v-for="item in childLevelMenus" :key="item.name" :item="item" :base-path="basePath" />
                        </ElMenu>
                    </SimpleScrollbar>
                </DarkModeContainer>
            </div>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import path from 'path-browserify'
import { computed, ref, watch } from 'vue'
import { type RouteRecordRaw, useRoute } from 'vue-router'
import { GLOBAL_SIDER_MENU_ID } from '@/constants/app'
import { useRouterPush } from '@/hooks/router'
import useAppStore from '@/store/modules/app'
import useRouteStore from '@/store/modules/route'
import useThemeStore from '@/store/modules/theme'
import { useMenu, useMixMenuContext } from '../../../context'

defineOptions({
    name: 'VerticalMixMenu'
})

const route = useRoute()
const appStore = useAppStore()
const themeStore = useThemeStore()
const routeStore = useRouteStore()
const { routerPushByKeyWithMetaQuery } = useRouterPush()

const { allMenus, activeFirstLevelMenuKey, setActiveFirstLevelMenuKey, getActiveFirstLevelMenuKey } = useMixMenuContext()
const { selectedKey } = useMenu()

const drawerVisible = ref(false)

const basePath = ref('')

const inverted = computed(() => !themeStore.isDark && themeStore.sider.inverted)

const childLevelMenus = computed<App.Global.Menu[]>(() => {
    return routeStore.menus.find(item => item.name === activeFirstLevelMenuKey.value)?.children || []
})

const hasChildMenus = computed(() => childLevelMenus.value.length > 0)

const showDrawer = computed(() => hasChildMenus.value && (drawerVisible.value || appStore.mixSiderFixed))

function handleSelectMixMenu(menu: RouteRecordRaw) {
    setActiveFirstLevelMenuKey(menu.name as string)
    const currentPath = path.resolve('/', menu.path as string)
    if (menu.children && menu.children.length > 0) {
        drawerVisible.value = true
        basePath.value = currentPath
    } else {
        routerPushByKeyWithMetaQuery(currentPath)
    }
}

function handleResetActiveMenu() {
    drawerVisible.value = false

    if (!appStore.mixSiderFixed) {
        getActiveFirstLevelMenuKey()
    }
}

const expandedKeys = ref<string[]>([])

function updateExpandedKeys() {
    if (appStore.siderCollapse || !selectedKey.value) {
        expandedKeys.value = []
        return
    }
    expandedKeys.value = routeStore.getSelectedMenuKeyPath(selectedKey.value)
}

watch(
    () => route.name,
    () => {
        updateExpandedKeys()
    },
    { immediate: true }
)
</script>

<style scoped></style>
