/** The union key namespace */
declare namespace UnionKey {
    /**
     * The login module
     *
     * - pwd-login: password login
     * - code-login: phone code login
     * - register: register
     * - reset-pwd: reset password
     * - bind-wechat: bind wechat
     */
    type LoginModule = 'pwd-login' | 'code-login' | 'register' | 'reset-pwd' | 'bind-wechat'

    /** Theme scheme */
    type ThemeScheme = 'light' | 'dark'

    /**
     * Reset cache strategy
     *
     * - close: re-cache when close page
     * - refresh: re-cache when refresh page
     */
    type ResetCacheStrategy = 'close' | 'refresh'

    /**
     * The layout mode
     *
     * - vertical: the vertical menu in left
     * - horizontal: the horizontal menu in top
     * - vertical-mix: two vertical mixed menus in left
     * - horizontal-mix: the vertical first level menus in left and horizontal child level menus in top
     */
    type ThemeLayoutMode = 'vertical' | 'horizontal' | 'vertical-mix' | 'horizontal-mix'

    /**
     * The scroll mode when content overflow
     *
     * - wrapper: the wrapper component's root element overflow
     * - content: the content component overflow
     */
    type ThemeScrollMode = 'wrapper' | 'content'

    /** Page animate mode */
    type ThemePageAnimateMode = 'fade' | 'fade-slide' | 'fade-bottom' | 'fade-scale' | 'zoom-fade' | 'zoom-out' | 'none'

    /**
     * Tab mode
     *
     * - chrome: chrome style
     * - button: button style
     */
    type ThemeTabMode = 'button' | 'chrome'

    type AlovaMethod = 'DELETE' | 'GET' | 'PATCH' | 'POST' | 'PUT'
}
