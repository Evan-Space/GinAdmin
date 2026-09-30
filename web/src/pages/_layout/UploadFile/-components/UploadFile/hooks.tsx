import { useRef, useState } from 'react'
import { Upload as UploadAPI } from '@src/request/Upload.ts'

export const useCustomUploadFile = () => {
    const filesRef = useRef<File[]>([])
    const [progress, setProgress] = useState<number>(0)

    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (files: FileList | []) => {
        if (files.length === 0) return

        const form = new FormData()
        form.append('file', files[0])
        setProgress(0)
        await UploadAPI('/uploadFile', form, (val) => {
            setProgress(val)
        })
    }

    return {
        handleFileChange,
        progress, // 上传进度条
    }
}
