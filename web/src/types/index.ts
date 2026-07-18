/**
 * 选项枚举类型 用于 Select 组件的 options 属性
 */
export interface OPTIONS_ENUM_TYPE {
    label: string
    value: string
}

/**
 * 分页
 */
export interface PaginationType {
    currentPage: number
    pageSize: number
    totalCount: number
}