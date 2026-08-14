<script lang="ts" setup>
import { computed } from 'vue'

defineOptions({ name: 'PinToggler' })

interface Props {
    pin?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits(['click'])

const icon = computed(() => (props.pin ? 'mdi:pin-off' : 'mdi:pin'))
</script>

<template>
    <ElTooltip :content="pin ? '固定' : '取消固定'" placement="bottom-start" :z-index="100">
        <ElButton text @click="emit('click')">
            <div class="flex items-center gap-8px text-lg">
                <SvgIcon :icon="icon" />
            </div>
        </ElButton>
    </ElTooltip>
</template>

<style scoped></style>
