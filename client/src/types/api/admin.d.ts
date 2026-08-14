declare namespace Api.Admin {
    /* api/method/admin.ts */

    /** 管理员登录请求参数 */
    interface LoginReq {
        userName: string
        password: string
    }

    /** 管理员登录响应参数 */
    interface LoginRes {
        token: string
    }

    /** 管理员用户信息响应参数 */
    interface UserInfoRes {
        ID: 1
        CreatedAt: '2025-11-29T21:08:16.036+08:00'
        UpdatedAt: '2025-11-29T21:08:16.044+08:00'
        DeletedAt: null
        uuid: '919147f3-2b18-46e0-9e6d-b89c662f5fdd'
        userName: 'admin'
        nickName: '超级管理员'
        avatar: 'https://qmplusimg.henrongyi.top/gva_header.jpg'
        introduction: ''
        authorityId: 200
        authority: {
            CreatedAt: '2025-11-29T21:08:15.568+08:00'
            UpdatedAt: '2025-11-29T21:08:15.642+08:00'
            DeletedAt: null
            authorityId: 200
            authorityName: '超级管理员'
            parentId: 0
            authorized_sub_roles: null
            children: null
            menus: null
            defaultRouter: 'dashboard'
        }
        authorities: null
        phone: '17611111111'
        email: '333333333@qq.com'
        isStaff: false
        isActive: true
        isSuperuser: false
        lastLogin: '2025-11-30 22:01:24'
    }

    /** 管理员密码修改请求参数 */
    interface ChangePasswordReq {
        password: string
        newPassword: string
    }

    /** 管理员操作进度响应参数 */
    interface ProcessRes {
        progress: number
        total: number
        totalErr: number
        totalSuccess: number
        respErr: Array<Commons.respErr>
        status: string
    }

    /** 管理员权限列表响应参数 */
    interface AuthorityListRes {
        authorityId: string
        authorityName: string
    }
}
