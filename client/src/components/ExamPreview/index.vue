<template>
    <div
        class="preview"
        :class="{
            'sub-no-items': onlyDigits() && (!question?.items || (question?.items && question?.items.length === 0))
        }">
        <!-- 题干 -->
        <div v-if="question?.title" v-html="question.title" ref="titleRef" class="title content"></div>

        <!-- 选项 -->
        <div class="preview-items" v-if="question?.items && question.items.length > 0">
            <template v-for="item in question?.items">
                <div
                    v-if="item.content !== null && item.content !== ''"
                    :key="item.prefix"
                    class="item-option item-key"
                    :class="{ 'is-active': question?.correct === item.prefix }"
                    :ref="'item' + item.prefix">
                    <label class="option-item-label">{{ item.prefix }}.</label>
                    <div v-html="item.content" class="option-item-content content"></div>
                </div>
            </template>
        </div>

        <div class="sub-content">
            <!-- 解析 -->
            <div class="item-key" v-if="!question?.subQuestion || question?.subQuestion?.length === 0">
                <ElTag class="item-label" type="primary">解析</ElTag>
                <div class="item-content" :class="{ 'null-content': question!.analytical === '' }">
                    <div v-if="question!.analytical !== ''" v-html="question!.analytical"></div>
                    <span v-else>空</span>
                </div>
            </div>

            <!-- 答案 -->
            <div class="item-key" v-if="!question?.subQuestion || question?.subQuestion?.length === 0">
                <ElTag class="item-label" type="primary">答案</ElTag>
                <div v-html="question!.correct" class="item-content"></div>
                <ElTag v-if="scoreDetail && scoreDetail !== ''" class="item-label score" type="primary">答题选择</ElTag>
                <div v-if="scoreDetail && scoreDetail !== ''" class="item-content">
                    {{ scoreDetail }}
                </div>
            </div>

            <!-- 分数 -->
            <div class="item-key" v-if="question?.score && (!question.subQuestion || question?.subQuestion?.length === 0)">
                <ElTag class="item-label" type="primary">分数</ElTag>
                <div class="item-content" :class="{ 'null-content': question.score === 0 }">
                    {{ question.score === 0 ? '空' : question.score }}
                </div>
            </div>

            <div class="item-key" v-for="row in question?.subQuestion" :key="row.uid">
                <exam-preview :question="row" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts" name="ExamPreview">
import { type PropType, useTemplateRef } from 'vue'

defineProps({
    question: {
        type: Object as PropType<Api.Examination.QuestionMetaData>
    },
    scoreDetail: {
        type: String,
        default: undefined
    }
})

defineOptions({
    name: 'ExamPreview'
})

const titleRef = useTemplateRef<HTMLDivElement>('titleRef')

function onlyDigits(): boolean {
    if (!titleRef.value || !titleRef.value?.textContent) return true

    const regex = /^\d+[.、]?$/g
    return regex.test(titleRef.value.textContent.trim())
}
</script>

<style lang="scss">
@use './style.scss';
</style>
