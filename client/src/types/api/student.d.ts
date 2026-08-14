declare namespace Api.Student {
    interface Info {
        ID?: number
        uid?: string
        grade: string
        class: string
        name: string
        phone: string
    }

    type ListReq = Api.Commons.RecordNullable<Pick<Student.Info, 'grade' | 'class' | 'name'> & Api.Commons.PaginationRequestQuery>

    type ListRes = Api.Common.PaginatingResponseRecord<Omit<Student.Info, 'ID' | 'uid'> & { ID: number; uid: string }>
}
