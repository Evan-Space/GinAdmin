import { useState } from 'react'

import { getUserNameListOptionsAPI, getUserListAPI } from '@src/request/userList'
import { OPTIONS_ENUM_TYPE, PaginationTypeResponse, PaginationTypeQuery } from '@src/types'
import { useForm } from 'antd/es/form/Form'
import { FieldType } from './types'

export const useUserList = () => {
    const [form] = useForm<FieldType>()
    const [pagination, setPagination] = useState<PaginationTypeResponse>({
        currentPage: 1,
        pageSize: 10,
        total: 0,
    })
    /**
     * 获取用户名称枚举值
     */
    const handleGetUserNameListOptions = async () => {
        const res = await getUserNameListOptionsAPI()
        if (res.code !== 0) {
            throw new Error(res.msg)
        }
        return res.data.map((item: OPTIONS_ENUM_TYPE) => ({
            label: item.label,
            value: item.label,
        }))
    }

    /**
     * 获取用户列表
     */
    const handleGetUserList = async (params: Partial<FieldType> & PaginationTypeQuery) => {
        const res = await getUserListAPI({
            ...params,
        })
        if (res.code !== 0) {
            throw new Error(res.msg)
        }
        setPagination({
            currentPage: res.data.currentPage,
            pageSize: res.data.pageSize,
            total: res.data.total,
        })
        return res.data.list || []
    }

    return {
        form,
        pagination,
        setPagination,
        handleGetUserNameListOptions,
        handleGetUserList,
    }
}
