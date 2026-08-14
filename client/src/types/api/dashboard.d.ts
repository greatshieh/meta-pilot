declare namespace Api.Dashboard {
    /** api/method/dashBoard.ts */

    /** 仪表盘分类统计响应参数 */
    interface CategorySummary {
        examTotal: number
        paperTotal: number
        studentTotal: number
        teacherTotal: number
    }

    /** 仪表盘班级统计响应参数 */
    interface RaddarRes {
        grade: string
        subject: string
        avg: number
        max: number
        min: number
    }

    /** 仪表盘班级统计响应参数 */
    interface HistoricalAverages {
        avgScore: number
        examName: string
        grade: string
    }

    /** 仪表盘班级统计响应参数 */
    interface SubjectEntryProgress {
        examName: string
        grade: string
        total: number
        subjects: Record<string, number>
    }
}
