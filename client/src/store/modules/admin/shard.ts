import { localStg } from '@/utils/storage'

export function getLocalToken() {
    return localStg.get('token') || ''
}

export function removeLocalToken() {
    localStg.remove('token')
}
