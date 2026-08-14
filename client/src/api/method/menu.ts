import alova from '@/utils/alova'

const baseURL = '/menus'

/**
 * 获取异步菜单数据
 * @returns 返回一个Promise，resolve时返回菜单元数据数组
 */
function getAsyncMenu(): Promise<Api.Menu.MenuMeta[]> {
    // 发送GET请求获取异步菜单数据，禁用缓存
    return alova.Get<Api.Menu.MenuMeta[]>(`${baseURL}/authorized`, { cacheFor: 0 })
}

/**
 * 获取菜单列表数据
 * 发送GET请求到'/menus'端点获取菜单元信息数组
 * @returns Promise<Api.Menu.MenuMeta[]> 返回菜单元信息数组的Promise对象
 */
function getMenuList(): Promise<Api.Menu.MenuMeta[]> {
    // 使用alova实例发送GET请求获取菜单列表
    // 请求地址: /menus
    // 返回类型: Api.Menu.MenuMeta[] (菜单元信息数组)
    return alova.Get<Api.Menu.MenuMeta[]>(baseURL)
}

/**
 * 更新菜单信息
 * @param data 菜单元数据对象，包含需要更新的菜单信息
 * @returns 返回一个Promise，resolve时返回更新操作的结果字符串
 */
function updateMenu(data: Api.Menu.MenuMeta): Promise<string> {
    return alova.Put<string>(baseURL, data)
}

/**
 * 添加新菜单
 * @param data 菜单元数据对象，包含菜单的基本信息
 * @returns 返回一个Promise，resolve时代表菜单添加成功
 */
function addMenu(data: Api.Menu.MenuMeta) {
    // 发送POST请求添加新菜单
    return alova.Post(baseURL, data)
}

function deleteMenu(menuId: number): Promise<string> {
    return alova.Delete<string>(`${baseURL}/${menuId}`)
}

export { getAsyncMenu, getMenuList, updateMenu, addMenu, deleteMenu }
