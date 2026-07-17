import { createFileRoute } from '@tanstack/react-router'
import { Table, Form, Select, Input, Space } from 'antd'
import { TableColumns, dataSource, NameSelectOptions } from './constant'
import { FieldType } from './types'

export const Route = createFileRoute('/userList/')({
    component: RouteComponent,
})




function RouteComponent() {
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
                        options={NameSelectOptions}
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
