import { GET } from '@src/request/request'
import { OPTIONS_ENUM_TYPE } from '@src/types'



/**
 * 获取用户名称枚举值
 * */ 
export const getUserNameListOptionsAPI = async () => {
    return GET<OPTIONS_ENUM_TYPE[]>('/admin-user/userNameOptions')
}


/**
 * 获取用户列表
*/
export const getUserListAPI = async () => {
    return GET('/admin-user/list')
}