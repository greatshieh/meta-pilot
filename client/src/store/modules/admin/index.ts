import { defineStore } from 'pinia'
import { localStg } from '@/utils/storage'
import { SetupStoreId } from '../../enum'

import { getLocalToken, removeLocalToken } from './shard'
import useRouteStore from '../route'
import useTabStore from '../tab'
import { sysUserApi } from '@/api'
import { computed, ref } from 'vue'

const useAdminStore = defineStore(SetupStoreId.Admin, () => {
    const userInfo = ref<App.SysUser.UserInfo>({
        token: '',
        userName: '',
        nickName: '',
        avatar: '',
        phone: '',
        email: '',
        authority: { authorityId: 0, authorityName: '' },
        authorities: []
    })

    const getToken = computed(() => {
        return userInfo.value.token || getLocalToken()
    })

    // admin login
    function login(adminInfo: Api.Admin.LoginReq) {
        return new Promise<void>((resolve, reject) => {
            sysUserApi
                .login(adminInfo)
                .then(resp => {
                    userInfo.value.token = resp.token
                    localStg.set('token', resp.token)
                    resolve()
                })
                .catch(error => {
                    reject(error)
                })
        })
    }

    function getInfo() {
        return new Promise<Omit<App.SysUser.UserInfo, 'token' | 'userName'>>((resolve, reject) => {
            sysUserApi
                .getSelfInfo()
                .then(resp => {
                    if (resp) {
                        userInfo.value.userName = resp.userName
                        userInfo.value.nickName = resp.nickName
                        userInfo.value.avatar = resp.avatar
                        userInfo.value.email = resp.email
                        userInfo.value.phone = resp.phone
                        userInfo.value.authority = resp.authority
                        userInfo.value.authorities = resp.authorities || []
                    }
                    resolve(userInfo.value)
                })
                .catch(err => {
                    reject(err)
                })
        })
    }

    // admin logout
    function logout() {
        userInfo.value = {
            token: '',
            userName: '',
            nickName: '',
            avatar: '',
            phone: '',
            email: '',
            authority: { authorityId: 0, authorityName: '' },
            authorities: []
        }
        // 清除所有状态
        removeToken()

        const routeStore = useRouteStore()
        routeStore.resetStore()
        routeStore.isInitAuthRoute = false
        const tabStore = useTabStore()
        tabStore.clearTabs()
    }

    // 修改密码
    function changePassword(data: Api.SysUser.ChangePasswordReq) {
        return new Promise(resolve => {
            sysUserApi.changePassword(data).then(() => {
                resolve(0)
            })
        })
    }

    function removeToken() {
        removeLocalToken()
    }

    return {
        userInfo,
        login,
        getToken,
        getInfo,
        logout,
        changePassword,
        removeToken
    }
})

export default useAdminStore
