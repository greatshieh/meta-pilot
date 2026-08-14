import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import type { PluginOption } from 'vite'
import progress from 'vite-plugin-progress'
import VueDevtools from 'vite-plugin-vue-devtools'
import { setupHtmlPlugin } from './html'
import Unocss from 'unocss/vite'
import { setupUnplugin } from './unplugin'
import { envParse } from 'vite-plugin-env-parse'
import { setupMockServer } from './mock'

export function setupVitePlugins(viteEnv: Env.ImportMeta, buildTime: string) {
    const plugins: PluginOption = [
        vue(),
        vueJsx(),
        VueDevtools(),
        envParse(),
        Unocss(),
        ...setupUnplugin(viteEnv),
        progress(),
        setupHtmlPlugin(buildTime)
    ]

    if (viteEnv.nodeEnv === 'development') {
        plugins.push(setupMockServer())
    }

    return plugins
}
