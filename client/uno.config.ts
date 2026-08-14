import type { Theme } from '@unocss/preset-uno'
import transformerDirectives from '@unocss/transformer-directives'
import transformerVariantGroup from '@unocss/transformer-variant-group'
import { defineConfig } from '@unocss/vite'
import { presetWind3 } from 'unocss'
import { presetUniAdmin, presetIcon } from './src/plugins/unopreset'

export default defineConfig<Theme>({
    content: {
        pipeline: {
            exclude: ['node_modules', 'dist']
        }
    },
    theme: {
        colors: {
            nprogress: 'rgb(var(--nprogress-color))',
            container: 'rgb(var(--bg-color-muted)',
            base: 'var(--el-bg-color)',
            page: 'var(--el-bg-color-page)',
            muted: 'var(--bg-color-muted)',
            inverted: 'rgb(var(--inverted-bg-color))',
            primary: 'var(--el-color-primary)',
            success: 'var(--el-color-success)',
            danger: 'var(--el-color-danger)',
            'text-primary': 'var(--el-text-color-primary)',
            'text-regular': 'var(--el-text-color-regular))',
            'text-secondary': 'var(--el-text-color-secondary)',
            'text-placeholder': 'var(--el-text-color-placeholder)',
            'text-disabled': 'var(--el-text-color-disabled)'
        },
        boxShadow: {
            header: 'var(--header-box-shadow)',
            sider: 'var(--sider-box-shadow)',
            tab: 'var(--layout-box-shadow)'
        },
        fontSize: {
            'icon-xs': '0.875rem',
            'icon-small': '1rem',
            icon: '1.125rem',
            'icon-large': '1.5rem',
            'icon-xl': '2rem'
        }
    },
    transformers: [transformerDirectives(), transformerVariantGroup()],
    presets: [presetWind3({ dark: 'class' }), presetUniAdmin(), presetIcon()]
})
