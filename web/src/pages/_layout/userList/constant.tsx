import { UserListItemType } from './types'
import { ColumnsType } from 'antd/es/table'


export const TableColumns: ColumnsType<UserListItemType> = [
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
        title: 'username',
        dataIndex: 'username',
        width: 100,
    },
    {
        title: 'age',
        dataIndex: 'age',
        width: 100,
    },
    {
        title: 'address',
        dataIndex: 'address',
        key: 'address',
        render: (text: string) => {
            return text ? text : '-'
        },
        width: 100,
    },
    {
        title: 'email',
        dataIndex: 'email',
        render: (text: string) => {
            return text ? text : '-'
        },
        width: 100,
    }
]