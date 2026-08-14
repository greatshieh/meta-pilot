import { computed, ref, watch } from 'vue'
import { type RouteRecordRaw, useRoute, useRouter } from 'vue-router'
import useContext from '@/hooks/use-context'
import useRouteStore from '@/store/modules/route'

export const { setupStore: setupMixMenuContext, useStore: useMixMenuContext } = useContext('mix-menu', useMixMenu)

function useMixMenu() {
    const route = useRoute()
    const routeStore = useRouteStore()
    const { selectedKey } = useMenu()

    const activeFirstLevelMenuKey = ref('')

    function setActiveFirstLevelMenuKey(key: string) {
        activeFirstLevelMenuKey.value = key
    }

    function getActiveFirstLevelMenuKey() {
        const [firstLevelRouteName] = selectedKey.value.split('_')
        setActiveFirstLevelMenuKey(firstLevelRouteName)
    }

    const allMenus = computed<RouteRecordRaw[]>(() => {
        return routeStore.menus.map(menu => {
            return {
                name: menu.name,
                path: menu.routePath,
                meta: { affix: menu.affix, icon: menu.icon, title: menu.title },
                children: menu.children?.map(child => ({
                    name: child.name,
                    path: child.routePath,
                    meta: { affix: child.affix, icon: child.icon, title: child.title }
                }))
            } as RouteRecordRaw
        })
    })

    const childLevelMenus = computed<App.Global.Menu[]>(
        () => routeStore.menus.find(menu => menu.name === activeFirstLevelMenuKey.value)?.children || []
    )

    const isActiveFirstLevelMenuHasChildren = computed(() => {
        if (!activeFirstLevelMenuKey.value) {
            return false
        }

        const findItem = allMenus.value.find(item => item.name === activeFirstLevelMenuKey.value)

        return Boolean(findItem?.children?.length)
    })

    watch(
        () => route.name,
        () => {
            getActiveFirstLevelMenuKey()
        },
        { immediate: true }
    )

    return {
        allMenus,
        childLevelMenus,
        isActiveFirstLevelMenuHasChildren,
        activeFirstLevelMenuKey,
        setActiveFirstLevelMenuKey,
        getActiveFirstLevelMenuKey
    }
}

export function useMenu() {
    const route = useRoute()
    const router = useRouter()

    const selectedKey = computed(() => {
        const { hidden, activeName } = route.meta
        const name = route.path as string

        const routeName = (hidden ? router.getRoutes().find(item => item.name === activeName)?.path : name) || name
        return routeName
    })

    return {
        selectedKey
    }
}
