import { GET, POST } from '@src/request/request'
import { OPTIONS_ENUM_TYPE, PaginationType } from '@src/types'
import { UserListItemType } from '@src/pages/_layout/userList/types'
import { FieldType as UserListFormType } from '@src/pages/_layout/userList/types'



/**
 * 获取用户名称枚举值
 * */ 
export const getUserNameListOptionsAPI = async () => {
    return GET<OPTIONS_ENUM_TYPE[]>('/admin-user/userNameOptions')
}


/**
 * 获取用户列表
*/
export const getUserListAPI = async (params: Partial<UserListFormType> & PaginationType) => {
    return POST<{ list: UserListItemType[] } & PaginationType>('/admin-user/list', params)
}