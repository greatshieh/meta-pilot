declare namespace Api.Authority {
    interface Auth {
        createAt: string
        updateAt: string
        authorityId: number
        parentId: number
        authorityName: string
        defaultRoute: string
        children?: Auth[]
        menus: Api.Menu.MenuMeta[]
        users: Api.SysUser.User[]
    }

    type AuthListRes = Auth[]

    type AddAuthorityReq = Api.Commons.RecordNullable<{
        parentId: number
        authorityId: number
        authorityName: string
    }>

    type UpdateAuthorityReq = AddAuthorityReq
}
