<template>
    <!-- 根据 icon 类型决定使用的不同类型的图标组件 -->
    <ElIcon v-if="icon !== '' && icon.startsWith('ElIcon')" class="sub-ElIcon">
        <component :is="icon.replace('ElIcon-', '')" />
    </ElIcon>
    <ElIcon v-else-if="icon !== ''">
        <SvgIcon :icon="icon" />
    </ElIcon>
    <!-- 菜单标题 -->
    <span v-if="title" class="menu-title">{{ title }}</span>
</template>

<script lang="ts" setup>
defineOptions({ name: 'Item' })

defineProps({
    icon: {
        type: String,
        default: ''
    },
    title: {
        type: String,
        default: ''
    },
    isSubMenuItem: {
        type: Boolean,
        default: false
    }
})
</script>

<style lang="scss" scoped>
.sub-el-icon {
    width: 1em;
    height: 1em;
    color: currentcolor;
}

.menu-title {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    position: relative;
}
</style>
