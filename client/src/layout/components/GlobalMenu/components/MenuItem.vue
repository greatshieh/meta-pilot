<template>
    <!-- 显示含多个子节点的父菜单，或始终显示的单子节点 -->
    <ElSubMenu v-if="hasChildren" ref="subMenu" :index="resolvePath(item.routePath)" teleported>
        <template #title>
            <Item :icon="item.icon" :title="item.title" isSubMenuItem />
        </template>
        <MenuItem
            v-for="child in item.children"
            :key="resolvePath(item.routePath) + '/' + child.routePath"
            :item="child"
            :base-path="resolvePath(item.routePath)" />
    </ElSubMenu>

    <!-- 显示叶子节点或唯一子节点且父节点未配置始终显示 -->

    <ElMenuItem v-else :index="resolvePath(item.routePath)">
        <Item :icon="item.icon" :title="item.title" />
    </ElMenuItem>
</template>

<script setup lang="ts" name="MenuItem">
import path from 'path-browserify'
import { type PropType, ref } from 'vue'
import { isExternal } from '@/utils/validate'

defineOptions({ name: 'MenuItem' })

const { item, basePath } = defineProps({
    // route object
    item: {
        type: Object as PropType<App.Global.Menu>,
        required: true
    },
    basePath: {
        type: String,
        default: ''
    }
})

const hasChildren = item.children && item.children.length > 0

/**
 * 获取完整路径，适配外部链接
 *
 * @param routePath 路由路径
 * @returns 绝对路径
 */
function resolvePath(routePath: string) {
    if (isExternal(routePath)) return routePath
    // 拼接父路径和当前路径
    return path.resolve(basePath, routePath)
}
</script>
