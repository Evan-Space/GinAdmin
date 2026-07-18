const BASE_URL = 'http://localhost:8080/api/v1'
const TOKEN_KEY = 'auth_token'

// 不需要 token 的接口
const PUBLIC_URLS = ['/login']

export interface fetchResponse<T> {
    code: number
    msg: string
    cost: number
    request_id: string
    data: T
}

export const getToken = () => {
    return localStorage.getItem(TOKEN_KEY) || ''
}
export const redirectToLogin = () => {
    localStorage.removeItem(TOKEN_KEY)
    window.location.href = '/login'
}

const GET = async <T>(url: string, config?: RequestInit): Promise<fetchResponse<T>> => {
    const isPublic = PUBLIC_URLS.some((path) => url.startsWith(path))
    const token = getToken()

    // 无 token 且不是公开接口 → 跳转登录
    if (!token && !isPublic) {
        redirectToLogin()
        return Promise.reject(new Error('未登录'))
    }

    return fetch(`${BASE_URL}${url}`, {
        method: 'get',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
        cache: 'no-store',
        ...config,
    })
        .then(async (responseData) => {
            const result = await responseData.json()
            if (result.code === 401) {
                redirectToLogin()
            }
            return result
        })
        .catch((err) => {
            return new Error(err)
        })
}

const POST = async <T>(url: string, data: any, config?: RequestInit): Promise<fetchResponse<T>> => {
    const isPublic = PUBLIC_URLS.some((path) => url.startsWith(path))
    const token = getToken()

    if (!token && !isPublic) {
        redirectToLogin()
        return Promise.reject(new Error('未登录'))
    }

    return fetch(`${BASE_URL}${url}`, {
        method: 'post',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
        cache: 'no-store',
        body: JSON.stringify(data),
        ...config,
    })
        .then(async (responseData) => {
            const result = await responseData.json()
            if (result.code === 401) {
                redirectToLogin()
            }
            return result
        })
        .catch((err) => {
            return new Error(err)
        })
}

export { GET, POST }
