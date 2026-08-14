declare namespace Api.Works {
    interface WorksQueryParams extends PageQuery {
        grade: string
        class: string
        name: string
        artType: string
        category: string
        academic: string
        isGood: ''
    }

    interface WorksMetaData {
        id: number
        grade: string
        class: string
        name: string
        picName: string
        category: string
        remark: string
        isGood: string
        examName: string
        artType: string
        academic: string
        uid: string
        url: string
        [label: string]: WorksMetaData[]
    }

    interface WorksTableResponse extends PageQuery {
        total: number
        list: WorksMetaData[]
    }
}
