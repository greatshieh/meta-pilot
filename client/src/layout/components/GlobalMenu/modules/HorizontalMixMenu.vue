<template>
    <Teleport defer :to="`#${GLOBAL_HEADER_MENU_ID}`">
        <ElMenu
            ellipsis
            class="w-full"
            mode="horizontal"
            :default-active="selectedKey"
            router
            unique-opened
            @select="val => routerPushByKeyWithMetaQuery(val)">
            <MenuItem v-for="item in childLevelMenus" :key="item.name" :item="item" :base-path="basePath" />
        </ElMenu>
    </Teleport>
    <Teleport defer :to="`#${GLOBAL_SIDER_MENU_ID}`">
        <FirstLevelMenu
            :menus="allMenus"
            :active-menu-key="activeFirstLevelMenuKey"
            :sider-collapse="appStore.siderCollapse"
            :dark-mode="themeStore.isDark"
            @select="handleSelectMixMenu"
            @toggle-sider-collapse="appStore.toggleSidebar" />
    </Teleport>
</template>

<script setup lang="ts">
import path from 'path-browserify'
import { computed, ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { GLOBAL_HEADER_MENU_ID, GLOBAL_SIDER_MENU_ID } from '@/constants/app'
import { useRouterPush } from '@/hooks/router'
import useAppStore from '@/store/modules/app'
import useRouteStore from '@/store/modules/route'
import useThemeStore from '@/store/modules/theme'
import { useMenu, useMixMenuContext } from '../../../context'

defineOptions({
    name: 'HorizontalMixMenu'
})

const appStore = useAppStore()
const themeStore = useThemeStore()
const routeStore = useRouteStore()
const { routerPushByKeyWithMetaQuery } = useRouterPush()
const { allMenus, activeFirstLevelMenuKey, setActiveFirstLevelMenuKey } = useMixMenuContext()
const { selectedKey } = useMenu()

const basePath = ref('')

const childLevelMenus = computed<App.Global.Menu[]>(() => {
    return routeStore.menus.find(item => item.name === activeFirstLevelMenuKey.value)?.children || []
})

function handleSelectMixMenu(menu: RouteRecordRaw) {
    setActiveFirstLevelMenuKey(menu.name as string)
    const currentPath = path.resolve('/', menu.path as string)
    if (!menu.children || menu.children?.length === 0) {
        routerPushByKeyWithMetaQuery(currentPath)
    } else {
        basePath.value = currentPath
    }
}
</script>

<style scoped></style>
