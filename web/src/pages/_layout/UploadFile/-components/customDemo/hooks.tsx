import { useRef } from 'react'
import { getToken } from '@src/request/request'

export const useCustomUploadFile = () => {
    const filesRef = useRef<File[]>([])


    /**
     * 拿到上传文件
     * */
    const handleFileChange = async (files: FileList | []) => {
        if (files.length === 0) return

        const form = new FormData()
        form.append('file', files[0])

        await fetch('http://localhost:8080/api/v1/upload/uploadFile', {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
            method: 'POST',
            body: form,
        })

    }

    return {
        handleFileChange,
    }
}
