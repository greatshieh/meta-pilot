/**
 * 自定义图片上传adapter
 */

import type { FileLoader, UploadResponse } from 'ckeditor5'
import questionApi from '@/api/method/question'

export class UploadAdapter {
    loader: FileLoader
    constructor(loader: FileLoader) {
        this.loader = loader
    }

    async upload() {
        const file = await this.loader.file

        const data = new FormData()
        data.append('img', file as File)

        return new Promise<UploadResponse>(resolve => {
            questionApi.uploadImage(data).then(resp => {
                resolve({ default: resp.url } as UploadResponse)
            })
        })
    }

    abort() {}
}
