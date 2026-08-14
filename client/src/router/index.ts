import type { App } from 'vue'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { createRouterGuard } from './guard'

const Layout = () => import('@/layout/index.vue')
/**
 * constantRoutes
 * 所有权限都能访问的页面
 */
export const constantRoutes: RouteRecordRaw[] = [
    {
        path: '/redirect',
        component: Layout,
        meta: { hidden: true },
        children: [
            {
                path: '/redirect/:path(.*)',
                component: () => import('@/views/redirect/index.vue')
            }
        ]
    },
    {
        name: 'not-found',
        path: '/:pathMatch(.*)*',
        component: () => import('@/views/exception/404.vue'),
        meta: { hidden: true }
    },
    {
        name: 'no-auth',
        path: '/403',
        component: () => import('@/views/exception/403.vue'),
        meta: { hidden: true }
    },
    {
        path: '/',
        redirect: '/dashboard',
        meta: { hidden: true }
    },
    {
        path: '/login',
        name: 'login',
        component: () => import('@/views/login/index.vue'),
        meta: { hidden: true }
    }
]

/**
 * 创建路由
 */
export const router = createRouter({
    history: createWebHistory(import.meta.env.VITE_BASE_URL),
    routes: constantRoutes,
    scrollBehavior: () => {
        return {
            top: 0,
            left: 0
        }
    }
})

export async function setupRouter(app: App) {
    app.use(router)
    createRouterGuard(router)
    await router.isReady()
}

export default router
