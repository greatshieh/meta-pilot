<template>
    <div style="height: 100%; position: relative">
        <div class="exception-conatiner flex-col-center size-full min-h-520px gap-24px overflow-hidden">
            <div class="svg-container flex text-size-400px text-primary">
                <SvgIcon :icon="icon" />
            </div>
            <div>{{ second }}秒后自动返回首页</div>
            <ElButton type="primary" @click="goTODahboard">返回首页</ElButton>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useIntervalFn } from '@vueuse/core'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({ name: 'ExceptionBase' })

const $router = useRouter()

type ExceptionType = '403' | '404' | '500'

interface Props {
    /**
     * Exception type
     *
     * - 403: no permission
     * - 404: not found
     * - 500: service error
     */
    type: ExceptionType
}

const props = defineProps<Props>()

const second = ref(5)

const iconMap: Record<ExceptionType, string> = {
    '403': 'no-permission',
    '404': 'not-found',
    '500': 'service-error'
}

const icon = computed(() => iconMap[props.type])

const { resume, pause } = useIntervalFn(
    () => {
        second.value = second.value - 1
        if (second.value == 0) {
            // 停止计时器
            goTODahboard()
            pause()
        }
    },
    1000,
    { immediate: false }
)

onMounted(() => {
    resume()
})

const goTODahboard = () => {
    $router.push('/dashboard')
    pause()
}
</script>

<style lang="scss" scoped></style>
