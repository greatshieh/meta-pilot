declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<{}, {}, ayn>
    export default component
}

declare module 'virtual:svg-icons-register' {
    import { Plugin } from 'vite'

    // 根据实际使用的 svg 图标库进行调整
    const svgIconsRegister: Plugin
    export default svgIconsRegister
}

declare module '@wiris/mathtype-ckeditor5/dist/index.js'
declare module '@/vendor/Export2Excel.js'
// oxlint-disable-next-line no-explicit-any
declare const mammoth: any
