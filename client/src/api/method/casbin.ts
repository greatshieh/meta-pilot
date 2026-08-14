import alova from '@/utils/alova'

const baseURL = '/policies/authorities'

function setPolicy(authorityId: number, policy: Api.Cabin.PolicyPathRequest) {
    return alova.Post(`${baseURL}/${authorityId}`, policy)
}

function getPolicyPathByAuthorityId(authorityId: number): Promise<Api.Cabin.PolicyPathResponse> {
    return alova.Get<Api.Cabin.PolicyPathResponse>(`${baseURL}/${authorityId}`, { cacheFor: 0 })
}

export { setPolicy, getPolicyPathByAuthorityId }
