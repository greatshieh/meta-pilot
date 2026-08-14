declare namespace Api.User {
    interface UserLoginRequest {
        name: string
        phone: string
        code: string
    }

    interface UserLoginResponse {
        grade: string
        class: string
        name: string
        token: string
        uid: string
    }

    interface UserDetailInfo {
        name: string
        grade: string
        class: string
    }

    interface ReportData {
        classRanking: string
        gradeRanking: string
        id: number
        uid: string
        examName: string
        examUid: string
        grade: string
        class: string
        subject: string
        score: number
        detail: string[]
    }

    interface UserReportResponse {
        es: ReportData[]
        total: number
    }

    interface ExamInfo {
        uid: string
        name: string
        createdAt: string
    }

    interface UserExamListResponse {
        names: ExamInfo[]
    }

    interface WorksResponse {
        list: Api.Works.WorksMetaData[]
    }

    interface GoodWorkResponse {
        id: number
        url: string
        remark: string
        category: string
        isGood: string
    }

    interface GoodWorkListRequest extends PageQuery {
        grade: string
        class: string
        category: string
        academic: number
    }

    interface GoodWorkList extends PageQuery {
        list: GoodWorkResponse[]
        total: number
    }

    interface userState {
        name: string
        grade: string
        class: string
        code: string
        artType: string
        es: Record<string, string[]>
        token: string
    }
}
