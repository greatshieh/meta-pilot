import type { Router } from 'vue-router'
import 'nprogress/nprogress.css' // 引入NProgress进度条样式
import NProgress from 'nprogress' // 引入NProgress库

/**
 * 创建路由进度条守卫
 * 功能：在路由切换时显示进度条，提升用户体验
 * @param router - Vue Router实例
 */
export function createProgressGuard(router: Router) {
    // 配置NProgress：禁用旋转动画
    NProgress.configure({ showSpinner: false })

    // 路由切换前回调：开始显示进度条
    router.beforeEach((_to, _from, next) => {
        NProgress.start() // 启动进度条动画
        next() // 继续路由导航
    })

    // 路由切换后回调：完成进度条
    router.afterEach(() => {
        NProgress.done() // 结束进度条动画
    })
}
