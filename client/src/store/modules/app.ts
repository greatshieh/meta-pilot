import { defineStore } from 'pinia'
import { SetupStoreId } from '../enum'
import useThemeStore from './theme'
import { ref } from 'vue'

const useAppStore = defineStore(SetupStoreId.App, () => {
    const siderCollapse = ref(false)
    const device = ref('desktop')
    const size = ref('medium')
    const themeDrawerVisible = ref(false)
    const mixSiderFixed = ref(false)
    const contentXScrollable = ref(false)
    const fullContent = ref(false)
    const isMobile = ref(false)
    const reloadFlag = ref(true)

    function toggleSidebar() {
        siderCollapse.value = !siderCollapse.value
    }

    function toggleMixSiderFixed() {
        mixSiderFixed.value = !mixSiderFixed.value
    }

    function setContentXScrollable(val: boolean) {
        contentXScrollable.value = val
    }

    function toggleFullContent() {
        fullContent.value = !fullContent.value
    }

    async function reloadPage(duration = 300) {
        reloadFlag.value = false

        const themestore = useThemeStore()

        const d = themestore.page.animate ? duration : 40

        await new Promise(resolve => {
            setTimeout(resolve, d)
        })

        reloadFlag.value = true
    }

    return {
        siderCollapse,
        device,
        size,
        themeDrawerVisible,
        mixSiderFixed,
        contentXScrollable,
        fullContent,
        isMobile,
        reloadFlag,
        toggleSidebar,
        toggleMixSiderFixed,
        setContentXScrollable,
        toggleFullContent,
        reloadPage
    }
})

export default useAppStore
