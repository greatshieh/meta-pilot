<template>
    <div class="bg-page c-text-primary transition-300" :class="{ 'bg-inverted text-#1f1f1f': inverted }">
        <slot></slot>
    </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'DarkModeContainer' })

interface Props {
    inverted?: boolean
}

defineProps<Props>()
</script>

<style scoped></style>
