import alova from '@/utils/alova'

function fetchCategorySummary(): Promise<Api.Dashboard.CategorySummary> {
    return alova.Get<Api.Dashboard.CategorySummary>('/dashboard/summary/categories')
}

/** 仪表盘班级统计响应参数 */
function fetchGradeSubjectAverages(): Promise<Array<Api.Dashboard.RaddarRes>> {
    return alova.Get<Array<Api.Dashboard.RaddarRes>>('/dashboard/analytics/grade-subject-averages')
}

/** 仪表盘班级统计响应参数 */
function fetchGradeHistoricalAverages(): Promise<Array<Api.Dashboard.HistoricalAverages>> {
    return alova.Get<Array<Api.Dashboard.HistoricalAverages>>('/dashboard/analytics/grade-historical-averages')
}

/** 仪表盘班级统计响应参数 */
function fetchGradeSubjectEntryProgress(): Promise<Array<Api.Dashboard.SubjectEntryProgress>> {
    return alova.Get<Array<Api.Dashboard.SubjectEntryProgress>>('/dashboard/progress/grade-subject-entry')
}

export default {
    fetchCategorySummary,
    fetchGradeSubjectAverages,
    fetchGradeHistoricalAverages,
    fetchGradeSubjectEntryProgress
}
