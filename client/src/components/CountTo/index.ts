import { buildProps } from 'element-plus/es/utils/index.mjs'
import type { PropType } from 'vue'

export const CountToProps = buildProps({
    startVal: {
        type: Number as PropType<number>,
        default: 0
    },
    endVal: {
        type: Number as PropType<number>,
        default: 2025
    },
    duration: {
        type: Number as PropType<number>,
        default: 3000
    },
    autoplay: {
        type: Boolean as PropType<boolean>,
        default: false
    },
    decimals: {
        type: Number as PropType<number>,
        default: 0
    },
    decimal: {
        type: String as PropType<string>,
        default: '.'
    },
    separator: {
        type: String as PropType<string>,
        default: ','
    },
    prefix: {
        type: String as PropType<string>,
        default: ''
    },
    suffix: {
        type: String as PropType<string>,
        default: ''
    },
    useEasing: {
        type: Boolean as PropType<boolean>,
        default: true
    },
    easingFn: {
        type: Function as PropType<(t: number, b: number, c: number, d: number) => number>,
        default(t: number, b: number, c: number, d: number) {
            return (c * (-(2 ** ((-10 * t) / d)) + 1) * 1024) / 1023 + b
        }
    }
})
