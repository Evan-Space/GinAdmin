import { UserListItemType } from './types'
import { ColumnsType } from 'antd/es/table'
import { Tag } from 'antd'
import { Space, Button, Modal } from 'antd'

export const getTableColumns = (
    handleDeleteAccount: (params: { id: number; type: '0' | '1' }) => void,
    handleChangeStatus: (id: number) => void,
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
                return (
                    <Tag color={text === 1 ? 'green' : 'red'}>{text === 1 ? '启用' : '禁用'}</Tag>
                )
            },
        },
        {
            title: '操作',
            dataIndex: 'action',
            width: 120,
            render: (_: any, record: UserListItemType) => {
                return (
                    <Space>
                        <Button
                            type="text"
                            danger
                            onClick={() => {
                                Modal.confirm({
                                    title: '确认删除',
                                    content: `确定要删除账号 "${record.name}" 吗？此操作不可恢复。`,
                                    onOk: () => {
                                        handleDeleteAccount({
                                            id: Number(record.id),
                                            type: '1', // 1: admin 0: 普通成员
                                        })
                                    },
                                })
                            }}
                        >
                            删除
                        </Button>

                        <Button
                            color="purple"
                            variant="text"
                            onClick={() => {
                                handleChangeStatus(Number(record.id))
                            }}
                        >
                            {record.status === 1 ? '禁用' : '启用'}
                        </Button>

                        
                    </Space>
                )
            },
        },
    ]
}
