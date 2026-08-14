<template>
    <ElCard class="card-wrapper" shadow="never">
        <div ref="chartRef" class="h-280px w-full" />
    </ElCard>
</template>

<script setup lang="ts" name="Chart">
import { type PropType, watch } from 'vue'
import { type ECOption, useEcharts } from '@/hooks/echarts'

defineOptions({ name: 'Chart' })

const props = defineProps({
    chartOptions: {
        type: Object as PropType<ECOption>,
        default: () => {
            return {}
        }
    }
})

const { domRef: chartRef, setOptions } = useEcharts(() => props.chartOptions, {
    onRender() {}
})

watch(
    () => props.chartOptions,
    val => {
        setOptions(val)
    }
)
</script>
