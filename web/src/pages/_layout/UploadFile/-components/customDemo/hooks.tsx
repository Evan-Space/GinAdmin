import { useState } from 'react'
import { xhrRequest } from '@src/request/UploadFile'
import { POST} from '@src/request'
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
        if (!file) return
        setOnProgress(0)

        const { fileHash } = await computeFileHash(file, CHUNK_SIZE, (value) => {
            setOnProgress(Math.round(value * 100))
        })


        // const total = Math.ceil(file.size / CHUNK_SIZE)
        const initRes = await POST<InitData>('/upload/init', {
            file_name: file.name,
            file_size: file.size,
            file_hash: fileHash,
        })
        if (initRes.code !== 0) return
        if (initRes.data.finished) { // 秒传
            setOnProgress(100)
            return
        }

        const { upload_id: uploadId, chunk_size: chunkSize, chunk_total: total, uploaded } = initRes.data
        const doneSet = new Set(uploaded ?? [])

        const loaded = Array.from({ length:total }, (_, index) => {
            if (!doneSet.has(index)) return 0
            const start = index * chunkSize
            return Math.min(chunkSize, file.size - start)
        })

        const report = () => {
            const sum = loaded.reduce((a, b) => a + b, 0)
            setOnProgress(Math.min(99, Math.round((sum / file.size) * 100)))
        }

        report()

        
        const taskArray = Array.from({ length: total }, (_, index) => index)
            .filter(item => !doneSet.has(item))
            .map((index) => {
                return () => {
                    const startIndex = index * chunkSize
                    const blob = file.slice(startIndex, startIndex + chunkSize) // 截取某段 文件
                    const queryParams = new URLSearchParams({
                        upload_id: uploadId,
                        index: String(index),
                    })
                    return xhrRequest({
                        url: `/upload/chunk?${queryParams}`,
                        method: 'POST',
                        headers: { 'Content-Type': 'application/octet-stream' },
                        body: blob,
                        onProgress: (percent) => {
                            const bytes = (blob.size * percent) / 100
                            if (bytes > loaded[index]) {
                                loaded[index] = bytes
                                report()
                            }
                            // const loaded = startIndex + (blob.size * percent) / 100
                            // setOnProgress(Math.round((loaded / file.size) * 100))
                        },
                    })
                }
            })

        try {
            const result = await limitPromise(taskArray, 3)
            if (result.some((item) => item instanceof Error || item.code !== 0)) return
            const done = await POST<{ path: string }>('/upload/complete', {
                upload_id: initRes.data.upload_id,
            })

            if (done.code !== 0) return
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
