import { useRef, useState } from 'react'
import { Upload as UploadAPI } from '@src/request/Upload.ts'


const CHUNK_SIZE = 5 * 1024 * 1024

export const useCustomUploadFile = () => {
    const filesRef = useRef<File[]>([])
    const [progress, setProgress] = useState<number>(0)

    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (files: FileList | []) => {
        if (files.length === 0) return
        const file = files[0] // 先只上传一个
        const total = Math.ceil(file.size / CHUNK_SIZE)
        setProgress(0)

        for (let i = 0; i < total; i++) {
            const blob = file.slice(i * CHUNK_SIZE, (i + 1) * CHUNK_SIZE)
            await UploadAPI(`/chunk?index=${i}&total=${total}`, blob, (val) => {
                setProgress(Math.round((i + val / 100) / total) * 100)
            })
        }

        await UploadAPI(`/merge?total=${total}&filename=${encodeURIComponent(file.name)}`)
        setProgress(100)

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
