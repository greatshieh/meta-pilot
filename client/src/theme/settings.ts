/** Default theme settings */
export const themeSettings: App.Theme.ThemeSetting = {
    resetCacheStrategy: 'close',
    layout: {
        mode: 'vertical',
        scrollMode: 'content',
        reverseHorizontalMix: false
    },
    page: {
        animate: true,
        animateMode: 'fade-scale'
    },
    header: {
        height: 56,
        breadcrumb: {
            visible: true,
            showIcon: true
        },
        multilingual: {
            visible: true
        }
    },
    tab: {
        visible: true,
        cache: true,
        height: 44,
        mode: 'chrome'
    },
    fixedHeaderAndTab: true,
    sider: {
        inverted: false,
        width: 220,
        collapsedWidth: 64,
        mixWidth: 90,
        mixCollapsedWidth: 64,
        mixChildMenuWidth: 200
    },
    footer: {
        visible: true,
        fixed: false,
        height: 48,
        right: true
    },
    watermark: {
        visible: false,
        text: '成绩管理系统'
    }
}

/**
 * Override theme settings
 *
 * If publish new version, use `overrideThemeSettings` to override certain theme settings
 */
export const overrideThemeSettings: Partial<App.Theme.ThemeSetting> = {}
