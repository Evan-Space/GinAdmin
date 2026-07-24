import { createFileRoute } from '@tanstack/react-router'
import { Table, Form, Select, Input, Space, Button } from 'antd'
import { TableColumns } from './constant'
import { FieldType } from './types'
import { useUserList } from './hooks'
import { UserListItemType } from './types'

export const Route = createFileRoute('/_layout/userList/')({
    component: RouteComponent,
})

function RouteComponent() {
    const { form, pagination, USER_NAME_LIST_OPTIONS, user_list_data, handleSearch } = useUserList()

    return (
        <Space orientation="vertical" size="medium" style={{ display: 'flex' }}>
            <Form
                form={form}
                layout="inline"
                wrapperCol={{ span: 8, style: { minWidth: '200px' } }}
            >
                <Form.Item<FieldType> label="姓名" name="nickname">
                    <Select allowClear options={USER_NAME_LIST_OPTIONS} />
                </Form.Item>

                <Form.Item<FieldType> label="年龄" name="age">
                    <Input />
                </Form.Item>
                <Form.Item>
                    <Button
                        type="primary"
                        htmlType="submit"
                        onClick={() => handleSearch({ currentPage: 1, pageSize: 10 })}
                    >
                        Search
                    </Button>
                </Form.Item>
            </Form>
            <div className="my-12 bg-[#ccc]" />
            <Table<UserListItemType>
                bordered
                columns={TableColumns}
                dataSource={user_list_data}
                rowKey="id"
                pagination={{
                    current: pagination.currentPage,
                    pageSize: pagination.pageSize,
                    total: pagination.total,
                    onChange: (page, pageSize) => {
                        handleSearch({
                            currentPage: page,
                            pageSize: pageSize,
                        })
                    },
                }}
            />
        </Space>
    )
}
