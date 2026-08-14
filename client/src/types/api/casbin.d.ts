declare namespace Api.Cabin {
    interface PolicyPath {
        path: string
        method: string
    }

    type PolicyPathResponse = PolicyPath[]

    interface PolicyPathRequest {
        casbinInfos: PolicyPath[]
    }
}
