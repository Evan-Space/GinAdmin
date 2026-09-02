import { useForm } from 'antd/es/form/Form'
import { OPTIONS_ENUM_TYPE } from '@src/types'
import { getMemberUserNameListOptionsAPI } from '@src/request/userList'
import { useRequest } from 'ahooks'

export const useMemberList = () => {
    const [form] = useForm()

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

    return {
        form,
        USER_NAME_LIST_OPTIONS,
    }
}
