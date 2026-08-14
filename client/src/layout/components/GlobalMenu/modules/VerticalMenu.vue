<template>
    <Teleport defer :to="`#${GLOBAL_SIDER_MENU_ID}`">
        <SimpleScrollbar>
            <ElMenu mode="vertical" :default-active="selectedKey" :collapse="isSidebarCollapsed" router unique-opened>
                <MenuItem v-for="item in routeStore.menus" :key="item.routePath" :item="item" base-path="/" />
            </ElMenu>
        </SimpleScrollbar>
    </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { GLOBAL_SIDER_MENU_ID } from '@/constants/app'
import useAppStore from '@/store/modules/app'
import useRouteStore from '@/store/modules/route/index'
import { useMenu } from '../../../context'

defineOptions({ name: 'VerticalMenu' })

const appStore = useAppStore()
const routeStore = useRouteStore()

const isSidebarCollapsed = computed(() => appStore.siderCollapse)

const { selectedKey } = useMenu()
</script>

<style scoped></style>
