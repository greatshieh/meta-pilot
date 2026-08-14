<template>
    <ElBreadcrumb v-if="themeStore.header.breadcrumb.visible">
        <!-- define component start: BreadcrumbContent -->
        <DefineBreadcrumbContent v-slot="{ breadcrumb }"></DefineBreadcrumbContent>

        <!-- define component end: BreadcrumbContent -->
        <ElBreadcrumbItem v-for="item in routeStore.breadcrumbs" :key="item.name">
            <div class="i-flex-y-center align-middle text-14px">
                <!-- <SvgIcon
                    v-if="themeStore.header.breadcrumb.showIcon && breadcrumb.icon"
                    :icon="breadcrumb.icon"
                    class="mr-4px text-icon"></SvgIcon> -->
                {{ item.title }}
            </div>
        </ElBreadcrumbItem>
    </ElBreadcrumb>
</template>

<script setup lang="ts">
import * as pathToRegexp from 'path-to-regexp'
import { ref, watch } from 'vue'
import type { RouteLocationMatched, RouteLocationRaw } from 'vue-router'
import { useRoute, useRouter } from 'vue-router'
import useThemeStore from '@/store/modules/theme'
import useRouteStore from '@/store/modules/route'
import { useRouterPush } from '@/hooks/router'
import { createReusableTemplate } from '@vueuse/core'

defineOptions({ name: 'GlobalBreadcrumb' })

const themeStore = useThemeStore()
const routeStore = useRouteStore()

const { routerPushByKey } = useRouterPush()

interface BreadcrumbContentProps {
    breadcrumb: App.Global.Menu
}

const [DefineBreadcrumbContent, BreadcrumbContent] = createReusableTemplate<BreadcrumbContentProps>()

function handleClickMenu(key: string) {
    routerPushByKey(key)
}
</script>

<style lang="scss" scoped>
.app-breadcrumb.el-breadcrumb {
    display: inline-block;
    font-size: 14px;
    line-height: 60px;
    margin-left: 8px;

    .no-redirect {
        color: var(--ep-text-color-disabled);
        cursor: text;
    }
}
</style>
