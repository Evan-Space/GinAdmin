import { useState } from 'react'
import { xhrRequest } from '@src/request/UploadFile'
import { POST } from '@src/request'

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
        const uploadId = crypto.randomUUID() // 生成一个 uploadId

        try {
            for (let index = 0; index < total; index++) {
                const startIndex = index * CHUNK_SIZE
                const blob = file.slice(startIndex, startIndex + CHUNK_SIZE) // 截取某段 文件
                const queryParams = new URLSearchParams({
                    upload_id: uploadId,
                    index: String(index),
                })
                const res = await xhrRequest({
                    url: `/upload/chunk?${queryParams}`,
                    method: 'POST',
                    headers: { 'Content-Type': 'application/octet-stream' },
                    body: blob,
                    onProgress: (percent) => {
                        const loaded = startIndex + (blob.size * percent) / 100
                        setOnProgress(Math.round((loaded / file.size) * 100))
                    },
                })

                if (res.code !== 0) return
            }

            const done = await POST<{ path: string }>('/upload/complete', {
                upload_id: uploadId,
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
