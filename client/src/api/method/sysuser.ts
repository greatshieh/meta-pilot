import alova from '@/utils/alova'

const baseURL = '/users'

/**
 * 管理员登录接口
 * @param data 登录请求参数，包含管理员账号和密码等信息
 * @returns 返回Promise，包含登录响应数据(如token等)
 */
function login(data: Api.SysUser.LoginReq): Promise<Api.SysUser.LoginRes> {
    return alova.Post<Api.SysUser.LoginRes>(`${baseURL}/login`, data)
}

/**
 * 管理员登出接口
 * @returns 返回Promise，执行登出操作
 */
function logout() {
    return alova.Post(`${baseURL}/logout`)
}

/**
 * 修改管理员密码接口
 * @param data 密码修改请求参数，包含原密码和新密码
 * @returns 返回Promise，执行密码修改操作
 */
function changePassword(data: Api.Admin.ChangePasswordReq) {
    return alova.Patch(`${baseURL}/password`, data)
}

/**
 * 获取登陆用户信息接口
 * @returns 返回Promise，包含管理员用户信息(禁用缓存)
 */
function getSelfInfo(): Promise<Api.SysUser.User> {
    return alova.Get<Api.SysUser.User>(`${baseURL}/me`, {
        cacheFor: 0
    })
}

/**
 * 获取系统用户列表
 * @param params 查询参数对象，包含分页和搜索条件
 * @returns 返回一个Promise，resolve时返回用户列表的响应数据
 */
function listSysUser(params: Api.SysUser.ListReq): Promise<Api.SysUser.ListRes> {
    // 发送GET请求获取用户列表数据，禁用缓存
    return alova.Get<Api.SysUser.ListRes>(baseURL, { params, cacheFor: 0 })
}

/**
 * 注册系统用户
 * @param data 用户注册请求数据，包含用户名、密码等信息
 * @returns 返回一个Promise，resolve时返回创建的用户信息
 */
function registerSysUser(data: Api.SysUser.RegisterReq): Promise<Api.SysUser.User> {
    // 发送POST请求注册新用户
    return alova.Post<Api.SysUser.User>(baseURL, data)
}

/**
 * 删除系统用户
 * @param id 要删除的用户ID
 * @returns 返回一个Promise，resolve时代表删除成功（无返回值）
 */
function deleteSysUser(id: number): Promise<void> {
    // 发送DELETE请求删除指定ID的用户
    return alova.Delete<void>(`${baseURL}/${id}`)
}

/**
 * 更新用户信息
 * @param data 用户信息更新请求数据，包含需要修改的用户信息
 * @returns 返回一个Promise，resolve时返回更新后的用户信息
 */
function setUserInfo(data: Api.SysUser.RegisterReq): Promise<Api.SysUser.User> {
    // 发送PUT请求更新用户信息，禁用缓存
    return alova.Put<Api.SysUser.User>(`${baseURL}/${data.id}`, data, { cacheFor: 0 })
}

/**
 * 重置用户密码
 * @param id 用户ID，用于指定需要重置密码的用户
 * @returns 返回一个Promise，resolve时代表密码重置成功
 */
function resetPassword(id: number): Promise<void> {
    // 发送POST请求重置指定用户的密码
    return alova.Post<void>(`${baseURL}/password/reset`, { id })
}

/**
 * 设置用户权限
 * @param data 包含权限ID的对象
 * @param data.authorityId 要设置的权限ID
 * @returns 返回一个Promise，resolve时代表用户权限设置成功
 */
function setUserAuthority(data: { authorityId: number }): Promise<void> {
    return alova.Patch<void>(`${baseURL}/me/authority`, data)
}

export {
    login,
    logout,
    changePassword,
    getSelfInfo,
    listSysUser,
    registerSysUser,
    deleteSysUser,
    setUserInfo,
    resetPassword,
    setUserAuthority
}
