import { computed, ref } from 'vue'
import { useEventListener } from '@vueuse/core'
import { defineStore } from 'pinia'
import { router } from '@/router'
import { localStg } from '@/utils/storage'

import {
    extractTabsByAllRoutes,
    findTabByRouteName,
    getAllTabs,
    getDefaultHomeTab,
    getAffixTabNames,
    getTabByRoute,
    isTabInTabs,
    filterTabsByNames
} from './shared'
import { SetupStoreId } from '@/store/enum'
import useThemeStore from '../theme'
import { useRouterPush } from '@/hooks/router'
import useRouteStore from '../route'

const useTabStore = defineStore(SetupStoreId.Tab, () => {
    const routeStore = useRouteStore()
    const themeStore = useThemeStore()
    const { routerPush } = useRouterPush(false)

    /** Tabs */
    const tabs = ref<App.Global.Tab[]>([])

    /** Get active tab */
    const homeTab = ref<App.Global.Tab>()

    /** Init home tab */
    function initHomeTab() {
        homeTab.value = getDefaultHomeTab(router, routeStore.routeHome)
    }

    /** Get all tabs */
    const allTabs = computed(() => getAllTabs(tabs.value, homeTab.value))

    /** Active tab id */
    const activeTabName = ref<string>('')

    /**
     * Set active tab id
     *
     * @param id Tab id
     */
    function setActiveTabName(id: string) {
        activeTabName.value = id
    }

    /**
     * Init tab store
     *
     * @param currentRoute Current route
     */
    function initTabStore(currentRoute: App.Global.TabRoute) {
        const storageTabs = localStg.get('globalTabs')

        if (themeStore.tab.cache && storageTabs) {
            const extractedTabs = extractTabsByAllRoutes(router, storageTabs)
            tabs.value = extractedTabs
        }

        addTab(currentRoute)
    }

    /**
     * Add tab
     *
     * @param route Tab route
     * @param active Whether to activate the added tab
     */
    function addTab(route: App.Global.TabRoute, active = true) {
        const tab = getTabByRoute(route)

        const isHomeTab = tab.name === homeTab.value?.name

        if (!isHomeTab && !isTabInTabs(tab.name, tabs.value)) {
            tabs.value.push(tab)
        }

        if (active) {
            setActiveTabName(tab.name)
        }
    }

    /**
     * Remove tab when remove active tab, switch to next tab or home tab
     *
     * @param tabName Tab name
     */
    async function removeTab(tabName: string) {
        const removeTabIndex = tabs.value.findIndex(tab => tab.name === tabName)
        if (removeTabIndex === -1) return

        const isRemoveActiveTab = activeTabName.value === tabName

        // if remove the last tab, then switch to the second last tab
        const nextTab = tabs.value[removeTabIndex + 1] || tabs.value[removeTabIndex - 1] || homeTab.value

        // remove tab
        tabs.value.splice(removeTabIndex, 1)

        // if current tab is removed, then switch to next tab
        if (isRemoveActiveTab && nextTab) {
            await switchRouteByTab(nextTab)
        }
    }

    /** remove active tab */
    async function removeActiveTab() {
        await removeTab(activeTabName.value)
    }

    /**
     * remove tab by route name
     *
     * @param routeName route name
     */
    async function removeTabByRouteName(routeName: string) {
        const tab = findTabByRouteName(routeName, tabs.value)
        if (!tab) return

        await removeTab(tab.name)
    }

    /**
     * Clear tabs
     *
     * @param excludes Exclude tab ids
     */
    async function clearTabs(excludes: string[] = []) {
        const remainTabNames = [...getAffixTabNames(tabs.value), ...excludes]

        // Identify tabs to be removed and collect their routeKeys if strategy is 'close'
        const tabsToRemove = tabs.value.filter(tab => !remainTabNames.includes(tab.name))
        const routeKeysToReset: string[] = []

        const removedTabsNames = tabsToRemove.map(tab => tab.name)

        // If no tabs are actually being removed based on excludes and fixed tabs, exit
        if (removedTabsNames.length === 0) {
            return
        }

        const isRemoveActiveTab = removedTabsNames.includes(activeTabName.value)
        // filterTabsByIds returns tabs NOT in removedTabsNames, so these are the tabs that will remain
        const updatedTabs = filterTabsByNames(removedTabsNames, tabs.value)

        function update() {
            tabs.value = updatedTabs
        }

        if (isRemoveActiveTab) {
            const activeTabCandidate = updatedTabs[updatedTabs.length - 1] || homeTab.value

            if (activeTabCandidate) {
                // Ensure there's a tab to switch to
                await switchRouteByTab(activeTabCandidate)
            }
        }
        // Update the tabs array regardless of switch success or if a candidate was found
        update()

        // After tabs are updated and route potentially switched, reset cache for removed tabs
        for (const routeKey of routeKeysToReset) {
            routeStore.resetRouteCache(routeKey)
        }
    }

    /**
     * Switch route by tab
     *
     * @param tab
     */
    async function switchRouteByTab(tab: App.Global.Tab) {
        const fail = await routerPush(tab.fullPath)
        if (!fail) {
            setActiveTabName(tab.name)
        }
    }

    /**
     * Clear left tabs
     *
     * @param tabId
     */
    async function clearLeftTabs(tabId: string) {
        const tabNames = tabs.value.map(tab => tab.name)
        const index = tabNames.indexOf(tabId)
        if (index === -1) return

        const excludes = tabNames.slice(index)
        await clearTabs(excludes)
    }

    /**
     * Clear right tabs
     *
     * @param tabId
     */
    async function clearRightTabs(tabId: string) {
        const isHomeTab = tabId === homeTab.value?.name
        if (isHomeTab) {
            clearTabs()
            return
        }

        const tabNames = tabs.value.map(tab => tab.name)
        const index = tabNames.indexOf(tabId)
        if (index === -1) return

        const excludes = tabNames.slice(0, index + 1)
        await clearTabs(excludes)
    }

    /**
     * Set new title of tab
     *
     * @default activeTabName
     * @param title New tab title
     * @param tabName Tab name
     */
    function setTabTitle(title: string, tabName?: string) {
        const name = tabName || activeTabName.value

        const tab = tabs.value.find(item => item.name === name)
        if (!tab) return

        tab.title = title
    }

    /**
     * Is tab retain
     *
     * @param tabName
     */
    function isTabRetain(tabName: string) {
        if (tabName === homeTab.value?.name) return true

        const fixedtabNames = getAffixTabNames(tabs.value)

        return fixedtabNames.includes(tabName)
    }

    /** Cache tabs */
    function cacheTabs() {
        if (!themeStore.tab.cache) return

        localStg.set('globalTabs', tabs.value)
    }

    // cache tabs when page is closed or refreshed
    useEventListener(window, 'beforeunload', () => {
        cacheTabs()
    })

    return {
        /** All tabs */
        tabs: allTabs,
        activeTabName,
        initHomeTab,
        initTabStore,
        addTab,
        removeTab,
        removeActiveTab,
        removeTabByRouteName,
        clearTabs,
        clearLeftTabs,
        clearRightTabs,
        setTabTitle,
        switchRouteByTab,
        isTabRetain,
        cacheTabs
    }
})

export default useTabStore
