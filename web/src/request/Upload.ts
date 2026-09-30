import { fetchResponse, getToken, redirectToLogin } from '@src/request/request.ts'

const UPLOAD_BASE_URL = 'http://localhost:8080/api/v1/upload'

export const Upload = <T>(
    url: string,
    body?: XMLHttpRequestBodyInit | null,
    onProgress?: (percent: number) => void,
) => {
    const token = getToken()
    if (!token) {
        redirectToLogin()
        return Promise.reject(new Error('未登陆'))
    }
    return new Promise<fetchResponse<T>>((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('POST', `${UPLOAD_BASE_URL}${url}`)
        xhr.setRequestHeader('Authorization', `Bearer ${token}`)

        xhr.upload.onprogress = (event) => {
            if (!event.lengthComputable || !onProgress) return
            onProgress(Math.round((event.loaded / event.total) * 100))
        }

        xhr.onload = () => {
            const result = JSON.parse(xhr.responseText)
            if (result.code === 401) {
                redirectToLogin()
                return
            }
            resolve(result)
        }

        xhr.onerror = () => reject(new Error('上传失败，请稍后重试'))
        xhr.send(body)
    })
}
