import { createApp } from 'vue'
import './plugins/assets'
import { setupStore } from '@/store'
import App from './App.vue'
import { setupIconifyOffline, setupLoading } from './plugins'
import { setupRouter } from './router'

/**
 * 异步初始化Vue应用
 * 功能：按顺序执行应用初始化流程，包括加载动画、插件设置、状态管理和路由配置
 */
async function setupApp() {
    // 1. 显示加载动画
    setupLoading()

    // 2. 配置Iconify图标库离线模式
    setupIconifyOffline()

    // 3. 创建Vue应用实例
    const app = createApp(App)

    // 4. 设置UI组件库
    //   setupUI(app)

    // 5. 初始化状态管理
    setupStore(app)

    // 6. 异步设置路由（等待路由准备完成）
    await setupRouter(app)

    // 7. 将应用挂载到DOM
    app.mount('#app')
}

// 启动应用初始化流程
setupApp()
