import type { PageTabCssVars } from '../../types'

/** The active color of the tab */
export const ACTIVE_COLOR = '#1890ff'

export function createTabCssVars() {
    const cssProps: PageTabCssVars = {
        '--tab-active-text-color': '#fff',
        '--tab-active-bg-color': 'var(--el-color-primary)',
        '--tab-border-color': 'var(--el-color-primary-light-9)',
        '--tab-hover-text-color': 'var(--el-text-color-primary)',
        '--tab-hover-bg-color': 'var(--el-color-primary-light-7)'
    }

    return cssProps
}
