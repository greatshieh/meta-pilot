<script setup lang="ts">
import { authorityApi, menuApi } from '@/api'
import { TreeInstance, TreeNode, TreeNodeData } from 'element-plus'
import { ref, shallowRef, watch } from 'vue'

defineOptions({ name: 'SetMenu' })

const { authorityId } = defineProps({
    maxHeight: {
        type: String,
        default: '400px'
    },
    authorityId: {
        type: Number,
        default: 0
    }
})

const filterTitle = ref('')
const filterName = ref('')

const authorityMenu = ref<Api.Menu.MenuMeta[]>()
const menuTreeRef = shallowRef<TreeInstance>()
const selectedMenuIds = ref<number[]>([])

const menuTreeProps = {
    label: (data: TreeNodeData, _node: TreeNode) => data.meta.title,
    children: 'children'
}

async function init() {
    // 获取所有菜单列表
    authorityMenu.value = await menuApi.getMenuList()
    // 获取当前角色已授权的菜单列表
    const res = await authorityApi.getAuthorityMenu(authorityId)
    selectedMenuIds.value = []
    res.forEach(item => {
        // 防止直接选中父级造成全选
        if (!res.some(some => some.parentId === item.id)) {
            selectedMenuIds.value.push(item.id)
        }
    })
}

init()

/**
 * 处理设置权限操作
 */
async function setMenuAuthority() {
    const checkedMenus = (menuTreeRef.value?.getCheckedNodes(false, true) || []) as Api.Menu.MenuMeta[]

    if (checkedMenus.length === 0) {
        ElMessage.error('请选择菜单')
        return
    }

    await authorityApi.setMenuAuthority(authorityId, checkedMenus)
}

function filterNode(value: [string, string], data: TreeNodeData) {
    if (value[0] === '' && value[1] === '') return true

    if (value[0] !== '' && data.meta.title.indexOf(value[0]) !== -1) return true
    if (value[1] !== '' && data.name.indexOf(value[1]) !== -1) return true

    return false
}

watch([filterTitle, filterName], val => {
    menuTreeRef.value?.filter(val)
})
</script>

<template>
    <div class="flex items-center justify-between gap-4 mb-2">
        <div class="flex items-center gap-2 flex-1">
            <ElInput v-model="filterTitle" size="small" placeholder="展示名称"></ElInput>
            <ElInput v-model="filterName" size="small" placeholder="路由name"></ElInput>
        </div>
        <ElButton type="primary" size="small" @click="setMenuAuthority">确定</ElButton>
    </div>
    <ElScrollbar :max-height="maxHeight">
        <ElTree
            ref="menuTreeRef"
            :data="authorityMenu"
            :props="menuTreeProps"
            node-key="id"
            default-expand-all
            :default-checked-keys="selectedMenuIds"
            :filter-node-method="filterNode"
            show-checkbox></ElTree>
    </ElScrollbar>
</template>
