<template>
    <ElCard class="card-wrapper filter-container">
        <ElForm inline size="small">
            <ElFormItem v-if="showItem.includes('grade')">
                <ElSelect v-model="queryParam.grade" placeholder="年级" clearable style="width: 80px" @change="handleSelectChange('grade')">
                    <ElOption v-for="item in gradeOptions" :key="item" :label="item" :value="item" />
                </ElSelect>
            </ElFormItem>

            <ElFormItem v-if="showItem.includes('class')">
                <ElSelect v-model="queryParam.class" placeholder="班级" clearable style="width: 80px" @change="handleSelectChange('class')">
                    <ElOption v-for="item in classOptions" :key="item" :label="item" :value="item" />
                </ElSelect>
            </ElFormItem>

            <ElFormItem v-if="showItem.includes('stuName')">
                <ElInput
                    v-model="queryParam.name"
                    placeholder="学生姓名"
                    style="width: 100px"
                    clearable
                    @input="handleSelectChange('stuName')"></ElInput>
            </ElFormItem>

            <ElFormItem v-if="showItem.includes('examName')">
                <ElSelect
                    :disabled="!queryParam.grade || queryParam.grade === ''"
                    v-model="queryParam.examName"
                    placeholder="请选择考试名称"
                    style="width: 240px"
                    @change="handleSelectChange('examName')">
                    <ElOption v-for="item in examOptions" :label="item" :value="item" :key="item"></ElOption>
                </ElSelect>
            </ElFormItem>

            <ElFormItem label="科目" v-if="showItem.includes('subject')">
                <ElSelect
                    v-model="queryParam.subject"
                    placeholder="科目"
                    clearable
                    style="width: 80px"
                    @change="handleSelectChange('subject')">
                    <ElOption v-for="s in subjectOptions" :key="s" :label="s" :value="s" />
                </ElSelect>
            </ElFormItem>

            <slot></slot>
        </ElForm>
    </ElCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import useSseStore from '@/store/modules/sse'
import { SUBJECT_OPTIONS } from '@/types/constants'
import type { SearchFilterProps } from './search-filter'

defineOptions({
    name: 'SearchFilter'
})

const sseStore = useSseStore()
const $emits = defineEmits(['update:queryParam', 'change'])
const props = withDefaults(defineProps<SearchFilterProps>(), {
    showItem: 'grade, class, stuName, examName, subject',
    queryParam: () => ({
        grade: '',
        class: '',
        examName: '',
        subject: '',
        name: ''
    })
})

const currentParams = computed({
    get: () => props.queryParam,
    set: val => $emits('update:queryParam', val)
})

const gradeOptions = sseStore.gradeOptions

const classOptions = computed(() => {
    if (!currentParams.value.grade) return []
    return sseStore.classOptions[currentParams.value.grade].sort(
        (a: string, b: string) => parseInt(a.match(/\d+/)?.[0] || '0', 10) - parseInt(b.match(/\d+/)?.[0] || '0', 10)
    )
})

const examOptions = computed(() => {
    if (!currentParams.value.grade) return []
    return sseStore.examNameOptions[currentParams.value.grade].map(e => e.examName)
})

const subjectOptions = computed(() => {
    if (!currentParams.value.grade) return []
    if (!currentParams.value.examName) return SUBJECT_OPTIONS
    return sseStore.examNameOptions[currentParams.value.grade]
        .find(e => e.examName === currentParams.value.examName)
        ?.subject.split(',')
        .sort((a, b) => SUBJECT_OPTIONS.indexOf(a) - SUBJECT_OPTIONS.indexOf(b))
})

let isSearchTriggered = false

const handleSelectChange = (key: string) => {
    if (isSearchTriggered) return

    isSearchTriggered = true

    if (key === 'grade') {
        currentParams.value.class = ''
        currentParams.value.examName = ''
        currentParams.value.subject = ''
        currentParams.value.name = ''
    } else if (key === 'class') {
        currentParams.value.examName = ''
        currentParams.value.subject = ''
        currentParams.value.name = ''
    }

    $emits('change')

    isSearchTriggered = false
}
</script>

<style lang="scss">
.filter-container {
    .el-card__body {
        padding: 16px 20px;
    }

    .el-form-item {
        margin-bottom: 0px;
    }
}
</style>
