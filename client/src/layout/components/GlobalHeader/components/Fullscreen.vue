<template>
    <div @click="handleToggle">
        <SvgIcon :icon="isFullscreen ? 'gridicons:fullscreen-exit' : 'gridicons:fullscreen'" />
    </div>
</template>

<script setup lang="ts">
import { useFullscreen } from '@vueuse/core'

defineOptions({ name: 'Fullscreen' })

const { isFullscreen, toggle } = useFullscreen()

function handleToggle() {
    toggle()
}
</script>
