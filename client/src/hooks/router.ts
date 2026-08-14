import type { RouteLocationRaw } from 'vue-router'
import { useRouter } from 'vue-router'
import { router as globalRouter } from '@/router'

export interface RouterPushOptions {
    query?: Record<string, string>
    params?: Record<string, string>
}

/**
 * Router push
 *
 * Jump to the specified route, it can replace function router.push
 *
 * @param inSetup Whether is in vue script setup
 */
export function useRouterPush(inSetup = true) {
    const router = inSetup ? useRouter() : globalRouter
    const route = globalRouter.currentRoute

    const routerPush = router.push

    const routerBack = router.back

    async function routerPushByKey(key: string, options?: RouterPushOptions) {
        const { query, params } = options || {}
        const routeLocation: RouteLocationRaw = {
            name: key
        }

        if (Object.keys(query || {}).length) {
            routeLocation.query = query
        }

        if (Object.keys(params || {}).length) {
            routeLocation.params = params
        }

        return routerPush(routeLocation)
    }

    function routerPushByKeyWithMetaQuery(key: string) {
        const allRoutes = router.getRoutes()
        const currentRoute = allRoutes.find(item => item.path === key)
        if (!currentRoute) {
            return Promise.reject(new Error(`Route ${key} not found`))
        }
        const name = currentRoute.name as string

        const query: Record<string, string> = {}

        return routerPushByKey(name, { query })
    }

    async function toHome() {
        return routerPushByKey('root')
    }

    /**
     * Navigate to login page
     *
     * @param loginModule The login module
     * @param redirectUrl The redirect url, if not specified, it will be the current route fullPath
     */
    async function toLogin(loginModule?: UnionKey.LoginModule, redirectUrl?: string) {
        const module = loginModule || 'pwd-login'

        const options: RouterPushOptions = {
            params: {
                module
            }
        }

        const redirect = redirectUrl || route.value.fullPath

        options.query = {
            redirect
        }

        return routerPushByKey('login', options)
    }

    /**
     * Toggle login module
     *
     * @param module
     */
    async function toggleLoginModule(module: UnionKey.LoginModule) {
        const query = route.value.query as Record<string, string>

        return routerPushByKey('login', { query, params: { module } })
    }

    /**
     * Redirect from login
     *
     * @param [needRedirect=true] Whether to redirect after login. Default is `true`
     */
    async function redirectFromLogin(needRedirect = true) {
        const redirect = route.value.query?.redirect as string

        if (needRedirect && redirect) {
            await routerPush(redirect)
        } else {
            await toHome()
        }
    }

    return {
        routerPush,
        routerBack,
        routerPushByKey,
        routerPushByKeyWithMetaQuery,
        toLogin,
        toggleLoginModule,
        redirectFromLogin
    }
}
