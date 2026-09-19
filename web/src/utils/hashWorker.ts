/// <reference lib="webworker" />

export const sha256Hex = async (data: BufferSource) => {
    const digest = await crypto.subtle.digest('SHA-256', data)
    return Array.from(new Uint8Array(digest))
        .map((byte) => byte.toString(16).padStart(2, '0'))
        .join('')
}

self.onmessage = async (event: MessageEvent<{ file: File; chunkSize: number }>) => {
    const { file, chunkSize } = event.data
    const total = Math.ceil(file.size / chunkSize)
    const chunkHashes: string[] = []

    for (let index = 0; index < total; index++) {
        const buffer = await file.slice(index * chunkSize, (index + 1) * chunkSize).arrayBuffer()
        chunkHashes.push(await sha256Hex(buffer))
        self.postMessage({ type: 'progress', value: (index + 1) / total })
    }

    // 组合哈希， 分片哈希按序合并后，再进行一次哈希，与服务器merger 的算法一直
    const fileHash = await sha256Hex(new TextEncoder().encode(chunkHashes.join('')))
    self.postMessage({ type: 'done', fileHash, chunkHashes })
}
