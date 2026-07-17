import { useRequest } from 'ahooks'
import { getUserNameListOptionsAPI, getUserListAPI } from '@src/request/userList'

export const useUserList = () => {
    console.log('useUserList')
    // 获取用户名称枚举值
    const resData1 = useRequest(getUserNameListOptionsAPI)

    // 获取用户列表
    const resData2 = useRequest(getUserListAPI)

    // console.log(resData1.data)
    // console.log(resData2.data)

    return {
        // USER_NAME_LIST_OPTIONS,
        // userListData,
        data: {},
    }
}
