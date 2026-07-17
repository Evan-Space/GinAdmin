import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, message, Typography } from 'antd'

const BASE_URL = 'http://localhost:8080/api/v1'

type LoginFormValues = {
    username: string
    password: string
}

export const Route = createFileRoute('/login/')({
    component: RouteComponent,
})

function RouteComponent() {
    const navigate = useNavigate()
    const [form] = Form.useForm<LoginFormValues>()
    const [messageApi, contextHolder] = message.useMessage()

    const handleSubmit = async (values: LoginFormValues) => {
        try {
            const response = await fetch(`${BASE_URL}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(values),
            })
            const result = await response.json()

            if (result.code !== 0) {
                messageApi.error(result.msg || '登录失败')
                return
            }

            localStorage.setItem('auth_token', result.data.access_token)
            messageApi.success('登录成功')
            navigate({ to: '/' })
        } catch {
            messageApi.error('网络异常，请稍后重试')
        }
    }

    return (
        <div className="flex min-h-screen items-center justify-center">
            {contextHolder}
            <Card style={{ width: 400 }}>
                <Typography.Title level={3} style={{ marginBottom: 24, textAlign: 'center' }}>
                    系统登录
                </Typography.Title>
                <Form<LoginFormValues>
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                    autoComplete="off"
                >
                    <Form.Item
                        label="用户名"
                        name="username"
                        rules={[
                            { required: true, message: '请输入用户名' },
                            { min: 3, max: 16, message: '用户名长度为 3-16 位' },
                        ]}
                    >
                        <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
                    </Form.Item>

                    <Form.Item
                        label="密码"
                        name="password"
                        rules={[
                            { required: true, message: '请输入密码' },
                            { min: 6, max: 18, message: '密码长度为 6-18 位' },
                        ]}
                    >
                        <Input.Password prefix={<LockOutlined />} placeholder="请输入密码" />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>
                            登录
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    )
}
