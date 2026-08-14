// 计算文件的eTag，来自七牛云开源算法
import CryptoJS from 'crypto-js'

// 以4M为单位分割
const BLOCKSIZE = 4 * 1024 * 1024

// sha1算法
const sha1 = (content: CryptoJS.lib.WordArray | string) => {
    const sha1 = CryptoJS.algo.SHA1.create()
    sha1.update(content)
    return sha1.finalize()
}

const calcEtag = (sha1WordArray: CryptoJS.lib.WordArray, blockCount: number) => {
    let prefix = '16'
    // 如果大于4M，则对各个块的sha1结果再次sha1
    if (blockCount > 1) {
        prefix = '96'
        sha1WordArray = sha1(sha1WordArray)
    }

    sha1WordArray = CryptoJS.enc.Hex.parse(prefix).concat(sha1WordArray)

    return sha1WordArray.toString(CryptoJS.enc.Base64).replace(/\//g, '_').replace(/\+/g, '-')
}

export const getEtag = (file: File) => {
    return new Promise<string>(resolve => {
        let sha1WordArray: CryptoJS.lib.WordArray = CryptoJS.lib.WordArray.create()
        const blockCount = Math.ceil(file.size / BLOCKSIZE)

        const reader = new FileReader()
        reader.onload = e => {
            const buffer = e.target!.result as ArrayBuffer

            for (let i = 0; i < blockCount; i++) {
                const wordArray = CryptoJS.lib.WordArray.create(buffer.slice(i * BLOCKSIZE, (i + 1) * BLOCKSIZE))
                sha1WordArray = sha1WordArray.concat(sha1(wordArray))
            }

            const hash = calcEtag(sha1WordArray, blockCount)
            resolve(hash)
        }

        reader.readAsArrayBuffer(file)
    })
}
