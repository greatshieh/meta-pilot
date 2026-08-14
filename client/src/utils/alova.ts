import { createAlova } from 'alova'
import adapterFetch from 'alova/fetch'
import VueHook from 'alova/vue'
import useAdminStore from '@/store/modules/admin'
import { localStg } from './storage'

const getAuthorization = (): string => {
    const adminStore = useAdminStore()
    return adminStore.getToken
}

const alova = createAlova({
    // 基础URL，根据环境变量和代理配置动态获取
    baseURL: import.meta.env.VITE_API_PREFIX,
    // 使用fetch适配器
    requestAdapter: adapterFetch(),
    // 超时时间
    timeout: 10 * 1000,
    statesHook: VueHook,
    // 请求前置拦截器，在发送请求前执行
    beforeRequest(method) {
        if (method.url?.includes('huaweicloud')) {
            method.baseURL = ''
        } else {
            // 从store中获取当前用户的授权token
            const Authorization = getAuthorization()
            // 将token添加到请求头的X-Token字段中
            method.config.headers['X-Token'] = Authorization
        }
    },
    responded: {
        onSuccess: async response => {
            let message: string = ''

            const newToken = response.headers.get('new-token')

            if (newToken) {
                const adminStore = useAdminStore()
                adminStore.userInfo.token = newToken
                localStg.set('token', newToken)
            }

            if (response.status >= 400) {
                message = `${response.status}: ${response.statusText}`
            } else {
                if (response.url?.includes('huaweicloud')) {
                    return response
                }

                const resp = (await response.json()) as ResponseData

                if (String(resp.code) === String(import.meta.env.VITE_SERVICE_SUCCESS_CODE)) {
                    if (resp.message.toLocaleLowerCase() !== 'success') {
                        ElMessage.success(resp.message)
                    }

                    return resp.data
                } else if (['100206', '110001'].includes(String(resp.code))) {
                    ElMessage({
                        type: 'error',
                        message: '请重新登录',
                        duration: 5 * 1000
                    })
                    const adminStore = useAdminStore()
                    adminStore.removeToken()
                    location.reload()
                } else {
                    message = resp.message
                }
            }

            if (message) {
                ElMessage({
                    type: 'error',
                    message,
                    duration: 5 * 1000
                })
                return Promise.reject(message)
            }
        },
        onError: error => {
            ElMessage({
                type: 'error',
                message: error,
                duration: 5 * 1000
            })
        }
    }
})

export default alova
