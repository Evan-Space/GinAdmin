import { UserListItemType } from './types'
import { ColumnsType } from 'antd/es/table'


export const TableColumns: ColumnsType<UserListItemType> = [
    {
        title: 'ID',
        dataIndex: 'id',
        align: 'center',
    },
    {
        title: 'Name',
        dataIndex: 'nickname',
        align: 'center',
    },
    {
        title: 'username',
        dataIndex: 'username',
        align: 'center',
    },
    {
        title: 'address',
        dataIndex: 'address',
        key: 'address',
        align: 'center',
        render: (text: string) => {
            return text ? text : '-'
        },
    },
    {
        title: 'email',
        dataIndex: 'email',
        align: 'center',
        render: (text: string) => {
            return text ? text : '-'
        },
    }
]