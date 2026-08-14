declare global {
    /**
     * 响应数据
     */
    interface ResponseData<T = unknown> {
        code: number
        data: T
        message: string
    }

    /**
     * 分页查询参数
     */
    interface PageQuery {
        page: number // 页码
        pageSize: number // 每页大小
        sortField?: string // 排序字段
        sortOrder?: string // 排序方式
    }

    /**
     * 分页响应对象
     */
    interface PageResult<T> {
        /** 数据列表 */
        list: T
        /** 总数 */
        total: number
        /** 页码 */
        page: number
        /** 每页数量 */
        pageSize: number
    }
}

export {}
