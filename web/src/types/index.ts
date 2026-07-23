/**
 * 选项枚举类型 用于 Select 组件的 options 属性
 */
export interface OPTIONS_ENUM_TYPE {
    label: string
    value: string
}

/**
 * 请求中的分页参数，不带 total
 */
export interface PaginationTypeQuery {
    currentPage: number
    pageSize: number
}


/**
 * 接口响应中的 分页数据，带 total
 */
export interface PaginationTypeResponse extends PaginationTypeQuery {
    total: number
}