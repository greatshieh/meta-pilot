declare namespace Api.Score {
    /** api/method/score.ts */

    interface UpdateScoreMetaData {
        name: string
        class: string
        score: number
        answer: string
    }

    interface UploadScoreReq {
        examName: string
        grade: string
        subject: string
        rows: UpdateScoreMetaData[]
    }

    interface UpdateScoreReq {
        examName: string
        detail: Record<string, number>
    }

    interface ExamScoreMetaData {
        uid: string
        name: string
        examUid: string
        examName: string
        grade: string
        class: string
        detail: Record<string, number>
        totalScore: number
    }

    interface CommonQuery extends PageQuery {
        grade: string
        class: string
        subject: string
        examName: string
    }

    interface SearchQuery extends Omit<CommonQuery, 'subject'> {
        name: string
    }

    interface ExamScoreRes extends CommonQuery {
        total: number
        list: ExamScoreMetaData[]
    }

    interface StatisticsReq {
        grade: string
        class: string
        subject: string
        examName: string
        interval: number
    }

    // 核心指标查询参数
    type MetricsQuery = Pick<StatisticsReq, 'grade' | 'examName'>

    interface MetricsicsRes {
        avgScore: number
        excellentRate: string
        maxScore: number
        maxScoreSource: string
        passRate: string
        studentCount: number
        totalCount: number
    }

    // 平均分请求
    type AvgScoreReq = Pick<StatisticsReq, 'grade' | 'class' | 'examName'>

    interface AvgScoreRes {
        avgScore: number
        class: string
        subject: string
    }

    interface PassRateResponse {
        passRate: string
        excellentRate: string
        class: string
        subject: string
    }

    interface ImprovementTrendResponse {
        class: string
        subject: string
        examName: string
        improvementRate: number
    }

    interface TopPerformersResponse {
        name: string
        class: string
        ranking: number
        examName: string
        examDate: string
        uid: string
        detail: string
        totalScore: number
        score: number
    }
}
