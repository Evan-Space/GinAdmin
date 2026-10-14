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

        worker.onmessage = (event: MessageEvent) => { // 监听 worker 的返回消息
            const data = event.data // 获取 worker 返回的消息
            if (data.type === 'progress') {
                onProgress?.(data.value) // 更新分割进度
                return
            }
            worker.terminate() // 终止 worker
            resolve({
                fileHash: data.fileHash, // 文件的 SHA-256 的十六进制字符串
                chunkHash: data.chunkHash, // 所有的分片的 hash
            })
        }
        worker.onerror = (event: ErrorEvent) => { // 错误处理
            worker.terminate()
            reject(new Error(event.message || '文件哈希值计算失败'))
        }
        worker.postMessage({ // 发送消息，告诉 worker，要分割的文件和 分割大小
            file,
            chunkSize,
        })
    })
}
