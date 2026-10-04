export interface FileHashResult {
    fileHash: string
    chunkHash: string[]
}

/**
 * 使用 webworker 计算hash ，避免页面卡死
 * 每个分片哈希 安顺序拼接后，再哈希一次，算法需要和后端一致
 * */
export const computeFileHash = (
    file: File,
    chunkSize: number,
    onProgress?: (val: number) => void,
): Promise<FileHashResult> => {
    return new Promise((resolve, reject) => {
        const worker = new Worker(new URL('./hashWorker.ts', import.meta.url), { type: 'module' })

        worker.onmessage = (event: MessageEvent) => {
            const data = event.data
            if (data.type === 'progress') {
                onProgress?.(data.value)
                return
            }
            worker.terminate()
            resolve({
                fileHash: data.fileHash,
                chunkHash: data.chunkHash,
            })
        }
        worker.onerror = (event: ErrorEvent) => {
            worker.terminate()
            reject(new Error(event.message || '文件哈希值计算失败'))
        }
        worker.postMessage({
            file,
            chunkSize,
        })
    })
}
