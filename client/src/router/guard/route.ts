import type { RouteLocationNormalized, Router } from 'vue-router'
import useAdminStore from '@/store/modules/admin'
import useRouteStore from '@/store/modules/route'

/**
 * 初始化路由配置
 * 功能：处理路由初始化逻辑，包括常量路由和权限路由的初始化，以及登录状态验证
 * @param to - 目标路由对象
 * @returns 返回重定向路径或null（null表示允许继续导航）
 */
async function initRoute(to: RouteLocationNormalized): Promise<string | null> {
    // 获取管理员和路由状态存储
    const adminStore = useAdminStore()
    const routeStore = useRouteStore()
    // 判断当前是否为404路由
    const isNotFoundRoute = to.name === 'not-found'

    // 1. 处理常量路由未初始化的情况
    if (!routeStore.isInitConstantRoute) {
        routeStore.initConstantRoute() // 初始化常量路由

        // 返回原始路径用于重定向（解决因路由未初始化导致的404问题）
        return to.fullPath
    }

    // 检查登录状态
    const isLogin = Boolean(adminStore.getToken)

    // 2. 处理未登录情况
    if (!isLogin) {
        // 如果不是404路由，返回null（由外部处理重定向）
        if (!isNotFoundRoute) return null

        // 处理带重定向参数的登录跳转
        if (to.query?.redirect) {
            return `/login?redirect=${to.query?.redirect}`
        }
        return '/login' // 默认登录页
    }

    // 3. 处理权限路由未初始化情况
    if (!routeStore.isInitAuthRoute) {
        await routeStore.initAuthRoute() // 异步初始化权限路由

        // 如果是404路由，返回原始路径重试
        if (isNotFoundRoute) return to.fullPath
    }

    // 触发登录状态下的路由切换回调
    // 4. 正常路由放行
    if (!isNotFoundRoute) return null

    // 5. 处理404路由的特殊情况
    const exist = await routeStore.getIsAuthRouteExist(to.path)
    return exist ? '/404' : null // 路由存在但无权限跳403，否则继续404
}

const handleRouteSwitch = (to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    // route with href
    if (to.meta.href) {
        window.open(to.meta.href, '_blank')
        return {
            path: from.fullPath,
            replace: true,
            query: from.query,
            hash: to.hash
        }
    }

    return true
}

export const createRouteGuard = (router: Router) => {
    router.beforeEach(async (to, _from) => {
        const location = await initRoute(to)

        if (location) return location

        const adminStore = useAdminStore()

        const loginRoute = 'login'
        const rootRoute = 'dashboard'
        const isLogin = Boolean(adminStore.getToken)
        const needLogin = Boolean(to.meta.roles)

        if (to.name === loginRoute && isLogin) {
            return { name: rootRoute }
        }

        if (!needLogin) {
            const location = handleRouteSwitch(to, _from)
            return location
        }

        if (!isLogin) {
            return {
                name: 'login',
                query: { redirect: to.query.redirect }
            }
        }

        return handleRouteSwitch(to, _from)
    })
}
