import { GET } from '@src/request/request'



/**
 * 获取用户名称枚举值
 * */ 
export const getUserNameListOptionsAPI = async () => {
    console.log('getUserNameListOptionsAPI')
    const res = await GET('/admin-user/userNameOptions')
    return res
}


/**
 * 获取用户列表
*/
export const getUserListAPI = async () => {
    const res = await GET('/admin-user/list')
    return res
}