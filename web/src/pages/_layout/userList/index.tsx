import { createFileRoute } from '@tanstack/react-router'
import { Table, Form, Select, Input, Space } from 'antd'
import { TableColumns } from './constant'
import { FieldType } from './types'
import { useUserList } from './hooks'
import { useRequest } from 'ahooks'
import { UserListItemType } from './types'



export const Route = createFileRoute('/_layout/userList/')({
    component: RouteComponent,
})




function RouteComponent() {

    const { handleGetUserNameListOptions, handleGetUserList } = useUserList()

    const { data: USER_NAME_LIST_OPTIONS = [] } = useRequest(handleGetUserNameListOptions)
    const { data: user_list_data = [] } = useRequest(handleGetUserList)


    return (
        <Space orientation="vertical" size="medium" style={{ display: 'flex' }}>
            <Form
                layout="inline"
                wrapperCol={{ span: 8, style: { minWidth: '200px' } }}
            >
                <Form.Item<FieldType>
                    label="姓名"
                    name="name"
                >
                    <Select
                        options={USER_NAME_LIST_OPTIONS}
                    />
                </Form.Item>

                <Form.Item<FieldType>
                    label="年龄"
                    name="age"
                >
                    <Input />
                </Form.Item>
            </Form>
            <div className="my-12 bg-[#ccc]" />
            <Table<UserListItemType> bordered columns={TableColumns} dataSource={user_list_data} rowKey="id" />
        </Space>
    )
}
