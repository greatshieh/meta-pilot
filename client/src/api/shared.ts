import useAdminStore from '@/store/modules/admin'

export function getAuthorization() {
    const adminStore = useAdminStore()
    return adminStore.userInfo.token
}
