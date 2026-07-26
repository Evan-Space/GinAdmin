export interface AddUserFormFields {
    nickname: string
    userName: string
    age: number
    email: string
    password: string
    phone: string
    address: string
    gender: string
    birthday: string
    role: string
    status: StatusEnum // 0: 禁用 1: 启用
}

/**
 * 账号是否启用的状态枚举
 */
export interface StatusEnum {
    DISABLED: 0
    ENABLED: 1
}
