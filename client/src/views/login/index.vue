<template>
    <div class="login-container min-h-screen w-full flex items-center justify-center overflow-hidden relative">
        <!-- 背景渐变 -->
        <div class="absolute inset-0 bg-gradient-to-br from-slate-900 via-indigo-900 to-purple-800"></div>

        <!-- 装饰元素 -->
        <div class="absolute top-0 left-0 w-full h-full overflow-hidden">
            <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-white/10 rounded-full blur-3xl animate-pulse"></div>
            <div
                class="absolute bottom-1/4 right-1/4 w-80 h-80 bg-purple-500/20 rounded-full blur-3xl animate-pulse"
                style="animation-delay: 1s"></div>
        </div>

        <!-- 登录卡片 -->
        <div class="relative z-10 w-full max-w-md mx-4">
            <div
                class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 shadow-2xl transition-all duration-500 hover:shadow-purple-500/20 transform hover:-translate-y-1">
                <!-- Logo -->
                <div class="flex justify-center mb-8">
                    <SystemLogo class="w-128px" />
                </div>

                <!-- 标题 -->
                <div class="text-center mb-8">
                    <h1 class="text-3xl font-bold text-white mb-2 tracking-tight">Meta Pilot Admin</h1>
                    <p class="text-white/70">欢迎回来，请登录您的账户</p>
                </div>

                <!-- 登录表单 -->
                <ElForm class="space-y-6 w-full" :model="formData" ref="form">
                    <!-- 用户名输入 -->
                    <ElFormItem prop="username">
                        <ElInput
                            v-model="formData.userName"
                            placeholder="请输入用户名"
                            class="w-full bg-white/10 border-white/20 text-white placeholder-white/50 rounded-lg"
                            :class="{ 'is-focused': focusState.username }"
                            @focus="focusState.username = true"
                            @blur="focusState.username = false">
                            <template #prefix>
                                <i class="i-local:user text-white/70"></i>
                            </template>
                        </ElInput>
                    </ElFormItem>

                    <!-- 密码输入 -->
                    <ElFormItem prop="passcode">
                        <ElInput
                            ref="password"
                            :type="passwordType"
                            v-model="formData.password"
                            placeholder="请输入密码"
                            class="w-full bg-white/10 border-white/20 text-white placeholder-white/50 rounded-lg"
                            :class="{ 'is-focused': focusState.password }"
                            @focus="focusState.password = true"
                            @blur="focusState.password = false"
                            @keyup="checkCapslock"
                            @keyup.enter="onSubmit">
                            <template #prefix>
                                <i class="i-local:password text-white/70"></i>
                            </template>
                            <template #suffix>
                                <i
                                    :class="passwordType === 'password' ? 'i-local:eye' : 'i-local:eye-open'"
                                    class="cursor-pointer select-none text-white/70 hover:text-white transition-colors"
                                    @click="showPwd" />
                            </template>
                        </ElInput>
                        <!-- 大写锁定提示 -->
                        <div v-if="capsTooltip" class="absolute right-12 top-1/2 -translate-y-1/2 text-xs text-yellow-300">
                            大写锁定已开启
                        </div>
                    </ElFormItem>

                    <!-- 登录按钮 -->
                    <ElButton
                        type="primary"
                        class="w-full h-12 rounded-lg text-base font-medium bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98]"
                        @click="onSubmit">
                        登录
                    </ElButton>

                    <!-- 底部链接 -->
                    <div class="flex justify-between items-center mt-6">
                        <div class="text-white/60 text-sm hover:text-white transition-colors cursor-pointer">忘记密码？</div>
                        <div class="text-white/60 text-sm hover:text-white transition-colors cursor-pointer">联系支持</div>
                    </div>
                </ElForm>
            </div>

            <!-- 版权信息 -->
            <div class="mt-8 text-center text-white/50 text-sm">© 2026 Meta Pilot. 保留所有权利</div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ElLoading } from 'element-plus'
import { nextTick, ref, watch } from 'vue'

import { type LocationQuery, type LocationQueryValue, useRoute, useRouter } from 'vue-router'
import useAdminStore from '@/store/modules/admin'

const adminStore = useAdminStore()
const $router = useRouter()
const $route = useRoute()

const formData = ref<Api.Admin.LoginReq>({ userName: '', password: '' })
const password = ref<HTMLInputElement>()
const capsTooltip = ref<boolean | string>(false)
const focusState = ref({ username: false, password: false })

const redirect = ref<LocationQueryValue | LocationQueryValue[]>(null)
const otherQuery = ref<LocationQuery>({})

const getOtherQuery = (query: LocationQuery): LocationQuery => {
    return Object.keys(query).reduce((acc: LocationQuery, cur: string) => {
        if (cur !== 'redirect') {
            acc[cur] = query[cur]
        }
        return acc
    }, {})
}

const checkCapslock = (e: KeyboardEvent) => {
    const { key } = e
    capsTooltip.value = key && key.length === 1 && key >= 'A' && key <= 'Z'
}

const passwordType = ref('password')
const showPwd = () => {
    if (passwordType.value === 'password') {
        passwordType.value = ''
    } else {
        passwordType.value = 'password'
    }
    nextTick(() => {
        password.value?.focus()
    })
}

watch(
    $route,
    route => {
        const query = route.query
        if (query) {
            redirect.value = query.redirect
            otherQuery.value = getOtherQuery(query)
        }
    },
    { immediate: true }
)

// 提交表单
const onSubmit = () => {
    const loadingInstance = ElLoading.service({
        lock: true,
        text: '正在登录',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
    })

    adminStore
        .login(formData.value)
        .then(() => {
            $router.push({
                path: (redirect.value as string) || '/dashboard',
                query: otherQuery.value
            })
        })
        .finally(() => {
            loadingInstance.close()
        })
}
</script>

<style lang="scss">
/* 全局样式重置 */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

/* 登录容器样式 */
.login-container {
    font-family:
        'Inter',
        'SF Pro Display',
        -apple-system,
        BlinkMacSystemFont,
        'Segoe UI',
        Roboto,
        sans-serif;
    overflow-x: hidden;

    /* 自定义滚动条 */
    &::-webkit-scrollbar {
        width: 6px;
    }

    &::-webkit-scrollbar-track {
        background: rgba(255, 255, 255, 0.1);
    }

    &::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.3);
        border-radius: 3px;
    }
}

/* Element Plus 样式覆盖 */
.login-container {
    .el-input {
        .el-input__wrapper {
            background: transparent !important;
            border: 1px solid rgba(255, 255, 255, 0.2) !important;
            border-radius: 0.5rem !important;
            box-shadow: none !important;
            transition: all 0.3s ease;

            &.is-focused {
                border-color: rgba(255, 255, 255, 0.5) !important;
                box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1) !important;
            }

            input {
                background-color: transparent !important;
                border: none !important;
                color: white !important;
                height: 2.5rem !important;
                caret-color: white !important;

                &::placeholder {
                    color: rgba(255, 255, 255, 0.5) !important;
                }

                &:-webkit-autofill {
                    box-shadow: 0 0 0px 1000px rgba(255, 255, 255, 0.1) inset !important;
                    -webkit-text-fill-color: white !important;
                }
            }

            .el-input__prefix,
            .el-input__suffix {
                color: rgba(255, 255, 255, 0.7) !important;
            }
        }
    }

    .el-button {
        border: none !important;
        font-weight: 500 !important;
        transition: all 0.3s ease !important;
    }

    .el-button--primary {
        background: linear-gradient(90deg, #4f46e5, #7c3aed) !important;

        &:hover {
            opacity: 0.9 !important;
        }

        &:active {
            opacity: 0.8 !important;
        }
    }

    .el-form-item {
        margin-bottom: 1.5rem !important;
    }
}

/* 动画效果 */
@keyframes pulse {
    0%,
    100% {
        opacity: 0.3;
        transform: scale(1);
    }
    50% {
        opacity: 0.6;
        transform: scale(1.05);
    }
}

.animate-pulse {
    animation: pulse 8s ease-in-out infinite;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .login-container {
        padding: 1rem;

        .max-w-md {
            max-width: 95% !important;
        }

        .p-8 {
            padding: 2rem !important;
        }

        h1 {
            font-size: 1.75rem !important;
        }
    }
}

@media (max-width: 480px) {
    .login-container {
        .p-8 {
            padding: 1.5rem !important;
        }

        .size-16 {
            width: 4rem !important;
            height: 4rem !important;
        }
    }
}
</style>
