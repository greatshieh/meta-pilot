import { defineStore } from 'pinia'
import type { QuestionMetaData } from '@/views/exambank/constants'
import { SetupStoreId } from '../enum'

export interface ExamStoreProp {
    question: QuestionMetaData | undefined
    status: string | undefined
}

const useExamStore = defineStore(SetupStoreId.Exam, {
    state: (): ExamStoreProp => {
        return {
            // 待编辑的试题
            question: undefined,
            // 是否是上传文件
            status: undefined
        }
    },

    actions: {
        setQuestion(q: QuestionMetaData) {
            this.question = q
        },

        clearQuestion() {
            this.question = undefined
        },

        setStatus(status: string) {
            this.status = status
        },

        clearStatus() {
            this.status = undefined
        },

        reset() {
            this.$reset()
        }
    }
})

export default useExamStore
