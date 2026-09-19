export interface FileHashResult {
    fileHash: string
    chunkHashes: string[]
}

export const computeFileHash = (
    file: File,
    chunkSize: number,
    onProgress?: (value: number) => void,
): Promise<FileHashResult> => {
    return new Promise((resolve, reject) => {
        const worker = new Worker(new URL('./hashWorker.ts', import.meta.url), { type: 'module' })

        worker.onmessage = (event) => {
            const data = event.data
            if (data.type === 'progress') {
                onProgress?.(data.value)
                return
            }

            worker.terminate()
            resolve({ fileHash: data.fileHash, chunkHashes: data.chunkHashes })
        }

        worker.onerror = (event) => {
            worker.terminate()
            reject(new Error(event.message || '文件哈希计算失败'))
        }
        worker.postMessage({ file, chunkSize })
    })
}
