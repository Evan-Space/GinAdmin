import { useState } from 'react'
import { xhrRequest } from '@src/request/UploadFile'

export const useCustomUploadFile = () => {
    const [onProgress, setOnProgress] = useState<number>(0)

    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (file: File | undefined) => {
        if (!file) return

        // const form = new FormData()
        // form.append('file', file)
        const blob = file.slice(0, file.size, file.type || "application/octet-stream")

        setOnProgress(0)
        try {
            const res = await xhrRequest({
                url: '/upload/uploadFile',
                method: 'POST',
                body: blob,
                onProgress: (val) => setOnProgress(val),
                headers: {
                    'X-File-Name': encodeURIComponent(file.name),
                },
            })
            if (res.code !== 0) return
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
