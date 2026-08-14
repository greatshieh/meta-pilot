<template>
    <div>
        <router-view v-slot="{ Component }">
            <keep-alive :include="routeStore.cacheRoutes">
                <component :is="Component" />
            </keep-alive>
        </router-view>
    </div>
</template>

<script setup lang="ts">
import useRouteStore from '@/store/modules/route'

const routeStore = useRouteStore()

defineOptions({
    name: 'SuperAdmin'
})
</script>
