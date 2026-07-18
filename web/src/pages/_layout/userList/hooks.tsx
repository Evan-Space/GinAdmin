import { useRequest } from 'ahooks'
import { getUserNameListOptionsAPI, getUserListAPI } from '@src/request/userList'
import { OPTIONS_ENUM_TYPE } from '@src/types'



// 获取用户名称枚举值
const handleGetUserNameListOptions = async () => {
    try {
        const res = await getUserNameListOptionsAPI()
        if (res.code !== 0) {
            throw new Error(res.msg)
        }
        return res.data.map((item: OPTIONS_ENUM_TYPE) => ({
            label: item.label,
            value: item.value,
        }))
    } catch (error) {
        return []
    }
}






export const useUserList = () => {
    // 获取用户名称枚举值
    const { data: USER_NAME_LIST_OPTIONS = [] } = useRequest(handleGetUserNameListOptions)

    // 获取用户列表
    const { data: user_list_data = [] } = useRequest(getUserListAPI, {
        manual: true,
    })

    return {
        USER_NAME_LIST_OPTIONS,
        user_list_data,
    }
}



