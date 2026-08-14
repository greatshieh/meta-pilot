// oss相关api

import type { RequestBody } from 'alova'
import alova from '@/utils/alova'

/**
 * @description: 从服务区获取obs临时授权
 * @param {string} method 需要请求obs服务器的方法
 * @param {string} key 文件名
 * @param {Record<string, string>} headers 需要请求obs服务器的请求头
 * @param {Record<string, string>} queryParams 需要请求obs服务器的查询参数
 * @param {string} fileType 需要请求obs服务器的文件类型，如image/jpeg，用于请求头的Content-Type
 * @return {ResponseData<{ActualSignedRequestHeaders: Record<string, string[]>, SignedUrl:string}>} 返回obs临时授权
 */
const requestOssAuthorization = (method: string, key?: string, headers?: Record<string, string>, queryParams?: Record<string, string>) => {
    return alova.Post<{
        ActualSignedRequestHeaders: Record<string, string[]>
        SignedUrl: string
    }>('/oss/token', { method, key, headers, queryParams })
}

/**
 * 初始化分段上传
 */
const ossApplyURL = (
    url: string,
    method: string,
    data: {
        data?: RequestBody
        headers?: Record<string, string[]>
    }
) => {
    if (method === 'post') {
        return alova.Post<Response>(url, data.data, {
            headers: Object.assign(
                {
                    'Content-Type': '',
                    Accept: 'application/json, text/plain, */*'
                },
                data.headers
            )
        })
    } else {
        return alova.Put<Response>(url, data.data, {
            headers: Object.assign({ Accept: 'application/json, text/plain, */*' }, data.headers),
            timeout: 5 * 60 * 1000
        })
    }
}

/**
 * 上传分片
 */
const uploadChunk = (url: string, data: RequestBody, headers: Record<string, string[]>) => {
    return alova.Put(url, data, {
        headers
    })
}

export default {
    requestOssAuthorization,
    ossApplyURL,
    uploadChunk
}
