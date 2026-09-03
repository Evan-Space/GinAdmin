import { useState } from 'react'

import { getUserNameListOptionsAPI, getUserListAPI } from '@src/request/userList'
import { OPTIONS_ENUM_TYPE, PaginationTypeResponse, PaginationTypeQuery } from '@src/types'
import { useForm } from 'antd/es/form/Form'
import { FieldType } from './types'
import { useRequest } from 'ahooks'
import { omitEmptyValues } from '@src/utils/utils'

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

    const { data: USER_NAME_LIST_OPTIONS = [] } = useRequest(handleGetUserNameListOptions)

    /**
     * 获取用户列表
     */
    const { data: { data: { list: user_list_data = [] } = {} } = {}, run: runGetUserList } = useRequest(getUserListAPI, {
        defaultParams: [{ currentPage: 1, pageSize: 10 }],
        onSuccess: (res) => {
            setPagination({
                currentPage: res.data.currentPage,
                pageSize: res.data.pageSize,
                total: res.data.total,
            })
        },
    })

    /**
     * 搜索
     */
    const handleSearch = (params: Partial<FieldType> & PaginationTypeQuery) => {
        const apiParams: Partial<FieldType> & PaginationTypeQuery = {
            ...omitEmptyValues({
                ...form.getFieldsValue(),
            }),
            ...pagination,
            ...params,
        }
        runGetUserList(apiParams)
    }

    return {
        form,
        pagination,
        USER_NAME_LIST_OPTIONS,
        user_list_data,
        handleSearch,
    }
}
