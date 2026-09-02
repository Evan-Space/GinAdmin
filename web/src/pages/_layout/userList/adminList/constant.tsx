import { UserListItemType } from './types'
import { ColumnsType } from 'antd/es/table'
import { Tag } from 'antd'


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
    },
    {
        title: 'status',
        dataIndex: 'status',
        width: 100,
        render: (text: number) => {
            return StatusContainer({ status: text })
        },
       
    }
]


/**
 * 账户状态标签展示
*/
export const StatusContainer = ({ status }: { status: number }) => {
    return (
        <Tag color={status === 1 ? 'green' : 'red'}>{status === 1 ? '启用' : '禁用'}</Tag>
    )
}