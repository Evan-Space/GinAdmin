import { POST } from '@src/request/request'
import { fetchResponse } from '../request'

// 登录接口
export const loginAPI = async (data: any): Promise<fetchResponse<{ access_token: string }>> => {
    return POST('/login', data)
}