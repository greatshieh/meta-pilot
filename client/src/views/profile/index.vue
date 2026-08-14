<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { InternalRuleItem } from 'async-validator'
import { reactive, ref, shallowRef } from 'vue'
import useAdminStore from '@/store/modules/admin'
import useThemeStore from '@/store/modules/theme'
import { sysUserApi } from '@/api'
import { mixColors } from '@/hooks/use-color'

const cssVar: Record<string, string> = {
    'text-1': '#eb3b5a',
    'text-2': '#a55eea',
    'text-3': '#0fb9b1',
    'text-4': '#fd9644',
    'text-5': '#778ca3',
    'bg-1': mixColors('#eb3b5a', '#ffffff', 0.9),
    'bg-2': mixColors('#a55eea', '#ffffff', 0.9),
    'bg-3': mixColors('#0fb9b1', '#ffffff', 0.9),
    'bg-4': mixColors('#fd9644', '#ffffff', 0.9),
    'bg-5': mixColors('#778ca3', '#ffffff', 0.9)
}

defineOptions({ name: 'Profile' })

const adminStore = useAdminStore()
const themeStore = useThemeStore()

const dialogTitle = ref('')
const showChangeDialog = ref(false)
const dialogFormRef = shallowRef<FormInstance>()
const dialogForm = reactive({
    content: '',
    code: ''
})

const showChangePWD = ref(false)
const pwdFormRef = shallowRef<FormInstance>()
const pwdForm = reactive({
    password: '',
    newPassword: '',
    checkPassword: ''
})

const activeTab = ref('notice')

const validatePass = (_rule: InternalRuleItem, value: string, callback: (error?: string | Error) => void) => {
    if (value === '') {
        callback(new Error('请输入密码'))
    } else {
        if (pwdForm.newPassword !== '') {
            if (!pwdFormRef.value) return
            pwdFormRef.value.validateField('checkPass')
        }
        callback()
    }
}
const validatePass2 = (_rule: InternalRuleItem, value: string, callback: (error?: string | Error) => void) => {
    if (value === '') {
        callback(new Error('请再次输入密码'))
    } else if (value !== pwdForm.newPassword) {
        callback(new Error('两次输入的密码不一致'))
    } else {
        callback()
    }
}

const rules = reactive<FormRules<typeof pwdForm>>({
    password: [{ required: true, validator: validatePass, trigger: 'blur' }],
    newPassword: [{ required: true, message: '新密码不能为空', trigger: 'blur' }],
    checkPassword: [{ required: true, validator: validatePass2, trigger: 'blur' }]
})

const handlePwdSubmit = (formInstance: FormInstance | undefined) => {
    if (!formInstance) return
    formInstance.validate(valid => {
        if (valid) {
            sysUserApi
                .changePassword({
                    password: pwdForm.password,
                    newPassword: pwdForm.newPassword
                })
                .then(() => {
                    ElMessage.success('密码修改成功')
                })
                .finally(() => (showChangePWD.value = false))
        }
    })
}

const handlePwdCancel = (formInstance: FormInstance | undefined) => {
    if (!formInstance) return
    formInstance.resetFields()
    showChangePWD.value = false
}

function handleOpenDialog(title: string) {
    dialogTitle.value = title
    showChangeDialog.value = true
}

function getCode() {
    ElMessageBox({ message: '示例程序，获取验证码功能未实现', type: 'info' })
}

function handleDialogSubmit(formInstance: FormInstance | undefined) {
    if (!formInstance) return
    formInstance.validate(valid => {
        if (valid) {
            ElMessageBox.confirm(`确认修改${dialogTitle.value}吗？`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                ElMessageBox({ message: `修改${dialogTitle.value}功能未实现`, type: 'info' }).then(() => (showChangeDialog.value = false))
            })
        }
    })
}

function handleDialogCancel(formInstance: FormInstance | undefined) {
    if (!formInstance) return
    formInstance.resetFields()
    showChangeDialog.value = false
}
</script>

<template>
    <div class="h-full mx-4">
        <ElCard class="card-wrapper mb-8">
            <div class="flex flex-col items-center justify-center gap-4">
                <img :src="adminStore.userInfo.avatar" width="128px" heigh="128px" class="rounded-full" />
                <div class="flex items-center gap-8 text-16px">
                    <div>
                        <SvgIcon icon="mdi:map-marker-circle" class="inline-block -mt-2px c-text-primary"></SvgIcon>
                        <span class="ml-2 c-text-secondary">西安·想回武汉</span>
                    </div>

                    <div>
                        <SvgIcon icon="mdi:office-building-outline" class="inline-block -mt-2px c-text-primary"></SvgIcon>
                        <span class="ml-2 c-text-secondary">自由职业</span>
                    </div>

                    <div>
                        <SvgIcon icon="hugeicons:new-job" class="inline-block -mt-2px c-text-primary"></SvgIcon>
                        <span class="ml-2 c-text-secondary">全栈工程师</span>
                    </div>
                </div>
            </div>
        </ElCard>

        <div class="grid grid-cols-12 gap-8">
            <div class="col-span-4 space-y-8">
                <div class="card-wrapper bg-[var(--el-fill-color-lighter)] p-4 gap-4 flex flex-col space-y-2">
                    <div class="text-18px flex items-center gap-2 text-color-primary">
                        <SvgIcon icon="ant-design:info-circle-filled"></SvgIcon>
                        <span>基本信息</span>
                    </div>
                    <div class="flex items-center text-14px gap-2">
                        <SvgIcon icon="ant-design:user-outlined" class="mt-1px text-blue-500"></SvgIcon>
                        <span>用户：</span>
                        <span>{{ adminStore.userInfo.nickName }}</span>
                        <ElButton type="primary" text size="small" class="ml-auto" disabled>修改</ElButton>
                    </div>
                    <div class="flex items-center text-14px gap-2">
                        <SvgIcon icon="tabler:brand-auth0" class="mt-1px text-blue-500"></SvgIcon>
                        <span>权限：</span>
                        <ElTag
                            v-for="(item, index) in adminStore.userInfo.authorities"
                            disable-transitions
                            :color="cssVar[`bg-${index + 1}`]"
                            :style="{ color: cssVar[`text-${index + 1}`] }"
                            size="small">
                            {{ item.authorityName }}
                        </ElTag>
                        <ElButton type="primary" text size="small" class="ml-auto" disabled>修改</ElButton>
                    </div>
                    <div class="flex items-center text-14px gap-2">
                        <SvgIcon icon="gg:phone" class="mt-1px text-purple-500"></SvgIcon>
                        <span>手机：</span>
                        <span>{{ adminStore.userInfo.phone }}</span>
                        <ElButton type="primary" text size="small" class="ml-auto" @click="handleOpenDialog('手机')">修改</ElButton>
                    </div>
                    <div class="flex items-center text-14px gap-2">
                        <SvgIcon icon="mdi:email-outline" class="mt-1px text-green-500"></SvgIcon>
                        <span>邮箱：</span>
                        <span>{{ adminStore.userInfo.email }}</span>
                        <ElButton type="primary" text size="small" class="ml-auto" @click="handleOpenDialog('邮箱')">修改</ElButton>
                    </div>
                    <div class="flex items-center text-14px gap-2">
                        <SvgIcon icon="mdi:form-textbox-password" class="mt-1px text-orange-500"></SvgIcon>
                        <span>密码：</span>
                        <span>已设置</span>
                        <ElButton type="primary" text size="small" class="ml-auto" @click="showChangePWD = true">修改</ElButton>
                    </div>
                </div>
                <div class="card-wrapper bg-[var(--el-fill-color-lighter)] p-4 gap-4 flex flex-col space-y-2">
                    <div class="text-18px flex items-center gap-2 text-color-primary">
                        <SvgIcon icon="simple-icons:hyperskill"></SvgIcon>
                        <span>技能</span>
                    </div>
                    <div class="flex items-center text-28px gap-4">
                        <SvgIcon icon="skill-icons:golang"></SvgIcon>
                        <SvgIcon :icon="themeStore.isDark ? 'skill-icons:vuejs-dark' : 'skill-icons:vuejs-light'"></SvgIcon>
                        <SvgIcon icon="skill-icons:typescript"></SvgIcon>
                        <SvgIcon icon="skill-icons:rust"></SvgIcon>
                        <SvgIcon :icon="themeStore.isDark ? 'skill-icons:mysql-dark' : 'skill-icons:mysql-light'"></SvgIcon>
                        <SvgIcon :icon="themeStore.isDark ? 'skill-icons:redis-dark' : 'skill-icons:redis-light'"></SvgIcon>
                        <SvgIcon :icon="themeStore.isDark ? 'skill-icons:python-dark' : 'skill-icons:python-light'"></SvgIcon>
                    </div>
                </div>
            </div>

            <div class="col-span-8 card-wrapper bg-[var(--el-fill-color-lighter)] p-4 gap-4 flex flex-col space-y-2">
                <el-tabs v-model="activeTab">
                    <el-tab-pane name="notice">
                        <template #label>
                            <div class="flex items-center gap-2 c-text-regular text-18px">
                                <SvgIcon icon="ant-design:info-circle-filled" class="mt-2px"></SvgIcon>
                                <span>通知</span>
                            </div>
                        </template>
                        <ElScrollbar view-class="space-y-8 w-full" height="500px" max-height="800px">
                            <ElCard>
                                <template #header>
                                    <span>人力资源部：</span>
                                    <span>关于全员涨薪的通知</span>
                                </template>
                                <p class="text-14px c-text-regular">
                                    所有员工，包括前端、后端、数据库、测试等，都将收到此通知。请在规定时间内填写期望涨薪幅度，否则将被视为未完成涨薪。
                                </p>
                                <template #footer>
                                    <span class="text-14px c-text-regular">2025-12-12</span>
                                </template>
                            </ElCard>
                            <ElCard>
                                <template #header>
                                    <span>产品管理部：</span>
                                    <span>关于新版本发布的通知</span>
                                </template>
                                <p class="text-14px c-text-regular">
                                    产品管理部已发布新版本，新功能包括但不限于：用户管理、角色管理、权限管理、日志管理等。请及时更新您的系统。
                                </p>
                                <template #footer>
                                    <span class="text-14px c-text-regular">2025-11-11</span>
                                </template>
                            </ElCard>
                            <ElCard>
                                <template #header>
                                    <span>人力资源部：</span>
                                    <span>关于全员放假的通知</span>
                                </template>
                                <p class="text-14px c-text-regular">
                                    国庆佳节，所有员工放假两周。放假期间，禁止工作，所有领导禁止以任何方式，包括但不限于电话、邮件、短信、企业微信、钉钉等，与员工进行任何沟通，除非是请员工吃饭和发红包。
                                </p>
                                <template #footer>
                                    <span class="text-14px c-text-regular">2025-9-30</span>
                                </template>
                            </ElCard>
                        </ElScrollbar>
                    </el-tab-pane>
                    <el-tab-pane name="project">
                        <template #label>
                            <div class="flex items-center gap-2 c-text-regular text-18px">
                                <SvgIcon icon="roentgen:milestone" class="mt-2px"></SvgIcon>
                                <span>项目</span>
                            </div>
                        </template>
                        <ElScrollbar view-class="space-y-8 w-full" height="400px" max-height="500px">
                            <el-timeline style="max-width: 100%">
                                <el-timeline-item timestamp="2025-08-30" placement="top">
                                    <div class="text-16px c-text-primary">后端框架完成</div>
                                    <div class="text-12px c-text-secondary mt-2">
                                        完成后端框架，包括登录页、注册页、首页、用户管理页、角色管理页、权限管理页、日志管理页等。
                                    </div>
                                </el-timeline-item>
                                <el-timeline-item timestamp="2025-08-30" placement="top">
                                    <div class="text-16px c-text-primary">后端框架完成</div>
                                    <div class="text-12px c-text-secondary mt-2">
                                        完成后端框架，包括登录页、注册页、首页、用户管理页、角色管理页、权限管理页、日志管理页等。
                                    </div>
                                </el-timeline-item>
                                <el-timeline-item timestamp="2025-07-30" placement="top">
                                    <div class="text-16px c-text-primary">前端框架完成</div>
                                    <div class="text-12px c-text-secondary mt-2">
                                        完成前端框架，包括登录页、注册页、首页、用户管理页、角色管理页、权限管理页、日志管理页等。
                                    </div>
                                </el-timeline-item>
                                <el-timeline-item timestamp="2025-06-30" placement="top">
                                    <div class="text-16px c-text-primary">新项目启动</div>
                                    <div class="text-12px c-text-secondary mt-2">通用vue3+typescript+gin后台模板</div>
                                </el-timeline-item>
                            </el-timeline>
                        </ElScrollbar>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </div>

        <ElDialog v-model="showChangeDialog" :title="`修改${dialogTitle}`" width="400">
            <ElForm ref="dialogFormRef" :model="dialogForm" size="small" label-width="80px">
                <ElFormItem :label="dialogTitle">
                    <ElInput v-model="dialogForm.content" autocomplete="off" :placeholder="`请输入${dialogTitle}`">
                        <template #prefix>
                            <SvgIcon v-if="dialogTitle === '手机'" icon="gg:phone" class="mt-1px"></SvgIcon>
                            <SvgIcon v-else icon="mdi:email-outline" class="mt-1px"></SvgIcon>
                        </template>
                    </ElInput>
                </ElFormItem>
                <ElFormItem label="验证码">
                    <ElInput v-model="dialogForm.code" autocomplete="off" style="width: 60%">
                        <template #prefix>
                            <SvgIcon icon="grommet-icons:validate" class="mt-1px"></SvgIcon>
                        </template>
                    </ElInput>
                    <ElButton type="primary" @click="getCode" class="ml-auto">获取验证码</ElButton>
                </ElFormItem>
                <ElFormItem>
                    <ElButton type="info" class="ml-auto" @click="handleDialogCancel(dialogFormRef)">取消</ElButton>
                    <ElButton type="primary" class="ml-auto" @click="handleDialogSubmit(dialogFormRef)">确认</ElButton>
                </ElFormItem>
            </ElForm>
        </ElDialog>

        <ElDialog v-model="showChangePWD" title="修改密码" width="400">
            <ElForm ref="pwdFormRef" :model="pwdForm" :rules="rules" label-position="top" size="small" label-width="80px">
                <ElFormItem label="原密码" prop="password">
                    <ElInput v-model="pwdForm.password" type="password" autocomplete="off" show-password />
                </ElFormItem>
                <ElFormItem label="新密码" prop="newPassword">
                    <ElInput v-model="pwdForm.newPassword" type="password" autocomplete="off" show-password />
                </ElFormItem>
                <ElFormItem label="确认密码" prop="checkPassword">
                    <ElInput v-model="pwdForm.checkPassword" type="password" autocomplete="off" show-password />
                </ElFormItem>
                <ElFormItem>
                    <ElButton type="info" class="ml-auto" @click="handlePwdCancel(pwdFormRef)">取消</ElButton>
                    <ElButton type="primary" class="ml-auto" @click="handlePwdSubmit(pwdFormRef)">确认</ElButton>
                </ElFormItem>
            </ElForm>
        </ElDialog>
    </div>
</template>
