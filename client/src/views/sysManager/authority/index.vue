<script setup lang="ts">
import { authorityApi } from '@/api'
import { reactive, ref, onMounted } from 'vue'
import { invalidateCache, Method } from 'alova'
import { AuthForm, useAuthForm } from './use-authority'
import { ElMessageBox } from 'element-plus'
import SetMenu from './components/SetMenu.vue'
import SetApi from './components/SetApi.vue'

defineOptions({ name: 'Authority' })

const tableData = ref<Api.Authority.Auth[]>([])

// 显示角色设置Dialog
const showSetAuthority = ref(false)
// 显示增加或修改角色Dialog
const showAddOrUpdate = ref(false)

// 需要操作的authorityId
const choosedAuthorityId = ref<number>()

const authFormData = reactive<AuthForm>({
    parentId: undefined,
    authorityId: undefined,
    authorityName: undefined
})

const { authorityOptions, dialogFlag, submitForm } = useAuthForm(showAddOrUpdate, authFormData)

function init() {
    authorityApi.getAuthorityList().then(resp => {
        tableData.value = resp
    })
}

function handleSetAuth(authorityId: number) {
    choosedAuthorityId.value = authorityId
    showSetAuthority.value = true
}

function handleAddOrUpdate(flag: 'add' | 'edit', parentId: number, authorityId?: number, authorithName?: string) {
    dialogFlag.value = flag

    authFormData.parentId = parentId
    if (authorityId) authFormData.authorityId = authorityId
    if (authorithName) authFormData.authorityName = authorithName

    showAddOrUpdate.value = true
}

function afterEdit() {
    if (tableData.value) {
        invalidateCache(authorityApi.getAuthorityList() as Method)
        init()
    }
}

function handleDeleteAuth(parentId: number, authorityId: number, authorityName: string) {
    ElMessageBox.confirm(`确认删除角色 ${authorityName} ?`, '确认删除', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        authorityApi.deleteAuthority(authorityId, authorityName).then(() => {
            invalidateCache(authorityApi.getAuthorityList() as Method)
            init()
        })
    })
}

onMounted(() => {
    init()
})
</script>

<template>
    <div class="p-6 bg-gradient-to-br from-slate-50 to-slate-200 min-h-[calc(100vh-64px)]">
        <!-- 页面头部 -->
        <div class="mb-8 flex justify-between items-center">
            <h1 class="text-28px font-700 text-slate-800">
                角色管理
                <span class="text-16px text-slate-500">管理系统角色权限，包括新增、编辑、删除和权限配置</span>
            </h1>
            <ElButton type="primary" size="default" @click="handleAddOrUpdate('add', 0)" class="add-btn">
                <i class="ant-design:plus-outlined mr-2"></i>
                新增角色
            </ElButton>
        </div>

        <!-- 角色列表 -->
        <div class="authority-card">
            <div class="px-6 py-5 border-b border-slate-100 bg-slate-50">
                <h2 class="text-18px font-600 text-slate-800 mb-1">
                    角色列表
                    <span class="text-14px text-slate-500 m-0">系统中所有角色的详细信息</span>
                </h2>
            </div>
            <div class="p-6">
                <ElTable :data="tableData" row-key="authorityId" class="authority-table"
                    :header-cell-style="{ backgroundColor: '#f8fafc', fontWeight: '600' }"
                    :row-style="{ transition: 'all 0.3s ease' }">
                    <ElTableColumn prop="authorityID" label="角色ID" width="120"></ElTableColumn>
                    <ElTableColumn prop="authorityName" label="角色名称" min-width="200" width="900"></ElTableColumn>
                    <ElTableColumn label="成员人数" width="120">
                        <template #default="{ row }">
                            <div class="member-count">
                                {{ row.users?.length || 0 }}
                            </div>
                        </template>
                    </ElTableColumn>
                    <ElTableColumn label="操作">
                        <template #default="{ row }">
                            <div class="action-buttons">
                                <ElButton type="primary" size="small" @click="handleSetAuth(row.authorityId)"
                                    class="action-btn action-btn-primary">
                                    <i class="ant-design:setting-outlined mr-1"></i>
                                    修改角色
                                </ElButton>
                                <ElButton type="success" size="small"
                                    @click="handleAddOrUpdate('add', row.parentId, row.authorityId)"
                                    class="action-btn action-btn-success">
                                    <i class="ant-design:folder-add-outlined mr-1"></i>
                                    添加子角色
                                </ElButton>
                                <ElButton type="info" size="small"
                                    @click="handleAddOrUpdate('edit', row.parentId, row.authorityId, row.authorityName)"
                                    class="action-btn action-btn-info">
                                    <i class="ant-design:edit-outlined mr-1"></i>
                                    编辑
                                </ElButton>
                                <ElButton type="danger" size="small"
                                    @click="handleDeleteAuth(row.parentId, row.authorityId, row.authorityName)"
                                    class="action-btn action-btn-danger">
                                    <i class="ant-design:delete-outlined mr-1"></i>
                                    删除
                                </ElButton>
                            </div>
                        </template>
                    </ElTableColumn>
                </ElTable>
            </div>
        </div>

        <!-- 新增/编辑角色对话框 -->
        <ElDialog v-model="showAddOrUpdate" :show-close="false" destroy-on-close @close="afterEdit" width="500px"
            class="auth-dialog">
            <template #header>
                <div class="flex justify-between items-center w-full">
                    <h3 class="text-18px font-600 text-slate-800 m-0">
                        {{ dialogFlag === 'add' ? '新增角色' : '编辑角色属性' }}
                    </h3>
                </div>
            </template>
            <template #default>
                <ElForm :model="authFormData" size="default" label-position="top" class="flex flex-col gap-5">
                    <ElFormItem label="父角色名称" prop="parentId" required>
                        <ElCascader v-model="authFormData.parentId" :disabled="dialogFlag === 'add'"
                            :options="authorityOptions" :show-all-levels="false" filterable :props="{
                                checkStrictly: true,
                                label: 'authorityName',
                                value: 'authorityId',
                                disabled: 'disabled',
                                emitPath: false
                            }" class="w-full"></ElCascader>
                    </ElFormItem>
                    <ElFormItem label="角色ID" prop="authorityId" required>
                        <ElInput v-model.number="authFormData.authorityId" :disabled="dialogFlag === 'edit'"
                            class="w-full"></ElInput>
                    </ElFormItem>
                    <ElFormItem label="角色名称" prop="authorityName" required>
                        <ElInput v-model="authFormData.authorityName" class="w-full"></ElInput>
                    </ElFormItem>
                </ElForm>
            </template>
            <template #footer>
                <div class="flex justify-end gap-3">
                    <ElButton type="info" @click="showAddOrUpdate = false">取消</ElButton>
                    <ElButton type="primary" @click="submitForm">确定</ElButton>
                </div>
            </template>
        </ElDialog>

        <!-- 角色配置对话框 -->
        <ElDialog v-model="showSetAuthority" destroy-on-close title="角色配置" width="800px" height="600px"
            class="set-auth-dialog">
            <div class="h-500px flex flex-col">
                <ElTabs default-value="menu" class="flex-1 flex flex-col">
                    <ElTabPane label="菜单权限" name="menu" class="p-6 flex-1 overflow-auto">
                        <SetMenu :authorityId="choosedAuthorityId"></SetMenu>
                    </ElTabPane>
                    <ElTabPane label="API 权限" name="api" class="p-6 flex-1 overflow-auto">
                        <SetApi :authorityId="choosedAuthorityId"></SetApi>
                    </ElTabPane>
                </ElTabs>
            </div>
        </ElDialog>
    </div>
</template>

<style lang="scss" scoped>
.authority-container {
    padding: 24px;
    background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
    min-height: calc(100vh - 64px);
}

/* 页面头部 */
.page-header {
    margin-bottom: 32px;
}

.page-title {
    font-size: 28px;
    font-weight: 700;
    color: #1e293b;
    margin: 0 0 8px 0;
}

.page-description {
    font-size: 16px;
    color: #64748b;
    margin: 0;
}

/* 操作栏 */
.action-bar {
    margin-bottom: 24px;
    display: flex;
    justify-content: flex-start;
}

.add-btn {
    background: linear-gradient(135deg, #3b82f6, #2563eb);
    border: none;
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 600;
    border-radius: 8px;
    transition: all 0.3s ease;

    &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
    }
}

/* 角色卡片 */
.authority-card {
    background: white;
    border-radius: 12px;
    box-shadow:
        0 4px 6px -1px rgba(0, 0, 0, 0.1),
        0 2px 4px -1px rgba(0, 0, 0, 0.06);
    overflow: hidden;
    margin-bottom: 24px;
}

.card-header {
    padding: 20px 24px;
    border-bottom: 1px solid #f1f5f9;
    background: #f8fafc;
}

.card-title {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
    margin: 0 0 4px 0;
}

.card-subtitle {
    font-size: 14px;
    color: #64748b;
    margin: 0;
}

.card-body {
    padding: 24px;
}

/* 角色表格 */
.authority-table {
    border-radius: 8px;
    overflow: hidden;

    :deep(.el-table__row) {
        &:hover {
            background-color: #f8fafc !important;
        }
    }

    :deep(.el-table__cell) {
        padding: 16px;
    }
}

.member-count {
    background: #e0f2fe;
    color: #0284c7;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
    display: inline-block;
}

/* 操作按钮 */
.action-buttons {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
}

.action-btn {
    border-radius: 6px;
    transition: all 0.3s ease;
    font-size: 12px;
    padding: 4px 12px;

    &:hover {
        transform: translateY(-1px);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    }

    &.action-btn-primary {
        background: linear-gradient(135deg, #3b82f6, #2563eb);
        border: none;
    }

    &.action-btn-success {
        background: linear-gradient(135deg, #10b981, #059669);
        border: none;
    }

    &.action-btn-info {
        background: linear-gradient(135deg, #6366f1, #4f46e5);
        border: none;
    }

    &.action-btn-danger {
        background: linear-gradient(135deg, #ef4444, #dc2626);
        border: none;
    }
}

/* 对话框 */
.auth-dialog {
    :deep(.el-dialog__header) {
        padding: 20px 24px;
        border-bottom: 1px solid #f1f5f9;
        background: #f8fafc;
    }

    :deep(.el-dialog__body) {
        padding: 24px;
    }

    :deep(.el-dialog__footer) {
        padding: 20px 24px;
        border-top: 1px solid #f1f5f9;
    }
}

.dialog-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
}

.dialog-title {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
    margin: 0;
}

.auth-form {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}

/* 角色配置对话框 */
.set-auth-dialog {
    :deep(.el-dialog__body) {
        padding: 0;
    }
}

.set-auth-container {
    height: 500px;
    display: flex;
    flex-direction: column;
}

.auth-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.auth-tab {
    padding: 24px;
    flex: 1;
    overflow: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .authority-container {
        padding: 16px;
    }

    .page-title {
        font-size: 24px;
    }

    .action-buttons {
        flex-direction: column;
        align-items: flex-start;
    }

    .action-btn {
        width: 100%;
        justify-content: center;
    }

    .authority-table {
        :deep(.el-table__cell) {
            padding: 12px 8px;
        }
    }

    .auth-dialog {
        width: 90% !important;
    }

    .set-auth-dialog {
        width: 95% !important;
        height: 80vh !important;
    }
}
</style>
