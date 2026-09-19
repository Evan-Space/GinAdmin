import { useRef, useState } from 'react'
import { message } from 'antd'

import { computeFileHash, FileHashResult } from '@src/utils/fileHash'
import {
    CHUNK_RETRY,
    CHUNK_SIZE,
    UPLOAD_CONCURRENCY,
    abortUploadAPI,
    completeUploadAPI,
    getUploadStatusAPI,
    initUploadAPI,
    uploadChunkAPI,
} from '@src/request/upload'

export type UploadStatus =
    'idle' | 'hashing' | 'uploading' | 'merging' | 'success' | 'paused' | 'error'

const STORAGE_KEY = 'upload_task_map'

// 本地记住 hash -> upload_id，刷新页面后能省掉重算哈希；丢了也不影响正确性
const readTaskMap = () => {
    try {
        return JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
    } catch {
        return {}
    }
}

const saveTaskId = (hash: string, uploadId: string) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ ...readTaskMap(), [hash]: uploadId }))
}

const dropTaskId = (hash: string) => {
    const map = readTaskMap()
    delete map[hash]
    localStorage.setItem(STORAGE_KEY, JSON.stringify(map))
}

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

const isAbortError = (error: unknown) => (error as Error)?.name === 'AbortError'

export const useChunkUpload = () => {
    const [status, setStatus] = useState<UploadStatus>('idle')
    const [progress, setProgress] = useState(0)
    const [hashProgress, setHashProgress] = useState(0)
    const [fileName, setFileName] = useState('')

    const fileRef = useRef<File | null>(null)
    const hashRef = useRef<FileHashResult | null>(null)
    const uploadIdRef = useRef('')
    const loadedRef = useRef<Record<number, number>>({}) // 每片已传字节，用于合成总进度
    const abortRef = useRef<AbortController | null>(null)
    const pausedRef = useRef(false)

    // 进度取历史最大值，避免失败重传导致进度条回退
    const refreshProgress = () => {
        const file = fileRef.current
        if (!file || file.size === 0) return
        const loaded = Object.values(loadedRef.current).reduce((sum, value) => sum + value, 0)
        setProgress((prev) => Math.max(prev, Math.min(99, Math.floor((loaded / file.size) * 100))))
    }

    /** 单片上传，失败按指数退避重试；暂停造成的中断不算失败 */
    const uploadChunk = async (index: number) => {
        const file = fileRef.current!
        const blob = file.slice(index * CHUNK_SIZE, (index + 1) * CHUNK_SIZE)
        for (let attempt = 0; attempt <= CHUNK_RETRY; attempt++) {
            try {
                await uploadChunkAPI(
                    {
                        upload_id: uploadIdRef.current,
                        index,
                        chunk_hash: hashRef.current!.chunkHashes[index],
                        blob,
                    },
                    {
                        signal: abortRef.current?.signal,
                        onProgress: (loaded) => {
                            loadedRef.current[index] = loaded
                            refreshProgress()
                        },
                    },
                )
                loadedRef.current[index] = blob.size
                refreshProgress()
                return
            } catch (error) {
                if (isAbortError(error)) throw error
                loadedRef.current[index] = 0
                if (attempt === CHUNK_RETRY) throw error
                await sleep(500 * 2 ** attempt)
            }
        }
    }

    /** 固定并发的工人池，浏览器同域连接数有限，一次性全发出去只会排队 */
    const runPool = async (indexes: number[]) => {
        const queue = [...indexes]
        const workers = Array.from(
            { length: Math.min(UPLOAD_CONCURRENCY, queue.length) },
            async () => {
                while (queue.length > 0) {
                    if (pausedRef.current) return
                    await uploadChunk(queue.shift()!)
                }
            },
        )
        await Promise.all(workers)
    }

    /** 从服务端给的已传清单出发，只补缺失分片，然后请求合并 */
    const runUpload = async (uploaded: number[], chunkTotal: number) => {
        setStatus('uploading')
        const done = new Set(uploaded)
        done.forEach((index) => {
            const start = index * CHUNK_SIZE
            loadedRef.current[index] = Math.min(CHUNK_SIZE, fileRef.current!.size - start)
        })
        refreshProgress()
        const pending = Array.from({ length: chunkTotal }, (_, index) => index).filter(
            (index) => !done.has(index),
        )
        await runPool(pending)
        if (pausedRef.current) {
            setStatus('paused')
            return
        }
        setStatus('merging')
        const res = await completeUploadAPI(uploadIdRef.current)
        if (res.code !== 0) throw new Error(res.msg)
        // 服务端点名发现缺片，补齐后再合并一次
        if (!res.data.finished) {
            await runPool(res.data.missing)
            const retry = await completeUploadAPI(uploadIdRef.current)
            if (retry.code !== 0 || !retry.data.finished)
                throw new Error(retry.msg || '文件合并失败')
        }
        dropTaskId(hashRef.current!.fileHash)
        setProgress(100)
        setStatus('success')
        message.success('上传完成')
    }

    /** 选择文件后的入口：算哈希 -> init -> 传分片 -> 合并 */
    const start = async (file: File) => {
        try {
            fileRef.current = file
            hashRef.current = null
            uploadIdRef.current = ''
            loadedRef.current = {}
            pausedRef.current = false
            abortRef.current = new AbortController()
            setFileName(file.name)
            setProgress(0)
            setHashProgress(0)
            setStatus('hashing')
            hashRef.current = await computeFileHash(file, CHUNK_SIZE, (value) =>
                setHashProgress(Math.floor(value * 100)),
            )
            const res = await initUploadAPI({
                file_name: file.name,
                file_size: file.size,
                file_hash: hashRef.current.fileHash,
                chunk_size: CHUNK_SIZE,
            })
            if (res.code !== 0) throw new Error(res.msg)
            // 秒传：服务端已有同内容文件
            if (res.data.finished) {
                setProgress(100)
                setStatus('success')
                message.success('文件已存在，秒传完成')
                return
            }
            uploadIdRef.current = res.data.upload_id
            saveTaskId(hashRef.current.fileHash, res.data.upload_id)
            await runUpload(res.data.uploaded, res.data.chunk_total)
        } catch (error) {
            if (isAbortError(error)) return
            setStatus('error')
            message.error((error as Error).message || '上传失败')
        }
    }
    /** 暂停：中断在途请求，已落盘的分片不会丢 */
    const pause = () => {
        pausedRef.current = true
        abortRef.current?.abort()
        setStatus('paused')
    }
    /** 继续：重新向服务端要一次已传清单，前端不用自己记进度 */
    const resume = async () => {
        if (!fileRef.current || !uploadIdRef.current) return
        pausedRef.current = false
        abortRef.current = new AbortController()
        try {
            const res = await getUploadStatusAPI(uploadIdRef.current)
            if (res.code !== 0) throw new Error(res.msg)
            loadedRef.current = {}
            await runUpload(res.data.uploaded, res.data.chunk_total)
        } catch (error) {
            if (isAbortError(error)) return
            setStatus('error')
            message.error((error as Error).message || '续传失败')
        }
    }
    /** 取消：服务端删临时分片，本地清状态 */
    const cancel = async () => {
        pausedRef.current = true
        abortRef.current?.abort()
        if (uploadIdRef.current) {
            await abortUploadAPI(uploadIdRef.current)
            if (hashRef.current) dropTaskId(hashRef.current.fileHash)
        }
        uploadIdRef.current = ''
        loadedRef.current = {}
        setProgress(0)
        setStatus('idle')
    }
    return { status, progress, hashProgress, fileName, start, pause, resume, cancel }
}
