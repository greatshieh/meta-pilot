/** The storage namespace */
declare namespace StorageType {
    interface Session {
        /**
         * the theme settings
         */
        themeSettings: App.Theme.ThemeSetting
    }

    interface Local {
        /** The token */
        token: string
        /** Fixed sider with mix-menu */
        mixSiderFixed: CommonType.YesOrNo
        /** The refresh token */
        refreshToken: string
        /** The dark mode */
        darkMode: boolean
        /** The theme settings */
        themeSettings: App.Theme.ThemeSetting
        /**
         * The override theme flags
         *
         * The value is the build time of the project
         */
        overrideThemeFlag: string
        /** The global tabs */
        globalTabs: App.Global.Tab[]
        /** The backup theme setting before is mobile */
        backupThemeSettingBeforeIsMobile: {
            layout: UnionKey.ThemeLayoutMode
            siderCollapse: boolean
        }
    }
}
