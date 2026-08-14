declare namespace Api {
    interface ApiMetaData {
        id: number
        path: string
        description: string
        apiGroup: string
        required: boolean | undefined
        method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
    }

    type ListReq = Api.Commons.RecordNullable<ApiMetaData> & Api.Commons.PaginationRequestQuery

    type ListRes = Api.Commons.PaginatingResponseRecord<ApiMetaData>

    type groupsRes = string[]

    interface syncRes {
        newApis: ApiMetaData[]
        deleteApis: ApiMetaData[]
        ignoreApis: ApiMetaData[]
    }
}
