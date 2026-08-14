/** The global namespace for the app */
declare namespace App {
    namespace SysUser {
        interface userAuth {
            authorityId: number
            authorityName: string
        }
        interface UserInfo {
            token: string
            userName: string
            nickName: string
            avatar: string
            phone: string
            email: string
            authority: userAuth
            authorities: userAuth[]
        }
    }

    /** Admin namespace */
    namespace Admin {
        interface AdminState {
            token: string
            nickName: string
            userName: string
            avatar: string
        }
    }

    namespace User {
        interface UserState {
            name: string
            grade: string
            class: string
            code: string
            uid: string
            artType: string
            es: Record<string, string[]>
            token: string
        }
    }

    /** Theme namespace */
    namespace Theme {
        /** Theme setting */
        interface ThemeSetting {
            /** Theme scheme */
            themeScheme: UnionKey.ThemeScheme
            /** Reset cache strategy */
            resetCacheStrategy: UnionKey.ResetCacheStrategy
            /** Layout */
            layout: {
                /** Layout mode */
                mode: UnionKey.ThemeLayoutMode
                /** Scroll mode */
                scrollMode: UnionKey.ThemeScrollMode
                /**
                 * Whether to reverse the horizontal mix
                 *
                 * if true, the vertical child level menus in left and horizontal first level menus in top
                 */
                reverseHorizontalMix: boolean
            }
            /** Page */
            page: {
                /** Whether to show the page transition */
                animate: boolean
                /** Page animate mode */
                animateMode: UnionKey.ThemePageAnimateMode
            }
            /** Header */
            header: {
                /** Header height */
                height: number
                /** Header breadcrumb */
                breadcrumb: {
                    /** Whether to show the breadcrumb */
                    visible: boolean
                    /** Whether to show the breadcrumb icon */
                    showIcon: boolean
                }
                /** Multilingual */
                multilingual: {
                    /** Whether to show the multilingual */
                    visible: boolean
                }
            }
            /** Tab */
            tab: {
                /** Whether to show the tab */
                visible: boolean
                /**
                 * Whether to cache the tab
                 *
                 * If cache, the tabs will get from the local storage when the page is refreshed
                 */
                cache: boolean
                /** Tab height */
                height: number
                /** Tab mode */
                mode: UnionKey.ThemeTabMode
            }
            /** Fixed header and tab */
            fixedHeaderAndTab: boolean
            /** Sider */
            sider: {
                /** Inverted sider */
                inverted: boolean
                /** Sider width */
                width: number
                /** Collapsed sider width */
                collapsedWidth: number
                /** Sider width when the layout is 'vertical-mix' or 'horizontal-mix' */
                mixWidth: number
                /** Collapsed sider width when the layout is 'vertical-mix' or 'horizontal-mix' */
                mixCollapsedWidth: number
                /** Child menu width when the layout is 'vertical-mix' or 'horizontal-mix' */
                mixChildMenuWidth: number
            }
            /** Footer */
            footer: {
                /** Whether to show the footer */
                visible: boolean
                /** Whether fixed the footer */
                fixed: boolean
                /** Footer height */
                height: number
                /** Whether float the footer to the right when the layout is 'horizontal-mix' */
                right: boolean
            }
            /** Watermark */
            watermark: {
                /** Whether to show the watermark */
                visible: boolean
                /** Watermark text */
                text: string
            }
        }
    }

    /** Global namespace */
    namespace Global {
        type VNode = import('vue').VNode
        type RouteLocationNormalizedLoaded = import('vue-router').RouteLocationNormalizedLoaded
        type RouteQuery = import('vue-router').LocationQueryRaw

        /** The global header props */
        interface HeaderProps {
            /** Whether to show the logo */
            showLogo?: boolean
            /** Whether to show the menu toggler */
            showMenuToggler?: boolean
            /** Whether to show the menu */
            showMenu?: boolean
        }

        /** The global menu */
        type Menu = {
            /**
             * The menu key
             *
             * Equal to the route key
             */
            name: string
            /** 展示标题 */
            title: string
            /** The route key */
            routeKey: string
            /** The route path */
            routePath: string
            /** The menu icon */
            // icon?: () => VNode
            icon?: string
            affix?: boolean
            /** The menu children */
            children?: Menu[]
        }

        type Breadcrumb = Omit<Menu, 'children'> & {
            options?: Breadcrumb[]
        }

        /** Tab route */
        type TabRoute = Pick<RouteLocationNormalizedLoaded, 'name' | 'path' | 'meta'> &
            Partial<Pick<RouteLocationNormalizedLoaded, 'fullPath' | 'query' | 'matched'>>

        /** The global tab */
        type Tab = {
            /** 页签名称 */
            name: string
            /** 页签标题 */
            title: string
            /** 页签路由路径 */
            path: string
            /** 页签路由完整路径 */
            fullPath: string
            /** 页签图标 */
            icon?: string
            /** 是否固定页签 */
            affix?: boolean
            /** 是否开启缓存 */
            keepAlive?: boolean
            /** 路由查询参数 */
            query?: RouteQuery
        }

        /** Form rule */
        type FormRule = import('element-plus').FormItemRule

        /** The global dropdown key */
        type DropdownKey = 'closeCurrent' | 'closeOther' | 'closeLeft' | 'closeRight' | 'closeAll'
    }

    /** Service namespace */
    namespace Service {
        interface ServiceConfig {
            /** The backend service base url */
            baseURL: string
            /** The proxy pattern of the backend service base url */
            proxyPattern: string
        }

        /** The backend service response data */
        type Response<T = unknown> = {
            /** The backend service response code */
            code: string
            /** The backend service response message */
            msg: string
            /** The backend service response data */
            data: T
        }
    }
}
