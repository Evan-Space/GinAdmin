import { useRef, useState } from 'react'
import { Upload as UploadAPI } from '@src/request/Upload.ts'
import { POST } from '@src/request'

const CHUNK_SIZE = 5 * 1024 * 1024

export const useCustomUploadFile = () => {
    const filesRef = useRef<File[]>([])
    const [progress, setProgress] = useState<number>(0)

    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (files: FileList | []) => {
        const file = files[0]
        if (!file) return
        setProgress(0)
        const initRes = await POST<{
            upload_id: string
            chunk_total: number
        }>('/upload/init', {
            file_name: file.name,
            file_size: file.size,
        })

        if (initRes.code !== 0) return
        const { upload_id, chunk_total } = initRes.data
        for (let i = 0; i < chunk_total; i++) {
            const blob = file.slice(i * CHUNK_SIZE, (i + 1) * CHUNK_SIZE)
            const chunkRes = await UploadAPI(
                `/chunk?upload_id=${upload_id}&index=${i}`,
                blob,
                (val) => {
                    setProgress(Math.round(((i + val / 100) / chunk_total) * 100))
                },
            )
            if (chunkRes.code !== 0) return
        }

        const done = await UploadAPI<{ path: string }>(`/complete?upload_id=${upload_id}`)
        if (done.code !== 0) return
        setProgress(100)

        // if (files.length === 0) return
        // const file = files[0] // 先只上传一个
        // const total = Math.ceil(file.size / CHUNK_SIZE)
        // setProgress(0)
        //
        // for (let i = 0; i < total; i++) {
        //     const blob = file.slice(i * CHUNK_SIZE, (i + 1) * CHUNK_SIZE)
        //     await UploadAPI(`/chunk?index=${i}&total=${total}`, blob, (val) => {
        //         setProgress(Math.round((i + val / 100) / total) * 100)
        //     })
        // }
        //
        // await UploadAPI(`/merge?total=${total}&filename=${encodeURIComponent(file.name)}`)
        // setProgress(100)
        // const form = new FormData()
        // form.append('file', files[0])
        // setProgress(0)
        // await UploadAPI('/uploadFile', form, (val) => {
        //     setProgress(val)
        // })
    }

    return {
        handleFileChange,
        progress, // 上传进度条
    }
}
