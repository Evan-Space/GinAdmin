import type { FormRule } from 'antd'
import type { AddUserFormFields } from './types'

export const AddUserDrawerFormRules: Partial<Record<keyof AddUserFormFields, FormRule[]>> = {
    nickname: [
        { required: true, message: '请输入昵称' },
    ],
    userName: [
        { required: true, message: '请输入用户名' },
    ],
    age: [
        { required: true, message: '请输入年龄' },
    ],
    password: [
        { required: true, message: '请输入密码' },
    ],
}