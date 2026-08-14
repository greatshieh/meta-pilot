import type http from 'node:http'
import type { HttpProxy, ProxyOptions } from 'vite'
import { clearScreen, createColors } from './cli-helper'

const colors = createColors()

/**
 * Set http proxy
 *
 * @param env - The current env
 * @param enable - If enable http proxy
 */
export function createViteProxy(env: Env.ImportMeta, enable: boolean) {
    const isEnableHttpProxy = enable && env.VITE_HTTP_PROXY === 'Y'

    if (!isEnableHttpProxy || !env.VITE_API_PREFIX) return undefined

    const proxy: Record<string, ProxyOptions> = {}

    proxy[env.VITE_API_PREFIX] = {
        target: `${env.VITE_SERVICE_BASE_URL}:${env.VITE_SERVER_PORT}`, // 代理到目标路径
        changeOrigin: true,
        configure: (_proxy: HttpProxy.ProxyServer, options: ProxyOptions) => {
            _proxy.on('proxyReq', (_proxyReq: http.ClientRequest, req: http.IncomingMessage) => {
                clearScreen()
                console.log(colors.bgYellow(`  ${req.method}  `), colors.green(`${options.target}${req.url}`))
            })
            _proxy.on('error', (_err: Error, req: http.IncomingMessage) => {
                console.log(colors.bgRed(`Error：${req.method}  `), colors.green(`${options.target}${req.url}`))
            })
        }
        // rewrite: (path) => {
        //   console.log(path)
        //   path.replace(new RegExp(`^${env.VITE_API_PREFIX}`), '')
        // }
    }

    return proxy
}
