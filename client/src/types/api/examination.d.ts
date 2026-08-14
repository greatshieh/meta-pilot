declare namespace Api.Examination {
    /** api/method/examination.ts */

    interface PaperPageQuery extends PageQuery {
        grade?: string
        subject?: string
        name?: string
        isPublished?: boolean
    }

    interface ItemType {
        prefix: string
        content: string
    }

    interface QuestionMetaData {
        grade: string
        subject: string
        kindName: string
        uid: string
        title: string
        analytical?: string
        items?: ItemType[]
        correct?: string
        score?: number
        keyPoints?: string
        subQuestion?: QuestionMetaData[]
    }

    interface PaperMetaData {
        uid: string
        grade: string
        subject: string
        name: string
        questions: QuestionMetaData[]
        isPublished?: boolean
        edit?: boolean
        originalName?: string
    }

    interface PaperOptions {
        grade: string
        subject: string
        name: string
    }

    interface PaperResponse extends PageQuery {
        total: number
        list: PaperMetaData[]
    }

    interface UpdataPaperReq {
        uid: string
        name?: string
        isPublished?: boolean
    }
}
