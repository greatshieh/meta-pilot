import { defineStore } from 'pinia'
import type { RouteLocationNormalizedLoaded, RouteRecordRaw } from 'vue-router'

import router, { constantRoutes as staticRoutes } from '@/router'
import { SetupStoreId } from '../../enum'
import useAdminStore from '../admin'
import { computed, nextTick, ref, shallowRef } from 'vue'
import useTabStore from '../tab'
import { menuApi } from '@/api'

const views = import.meta.glob(['@/views/**/*.vue', '@/layout/index.vue'])

function getSelectedMenuKeyPathByKey(selectedKey: string, menus: App.Global.Menu[]) {
    const keyPath: string[] = []

    menus.some(menu => {
        const path = findMenuPath(selectedKey, menu)

        const find = Boolean(path?.length)

        if (find) {
            keyPath.push(...path!)
        }

        return find
    })

    return keyPath
}

/**
 * Find menu path
 *
 * @param targetKey Target menu key
 * @param menu Menu
 */
function findMenuPath(targetKey: string, menu: App.Global.Menu): string[] | null {
    const path: string[] = []

    function dfs(item: App.Global.Menu): boolean {
        path.push(item.name as string)

        if (item.name === targetKey) {
            return true
        }

        if (item.children) {
            for (const child of item.children) {
                if (dfs(child)) {
                    return true
                }
            }
        }

        path.pop()

        return false
    }

    if (dfs(menu)) {
        return path
    }

    return null
}

// const getRouteName = (path: string, routes: RouteRecordRaw[]): string | null => {
//     for (let i = 0; i < routes.length; i++) {
//         if (routes[i].path === path) {
//             return routes[i].name as string
//         }

//         if (routes[i].children) {
//             const childName = getRouteName(path, routes[i].children!)
//             return childName
//         }
//     }

//     return null
// }

/**
 * Get global menu by route
 *
 * @param route
 */
const getGlobalMenuByBaseRoute = (route: RouteLocationNormalizedLoaded | RouteRecordRaw) => {
    // const { SvgIconVNode } = useSvgIcon()

    const { name, path } = route
    const { title, icon, affix } = route.meta ?? {}

    const menu: App.Global.Menu = {
        name: name as string,
        title: title || '',
        routeKey: name as string,
        routePath: path,
        icon,
        affix
    }

    return menu
}

const getGlobalMenusByAuthRoutes = (routes: RouteRecordRaw[]) => {
    const menus: App.Global.Menu[] = []
    routes.forEach(route => {
        if (!route.meta?.hidden) {
            const menu = getGlobalMenuByBaseRoute(route)

            if (route.children?.some(child => !child.meta?.hide)) {
                menu.children = getGlobalMenusByAuthRoutes(route.children)
            }

            menus.push(menu)
        }
    })

    return menus
}

const getCacheRouteNames = (routes: RouteRecordRaw[]) => {
    const cacheNames: string[] = []

    routes.forEach(route => {
        // only get last two level route, which has component
        route.children?.forEach(child => {
            if (child.component && child.meta?.keepAlive) {
                cacheNames.push(child.name as string)
            }
        })
    })

    return cacheNames
}

function transformMenuToSearchMenus(menus: App.Global.Menu[], treeMap: App.Global.Menu[] = []) {
    if (menus && menus.length === 0) return []
    return menus.reduce((acc, cur) => {
        if (!cur.children) {
            acc.push(cur)
        }
        if (cur.children && cur.children.length > 0) {
            transformMenuToSearchMenus(cur.children, treeMap)
        }
        return acc
    }, treeMap)
}

function transformMenuToBreadcrumb(menu: App.Global.Menu) {
    const { children, ...rest } = menu

    const breadcrumb: App.Global.Breadcrumb = {
        ...rest
    }

    if (children?.length) {
        breadcrumb.options = children.map(transformMenuToBreadcrumb)
    }

    return breadcrumb
}

function getBreadcrumbsByRoute(route: RouteLocationNormalizedLoaded, menus: App.Global.Menu[]): App.Global.Breadcrumb[] {
    const name = route.name as string
    const activeKey = route.meta?.activeMenu

    for (const menu of menus) {
        if (menu.name === name) {
            return [transformMenuToBreadcrumb(menu)]
        }

        if (menu.name === activeKey) {
            const ROUTE_DEGREE_SPLITTER = '_'

            const parentKey = name.split(ROUTE_DEGREE_SPLITTER).slice(0, -1).join(ROUTE_DEGREE_SPLITTER)

            const breadcrumbMenu = getGlobalMenuByBaseRoute(route)
            if (parentKey !== activeKey) {
                return [transformMenuToBreadcrumb(breadcrumbMenu)]
            }

            return [transformMenuToBreadcrumb(menu), transformMenuToBreadcrumb(breadcrumbMenu)]
        }

        if (menu.children?.length) {
            const result = getBreadcrumbsByRoute(route, menu.children)
            if (result.length > 0) {
                return [transformMenuToBreadcrumb(menu), ...result]
            }
        }
    }

    return []
}

const getRouteName = (path: string, routes: RouteRecordRaw[]): string | null => {
    for (let i = 0; i < routes.length; i++) {
        if (routes[i].path === path) {
            return routes[i].name as string
        }

        if (routes[i].children) {
            const childName = getRouteName(path, routes[i].children!)
            return childName
        }
    }

    return null
}

const useRouteStore = defineStore(SetupStoreId.Route, () => {
    const tabStore = useTabStore()
    const isInitConstantRoute = ref(false)
    const isInitAuthRoute = ref(false)

    /** Home route key */
    const routeHome = ref(import.meta.env.VITE_ROUTE_HOME)

    /**
     * Set route home
     *
     * @param string Route key
     */
    function setRouteHome(string: string) {
        routeHome.value = string
    }

    /** constant routes */
    const constantRoutes = shallowRef<RouteRecordRaw[]>([])

    // function addConstantRoutes(routes: RouteRecordRaw[]) {
    //     const constantRoutesMap = new Map<string, RouteRecordRaw>([])

    //     routes.forEach(route => {
    //         constantRoutesMap.set(route.name as string, route)
    //     })

    //     constantRoutes.value = Array.from(constantRoutesMap.values())
    // }

    /** auth routes */
    const authRoutes = shallowRef<RouteRecordRaw[]>([])

    function addAuthRoutes(routes: Api.Menu.MenuMeta[]) {
        // 使用栈来模拟递归过程，避免深层递归可能导致的问题
        const stack: { routes: Api.Menu.MenuMeta[]; parentRoute?: RouteRecordRaw }[] = [{ routes }]
        const processedRoutes: RouteRecordRaw[] = []

        while (stack.length > 0) {
            const current = stack.pop()!
            const { routes, parentRoute } = current

            routes.forEach(item => {
                // 只保留 RouteRecordRaw 需要的属性
                const { path, name, component, meta, children, redirect } = item
                // 构造路由对象
                const route: RouteRecordRaw = {
                    path,
                    name,
                    redirect,
                    component: views[`/src/${component}`],
                    meta,
                    children: [] // 初始化为空数组，后续处理子路由
                }

                // 如果有子路由，将其加入处理栈
                if (children && children.length > 0) {
                    stack.push({ routes: children, parentRoute: route })
                }

                // 将处理好的路由添加到结果数组或父路由的children中
                if (parentRoute) {
                    // 如果有父路由，将当前路由添加到父路由的children中
                    if (!parentRoute.children) {
                        parentRoute.children = []
                    }
                    parentRoute.children.push(route)
                } else {
                    // 如果没有父路由，说明是顶级路由，直接添加到结果数组
                    processedRoutes.push(route)
                }
            })
        }

        // 更新authRoutes的值
        authRoutes.value = processedRoutes
    }

    const removeRouteFns: (() => void)[] = []

    /** Global menus */
    const menus = ref<App.Global.Menu[]>([])
    const searchMenus = computed(() => transformMenuToSearchMenus(menus.value))

    /** Get global menus */
    function getGlobalMenus(routes: RouteRecordRaw[]) {
        const globalMenu = getGlobalMenusByAuthRoutes(routes)
        if (globalMenu.length === 1) {
            menus.value = globalMenu[0].children || []
        } else {
            menus.value = globalMenu
        }
    }

    /** Cache routes */
    const cacheRoutes = ref<string[]>([])

    /**
     * Exclude cache routes
     *
     * for reset route cache
     */
    const excludeCacheRoutes = ref<string[]>([])

    /**
     * Get cache routes
     *
     * @param routes Vue routes
     */
    function getCacheRoutes(routes: RouteRecordRaw[]) {
        cacheRoutes.value = getCacheRouteNames(routes)
    }

    /**
     * Reset route cache
     *
     * @default router.currentRoute.value.name current route name
     * @param string
     */
    async function resetRouteCache(string?: string) {
        const routeName = string || (router.currentRoute.value.name as string)

        excludeCacheRoutes.value.push(routeName)

        await nextTick()

        excludeCacheRoutes.value = []
    }

    /** Global breadcrumbs */
    const breadcrumbs = computed(() => getBreadcrumbsByRoute(router.currentRoute.value, menus.value))

    /** Reset store */
    async function resetStore() {
        const routeStore = useRouteStore()

        routeStore.$reset()

        resetVueRoutes()

        // after reset store, need to re-init constant route
        await initConstantRoute()
    }

    /** Reset vue routes */
    function resetVueRoutes() {
        removeRouteFns.forEach(fn => fn())
        removeRouteFns.length = 0
    }

    /** init constant route */
    async function initConstantRoute() {
        if (isInitConstantRoute.value) return

        // if fetch constant routes failed, use static constant routes
        constantRoutes.value = staticRoutes

        handleConstantAndAuthRoutes()

        isInitConstantRoute.value = true

        tabStore.initHomeTab()
    }

    /** Init auth route */
    async function initAuthRoute() {
        const adminStore = useAdminStore()
        if (!adminStore.userInfo.userName) {
            await adminStore.getInfo()
        }
        await initDynamicAuthRoute()

        tabStore.initHomeTab()
    }

    /** Init dynamic auth route */
    async function initDynamicAuthRoute() {
        const anyncRoute = await menuApi.getAsyncMenu()

        const baseRouter: Api.Menu.MenuMeta[] = [
            {
                id: 0,
                path: '/',
                name: 'layout',
                component: 'layout/index.vue',
                meta: {
                    title: '底层layout'
                },
                children: []
            }
        ]

        baseRouter[0].children = anyncRoute

        addAuthRoutes(baseRouter)

        handleConstantAndAuthRoutes()

        setRouteHome('dashboard')

        // handleUpdateRootRouteRedirect('home')

        isInitAuthRoute.value = true
    }

    /** handle constant and auth routes */
    function handleConstantAndAuthRoutes() {
        const allRoutes = [...constantRoutes.value, ...authRoutes.value]

        resetVueRoutes()

        addRoutesToVueRouter(allRoutes)

        getGlobalMenus(allRoutes)

        getCacheRoutes(allRoutes)
    }

    /**
     * Add routes to vue router
     *
     * @param routes Vue routes
     */
    function addRoutesToVueRouter(routes: RouteRecordRaw[]) {
        routes.forEach(route => {
            const removeFn = router.addRoute(route)
            addRemoveRouteFn(removeFn)
        })
    }

    /**
     * Add remove route fn
     *
     * @param fn
     */
    function addRemoveRouteFn(fn: () => void) {
        removeRouteFns.push(fn)
    }

    /**
     * Update root route redirect when auth route mode is dynamic
     *
     * @param redirectKey Redirect route key
     */
    // function handleUpdateRootRouteRedirect(redirectKey: string) {
    //     const redirect = getRoutePath(redirectKey)

    //     if (redirect) {
    //         const rootRoute: CustomRoute = { ...ROOT_ROUTE, redirect }

    //         router.removeRoute(rootRoute.name)

    //         const [rootVueRoute] = getAuthVueRoutes([rootRoute])

    //         router.addRoute(rootVueRoute)
    //     }
    // }

    /**
     * Get is auth route exist
     *
     * @param routePath Route path
     */
    async function getIsAuthRouteExist(routePath: string) {
        const routeName = getRouteName(routePath, authRoutes.value)

        if (!routeName) {
            return false
        }

        return true
    }

    /**
     * Get selected menu key path
     *
     * @param selectedKey Selected menu key
     */
    function getSelectedMenuKeyPath(selectedKey: string) {
        return getSelectedMenuKeyPathByKey(selectedKey, menus.value)
    }

    async function onRouteSwitchWhenLoggedIn() {
        const adminStore = useAdminStore()
        await adminStore.getInfo()
    }

    async function onRouteSwitchWhenNotLoggedIn() {
        // some global init logic if it does not need to be logged in
    }

    return {
        resetStore,
        routeHome,
        menus,
        searchMenus,
        cacheRoutes,
        excludeCacheRoutes,
        resetRouteCache,
        breadcrumbs,
        initConstantRoute,
        isInitConstantRoute,
        initAuthRoute,
        isInitAuthRoute,
        initDynamicAuthRoute,
        getIsAuthRouteExist,
        getSelectedMenuKeyPath,
        onRouteSwitchWhenLoggedIn,
        onRouteSwitchWhenNotLoggedIn
    }
})

export default useRouteStore
