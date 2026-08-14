declare namespace Api.SSE {
    interface Response {
        grade: string
        class: string
        exam: Record<string, string>[]
    }
}
