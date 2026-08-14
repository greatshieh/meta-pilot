<template>
    <div :class="{ hidden: hidden }" class="pagination-container">
        <ElPagination
            :background="background"
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :layout="layout"
            :page-sizes="pageSizes"
            :total="total"
            v-bind="$attrs"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange" />
    </div>
</template>

<script setup lang="ts" name="Pagination">
import { computed } from 'vue'

defineOptions({ name: 'Pagination' })

const props = defineProps({
    total: {
        required: true,
        type: Number
    },
    page: {
        type: Number,
        default: 1
    },
    limit: {
        type: Number,
        default: 20
    },
    pageSizes: {
        type: Array<number>,
        default() {
            return [10, 15, 20, 30, 50]
        }
    },
    layout: {
        type: String,
        default: 'total, sizes, prev, pager, next, jumper'
    },
    background: {
        type: Boolean,
        default: false
    },
    autoScroll: {
        type: Boolean,
        default: false
    },
    hidden: {
        type: Boolean,
        default: false
    }
})

const $emits = defineEmits(['update:page', 'update:limit', 'pagination'])

const currentPage = computed({
    get: () => props.page,
    set: val => $emits('update:page', val)
})

const pageSize = computed({
    get: () => props.limit,
    set: val => $emits('update:limit', val)
})

function handleSizeChange(val: number) {
    $emits('pagination', { page: currentPage.value, limit: val })
}
function handleCurrentChange(val: number) {
    $emits('pagination', { page: val, limit: pageSize.value })
}
</script>

<style scoped>
.pagination-container {
    background-color: transparent;
    padding: 32px 16px;
    display: flex;
    justify-content: flex-end;
}

.pagination-container.hidden {
    display: none;
}
</style>
