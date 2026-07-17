
export interface fetchResponse {
    status: number;
    msg: string;
}


const BASE_URL = 'http://localhost:8080/api/v1';
const TOKEN = localStorage.getItem('auth_token') || ''

const GET = async <T>(url: string, config?: RequestInit): Promise<fetchResponse & T> => {
    return fetch(`${BASE_URL}${url}`, {
        method: 'get',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${TOKEN}`
        },
        cache: "no-store",
        ...config
    })
        .then(async (responseData) => {
            const result = await responseData.json()
            if (result.status === 401) {
                localStorage.removeItem('auth_token')
                window.location.href = '/login'
            }
            return result
        })
        .catch(err => {
            return new Error(err)
        })
}


const POST = async <T>(url: string, data: any, config?: RequestInit): Promise<fetchResponse & T> => {
    return fetch(`${BASE_URL}${url}`, {
        method: 'post',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${TOKEN}`
        },
        cache: "no-store",
        body: JSON.stringify(data),
        ...config
    })
        .then(async (responseData) => {
            const result = await responseData.json()
            if (result.status === 401) {
                localStorage.removeItem('auth_token')
                window.location.href = '/login'
            }
            return result
        })
        .catch(err => {
            return new Error(err)
        })
}

// // 流式请求
// const FetchStream = async (url: string, data: any, config?: RequestInit) => {
//     const baseUrl = await GetApiBaseUrl();
//     const authToken = localStorage.getItem('auth_token')

//     return fetch(`${baseUrl}${url}`, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json',
//             'Authorization': `Bearer ${authToken}`
//         },
//         cache: "no-store",
//         body: JSON.stringify(data),
//         ...config
//     })
// }


export {
    GET,
    POST,
    // FetchStream
}
