import { computed, getCurrentInstance, inject, ref } from 'vue'

import type { InjectionKey, Ref } from 'vue'

export const defaultNamespace = 'el'
const statePrefix = 'is-'

const _bem = (namespace: string, block: string, blockSuffix: string, element: string, modifier: string) => {
    let cls = namespace !== '' ? `${namespace}-${block}` : `${block}`
    if (blockSuffix) {
        cls += `-${blockSuffix}`
    }
    if (element) {
        cls += `__${element}`
    }
    if (modifier) {
        cls += `--${modifier}`
    }
    return cls
}

export const namespaceContextKey: InjectionKey<Ref<string | undefined>> = Symbol('namespaceContextKey')

export const useGetDerivedNamespace = (
    // 可选的命名空间覆盖值，允许组件传入自定义命名空间
    namespaceOverrides?: Ref<string | undefined>
) => {
    // 确定最终使用的命名空间：
    // 1. 如果传入了namespaceOverrides，则使用它
    // 2. 否则检查是否存在当前Vue实例：
    //    - 如果存在Vue实例，则从依赖注入中获取命名空间，如果获取不到则使用默认命名空间
    //    - 如果不存在Vue实例（例如在非组件环境中），则直接使用默认命名空间
    const derivedNamespace =
        namespaceOverrides || (getCurrentInstance() ? inject(namespaceContextKey, ref(defaultNamespace)) : ref(defaultNamespace))

    // 创建一个计算属性，确保返回有效的命名空间：
    // 1. 解包derivedNamespace的值
    // 2. 如果解包后的值为空，则使用默认命名空间
    const namespace = computed<string>(() => {
        return derivedNamespace.value !== undefined ? derivedNamespace.value : defaultNamespace
    })

    // 返回计算后的命名空间
    return namespace
}

export const useNamespace = (
    // BEM命名规则中的块（Block）名称，通常是组件的根元素名称
    // 例如：'button', 'input', 'form'
    block: string,
    // 可选的命名空间覆盖值，允许组件传入自定义命名空间
    // 例如：ref('my-custom-namespace')
    namespaceOverrides?: Ref<string | undefined>
) => {
    // 获取派生的命名空间，基于传入的overrides或上下文中的命名空间
    const namespace = useGetDerivedNamespace(namespaceOverrides)

    // 生成BEM格式的块（Block）类名
    // 示例：b() => 'el-button' （假设默认命名空间是'el'，block是'button'）
    // 示例：b('primary') => 'el-button-primary'
    const b = (blockSuffix = '') => _bem(namespace.value, block, blockSuffix, '', '')

    // 生成BEM格式的元素（Element）类名
    // 示例：e('icon') => 'el-button__icon'
    // 示例：e() => '' （没有元素时不生成类名）
    const e = (element?: string) => (element ? _bem(namespace.value, block, '', element, '') : '')

    // 生成BEM格式的修饰符（Modifier）类名
    // 示例：m('disabled') => 'el-button--disabled'
    // 示例：m() => '' （没有修饰符时不生成类名）
    const m = (modifier?: string) => (modifier ? _bem(namespace.value, block, '', '', modifier) : '')

    // 生成带有后缀的块和元素的BEM类名
    // 示例：be('group', 'item') => 'el-button-group__item'
    // 示例：be() => '' （参数不完整时不生成类名）
    const be = (blockSuffix?: string, element?: string) =>
        blockSuffix && element ? _bem(namespace.value, block, blockSuffix, element, '') : ''

    // 生成元素和修饰符的BEM类名
    // 示例：em('icon', 'loading') => 'el-button__icon--loading'
    // 示例：em() => '' （参数不完整时不生成类名）
    const em = (element?: string, modifier?: string) => (element && modifier ? _bem(namespace.value, block, '', element, modifier) : '')

    // 生成带后缀的块和修饰符的BEM类名
    // 示例：bm('group', 'disabled') => 'el-button-group--disabled'
    // 示例：bm() => '' （参数不完整时不生成类名）
    const bm = (blockSuffix?: string, modifier?: string) =>
        blockSuffix && modifier ? _bem(namespace.value, block, blockSuffix, '', modifier) : ''

    // 生成完整的BEM类名（块后缀、元素、修饰符）
    // 示例：bem('group', 'item', 'active') => 'el-button-group__item--active'
    // 示例：bem() => '' （参数不完整时不生成类名）
    const bem = (blockSuffix?: string, element?: string, modifier?: string) =>
        blockSuffix && element && modifier ? _bem(namespace.value, block, blockSuffix, element, modifier) : ''

    // 生成状态类名（通常用于表示组件的状态）
    // 示例：is('focus') => 'is-focus'
    // 示例：is('disabled', false) => '' （当状态为false时不生成类名）
    const is: {
        (name: string, state: boolean | undefined): string
        (name: string): string
    } = (name: string, ...args: [boolean | undefined] | []) => {
        const state = args.length >= 1 ? args[0]! : true
        return name && state ? `${statePrefix}${name}` : ''
    }

    // 生成CSS变量对象
    // 示例：cssVar({ 'font-size': '12px' }) => { '--el-font-size': '12px' }
    const cssVar = (object: Record<string, string>) => {
        const styles: Record<string, string> = {}
        for (const key in object) {
            if (object[key]) {
                styles[`--${namespace.value}-${key}`] = object[key]
            }
        }
        return styles
    }

    // 生成带块前缀的CSS变量对象
    // 示例：cssVarBlock({ 'font-size': '12px' }) => { '--el-button-font-size': '12px' }
    const cssVarBlock = (object: Record<string, string>) => {
        const styles: Record<string, string> = {}
        for (const key in object) {
            if (object[key]) {
                styles[`--${namespace.value}-${block}-${key}`] = object[key]
            }
        }
        return styles
    }

    // 生成CSS变量名
    // 示例：cssVarName('font-size') => '--el-font-size'
    const cssVarName = (name: string) => `--${namespace.value}-${name}`

    // 生成带块前缀的CSS变量名
    // 示例：cssVarBlockName('font-size') => '--el-button-font-size'
    const cssVarBlockName = (name: string) => `--${namespace.value}-${block}-${name}`

    // 返回包含所有方法的对象，供组件使用
    return {
        namespace, // 命名空间引用
        b, // 生成块类名
        e, // 生成元素类名
        m, // 生成修饰符类名
        be, // 生成块后缀+元素类名
        em, // 生成元素+修饰符类名
        bm, // 生成块后缀+修饰符类名
        bem, // 生成完整BEM类名
        is, // 生成状态类名
        cssVar, // 生成CSS变量对象
        cssVarName, // 生成CSS变量名
        cssVarBlock, // 生成带块前缀的CSS变量对象
        cssVarBlockName // 生成带块前缀的CSS变量名
    }
}

export type UseNamespaceReturn = ReturnType<typeof useNamespace>
