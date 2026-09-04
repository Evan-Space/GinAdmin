import { UserListItemType } from './types'
import { ColumnsType } from 'antd/es/table'
import { Tag } from 'antd'
import { Space, Button, Modal } from 'antd'

export const getTableColumns = (
    handleDeleteAccount: (params: { id: number; type: '0' | '1' }) => void
): ColumnsType<UserListItemType> => {
    return [
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
        },
        {
            title: "操作",
            dataIndex: "action",
            width: 120,
            render: (_: any, record: UserListItemType) => {
                return (
                    <Space>
                        <Button type="text" danger onClick={() => {
                            Modal.confirm({
                                title: "确认删除",
                                content: `确定要删除账号 "${record.name}" 吗？此操作不可恢复。`,
                                onOk: () => {
                                    handleDeleteAccount({
                                        id: Number(record.id),
                                        type: '1', // 1: admin 0: 普通成员
                                    })
                                }
                            })
                        }}>删除</Button>
                    </Space>
                )
            }
        }
    ]
}



/**
 * 账户状态标签展示
*/
export const StatusContainer = ({ status }: { status: number }) => {
    return (
        <Tag color={status === 1 ? 'green' : 'red'}>{status === 1 ? '启用' : '禁用'}</Tag>
    )
}