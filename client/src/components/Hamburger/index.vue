<!-- 汉堡按钮组件：展开/收缩菜单  -->
<template>
    <ElTooltip placement="bottom" :key="String(collapsed)" :content="collapsed ? '展开菜单' : '折叠菜单'" :z-index="zIndex">
        <ElButton text quaternary class="h-[36px] text-icon" v-bind="$attrs">
            <div class="flex-center gap-8px text-lg">
                <slot>
                    <SvgIcon :icon="icon" />
                </slot>
            </div>
        </ElButton>
    </ElTooltip>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({ name: 'Hamburger' })

interface Props {
    /** Show collapsed icon */
    collapsed?: boolean
    /** Arrow style icon */
    arrowIcon?: boolean
    zIndex?: number
}

const props = withDefaults(defineProps<Props>(), {
    arrowIcon: false,
    zIndex: 98
})

type NumberBool = 0 | 1

const icon = computed(() => {
    const icons: Record<NumberBool, Record<NumberBool, string>> = {
        0: {
            0: 'line-md:menu-fold-left',
            1: 'line-md:menu-fold-right'
        },
        1: {
            0: 'ph-caret-double-left-bold',
            1: 'ph-caret-double-right-bold'
        }
    }

    const arrowIcon = Number(props.arrowIcon || false) as NumberBool

    const collapsed = Number(props.collapsed || false) as NumberBool

    return icons[arrowIcon][collapsed]
})
</script>

<style scoped lang="scss">
.hamburger {
    display: inline-block;
    vertical-align: middle;
    width: 20px;
    height: 20px;
    color: var(--el-text-color-primary);
}

.hamburger.is-active {
    transform: rotate(180deg);
}
</style>
