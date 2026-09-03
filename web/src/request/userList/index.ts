import { GET, POST } from '@src/request/request'
import { OPTIONS_ENUM_TYPE, PaginationTypeQuery, PaginationTypeResponse } from '@src/types'
import { UserListItemType } from '@src/pages/_layout/userList/adminList/types'
import { FieldType as UserListFormType } from '@src/pages/_layout/userList/adminList/types'
import { FieldType as MemberListFormType } from '@src/pages/_layout/userList/memberList/types'
import { MemberListItemType } from '@src/pages/_layout/userList/memberList/types'



/**
 * 获取用户名称枚举值
 * */ 
export const getUserNameListOptionsAPI = async () => {
    return GET<OPTIONS_ENUM_TYPE[]>('/admin-user/userNameOptions')
}


/**
 * 获取用户列表
*/
export const getUserListAPI = async (params: Partial<UserListFormType> & PaginationTypeQuery) => {
    return POST<{ list: UserListItemType[] } & PaginationTypeResponse>('/admin-user/adminList', params)
}

/**
 * 新增用户
*/
export const createUserAPI = async (params: UserListFormType) => {
    return POST('/admin-user/create', params)
}

/**
 * 获取普通成员用户名枚举值
*/
export const getMemberUserNameListOptionsAPI = async () => {
    return GET<OPTIONS_ENUM_TYPE[]>('/admin-user/memberListName')
}
/*
获取普通成员列表
*/
export const getMemberListAPI = async (params: Partial<MemberListFormType> & PaginationTypeQuery) => {
    return POST<{list: MemberListItemType[] } & PaginationTypeResponse>('/admin-user/memberList', params)
}