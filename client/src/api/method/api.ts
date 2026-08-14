import alova from '@/utils/alova'

const baseURL = '/apis'

function getApiList(params: Api.ListReq): Promise<Api.ListRes> {
    return alova.Get<Api.ListRes>(baseURL, { params })
}

function addApi(api: Api.ApiMetaData): Promise<void> {
    return alova.Post('/apis', api)
}

function updateApi(api: Api.ApiMetaData): Promise<void> {
    return alova.Put(baseURL, api)
}

function deleteApi(id: number): Promise<void> {
    return alova.Delete(`${baseURL}/${id}`)
}

function getApiGroups(): Promise<Api.groupsRes> {
    return alova.Get<Api.groupsRes>(`${baseURL}/groups`)
}

/**
 * 获取API数据
 * @returns 返回一个Promise，resolve时包含新增，忽略和已删除的API数据
 */
function syncApi(): Promise<Api.syncRes> {
    // 发送GET请求到同步API接口，获取同步结果
    return alova.Get<Api.syncRes>(`${baseURL}/sync`)
}

function enterSyncApi(data: Api.syncRes): Promise<void> {
    return alova.Put(`${baseURL}/sync`, data)
}

function ignoreApi(api: Api.ApiMetaData & { flag: boolean }): Promise<void> {
    return alova.Post(`${baseURL}/ignore`, api)
}

function getAllApis(): Promise<Api.ApiMetaData[]> {
    return alova.Get<Api.ApiMetaData[]>(`${baseURL}/all`)
}

export { getApiList, addApi, updateApi, deleteApi, getApiGroups, syncApi, enterSyncApi, ignoreApi, getAllApis }
