import { mixColors } from '@/hooks/use-color'
import { Ref } from 'vue'

export interface AuthOptionProp {
    authorityId: number
    authorityName: string
    children?: AuthOptionProp[]
}

export interface SysUserProps extends Omit<Api.SysUser.User, 'authority' | 'authorities'> {
    authorityIds: number[]
}

export function useSysUser() {
    const defaultUser: SysUserProps = {
        id: 0, // 用户ID，初始为0
        userName: '', // 用户名，初始为空字符串
        nickName: '', // 昵称，初始为空字符串
        passWord: '', // 密码，初始为空字符串
        email: '', // 邮箱，初始为空字符串
        avatar: '', // 头像地址，初始为空字符串
        phone: '', // 手机号，初始为空字符串
        authorityId: 0, // 主要权限ID，初始为0
        authorityIds: [], // 权限ID数组，初始为空数组
        isActive: true, // 是否激活，初始为true
        isStaff: false, // 是否员工，初始为false
        lastLogin: '' // 最后登录时间，初始为空字符串
    }
    /**
     * 将后端返回的用户数据映射为前端所需的用户属性格式
     * @param user 后端返回的原始用户数据
     * @returns 映射后的用户属性对象，包含额外处理的权限字段
     */
    function mapToUserProps(user: Api.SysUser.User): SysUserProps {
        return {
            ...user,
            // 提取用户的主要权限ID，如果不存在则默认为0
            authorityId: user.authority?.authorityId || 0,
            // 提取用户的所有权限ID数组，如果不存在则默认为空数组
            authorityIds: user.authorities?.map(item => item.authorityId) || []
        }
    }

    /**
     * 重置用户信息到初始状态
     * @param userProps 用户信息的响应式引用对象
     */
    function resetUserInfo(userProps: Ref<SysUserProps>) {
        userProps.value = defaultUser
    }

    /**
     * 递归设置权限选项数据，将后端返回的权限树结构转换为前端使用的选项数组
     * @param AuthorityData 后端返回的权限数据数组（可能包含嵌套的子权限）
     * @param optionsData 用于存储转换后权限选项的数组（会递归填充）
     */
    function setAuthorityOptions(AuthorityData: Api.Authority.Auth[], optionsData: AuthOptionProp[]) {
        AuthorityData.forEach(item => {
            // 如果当前权限项有子权限且子权限数组不为空
            if (item.children && item.children.length) {
                // 创建包含children属性的选项对象
                const option = {
                    authorityId: item.authorityId, // 权限ID
                    authorityName: item.authorityName, // 权限名称
                    children: [] // 子权限选项数组
                }
                // 递归处理子权限数据
                setAuthorityOptions(item.children, option.children)
                // 将处理好的选项添加到选项数组中
                optionsData.push(option)
            } else {
                // 如果没有子权限，创建不包含children属性的选项对象
                const option = {
                    authorityId: item.authorityId, // 权限ID
                    authorityName: item.authorityName // 权限名称
                }
                // 将选项添加到选项数组中
                optionsData.push(option)
            }
        })
    }

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

    return {
        defaultUser,
        mapToUserProps,
        resetUserInfo,
        setAuthorityOptions,
        cssVar
    }
}
