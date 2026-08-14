<script setup lang="ts">
import { CSSProperties, onMounted, ref } from 'vue'
import { invalidateCache, Method } from 'alova'
import { menuApi } from '@/api'

defineOptions({
    name: 'MenuManagement'
})

const menuDefault: Api.Menu.MenuMeta = {
    id: 0,
    name: '',
    path: '',
    redirect: '',
    component: '',
    parentId: 0,
    children: [],
    meta: {}
}

const tableData = ref<Api.Menu.MenuMeta[]>()

const dialogVisible = ref(false)
const dialogTitle = ref('')
const choosedMenu = ref<Api.Menu.MenuMeta>(menuDefault)

function highLightChildren(data: { row: Api.Menu.MenuMeta; rowIndex: number }): CSSProperties {
    if (!data.row.parentId) {
        return {}
    }

    return { backgroundColor: 'var(--el-color-primary-light-9)' }
}

function deleteMenu(menu: Api.Menu.MenuMeta) {
    ElMessageBox.confirm('确认删除该菜单吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    })
        .then(() => {
            menuApi.deleteMenu(menu.id).then(() => {
                invalidateCache(menuApi.getMenuList() as Method)
                menuApi.getMenuList().then(res => {
                    tableData.value = res
                })
            })
        })
        .catch(() => {
            ElMessage.warning('取消删除菜单')
        })
}

function editMenu(menu: Api.Menu.MenuMeta) {
    dialogVisible.value = true
    dialogTitle.value = `编辑菜单-${menu.meta.title}`
    choosedMenu.value = menu
}

function addMenu() {
    dialogVisible.value = true
    dialogTitle.value = '新增菜单'
}

function addSubMenu(menu: Api.Menu.MenuMeta) {
    dialogVisible.value = true
    dialogTitle.value = `新增 ${menu.meta.title} 子菜单`
    choosedMenu.value = menuDefault
    choosedMenu.value.parentId = menu.id
    choosedMenu.value.component = menu.component.split('/').slice(0, 2).join('/')
    choosedMenu.value.sort = menu.children?.length + 1
}

onMounted(() => {
    invalidateCache(menuApi.getMenuList() as Method)
    menuApi.getMenuList().then(res => {
        tableData.value = res
    })
})
</script>

<template>
    <div class="h-full mx-4">
        <div class="flex items-center gap-8 my-4">
            <div class="text-2xl font-bold">菜单管理</div>
            <ElButton type="primary" size="small" plain @click="addMenu">新增菜单</ElButton>
        </div>

        <ElTable :data="tableData" row-key="id" :row-style="highLightChildren">
            <ElTableColumn prop="id" label="菜单ID" width="80"></ElTableColumn>
            <ElTableColumn prop="meta.title" label="展示名称" width="140"></ElTableColumn>
            <ElTableColumn prop="meta.icon" label="图标" width="80">
                <template #default="{ row }">
                    <ElIcon v-if="row.meta?.icon" size="24">
                        <SvgIcon :icon="row.meta.icon" />
                    </ElIcon>
                </template>
            </ElTableColumn>
            <ElTableColumn prop="meta.hidden" label="隐藏菜单" width="100">
                <template #default="{ row }">
                    <ElTag :type="row.meta?.hidden ? 'danger' : 'success'">
                        {{ row.meta?.hidden ? '是' : '否' }}
                    </ElTag>
                </template>
            </ElTableColumn>
            <ElTableColumn prop="name" label="路由name" width="220"></ElTableColumn>
            <ElTableColumn prop="path" label="路由path">
                <template #default="{ row }">
                    <span v-if="!row.parentId">/{{ row.path }}</span>
                    <span v-else>/{{ tableData?.find(e => e.id === row.parentId)?.path }}/{{ row.path }}</span>
                </template>
            </ElTableColumn>
            <ElTableColumn prop="component" label="组件路径"></ElTableColumn>
            <ElTableColumn prop="sort" label="排序" width="60">
                <template #default="{ row }">
                    {{ row.parentId ? `${tableData?.find(e => e.id === row.parentId)?.sort}-${row.sort}` : row.sort }}
                </template>
            </ElTableColumn>
            <ElTableColumn label="操作">
                <template #default="{ row }">
                    <ElButton type="primary" size="small" plain @click="addSubMenu(row)">新增子菜单</ElButton>
                    <ElButton type="primary" size="small" plain @click="editMenu(row)">编辑</ElButton>
                    <ElButton type="danger" size="small" plain @click="deleteMenu(row)">删除</ElButton>
                </template>
            </ElTableColumn>
        </ElTable>

        <MenuForm v-if="dialogVisible" v-model:show="dialogVisible" :menu="choosedMenu" :title="dialogTitle" />
    </div>
</template>
