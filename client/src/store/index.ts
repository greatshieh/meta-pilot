import { createPinia } from 'pinia'
import type { App } from 'vue'
import { resetSetupStore } from './plugins'

export const setupStore = (app: App) => {
    const store = createPinia()
    store.use(resetSetupStore)

    app.use(store)
}
