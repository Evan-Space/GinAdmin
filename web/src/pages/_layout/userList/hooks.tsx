
import { getUserNameListOptionsAPI, getUserListAPI } from '@src/request/userList'
import { OPTIONS_ENUM_TYPE } from '@src/types'

export const useUserList = () => {
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
                value: item.value,
            }))
        } catch (error) {
            return []
        }
    }

    /**
     * 获取用户列表
     */
    const handleGetUserList = async () => {
        try {
            const res = await getUserListAPI()
            if (res.code !== 0) {
                throw new Error(res.msg)
            }
            return res.data.list || []
        } catch (error) {
            return []
        }
    }

    return {
        handleGetUserNameListOptions,
        handleGetUserList
    }
}