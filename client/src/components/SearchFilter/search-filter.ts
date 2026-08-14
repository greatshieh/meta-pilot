export interface paramsProps {
    grade?: string
    class?: string
    name?: string
    examName?: string
    subject?: string
}

export interface SearchFilterProps {
    /**
     * @description 显示的选项
     */
    showItem?: string
    /**
     * form表单参数
     */
    queryParam?: paramsProps
}
