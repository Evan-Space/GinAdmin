import { OPTIONS_ENUM_TYPE } from '@src/types'
import { MemberListItemType } from './types'
import { ColumnsType } from 'antd/es/table'
import { Space, Button, Modal } from 'antd'

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


export const getTableColumns = (
    handleDeleteAccount: (params: { id: number; type: '0' | '1' }) => void
): ColumnsType<MemberListItemType> => {
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
            title: 'age',
            dataIndex: 'age',
            width: 100,
        },
        {
            title: 'status',
            dataIndex: 'status',
            width: 100,
        },
        {
            title: "操作",
            dataIndex: "action",
            width: 100,
            render: (_: any, record: MemberListItemType) => {
                return (
                    <Space>
                    <Button type="text" danger onClick={() => {
                        Modal.confirm({
                            title: "确认删除",
                            content: `确定要删除账号 "${record.nickname}" 吗？此操作不可恢复。`,
                            onOk: () => {
                                handleDeleteAccount({
                                    id: record.id,
                                    type: '0', // 0: 普通成员
                                })
                            }
                        })
                    }}>
                       删除</Button>
                    </Space>
                )
            }
        }
    ]
}