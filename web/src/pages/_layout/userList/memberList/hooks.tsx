import { useState } from 'react'
import { useForm } from 'antd/es/form/Form'
import { OPTIONS_ENUM_TYPE } from '@src/types'
import { getMemberUserNameListOptionsAPI, getMemberListAPI } from '@src/request/userList'
import { useRequest } from 'ahooks'
import { FieldType } from './types'
import { PaginationTypeResponse, PaginationTypeQuery } from '@src/types'
import { omitEmptyValues } from '@src/utils/utils'
import { deleteAccountAPI } from '@src/request/userList'

export const useMemberList = () => {
    const [form] = useForm()
    const [pagination, setPagination] = useState<PaginationTypeResponse>({
        currentPage: 1,
        pageSize: 10,
        total: 0,
    })

    /**
     * 获取普通成员用户名
     */
    const handleGetUserNameListOptions = async () => {
        const res = await getMemberUserNameListOptionsAPI()
        if (res.code !== 0) {
            throw new Error(res.msg)
        }
        return res.data.map((item: OPTIONS_ENUM_TYPE) => ({
            label: item.label,
            value: item.label,
        }))
    }

    const { data: USER_NAME_LIST_OPTIONS = [] } = useRequest(handleGetUserNameListOptions)

    const { run: runGetMemberList, data: { data: { list: member_list_data = [] } = {} } = {} } =
        useRequest(getMemberListAPI, {
            defaultParams: [{ currentPage: 1, pageSize: 10 }],
            onSuccess: (res) => {
                setPagination({
                    currentPage: res.data.currentPage,
                    pageSize: res.data.pageSize,
                    total: res.data.total,
                })
            },
        })

    // 获取成员列表,搜索方法
    const handleSearch = async (params: Partial<FieldType> & PaginationTypeQuery) => {
        const paramsObj = {
            ...omitEmptyValues(form.getFieldsValue()),
            ...pagination,
            ...params,
        }

        runGetMemberList(paramsObj)
    }

    /**
     * 删除用户
     */
    const handleDeleteAccount = async (params: { id: number }) => {
        const res = await deleteAccountAPI({
            id: params.id,
            type: '0',
        })
        if (res.code !== 0) {
            throw new Error(res.msg)
        }
        runGetMemberList({ currentPage: 1, pageSize: 10 })
    }
    return {
        form,
        USER_NAME_LIST_OPTIONS,
        handleSearch,
        member_list_data,
        pagination,
        handleDeleteAccount,
    }
}
