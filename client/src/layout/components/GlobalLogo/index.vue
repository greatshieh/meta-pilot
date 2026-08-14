<template>
    <RouterLink to="/" class="w-full flex justify-center items-center ws-nowrap overflow-hidden">
        <SystemLogo class="size-52px" />
        <h2 v-show="showTitle" class="pl-8px text-24px text-primary font-bold transition duration-300 ease-in-out">
            {{ title }}
        </h2>
    </RouterLink>
</template>

<script setup lang="ts">
defineOptions({ name: 'GlobalLogo' })

defineProps({
    showTitle: {
        type: Boolean,
        default: true
    }
})

const title = import.meta.env.VITE_APP_TITLE
</script>

<style lang="scss" scoped></style>
