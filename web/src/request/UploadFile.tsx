import { fetchResponse, getToken, redirectToLogin } from '@src/request/request.ts'

const BASE_URL = 'http://localhost:8080/api/v1'
export interface XHRRequestOptions {
    url: string
    method: 'POST' | 'GET' | 'PUT' | 'DELETE'
    body?: XMLHttpRequestBodyInit | null
    headers?: Record<string, string>
    onProgress?: (percent: number) => void
    signal?: AbortSignal
}

/**
 * 正式的 XMLHTTPRequest 封装
 *
 * XMLHTTP 可以通过 onPregressCallback 拿到上传进度
 * */
export function xhrRequest<T>(options: XHRRequestOptions): Promise<fetchResponse<T>> {
    const { url, method = 'POST', body = null, headers = {}, onProgress, signal } = options
    const token = getToken()

    return new Promise<fetchResponse<T>>((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open(method, `${BASE_URL}${url}`)
        if (!token) {
            redirectToLogin()
            return
        }

        xhr.setRequestHeader('Authorization', `Bearer ${token}`)
        Object.entries(headers).forEach(([key, value]) => xhr.setRequestHeader(key, value))

        if (onProgress) {
            xhr.upload.onprogress = (event) => {
                if (!event.lengthComputable) return
                onProgress(Math.round((event.loaded / event.total) * 100)) // Math.round 把一个数四舍五入成最近的一个整数
            }
        }

        xhr.onload = () => {
            let result: fetchResponse<T>
            try {
                result = JSON.parse(xhr.responseText) as fetchResponse<T>
            } catch {
                reject(new Error('响应解析失败'))
                return
            }
            if (result.code === 401) {
                redirectToLogin()
                return
            }
            resolve(result)
        }

        xhr.onerror = () => reject(new Error('网络错误，请求失败'))
        xhr.ontimeout = () => reject(new Error('请求超时'))
        xhr.onabort = () => reject(new Error('请求已取消'))

        if (signal) {
            signal.addEventListener('abort', () => xhr.abort(), { once: true })
        }

        xhr.send(body)
    })
}
