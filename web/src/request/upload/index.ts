// import axios from 'axios'
import { getToken, POST, GET, redirectToLogin, BASE_URL } from '../request'

// const uploadClient = axios.create({ baseURL: 'http://localhost:8080/api/v1/upload' })

// uploadClient.interceptors.request.use((config) => {
//     config.headers.Authorization = `Bearer ${getToken()}`
//     return config
// })

// 必须与服务端 storage.chunk_size 保持一致：文件哈希由分片哈希组合而成
export const CHUNK_SIZE = 5 * 1024 * 1024
export const UPLOAD_CONCURRENCY = 4
export const CHUNK_RETRY = 3

export interface UploadTaskResult {
    upload_id: string
    chunk_size: string
    chunk_total: number
    uploaded: number[]
    finished: boolean
    file_id: number
    url: string
}

export interface UploadCompleteResult {
    finished: boolean
    missing: number[]
    file_id: number
    url: string
    size: number
}

/**
 * 初始化上传
 * 同一个接口同时完成秒传探测与断点续传探测，前端不需要自己判断是妙传，还是断点续传
 * */
export const initUploadAPI = (params: {
    file_name: string
    file_size: number
    file_hash: string
    chunk_size: number
}) => POST<UploadTaskResult>('/upload/init', params)

/**
 * 查询上传进度
 * */
export const getUploadStatusAPI = (upload_id: string) => {
    return GET<UploadTaskResult>(`/upload/status?upload_id=${encodeURIComponent(upload_id)}`)
}

/**
 * 合并分片，missing 非空说明还有分片未上传
 * */
export const completeUploadAPI = (upload_id: string) => {
    return POST<UploadTaskResult>('/upload/complete', { upload_id })
}

/**
 * 取消上传，
 * 服务端收回临时分片
 * */
export const abortUploadAPI = (upload_id: string) => {
    return POST('/upload/abort', { upload_id })
}

/**
 * 上传单个分片
 * 用 XHR 而不是 fetch，因为只有 XHR 能拿到上传进度事件
 * 请求体是裸二进制，元信息放查询参数，与服务端 Chunk 接口对应
 */
export const uploadChunkAPI = (
    params: { upload_id: string; index: number; chunk_hash: string; blob: Blob },
    options: { signal?: AbortSignal; onProgress?: (loaded: number) => void } = {},
) =>
    new Promise<void>((resolve, reject) => {
        const query = new URLSearchParams({
            upload_id: params.upload_id,
            index: String(params.index),
            chunk_hash: params.chunk_hash,
        })
        const xhr = new XMLHttpRequest()
        xhr.open('POST', `${BASE_URL}/upload/chunk?${query.toString()}`)
        xhr.setRequestHeader('Content-Type', 'application/octet-stream')
        xhr.setRequestHeader('Authorization', `Bearer ${getToken()}`)
        const onAbort = () => xhr.abort()
        options.signal?.addEventListener('abort', onAbort)
        const cleanup = () => options.signal?.removeEventListener('abort', onAbort)
        xhr.upload.onprogress = (event) => options.onProgress?.(event.loaded)
        xhr.onload = () => {
            cleanup()
            try {
                const result = JSON.parse(xhr.responseText)
                if (result.code === 401) {
                    redirectToLogin()
                }
                result.code === 0 ? resolve() : reject(new Error(result.msg || '分片上传失败'))
            } catch {
                reject(new Error('分片上传响应异常'))
            }
        }
        xhr.onerror = () => {
            cleanup()
            reject(new Error('网络错误'))
        }
        xhr.onabort = () => {
            cleanup()
            reject(new DOMException('已暂停', 'AbortError'))
        }
        xhr.send(params.blob)
    })

