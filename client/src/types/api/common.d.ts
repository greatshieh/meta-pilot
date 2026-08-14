declare namespace Api.Commons {
    interface PaginatingCommonParams {
        // 页码
        page: number
        // 每页大小
        pageSize: number
    }

    interface PaginationRequestQuery extends PaginatingCommonParams {
        // 排序字段
        sortField?: string
        // 排序方式
        sortOrder?: string
    }

    interface PaginatingResponseRecord<T> extends PaginatingCommonParams {
        list: T[]
        total: number
    }

    /** add null to all properties */
    type RecordNullable<T> = {
        [K in keyof T]?: T[K] | undefined
    }

    interface respErr {
        grade: string
        class: string
        name: string
        err: string
    }
}
