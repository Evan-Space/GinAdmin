import { createFileRoute } from '@tanstack/react-router'
import { Card, Space, Form, Select, InputNumber } from 'antd'
import { useMemberList } from './hooks'
import { FieldType } from './types'
import { ACCOUNT_STATUS } from './constants'

export const Route = createFileRoute('/_layout/userList/memberList/')({
    component: RouteComponent,
})

function RouteComponent() {
    const { form, USER_NAME_LIST_OPTIONS } = useMemberList()

    
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


                </Form>
            </Card>
            <Card></Card>
        </Space>
    )
}
