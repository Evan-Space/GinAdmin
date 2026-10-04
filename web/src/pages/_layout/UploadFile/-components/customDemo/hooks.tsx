import { useState } from 'react'
import { xhrRequest } from '@src/request/UploadFile'
import { POST } from '@src/request'
import { limitPromise } from '@src/utils/utils.ts'
import { computeFileHash } from '@src/pages/_layout/UploadFile/-components/customDemo/utils/fileHash.ts'

const CHUNK_SIZE = 5 * 1024 * 1024

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

    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (file: File | undefined) => {
        if (!file) return // 如果文件不存在，则返回
        setOnProgress(0) // 设置上传进度为 0

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

        if (initRes.code !== 0) return
        if (initRes.data.finished) {
            // 秒传
            setOnProgress(100)
            return
        }

        // 如果文件没有上传过，则走后面分片上传逻辑
        const {
            upload_id: uploadId, // 上传 ID，由服务端生成，用于标识该文件上传任务
            chunk_size: chunkSize, // 分片规格大小
            chunk_total: total, // 真个文件被分片总数
            uploaded, // 服务端查询到的已经上传的分片序号
        } = initRes.data

        // 声明已经上传的分片，利用 Set 去重
        const doneSet = new Set(uploaded ?? [])


        /**
         * 构造需要上传的分片数组
         * 已经上传的分片，不重复上传
         * 未上传的分片，需要上传
         * loaded 数组，下标为分片序号，值为已经上传的字节数
        */
        const loaded = Array.from({ length: total }, (_, index) => {
            if (!doneSet.has(index)) return 0 // 如果某分片已经上传过，则本次不用再传
            const start = index * chunkSize // 计算某个分片的开始字节位置
            return Math.min(chunkSize, file.size - start) // 计算某个分片的结束字节位置
        })

        const report = () => {
            const sum = loaded.reduce((a, b) => a + b, 0) // 计算已经上传的总的字节数
            setOnProgress(Math.min(99, Math.round((sum / file.size) * 100))) // 设置上传进度
        }

        report()

        const taskArray = Array.from({ length: total }, (_, index) => index) // 构造需要上传的分片数组
            .filter((item) => !doneSet.has(item)) // 过滤掉已经上传的分片
            .map((index) => {
                return () => {
                    const startIndex = index * chunkSize // 计算某个分片的开始字节位置
                    const blob = file.slice(startIndex, startIndex + chunkSize) // 从整个文件字节中，截取当前分片需要上传的 blob 字节
                    const queryParams = new URLSearchParams({
                        upload_id: uploadId,
                        index: String(index), // 当前分片的序号
                    })
                    return xhrRequest({
                        url: `/upload/chunk?${queryParams}`,
                        method: 'POST',
                        headers: { 'Content-Type': 'application/octet-stream' },
                        body: blob,
                        onProgress: (percent) => {
                            const bytes = (blob.size * percent) / 100 // 计算当前分片已经上传的字节数
                            if (bytes > loaded[index]) {
                                loaded[index] = bytes // 更新已经上传的字节数
                                report() // 更新上传进度
                            }
                        },
                    })
                }
            })

        try {
            const result = await limitPromise(taskArray, 3) // 限制并发上传数量为 3
            if (result.some((item) => item instanceof Error || item.code !== 0)) return // 如果上传失败，则返回
            const done = await POST<{ path: string }>('/upload/complete', { // 调用 complete 接口，告知服务端可以合并文件，完成上传
                upload_id: initRes.data.upload_id,
            })

            if (done.code !== 0) return // 如果合并文件失败，则返回
            setOnProgress(100)
        } catch (err) {
            console.error(err)
        }
    }

    return {
        handleFileChange,
        onProgress,
    }
}
