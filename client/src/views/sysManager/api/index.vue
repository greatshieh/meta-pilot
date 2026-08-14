<script setup lang="ts">
import { apiApi } from '@/api'
import { AlovaMethodOptions, AlovaMethodRecord } from '@/constants/app'
import { transformArrayToOption } from '@/utils/common'
import { invalidateCache, Method } from 'alova'
import { FormInstance } from 'element-plus'
import { ref, shallowRef } from 'vue'

defineOptions({ name: 'ApiManager' })

// 新增 api dialog
const addApiDialogVisible = ref(false)
const dialogFormType = ref('add')
const dialogTitle = ref('')

const syncApiDialogVisible = ref(false)

const propDefault: Api.ApiMetaData = {
    id: 0,
    path: '',
    description: '',
    apiGroup: '',
    method: 'POST',
    required: undefined
}

const apiGroupOptions = ref<CommonType.Option<string>[]>()
const createFormRef = shallowRef<FormInstance>()
const menuProps = ref(propDefault)

const listParams = ref<Api.ListReq>({
    apiGroup: '',
    method: undefined,
    required: undefined,
    page: 1,
    pageSize: 10
})
const requiredOptions = ref<CommonType.Option<boolean>[]>([
    { label: '是', value: true },
    { label: '否', value: false }
])
const total = ref(0)

const tableData = ref<Api.ApiMetaData[]>()

// 同步 Api 数据
const syncApiData = ref<Api.syncRes>({
    newApis: [],
    deleteApis: [],
    ignoreApis: []
})

const openingSyncDialog = ref(false)

function openSyncDialog() {
    openingSyncDialog.value = true
    // 获取 api 内容
    apiApi
        .syncApi()
        .then(res => {
            syncApiData.value = res
            syncApiDialogVisible.value = true
        })
        .finally(() => {
            openingSyncDialog.value = false
        })
}

function openAddDialog(dialogType: string) {
    if (dialogType === 'edit') {
        dialogFormType.value = 'update'
        dialogTitle.value = '编辑 Api'
    } else {
        dialogFormType.value = 'add'
        dialogTitle.value = '新增 Api'
    }
    addApiDialogVisible.value = true
}

// 编辑 api，将 api 数据赋值给 dialog 表单
function handleEdit(row: Api.ApiMetaData) {
    menuProps.value = row
    openAddDialog('edit')
}

// 确认提交新增/编辑 api 表单
function closeAddDialog() {
    if (!createFormRef.value) return
    createFormRef.value.validate().then(valid => {
        if (!valid) return

        let callback = dialogFormType.value === 'add' ? apiApi.addApi : apiApi.updateApi

        callback(menuProps.value)
            .then(() => {
                search()
            })
            .finally(() => {
                addApiDialogVisible.value = false
            })
    })
}

// 关闭新增/编辑 api 对话框动画结束时清空表单
function onCloseAddDialog() {
    if (createFormRef.value) createFormRef.value.resetFields()
    menuProps.value = propDefault
}

const syncing = ref(false)
// 确认同步 api
function handleSyncApi() {
    if (syncApiData.value.newApis.some(item => !item.apiGroup || !item.description)) {
        ElMessage.error('api 描述和分组不能为空')
        return
    }

    syncing.value = true
    apiApi
        .enterSyncApi(syncApiData.value)
        .then(() => {
            search()
        })
        .finally(() => {
            syncing.value = false
            syncApiDialogVisible.value = false
        })
}

// 同步单条 api
function handleSyncClick(row: Api.ApiMetaData) {
    if (row.description == '') {
        ElMessage.error('api 描述不能为空')
        return
    }

    if (row.apiGroup == '') {
        ElMessage.error('api 分组不能为空')
        return
    }

    apiApi.addApi(row).then(() => {
        // 从新 api 列表中删除
        syncApiData.value.newApis = syncApiData.value.newApis.filter(item => !(item.path === row.path && item.method === row.method))
        search()
    })
}

function handleIgnoreClick(row: Api.ApiMetaData, flag: boolean) {
    apiApi.ignoreApi({ ...row, flag }).then(() => {
        if (flag) {
            // 从新 api 列表中删除
            syncApiData.value.newApis = syncApiData.value.newApis.filter(item => !(item.path === row.path && item.method === row.method))
            syncApiData.value.ignoreApis.push(row)
        } else {
            syncApiData.value.ignoreApis = syncApiData.value.ignoreApis.filter(
                item => !(item.path === row.path && item.method === row.method)
            )
            syncApiData.value.newApis.push(row)
        }
    })
}

function search() {
    invalidateCache(apiApi.getApiList(listParams.value) as Method)
    apiApi.getApiList(listParams.value).then(res => {
        tableData.value = res.list
        total.value = res.total
    })
}

function init() {
    apiApi.getApiGroups().then(res => {
        apiGroupOptions.value = transformArrayToOption<string>(res)
    })
    search()
}

init()
</script>

<template>
    <div class="h-full mx-4">
        <div class="search-box">
            <ElForm :model="listParams" size="small" inline>
                <ElFormItem label="路径" prop="path">
                    <ElInput v-model="listParams.path" clearable style="width: 180px"></ElInput>
                </ElFormItem>
                <ElFormItem label="描述" prop="description">
                    <ElInput v-model="listParams.description" clearable style="width: 180px"></ElInput>
                </ElFormItem>
                <ElFormItem label="分组" prop="apiGroup">
                    <ElSelect v-model="listParams.apiGroup" clearable style="width: 180px">
                        <ElOption v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value"></ElOption>
                    </ElSelect>
                </ElFormItem>
                <ElFormItem label="方法" prop="method">
                    <ElSelect v-model="listParams.method" clearable style="width: 180px">
                        <ElOption v-for="item in AlovaMethodOptions" :key="item.value" :label="item.label" :value="item.value"></ElOption>
                    </ElSelect>
                </ElFormItem>
                <ElFormItem label="是否必选" prop="required">
                    <ElSelect v-model="listParams.required" clearable style="width: 180px">
                        <ElOption v-for="item in requiredOptions" :key="item.label" :label="item.label" :value="item.value"></ElOption>
                    </ElSelect>
                </ElFormItem>

                <ElFormItem>
                    <ElButton type="primary" size="small" @click="search">查询</ElButton>
                </ElFormItem>
            </ElForm>
        </div>

        <div class="flex items-center my-4">
            <div class="text-2xl font-bold">Api 管理</div>
            <ElButton type="primary" size="small" plain class="ml-8" @click="openAddDialog('add')">新增 Api</ElButton>
            <ElButton type="primary" size="small" plain class="ml-8" @click="openSyncDialog" :loading="openingSyncDialog">
                同步 Api
            </ElButton>
            <ElButton type="primary" size="small" plain class="ml-8">上传 Api</ElButton>
            <ElButton type="primary" size="small" plain class="ml-8">导出 Api</ElButton>
            <ElButton type="primary" size="small" plain class="ml-8">下载模板</ElButton>
        </div>

        <ElTable :data="tableData" row-key="path">
            <ElTableColumn prop="path" label="Api 路径"></ElTableColumn>
            <ElTableColumn prop="description" label="Api 描述"></ElTableColumn>
            <ElTableColumn prop="apiGroup" label="Api 分组"></ElTableColumn>
            <ElTableColumn prop="method" label="Api 方法">
                <template #default="{ row }">
                    {{ AlovaMethodRecord[row.method as keyof typeof AlovaMethodRecord] }}
                </template>
            </ElTableColumn>
            <ElTableColumn prop="required" label="是否必选">
                <template #default="{ row }">
                    {{ row.required ? '是' : '否' }}
                </template>
            </ElTableColumn>
            <ElTableColumn label="操作">
                <template #default="{ row }">
                    <ElButton type="primary" size="small" plain class="ml-8" @click="handleEdit(row)">编辑</ElButton>
                    <ElButton type="danger" size="small" plain class="ml-8">删除</ElButton>
                </template>
            </ElTableColumn>
        </ElTable>

        <Pagination
            v-model:page="listParams.page"
            v-model:limit="listParams.pageSize"
            v-model:total="total"
            hide-on-single-page
            @pagination="search"></Pagination>

        <ElDialog v-model="addApiDialogVisible" width="400" :show-close="false" @close="onCloseAddDialog">
            <template #header>
                <div class="flex items-center justify-between">
                    <span>{{ dialogTitle }}</span>
                    <div class="gap-4">
                        <ElButton type="info" size="small" @click="addApiDialogVisible = false">取消</ElButton>
                        <ElButton type="primary" size="small" @click="closeAddDialog">确定</ElButton>
                    </div>
                </div>
            </template>
            <ElForm ref="createFormRef" :model="menuProps" size="small" label-position="top">
                <ElFormItem label="路径" prop="path" required>
                    <ElInput v-model="menuProps.path"></ElInput>
                </ElFormItem>
                <ElFormItem label="描述" prop="description">
                    <ElInput v-model="menuProps.description"></ElInput>
                </ElFormItem>
                <ElFormItem label="分组" prop="apiGroup">
                    <ElSelect v-model="menuProps.apiGroup" placeholder="请选择或新增" allow-create filterable>
                        <ElOption v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
                    </ElSelect>
                </ElFormItem>

                <div class="flex items-center gap-4">
                    <ElFormItem label="方法" prop="method">
                        <ElSelect v-model="menuProps.method" style="width: 240px">
                            <ElOption
                                v-for="item in AlovaMethodOptions"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value"></ElOption>
                        </ElSelect>
                    </ElFormItem>
                    <ElFormItem label="是否必选" prop="required">
                        <ElSwitch v-model="menuProps.required"></ElSwitch>
                    </ElFormItem>
                </div>
            </ElForm>
        </ElDialog>

        <ElDialog
            v-model="syncApiDialogVisible"
            :show-close="false"
            width="70%"
            class="h-90% mb-0 overflow-auto"
            style="--el-dialog-margin-top: 50px">
            <template #header>
                <div class="flex items-center justify-between">
                    <span class="font-bold text-28px">同步 Api</span>
                    <div class="gap-4">
                        <ElButton type="info" size="small" @click="syncApiDialogVisible = false">取消</ElButton>
                        <ElButton type="primary" size="small" @click="handleSyncApi" :loading="syncing">确定</ElButton>
                    </div>
                </div>
                <ElAlert
                    type="warning"
                    title="同步后, 新增api将写入数据表保存; 删除 api 将从数据表删除; 忽略的 api 将不会同步"
                    show-icon
                    :closable="false"></ElAlert>
            </template>

            <div class="mb-4 text-2xl c-primary flex items-center gap-4">
                <span>新增 Api</span>
                <ElAlert
                    class="flex-1"
                    type="info"
                    title="系统已注册, 但没有在 api 数据表中. 同步后会新增到数据表, 且加入鉴权"
                    show-icon
                    :closable="false"></ElAlert>
            </div>
            <ElTable :data="syncApiData.newApis" row-key="path" class="sync-table">
                <ElTableColumn prop="path" label="Api 路径"></ElTableColumn>
                <ElTableColumn prop="description" label="Api 描述">
                    <template #default="{ row }">
                        <ElInput v-model="row.description"></ElInput>
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="apiGroup" label="Api 分组" width="220">
                    <template #default="{ row }">
                        <ElSelect v-model="row.apiGroup" placeholder="请选择或新增" allow-create filterable>
                            <ElOption v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
                        </ElSelect>
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="method" label="Api 方法" width="220">
                    <template #default="{ row }">
                        {{ AlovaMethodRecord[row.method as keyof typeof AlovaMethodRecord] }}
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="required" label="是否必选" width="220">
                    <template #default="{ row }">
                        <ElSwitch v-model="row.required"></ElSwitch>
                    </template>
                </ElTableColumn>
                <ElTableColumn label="操作" width="220">
                    <template #default="{ row }">
                        <ElTooltip content="同步此 Api 到系统中注册">
                            <ElButton type="primary" plain size="small" @click="handleSyncClick(row)">同步</ElButton>
                        </ElTooltip>
                        <ElTooltip content="忽略此 Api, 不进行同步">
                            <ElButton type="warning" plain size="small" @click="handleIgnoreClick(row, true)">忽略</ElButton>
                        </ElTooltip>
                    </template>
                </ElTableColumn>
            </ElTable>

            <div class="my-4 text-2xl c-primary flex items-center gap-4">
                <span>已删除 Api</span>
                <ElAlert
                    class="flex-1"
                    type="info"
                    title="在 api 数据表中, 但没有在系统中注册. 同步后会从数据表中删除"
                    show-icon
                    :closable="false"></ElAlert>
            </div>
            <ElTable :data="syncApiData.deleteApis" row-key="path" class="sync-table">
                <ElTableColumn prop="path" label="Api 路径"></ElTableColumn>
                <ElTableColumn prop="description" label="Api 描述">
                    <template #default="{ row }">
                        <ElInput v-model="row.description"></ElInput>
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="apiGroup" label="Api 分组" width="220">
                    <template #default="{ row }">
                        <ElSelect v-model="row.apiGroup" placeholder="请选择或新增" allow-create filterable>
                            <ElOption v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
                        </ElSelect>
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="method" label="Api 方法" width="220">
                    <template #default="{ row }">
                        {{ AlovaMethodRecord[row.method as keyof typeof AlovaMethodRecord] }}
                    </template>
                </ElTableColumn>
            </ElTable>

            <div class="my-4 text-2xl c-primary flex items-center gap-4">
                <span>已忽略 Api</span>
                <ElAlert class="flex-1" type="info" title="不需要参与鉴权的 api" show-icon :closable="false" size="small"></ElAlert>
            </div>
            <ElTable :data="syncApiData.ignoreApis" row-key="path" class="sync-table">
                <ElTableColumn prop="path" label="Api 路径"></ElTableColumn>
                <ElTableColumn prop="description" label="Api 描述">
                    <template #default="{ row }">
                        <ElInput v-model="row.description"></ElInput>
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="apiGroup" label="Api 分组" width="220">
                    <template #default="{ row }">
                        <ElSelect v-model="row.apiGroup" placeholder="请选择或新增" allow-create filterable>
                            <ElOption v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
                        </ElSelect>
                    </template>
                </ElTableColumn>
                <ElTableColumn prop="method" label="Api 方法" width="220">
                    <template #default="{ row }">
                        {{ AlovaMethodRecord[row.method as keyof typeof AlovaMethodRecord] }}
                    </template>
                </ElTableColumn>
                <ElTableColumn label="取消忽略">
                    <template #default="{ row }">
                        <ElButton type="primary" size="small" @click="handleIgnoreClick(row, false)">恢复</ElButton>
                    </template>
                </ElTableColumn>
            </ElTable>
        </ElDialog>
    </div>
</template>

<style lang="scss">
.sync-table {
    .el-scrollbar__wrap {
        max-height: calc(15vh);
    }
}
</style>
