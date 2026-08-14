import { mockDevServerPlugin } from 'vite-plugin-mock-dev-server'

export function setupMockServer() {
    return mockDevServerPlugin({
        reload: true,
        prefix: '/api/v1',
        log: 'debug',
    })
}
