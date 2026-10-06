import { useRef, useState } from 'react'
import { xhrRequestWithRetry } from '@src/request/UploadFile'
import { POST, GET } from '@src/request'
import { limitPromise } from '@src/utils/utils.ts'
import { computeFileHash } from '@src/pages/_layout/UploadFile/-components/customDemo/utils/fileHash.ts'

const CHUNK_SIZE = 5 * 1024 * 1024

type UploadStatus =
    | 'idle'
    | 'hashing'
    | 'uploading'
    | 'paused'
    | 'merging'
    | 'done'
    | 'error'

type InitData = {
    upload_id: string
    chunk_size: number
    chunk_total: number
    uploaded: number[]
    finished: boolean
    path: string
}

export const useCustomUploadFile = () => {
    const [onProgress, setOnProgress] = useState<number>(0)
    const [status, setStatus] = useState<UploadStatus>('idle') // 默认状态是闲置状态

    const fileRef = useRef<File | null>(null)
    const uploadIdRef = useRef<string>('')
    const chunkSizeRef = useRef(CHUNK_SIZE)
    const totalRef = useRef<number>(0)
    const abortRef = useRef<AbortController | null>(null)
    const pausedRef = useRef<boolean>(false) // 暂停
    const loadedRef = useRef<number[]>([]) // 已经上传切片的进度

    /**
     * 通用的，计算当前所有切片，已经上传的字节数，并更新上传进度
     *  */
    const reportProgress = () => {
        const file = fileRef.current
        if (!file) return
        const sum = loadedRef.current.reduce((a, b) => a + b, 0)
        setOnProgress(Math.min(99,
            Math.round((sum / file.size) * 100)
            )
        )
    }

    /**
     * 在断点续传场景中，根据已经上传的分片序号，计算已经上传的字节数，并更新上传进度
     * param uploaded: 已经上传的分片序号
     * return: 上传结果
     */
    const syncLoadedFromUploaded = (uploaded: number[]) => {
        const file = fileRef.current
        if (!file) return
        const chunkSize = chunkSizeRef.current
        const doneSet = new Set(uploaded ?? []) // 对 uploaded 去重，doneSet 里面存储的是已经上传 分片的 index

        loadedRef.current = Array.from({ length: totalRef.current }, (_, index) => {
            // 计算已经上传的分片进度
            if (!doneSet.has(index)) return 0
            const start = index * chunkSize
            return Math.min(chunkSize, file.size - start)
        })
        reportProgress()
    }

    /**
     * 上传切片
     * param uploaded: 已经上传的分片序号
     * return: 上传结果
     */
    const runUpload = async (uploaded: number[]) => {
        const file = fileRef.current
        const uploadId = uploadIdRef.current
        if (!file || !uploadId) return

        setStatus('uploading')
        pausedRef.current = false
        abortRef.current?.abort() // 如果之前有任务没有中断，先中断之前的上传任务
        abortRef.current = new AbortController()
        const signal = abortRef.current.signal

        const chunkSize = chunkSizeRef.current
        const total = totalRef.current
        syncLoadedFromUploaded(uploaded) // 计算当前进度, 正常情况下，刚开始的 uploaded 是[]

        const doneSet = new Set(uploaded ?? []) // 对 uploaded 去重，doneSet 里面存储的是已经上传 分片的 index
        const pending = Array.from({ length: total }, (_, i) => i).filter((i) => !doneSet.has(i)) // 过滤掉已经上传的分片, 避免重复上传

        if (pending.length > 0) {
            const taskArray = pending.map((index) => {
                // 构造需要上传的分片 Promise 数组
                return () => {
                    const startIndex = index * chunkSize // 计算某个分片的开始字节位置
                    const blob = file.slice(startIndex, startIndex + chunkSize) // 从整个文件字节中，截取当前分片需要上传的 blob 字节
                    const queryParams = new URLSearchParams({
                        // 构造请求参数
                        upload_id: uploadId,
                        index: String(index),
                    })

                    return xhrRequestWithRetry({
                        url: `/upload/chunk?${queryParams}`,
                        method: 'POST',
                        headers: { 'Content-Type': 'application/octet-stream' },
                        body: blob,
                        retryNum: 3,
                        retryDelaySecond: 2, // 2s 后重试
                        signal,
                        onProgress: (percent) => {
                            const bytes = (blob.size * percent) / 100 // 计算当前分片已经上传的字节数
                            if (bytes > loadedRef.current[index]) {
                                loadedRef.current[index] = bytes // 更新已经上传的字节数
                                reportProgress() // 更新上传进度
                            }
                        },
                    })
                }
            })

            const result = await limitPromise(taskArray, 3) // 限制并发上传数量为 3
            if (pausedRef.current || signal.aborted) {
                // 如果暂停或被中断，则设置状态为暂停
                setStatus('paused')
                return
            }
            if (result.some((item) => item instanceof Error || item.code !== 0)) {
                // 如果上传中有错误，则设置状态为错误
                setStatus('error')
                return
            }
        }

        setStatus('merging') // 设置状态为正在合并中
        const done = await POST<{ path: string }>('/upload/complete', {
            // 上传完成后，调用 complete 接口，告知服务端可以合并文件，完成上传
            upload_id: uploadId,
        })
        if (done.code !== 0) {
            // 如果合并文件失败，则设置状态为错误
            setStatus('error')
            return
        }
        setOnProgress(100) // 设置上传进度为 100%
        setStatus('done') // 设置状态为完成
    }

    /**
     *
     * 处理文件 Hash 计算,
     * 初始化上传任务，调用 init 接口
     * 如果文件已经上传过，则直接返回，实现秒传
     * 如果文件未上传过，则调用上传 Chunk 接口，上传分片
     *
     * param file: 文件
     * */
    const handleFileChange = async (file: File | undefined) => {
        if (!file) return // 如果文件不存在，则返回
        abortRef.current?.abort()
        fileRef.current = file
        pausedRef.current = false
        setOnProgress(0)
        setStatus('hashing')

        /**
         * 计算文件的 SHA-256 的十六进制字符串
         * 使用 webworker 计算hash ，避免页面卡死
         * 每个分片哈希 安顺序拼接后，再哈希一次，算法需要和后端一致
         */
        const { fileHash } = await computeFileHash(file, CHUNK_SIZE, (value) => {
            setOnProgress(Math.round(value * 100))
        })

        /**
         * 初始化上传
         * 先调用 init 接口，查看当前文件是否已经上传过。
         * 如果已经上传过，则直接返回，实现秒传
         * 如果未上传过，则走后面分片上传逻辑
         */
        const initRes = await POST<InitData>('/upload/init', {
            file_name: file.name,
            file_size: file.size,
            file_hash: fileHash,
        })

        if (initRes.code !== 0) {
            setStatus('error')
            return
        }
        if (initRes.data.finished) {
            setOnProgress(100) // 秒传
            setStatus('done')
            return
        }

        uploadIdRef.current = initRes.data.upload_id
        chunkSizeRef.current = initRes.data.chunk_size
        totalRef.current = initRes.data.chunk_total

        await runUpload(initRes.data.uploaded ?? [])
    }

    /**
     * 暂停上传
     */
    const handlePause = () => {
        if (status !== 'uploading') return
        pausedRef.current = true
        abortRef.current?.abort()
        setStatus('paused')
    }

    /**
     * 暂停后恢复上传
     * */
    const handleResume = async () => {
        if (!uploadIdRef.current || !fileRef.current) return
        // 查询之前某个上传任务的最新状态
        const statusRes = await GET<InitData>(
            `/upload/status?upload_id=${encodeURIComponent(uploadIdRef.current)}`,
        )
        if (statusRes.code !== 0) {
            setStatus('error')
            return
        }

        await runUpload(statusRes.data.uploaded ?? [])
    }

    return {
        handleFileChange,
        onProgress,
        handlePause,
        handleResume,
        status,
    }
}
