import alova from '@/utils/alova'

const baseURL = '/authorities'

function getAuthorityList(): Promise<Api.Authority.AuthListRes> {
    return alova.Get<Api.Authority.AuthListRes>(baseURL)
}

function addAuthority(data: Api.Authority.AddAuthorityReq): Promise<Api.Authority.Auth> {
    return alova.Post<Api.Authority.Auth>(baseURL, data)
}

function updateAuthority(data: Api.Authority.UpdateAuthorityReq): Promise<Api.Authority.Auth> {
    return alova.Put<Api.Authority.Auth>(baseURL, data)
}

function deleteAuthority(authorityId: number, authorityName: string): Promise<void> {
    return alova.Delete<void>(baseURL, {
        authorityId,
        authorityName
    })
}

/**
 * 获取指定角色对应的菜单列表
 * @returns 返回一个Promise，resolve时包含指定角色有权限访问的菜单元数据数组
 */
function getAuthorityMenu(authorityId: number) {
    // 发送GET请求获取指定角色权限对应的菜单列表
    return alova.Get<Api.Menu.MenuMeta[]>(`${baseURL}/${authorityId}/menus`, { cacheFor: 0 })
}

function setMenuAuthority(authorityId: number, menus: Api.Menu.MenuMeta[]): Promise<void> {
    return alova.Post<void>(`${baseURL}/${authorityId}/menus`, { menus }, { cacheFor: 0 })
}

export { getAuthorityList, addAuthority, updateAuthority, deleteAuthority, getAuthorityMenu, setMenuAuthority }
