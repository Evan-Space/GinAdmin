import { useState } from 'react'

import { getUserNameListOptionsAPI, getUserListAPI } from '@src/request/userList'
import { OPTIONS_ENUM_TYPE } from '@src/types'
import { useForm } from 'antd/es/form/Form'
import { FieldType } from './types'

import { PaginationType } from '@src/types'


export const useUserList = () => {
    const [form] = useForm<FieldType>()
    const [paginations, setPaginations] = useState<PaginationType>({
        currentPage: 1,
        pageSize: 10,
        totalCount: 0
    })
    /**
     * 获取用户名称枚举值
     */
    const handleGetUserNameListOptions = async () => {
        try {
            const res = await getUserNameListOptionsAPI()
            if (res.code !== 0) {
                throw new Error(res.msg)
            }
            return res.data.map((item: OPTIONS_ENUM_TYPE) => ({
                label: item.label,
                value: item.label,
            }))
        } catch (error) {
            return []
        }
    }

    /**
     * 获取用户列表
     */
    const handleGetUserList = async (params: Partial<FieldType>) => {
        try {
            const res = await getUserListAPI({
                ...params,
                ...paginations,
            })
            if (res.code !== 0) {
                throw new Error(res.msg)
            }
            return res.data.list || []
        } catch (error) {
            return []
        }
    }

    return {
        form,
        handleGetUserNameListOptions,
        handleGetUserList,
    }
}