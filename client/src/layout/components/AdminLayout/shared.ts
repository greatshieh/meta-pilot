import type { AdminLayoutProps, LayoutCssVars, LayoutCssVarsProps } from '../../types'

/** The id of the scroll element of the layout */
export const LAYOUT_SCROLL_EL_ID = '__SCROLL_EL_ID__'

/** The max z-index of the layout */
export const LAYOUT_MAX_Z_INDEX = 100

/**
 * Create layout css vars by css vars props
 *
 * @param props Css vars props
 */
function createLayoutCssVarsByCssVarsProps(props: LayoutCssVarsProps) {
    const cssVars: LayoutCssVars = {
        '--layout-header-height': `${props.headerHeight}px`,
        '--layout-header-z-index': props.headerZIndex,
        '--layout-tab-height': `${props.tabHeight}px`,
        '--layout-tab-z-index': props.tabZIndex,
        '--layout-sider-width': `${props.siderWidth}px`,
        '--layout-sider-collapsed-width': `${props.siderCollapsedWidth}px`,
        '--layout-sider-z-index': props.siderZIndex,
        '--layout-mobile-sider-z-index': props.mobileSiderZIndex,
        '--layout-footer-height': `${props.footerHeight}px`,
        '--layout-footer-z-index': props.footerZIndex
    }

    return cssVars
}
/**
 * Create layout css vars
 *
 * @param props
 */
export function createLayoutCssVars(props: AdminLayoutProps) {
    const { mode, isMobile, maxZIndex = LAYOUT_MAX_Z_INDEX, headerHeight, tabHeight, siderWidth, siderCollapsedWidth, footerHeight } = props

    const headerZIndex = maxZIndex - 3
    const tabZIndex = maxZIndex - 5
    const siderZIndex = mode === 'vertical' || isMobile ? maxZIndex - 1 : maxZIndex - 4
    const mobileSiderZIndex = isMobile ? maxZIndex - 2 : 0
    const footerZIndex = maxZIndex - 5

    const cssProps: LayoutCssVarsProps = {
        headerHeight,
        headerZIndex,
        tabHeight,
        tabZIndex,
        siderWidth,
        siderZIndex,
        mobileSiderZIndex,
        siderCollapsedWidth,
        footerHeight,
        footerZIndex
    }

    return createLayoutCssVarsByCssVarsProps(cssProps)
}
