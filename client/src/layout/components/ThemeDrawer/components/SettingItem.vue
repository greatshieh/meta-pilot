<template>
    <div class="w-full flex items-center justify-between">
        <div>
            <span class="pr-8px text-text-primary">{{ label }}</span>
            <slot name="suffix"></slot>
        </div>
        <slot></slot>
    </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'SettingItem' })

interface Props {
    /** Label */
    label: string
}

defineProps<Props>()
</script>

<style scoped></style>
