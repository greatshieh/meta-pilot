import type { Component } from 'vue'
import { h } from 'vue'
import SvgIcon from '@/components/SvgIcon/index.vue'

interface IconConfig {
    /** Iconify icon name */
    icon?: string
    /** Local icon name */
    localIcon?: string
    /** Icon color */
    color?: string
    /** Icon size */
    fontSize?: number
}

/**
 * Svg icon render hook
 *
 * @param SvgIcon Svg icon component
 */
const useSvgIconRender = (SvgIcon: Component) => {
    type IconStyle = Partial<Pick<CSSStyleDeclaration, 'color' | 'fontSize'>>

    /**
     * Svg icon VNode
     *
     * @param config
     */
    const SvgIconVNode = (config: IconConfig) => {
        const { color, fontSize, icon, localIcon } = config

        const style: IconStyle = {}

        if (color) {
            style.color = color
        }
        if (fontSize) {
            style.fontSize = `${fontSize}px`
        }

        if (!icon && !localIcon) {
            return undefined
        }

        return () => h(SvgIcon, { icon, localIcon, style })
    }

    return {
        SvgIconVNode
    }
}

export const useSvgIcon = () => {
    const { SvgIconVNode } = useSvgIconRender(SvgIcon)

    return {
        SvgIconVNode
    }
}
