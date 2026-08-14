export interface RaddarChartProp {
    grade: string
    subjects: string[]
    series: Array<{
        type: 'radar'
        areaStyle: {
            opacity: number
        }
        z: number
        data: Array<{
            name: string
            value: number[]
        }>
    }>
}

export interface TrendDataProp {
    title: string
    data: number[]
    xAxis: string[]
}
