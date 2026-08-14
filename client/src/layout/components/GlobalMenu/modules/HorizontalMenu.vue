<template>
    <Teleport defer :to="`#${GLOBAL_HEADER_MENU_ID}`">
        <ElMenu ellipsis class="w-full" mode="horizontal" :default-active="selectedKey" router unique-opened>
            <MenuItem v-for="item in routeStore.menus" :key="item.name" :item="item" base-path="/" />
        </ElMenu>
    </Teleport>
</template>

<script setup lang="ts">
import { GLOBAL_HEADER_MENU_ID } from '@/constants/app'
import useRouteStore from '@/store/modules/route'
import { useMenu } from '../../../context'

defineOptions({ name: 'HorizontalMenu' })

const routeStore = useRouteStore()
const { selectedKey } = useMenu()
</script>

<style scoped></style>
