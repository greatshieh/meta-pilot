import { useTitle } from '@vueuse/core' // 引入VueUse的标题工具
import type { Router } from 'vue-router' // 引入Vue Router类型

/**
 * 创建文档标题守卫
 * 功能：在路由切换后自动更新页面标题
 * @param router - Vue Router实例
 */
export function createDocumentTitleGuard(router: Router) {
    // 路由切换后回调：更新页面标题
    router.afterEach(to => {
        // 使用路由元信息中的标题，若不存在则使用环境变量中的默认标题
        useTitle(to.meta.title || import.meta.env.VITE_APP_TITLE)
    })
}
