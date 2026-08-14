<template>
    <!-- define component: MixMenuItem -->
    <DefineMixMenuItem v-slot="{ label, icon, active, isMini }">
        <div
            class="mx-4px mb-6px flex flex-col items-center cursor-pointer rounded-8px bg-transparent px-4px py-8px transition-300 hover:bg-[rgb(0,0,0,0.08)]"
            :class="{
                'text-primary selected-mix-menu': active,
                'text-white:65 hover:text-white': inverted,
                '!text-white !bg-primary': active && inverted
            }">
            <SvgIcon :icon="icon" :class="[isMini ? 'text-icon-small' : 'text-icon-large']" />
            <p class="w-full ellipsis-text text-center text-12px transition-height-300" :class="[isMini ? 'h-0 pt-0' : 'h-20px pt-4px']">
                {{ label }}
            </p>
        </div>
    </DefineMixMenuItem>
    <!-- define component end: MixMenuItem -->

    <div class="h-full flex flex-col items-stretch flex-1 overflow-hidden">
        <slot></slot>
        <SimpleScrollbar>
            <MixMenuItem
                v-for="menu in topLevelMenus"
                :key="menu.name"
                :label="menu.meta?.title || ''"
                :icon="menu.meta?.icon || ''"
                :active="menu.name === activeMenuKey"
                :is-mini="siderCollapse"
                @click="handleClickMixMenu(menu)" />
        </SimpleScrollbar>
        <Hamburger
            arrow-icon
            :collapsed="siderCollapse"
            :z-index="99"
            :class="{ 'text-white:88 !hover:text-white': inverted }"
            @click="toggleSiderCollapse" />
    </div>
</template>

<script setup lang="ts">
import { createReusableTemplate } from '@vueuse/core'
import { computed } from 'vue'
import type { RouteRecordRaw } from 'vue-router'

defineOptions({ name: 'FirstLevelMenu' })

interface Props {
    menus: RouteRecordRaw[]
    activeMenuKey?: string
    inverted?: boolean
    siderCollapse?: boolean
    darkMode?: boolean
}

const props = defineProps<Props>()

interface Emits {
    (e: 'select', menu: RouteRecordRaw): boolean
    (e: 'toggleSiderCollapse'): void
}

const emit = defineEmits<Emits>()

interface MixMenuItemProps {
    /** Menu item label */
    label: string
    /** Menu item icon */
    icon: string
    /** Active menu item */
    active: boolean
    /** Mini size */
    isMini?: boolean
}
const [DefineMixMenuItem, MixMenuItem] = createReusableTemplate<MixMenuItemProps>()

function handleClickMixMenu(menu: RouteRecordRaw) {
    emit('select', menu)
}

function toggleSiderCollapse() {
    emit('toggleSiderCollapse')
}

const topLevelMenus = computed<RouteRecordRaw[]>(() => {
    return props.menus.map(item => {
        if (item.children && item.children.length > 0) {
            return {
                ...item,
                children: item.children.filter(child => !child.meta?.hidden)
            }
        } else {
            return item
        }
    })
})
</script>

<style scoped>
.selected-mix-menu {
    background-color: var(--el-color-primary-light-7);
}
</style>
