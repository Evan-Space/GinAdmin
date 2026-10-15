/// <reference lib="webworker" /> // 告诉编译器，这里是 webworker 环境


/**
 * 把任意二进制数据算成 SHA-256 的十六进制字符串
 * */
export const sha256Hex = async (data: BufferSource) => {
    const digest = await crypto.subtle.digest('SHA-256', data) // 将 data 做 SHA-256  转换，返回 32 字节的ArrayBuffer 
    return Array.from(new Uint8Array(digest)) // 将 ArrayBuffer 转换为 Uint8Array，然后转换为字节数组
        .map((byte) => byte.toString(16).padStart(2, '0')) // 将字节数组转换为十六进制字符串，每个字节转换为 2 位，不足 2 位的前面补 0
        .join('') // 将字节数组转换为十六进制字符串，每个字节转换为 2 位，不足 2 位的前面补 0
} // 返回 SHA-256 的十六进制字符串

self.onmessage = async (event: MessageEvent<{ file: File; chunkSize: number }>) => {
    const { file, chunkSize } = event.data // worker 获取 要分割的文件和 分割大小
    const total = Math.ceil(file.size / chunkSize) // 计算分割数量
    const chunkHashes: string[] = [] // 声明变量，存储分割后的 SHA-256 的十六进制字符串
    for (let index = 0; index < total; index++) { // 开始分割
        const buffer = await file.slice(index * chunkSize, (index + 1) * chunkSize).arrayBuffer() // 分割文件，返回 ArrayBuffer
        chunkHashes.push(await sha256Hex(buffer)) // 计算分割后的 SHA-256 的十六进制字符串
        self.postMessage({ // 发送消息，告诉主线程，分割进度
            type: 'progress',
            value: (index + 1) / total, // 分割进度
        })
    }
    const fileHash = await sha256Hex(new TextEncoder().encode(chunkHashes.join(''))) // 将所有的分片的 hash 合并后，再计算文件的 SHA-256 的十六进制字符串
    self.postMessage({
        type: 'done',
        fileHash: fileHash,
        chunkHashes: chunkHashes, // 发送消息，告诉主线程，所有的分片的 hash
    })
}
