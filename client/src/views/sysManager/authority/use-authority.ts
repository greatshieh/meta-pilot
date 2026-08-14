import { authorityApi } from '@/api'
import { effectScope, Reactive, Ref, ref, watch } from 'vue'

export interface AuthForm {
    parentId: number | undefined
    authorityId: number | undefined
    authorityName: string | undefined
}

interface AuthOptionProp {
    authorityId: number
    authorityName: string
    children?: AuthOptionProp[]
}

export function useAuthForm(visible: Ref<boolean>, formData: Reactive<AuthForm>) {
    const dialogFlag = ref<'add' | 'edit'>('add')
    // 创建effect作用域，用于管理副作用
    const scope = effectScope()
    const authorityOptions = ref<AuthOptionProp[]>([{ authorityId: 0, authorityName: '根用户/严格模式下为当前用户' }])

    function setAuthorityOptions(AuthorityData: Api.Authority.Auth[], optionsData: AuthOptionProp[], disabled: boolean) {
        AuthorityData.forEach(item => {
            // 如果当前权限项有子权限且子权限数组不为空
            if (item.children && item.children.length) {
                // 创建包含children属性的选项对象
                const option = {
                    authorityId: item.authorityId, // 权限ID
                    authorityName: item.authorityName, // 权限名称
                    disabled: disabled || item.authorityId === formData.authorityId,
                    children: [] // 子权限选项数组
                }
                // 递归处理子权限数据
                setAuthorityOptions(item.children, option.children, disabled || item.authorityId === formData.authorityId)
                // 将处理好的选项添加到选项数组中
                optionsData.push(option)
            } else {
                // 如果没有子权限，创建不包含children属性的选项对象
                const option = {
                    authorityId: item.authorityId, // 权限ID
                    authorityName: item.authorityName, // 权限名称
                    disabled: disabled || item.authorityId === formData.authorityId
                }
                // 将选项添加到选项数组中
                optionsData.push(option)
            }
        })
    }

    async function submitForm() {
        const callback = dialogFlag.value === 'add' ? authorityApi.addAuthority(formData) : authorityApi.updateAuthority(formData)
        const msg = dialogFlag.value === 'add' ? `增加权限成功` : `编辑权限属性成功`

        callback
            .then(() => {
                ElMessage.success(msg)
            })
            .finally(() => (visible.value = false))
    }

    // 在effect作用域内设置状态监听
    scope.run(() => {
        // 监听暗黑模式状态变化
        watch(
            () => visible.value, // 监听的暗黑模式计算属性
            val => {
                if (val) {
                    // 初始化权限选项
                    authorityApi.getAuthorityList().then(resp => {
                        if (resp) setAuthorityOptions(resp, authorityOptions.value, false)
                    })
                } else {
                    formData.parentId = undefined
                    formData.authorityId = undefined
                    formData.authorityName = undefined
                }
            }
        )
    })

    return { authorityOptions, dialogFlag, submitForm }
}
