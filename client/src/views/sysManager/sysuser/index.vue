<script setup lang="ts">
import { invalidateCache, type Method } from 'alova'
import type { FormInstance, FormRules } from 'element-plus'
import { computed, reactive, ref, watch } from 'vue'
import useAdminStore from '@/store/modules/admin'
import { SysUserProps, useSysUser } from './shard'
import { authorityApi, sysUserApi } from '@/api'

defineOptions({ name: 'SysUser' })

const adminStore = useAdminStore()
const tableData = ref<Api.SysUser.User[]>([])

const query = reactive<Api.SysUser.ListReq>({
    page: 1,
    pageSize: 15
})

const total = ref<number>(0)

const dialogVisible = ref<boolean>(false)
const dialogFlag = ref<string>('create')
const dialogTitle = computed<string>(() => (dialogFlag.value === 'create' ? '注册新用户' : '更新用户'))

const { defaultUser, mapToUserProps, resetUserInfo, setAuthorityOptions, cssVar } = useSysUser()

const userProps = ref<SysUserProps>(defaultUser)

const authorityOptions = ref<Api.Authority.Auth[]>([])

// 初始化权限选项

const createFormRef = ref<FormInstance>()

// 表单校验规则
const createFormRules = reactive<FormRules<SysUserProps>>({
    userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    nickName: [{ required: true, message: '请输入员工姓名', trigger: 'blur' }],
    passWord: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    authorityId: [{ required: true, message: '请选择角色', trigger: 'blur' }]
})

function apiCallback<T>(promise: Promise<T>, successMsg: string) {
    promise
        .then(() => {
            ElMessage.success(successMsg)
            invalidateCache(sysUserApi.listSysUser(query) as Method)
            handleSearch()
        })
        .finally(() => {
            dialogVisible.value = false
        })
}

function handleCreate(formEl: FormInstance | undefined) {
    if (!formEl) return
    formEl.validate(valid => {
        if (valid) {
            let callback = dialogFlag.value === 'create' ? sysUserApi.registerSysUser : sysUserApi.setUserInfo
            userProps.value.authorityId = userProps.value.authorityIds[0]
            apiCallback(callback(userProps.value), dialogFlag.value === 'create' ? '注册成功' : '更新成功')
        }
    })
}

function handelDelete(id: number, idx: number) {
    ElMessageBox.confirm('确认删除该用户吗？', '删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        sysUserApi.deleteSysUser(id).then(() => {
            ElMessage.success('删除成功')
            tableData.value.splice(idx, 1)
        })
    })
}

function resetPassword(id: number) {
    ElMessageBox.confirm('确认重置该用户密码吗？', '重置确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    })
        .then(() => {
            sysUserApi.resetPassword(id).then(() => {
                ElMessage.success('重置密码成功')
            })
        })
        .catch(() => {
            ElMessage.warning('取消重置密码')
        })
}

function handleSearch() {
    sysUserApi.listSysUser(query).then(resp => {
        tableData.value = resp.list
        total.value = resp.total
    })
}

function handleUpdateShow(row: Api.SysUser.User) {
    userProps.value = mapToUserProps(row)
    dialogFlag.value = 'update'
    dialogVisible.value = true
}

function handleRegisterShow() {
    dialogFlag.value = 'create'
    dialogVisible.value = true
}

function handleDialogClose() {
    if (createFormRef.value) createFormRef.value.resetFields()
    resetUserInfo(userProps)
}

function init() {
    // 获取人员列表
    handleSearch()
    // 初始化权限选项
    authorityApi.getAuthorityList().then(resp => {
        if (resp) setAuthorityOptions(resp, authorityOptions.value)
    })
}

init()

watch(query, () => {
    invalidateCache(sysUserApi.listSysUser(query) as Method)
    handleSearch()
})
</script>

<template>
    <div class="h-full mx-4">
        <ElCard class="card-wrapper mb-8">
            <ElForm :model="query" size="small" inline>
                <ElFormItem label="姓名">
                    <ElInput v-model="query.userName" />
                </ElFormItem>

                <ElFormItem label="电话">
                    <ElInput v-model="query.phone" />
                </ElFormItem>

                <!-- <ElFormItem label="角色">
                    <ElSelect v-model="query.authorityId" style="width: 160px" clearable>
                        <ElOption
                            v-for="item in roleOptions"
                            :key="item.authorityId"
                            :value="item.authorityId"
                            :label="item.authorityName" />
                    </ElSelect>
                </ElFormItem> -->

                <ElFormItem label="状态">
                    <ElSelect v-model="query.isActive" style="width: 160px" clearable>
                        <ElOption key="有效" label="有效" :value="true" />
                        <ElOption key="冻结" label="冻结" :value="false" />
                    </ElSelect>
                </ElFormItem>
            </ElForm>
        </ElCard>

        <div class="flex items-center gap-8 my-4">
            <div class="text-2xl font-bold">用户管理</div>
            <ElButton type="primary" size="small" plain @click="handleRegisterShow">注册新用户</ElButton>
        </div>

        <ElTable :data="tableData">
            <ElTableColumn label="头像" prop="avatar" width="120">
                <template #default="{ row }">
                    <img :src="row.avatar" alt="avatar" class="w-24px h-24px rounded-full" />
                </template>
            </ElTableColumn>
            <ElTableColumn label="用户名" prop="userName" width="180"></ElTableColumn>
            <ElTableColumn label="姓名" prop="nickName" width="180"></ElTableColumn>
            <ElTableColumn label="电话" prop="phone" width="160"></ElTableColumn>
            <ElTableColumn label="邮箱" prop="email" width="240"></ElTableColumn>
            <ElTableColumn label="角色">
                <template #default="{ row }">
                    <div class="flex gap-2">
                        <ElTag
                            v-for="(item, index) in row.authorities"
                            disable-transitions
                            :color="cssVar[`bg-${(index as number) + 1}`]"
                            :style="{ color: cssVar[`text-${(index as number) + 1}`] }"
                            size="small">
                            {{ item.authorityName }}
                        </ElTag>
                    </div>
                </template>
            </ElTableColumn>

            <ElTableColumn label="状态" width="120">
                <template #default="{ row }">
                    <ElTag disable-transitions :type="row.isActive ? 'success' : 'warning'" size="small" round>
                        <SvgIcon
                            class="inline-block size-1.25em mr-2"
                            :icon="row.isActive ? 'ep:success-filled' : 'fa-solid:dot-circle'"></SvgIcon>
                        {{ row.isActive ? '有效' : '冻结' }}
                    </ElTag>
                </template>
            </ElTableColumn>

            <ElTableColumn label="是否员工" width="120">
                <template #default="{ row }">
                    <ElTag disable-transitions :type="row.isStaff ? 'success' : 'warning'" size="small" round>
                        <SvgIcon
                            class="inline-block size-1.25em mr-2"
                            :icon="row.isStaff ? 'ep:success-filled' : 'fa-solid:dot-circle'"></SvgIcon>
                        {{ row.isStaff ? '是' : '否' }}
                    </ElTag>
                </template>
            </ElTableColumn>

            <ElTableColumn label="操作">
                <template #default="scope">
                    <ElButton
                        type="danger"
                        size="small"
                        :disabled="scope.row.userName === adminStore.userInfo.userName"
                        plain
                        @click="handelDelete(scope.row.id, scope.$index)">
                        删除
                    </ElButton>
                    <ElButton
                        type="primary"
                        size="small"
                        :disabled="adminStore.userInfo.authority.authorityId !== 200 && scope.row.userName === adminStore.userInfo.userName"
                        plain
                        @click="handleUpdateShow(scope.row)">
                        修改
                    </ElButton>
                    <ElButton
                        type="danger"
                        size="small"
                        :disabled="scope.row.userName === adminStore.userInfo.userName"
                        plain
                        @click="resetPassword(scope.row.id)">
                        重置密码
                    </ElButton>
                </template>
            </ElTableColumn>
        </ElTable>

        <pagination
            v-model:page="query.page"
            v-model:page-size="query.pageSize"
            v-model:total="total"
            hide-on-single-page
            @pagination="handleSearch"></pagination>

        <ElDialog v-model="dialogVisible" :show-close="false" @close="handleDialogClose">
            <template #header>
                <div class="flex items-center justify-between">
                    <span>{{ dialogTitle }}</span>
                    <div class="gap-4">
                        <ElButton type="info" size="small" @click="dialogVisible = false">取消</ElButton>
                        <ElButton type="primary" size="small" @click="handleCreate(createFormRef)">确定</ElButton>
                    </div>
                </div>
            </template>
            <ElForm ref="createFormRef" :model="userProps" size="small" :rules="createFormRules" label-position="top">
                <div class="flex gap-4">
                    <ElFormItem label="用户名" prop="userName">
                        <ElInput v-model="userProps.userName"></ElInput>
                    </ElFormItem>
                    <ElFormItem label="姓名" prop="nickName">
                        <ElInput v-model="userProps.nickName"></ElInput>
                    </ElFormItem>
                    <ElFormItem label="电话" prop="phone">
                        <ElInput v-model="userProps.phone"></ElInput>
                    </ElFormItem>
                    <ElFormItem label="邮箱" prop="email">
                        <ElInput v-model="userProps.email"></ElInput>
                    </ElFormItem>
                    <ElFormItem label="密码" prop="passWord" v-if="dialogTitle.includes('注册')">
                        <ElInput v-model="userProps.passWord"></ElInput>
                    </ElFormItem>

                    <ElFormItem label="状态" prop="isActive">
                        <ElSwitch v-model="userProps.isActive"></ElSwitch>
                    </ElFormItem>

                    <ElFormItem label="是否员工" prop="isStaff">
                        <ElSwitch v-model="userProps.isStaff"></ElSwitch>
                    </ElFormItem>
                </div>

                <ElFormItem label="角色" prop="authorityId">
                    <ElCascader
                        class="w-50%"
                        v-model="userProps.authorityIds"
                        :options="authorityOptions"
                        :show-all-levels="false"
                        :props="{
                            multiple: true,
                            checkStrictly: true,
                            label: 'authorityName',
                            value: 'authorityId',
                            emitPath: false
                        }"
                        :clearable="false"></ElCascader>
                </ElFormItem>
            </ElForm>
        </ElDialog>
    </div>
</template>

<style lang="scss" scoped>
.pagination-container {
    width: 80%;
    margin: auto;
}
</style>
