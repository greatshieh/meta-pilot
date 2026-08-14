declare namespace Api.SysUser {
    interface LoginReq {
        userName: string
        password: string
    }

    /** 管理员登录响应参数 */
    interface LoginRes {
        token: string
    }

    interface ChangePasswordReq {
        password: string
        newPassword: string
    }

    interface User {
        id?: number
        userName: string
        nickName: string
        passWord?: string
        email: string
        avatar: string
        phone: string
        authorityId: number
        authority: Api.Authority.Auth
        authorities: Api.Authority.Auth[]
        isActive: boolean
        isStaff: boolean
        lastLogin: string
    }

    type ListReq = Api.Commons.RecordNullable<Omit<User, 'passWord'> & Api.Commons.PaginatingCommonParams>

    type ListRes = Api.Commons.PaginatingResponseRecord<Omit<User, 'passWord'>>

    type RegisterReq = Api.Commons.RecordNullable<User>
}
