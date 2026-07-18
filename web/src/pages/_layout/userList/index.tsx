import { createFileRoute } from '@tanstack/react-router'
import { Table, Form, Select, Input, Space } from 'antd'
import { TableColumns, dataSource, NameSelectOptions } from './constant'
import { FieldType } from './types'
import { useUserList } from './hooks'



export const Route = createFileRoute('/_layout/userList/')({
    component: RouteComponent,
})




function RouteComponent() {

    const { USER_NAME_LIST_OPTIONS } = useUserList()

    console.log(USER_NAME_LIST_OPTIONS)
    console.log(111)


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
            <Table bordered columns={TableColumns} dataSource={dataSource} />
        </Space>
    )
}
