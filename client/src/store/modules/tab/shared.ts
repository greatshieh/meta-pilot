import type { Router } from 'vue-router'

/**
 * Get all tabs
 *
 * @param tabs Tabs
 * @param homeTab Home tab
 */
export function getAllTabs(tabs: App.Global.Tab[], homeTab?: App.Global.Tab) {
    if (!homeTab) {
        return []
    }

    const filterHomeTabs = tabs.filter(tab => tab.name !== homeTab.name)

    const fixedTabs = filterHomeTabs.filter(isAffixTab)

    const remainTabs = filterHomeTabs.filter(tab => !isAffixTab(tab))

    const allTabs = [homeTab, ...fixedTabs, ...remainTabs]

    return allTabs
}

/**
 * Is fixed tab
 *
 * @param tab
 */
function isAffixTab(tab: App.Global.Tab) {
    return tab.affix
}

/**
 * Get tab id by route
 *
 * @param route
 */
export function getTabNameByRoute(route: App.Global.TabRoute) {
    const { path, query = {}, meta } = route

    let id = path

    if (meta.multiTab) {
        const queryKeys = Object.keys(query).sort()
        const qs = queryKeys.map(key => `${key}=${query[key]}`).join('&')

        id = `${path}?${qs}`
    }

    return id
}

/**
 * Get tab by route
 *
 * @param route
 */
export function getTabByRoute(route: App.Global.TabRoute) {
    const { name, path, fullPath, meta } = route

    const { title, affix, icon } = meta

    const tab: App.Global.Tab = {
        name: name as string,
        title: title || '',
        path,
        fullPath: fullPath || path,
        icon,
        affix
    }

    return tab
}

/**
 * The vue router will automatically merge the meta of all matched items, and the icons here may be affected by other
 * matching items, so they need to be processed separately
 *
 * @param route
 */
export function getRouteIcons(route: App.Global.TabRoute) {
    // Set default value for icon at the beginning
    let icon: string = route?.meta?.icon || import.meta.env.VITE_MENU_ICON

    // Route.matched only appears when there are multiple matches,so check if route.matched exists
    if (route.matched) {
        // Find the meta of the current route from matched
        const currentRoute = route.matched.find(r => r.name === route.name)
        // If icon exists in currentRoute.meta, it will overwrite the default value
        icon = currentRoute?.meta?.icon || icon
    }

    return icon
}

/**
 * Get default home tab
 *
 * @param router
 * @param homeRouteName routeHome in useRouteStore
 */
export function getDefaultHomeTab(router: Router, homeRouteName: string) {
    let homeTab: App.Global.Tab = {
        title: '首页',
        name: 'home',
        path: '/dashboard',
        fullPath: '/dashboard'
    }

    const routes = router.getRoutes()
    const homeRoute = routes.find(route => route.name === homeRouteName)
    if (homeRoute) {
        homeTab = getTabByRoute(homeRoute)
    }

    return homeTab
}

/**
 * Is tab in tabs
 *
 * @param tab
 * @param tabs
 */
export function isTabInTabs(tabName: string, tabs: App.Global.Tab[]) {
    return tabs.some(tab => tab.name === tabName)
}

/**
 * Filter tabs by id
 *
 * @param tabName
 * @param tabs
 */
export function filterTabsByName(tabName: string, tabs: App.Global.Tab[]) {
    return tabs.filter(tab => tab.name !== tabName)
}

/**
 * Filter tabs by ids
 *
 * @param tabNames
 * @param tabs
 */
export function filterTabsByNames(tabNames: string[], tabs: App.Global.Tab[]) {
    return tabs.filter(tab => !tabNames.includes(tab.name))
}

/**
 * extract tabs by all routes
 *
 * @param router
 * @param tabs
 */
export function extractTabsByAllRoutes(router: Router, tabs: App.Global.Tab[]) {
    const routes = router.getRoutes()

    const routeNames = routes.map(route => route.name)

    return tabs.filter(tab => routeNames.includes(tab.name))
}

/**
 * Get fixed tabs
 *
 * @param tabs
 */
export function getAffixTabs(tabs: App.Global.Tab[]) {
    return tabs.filter(isAffixTab)
}

/**
 * Get fixed tab ids
 *
 * @param tabs
 */
export function getAffixTabNames(tabs: App.Global.Tab[]) {
    const fixedTabs = getAffixTabs(tabs)

    return fixedTabs.map(tab => tab.name)
}

/**
 * find tab by route name
 *
 * @param name
 * @param tabs
 */
export function findTabByRouteName(name: string, tabs: App.Global.Tab[]) {
    const tabName = name
    const multiTabName = `${name}?`

    return tabs.find(tab => tab.name === tabName || tab.name.startsWith(multiTabName))
}
