import { OPTIONS_ENUM_TYPE } from '@src/types'
import { MemberListItemType } from './types'
import { ColumnsType } from 'antd/es/table'

export const ACCOUNT_STATUS: OPTIONS_ENUM_TYPE[] = [
    {
        label: '启用',
        value: 1,
    },
    {
        label: '禁用',
        value: 0,
    },
]



export const TableColumns: ColumnsType<MemberListItemType> = [

    {
        title: 'ID',
        dataIndex: 'id',
        width: 100,
    },
    {
        title: 'Nickname',
        dataIndex: 'nickname',
        width: 100,
    },
    
    {
        title: 'age',
        dataIndex: 'age',
        width: 100,
    },
    {
        title: 'status',
        dataIndex: 'status',
        width: 100,
    },
]