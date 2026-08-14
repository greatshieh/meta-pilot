import type { Router } from 'vue-router'
import { createProgressGuard } from './progress'
import { createRouteGuard } from './route'
import { createDocumentTitleGuard } from './title'

/**
 * 创建路由守卫整合函数
 * 功能：集中注册所有路由守卫，包括：
 * 1. 路由切换进度条
 * 2. 路由权限验证
 * 3. 页面标题自动更新
 * @param router - Vue Router实例
 */
export function createRouterGuard(router: Router) {
    // 注册路由切换进度条守卫
    createProgressGuard(router)

    // 注册路由权限验证守卫（核心业务逻辑）
    createRouteGuard(router)

    // 注册页面标题自动更新守卫
    createDocumentTitleGuard(router)
}
