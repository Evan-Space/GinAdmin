import { useState } from 'react'
import { xhrRequest } from '@src/request/UploadFile'
import { POST } from '@src/request'
import { limitPromise } from '@src/utils/utils.ts'

const CHUNK_SIZE = 5 * 1024 * 1024 // 分片大小

export const useCustomUploadFile = () => {
    const [onProgress, setOnProgress] = useState<number>(0)

    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (file: File | undefined) => {
        if (!file) return

        const blob = file.slice(0, file.size, file.type || 'application/octet-stream')

        setOnProgress(0)
        const total = Math.ceil(file.size / CHUNK_SIZE)
        const initRes = await POST<{ upload_id: string; chunk_size: number; chunk_total: number }>(
            '/upload/init',
            {
                file_name: file.name,
                file_size: file.size,
            },
        )
        if (initRes.code !== 0) return

        
        const taskArray = Array.from({ length: total }, (_, index) => {
            return () => {
                const startIndex = index * CHUNK_SIZE
                const blob = file.slice(startIndex, startIndex + CHUNK_SIZE) // 截取某段 文件
                const queryParams = new URLSearchParams({
                    upload_id: initRes.data.upload_id,
                    index: String(index),
                })
                return xhrRequest({
                    url: `/upload/chunk?${queryParams}`,
                    method: 'POST',
                    headers: { 'Content-Type': 'application/octet-stream' },
                    body: blob,
                    onProgress: (percent) => {
                        const loaded = startIndex + (blob.size * percent) / 100
                        setOnProgress(Math.round((loaded / file.size) * 100))
                    },
                })
            }
        })

        try {
            const result = await limitPromise(taskArray, 3)
            if (result.some((item) => item instanceof Error || item.code !== 0)) return
            const done = await POST<{ path: string }>('/upload/complete', {
                upload_id: initRes.data.upload_id,
                file_name: file.name,
                chunk_total: total,
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
