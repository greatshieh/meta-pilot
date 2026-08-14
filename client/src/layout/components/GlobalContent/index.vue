<template>
    <router-view v-slot="{ Component, route }">
        <transition
            :name="transitionName"
            mode="out-in"
            @before-leave="appStore.setContentXScrollable(true)"
            @after-leave="resetScroll"
            @after-enter="appStore.setContentXScrollable(false)">
            <keep-alive :include="routeStore.cacheRoutes" :exclude="routeStore.excludeCacheRoutes">
                <component
                    v-if="appStore.reloadFlag"
                    :is="Component"
                    :key="route.name"
                    class="flex-grow bg-[--el-fill-color] transition-300 p-16px" />
            </keep-alive>
        </transition>
    </router-view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import useAppStore from '@/store/modules/app'
import useThemeStore from '@/store/modules/theme'
import useTabStore from '@/store/modules/tab'
import useRouteStore from '@/store/modules/route'
import { LAYOUT_SCROLL_EL_ID } from '../AdminLayout/shared'

defineOptions({ name: 'GlobalContent' })

const appStore = useAppStore()
const themeStore = useThemeStore()
const routeStore = useRouteStore()

const tabStore = useTabStore()

const transitionName = computed(() => (themeStore.page.animate ? themeStore.page.animateMode : 'fade-scale'))

function resetScroll() {
    const el = document.querySelector(`#${LAYOUT_SCROLL_EL_ID}`)

    el?.scrollTo({ left: 0, top: 0 })
}
</script>

<style lang="scss" scoped></style>
