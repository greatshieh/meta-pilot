import { blake2b } from 'hash-wasm'

// 计算单个块的哈希值
async function calculateBlockHash(block: Uint8Array): Promise<string> {
    return blake2b(block, 128) // 128 字节输出
}

// 流式读取文件并并行计算哈希值
export async function calculateFileHash(file: File, concurrency: number = 4): Promise<string> {
    return new Promise((resolve, reject) => {
        const blockHashes: string[] = [] // 存储每个块的哈希值
        // 使用 ReadableStream 逐块读取文件
        const stream = file.stream()
        const reader = stream.getReader()

        // 并行计算块的哈希值
        const processBlocks = async () => {
            const blockPromises: Promise<string>[] = []

            try {
                while (true) {
                    const { done, value } = await reader.read()
                    if (done) {
                        // 文件读取完成，等待所有块计算完成
                        const hashes = await Promise.all(blockPromises)
                        blockHashes.push(...hashes)

                        // 计算最终哈希值
                        if (blockHashes.length > 1) {
                            // 将每个块的哈希值合并后再计算一次哈希值
                            const combinedHashes = blockHashes.join('')
                            const combinedBuffer = new TextEncoder().encode(combinedHashes)
                            const finalHash = await blake2b(combinedBuffer, 128)
                            const base64Hash = btoa(
                                finalHash
                                    .match(/\w{2}/g)!
                                    .map(byte => String.fromCharCode(parseInt(byte, 16)))
                                    .join('')
                            )
                            resolve(base64Hash)
                        } else {
                            // 只有一个块，直接返回其哈希值的 Base64 编码
                            const base64Hash = btoa(
                                blockHashes[0]
                                    .match(/\w{2}/g)!
                                    .map(byte => String.fromCharCode(parseInt(byte, 16)))
                                    .join('')
                            )
                            resolve(base64Hash)
                        }
                        return
                    }

                    // 将当前块的哈希计算任务加入队列
                    const blockPromise = calculateBlockHash(value)
                    blockPromises.push(blockPromise)

                    // 如果达到并发数，等待一个任务完成
                    if (blockPromises.length >= concurrency) {
                        const index = await Promise.race(
                            blockPromises.map((p, i) => p.then(() => i)) // 返回完成的 Promise 的索引
                        )
                        blockPromises.splice(index, 1) // 移除已完成的任务
                    }
                }
            } catch (error) {
                reject(error)
            }
        }

        // 开始处理块
        processBlocks()
    })
}
