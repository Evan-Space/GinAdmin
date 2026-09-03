import { createFileRoute } from '@tanstack/react-router'
import { Card, Space, Form, Select, InputNumber, Button, Table } from 'antd'
import { useMemberList } from './hooks'
import { FieldType, MemberListItemType } from './types'
import { ACCOUNT_STATUS, getTableColumns } from './constants'

export const Route = createFileRoute('/_layout/userList/memberList/')({
    component: RouteComponent,
})

function RouteComponent() {
    const { form, USER_NAME_LIST_OPTIONS, handleSearch, member_list_data, pagination, handleDeleteAccount } = useMemberList()

    return (
        <Space orientation="vertical" size="medium" className="flex w-full">
            <Card>
                <Form
                    form={form}
                    layout="inline"
                    wrapperCol={{ span: 8, style: { minWidth: '200px' } }}
                >
                    <Form.Item<FieldType> label="姓名" name="nickname">
                        <Select allowClear options={USER_NAME_LIST_OPTIONS} />
                    </Form.Item>

                    <Form.Item<FieldType> label="年龄" name="age">
                        <InputNumber min={1} max={100} />
                    </Form.Item>

                    <Form.Item<FieldType> label="账号状态" name="status">
                        <Select allowClear options={ACCOUNT_STATUS} />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" onClick={() => handleSearch({ currentPage: 1, pageSize: 10 })}>搜索</Button>
                    </Form.Item>
                </Form>
            </Card>
            <Card>
                <Space orientation="vertical" size="medium" className="w-full mb-4">
                    {/* <Button
                        type="primary"
                        // onClick={() => {
                        //     setAddUserDrawerOpen(true)
                        // }}
                    >
                        添加账号
                    </Button> */}

                    <Table<MemberListItemType>
                        bordered
                        columns={getTableColumns(handleDeleteAccount)}
                        dataSource={member_list_data}
                        rowKey="id"
                        pagination={{
                            current: pagination.currentPage,
                            pageSize: pagination.pageSize,
                            total: pagination.total,
                            showSizeChanger: true,
                            onChange: (page, pageSize) => {
                                handleSearch({
                                    currentPage: page,
                                    pageSize: pageSize,
                                })
                            },
                        }}
                        />
                </Space>
            </Card>
        </Space>
    )
}
