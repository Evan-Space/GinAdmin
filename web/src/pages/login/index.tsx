import { createFileRoute } from '@tanstack/react-router'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, Typography } from 'antd'
import { useLogin } from './hooks'

export const Route = createFileRoute('/login/')({
    component: LoginPage,
})



type LoginFormValues = {
    username: string
    password: string
}




function LoginPage() {
    
    const { form, handleSubmit } = useLogin()
 
    return (
        <div className="flex min-h-screen items-center justify-center">
            <Card style={{ width: 400 }}>
                <Typography.Title level={3} style={{ marginBottom: 24, textAlign: 'center' }}>
                    系统登录
                </Typography.Title>
                <Form<LoginFormValues>
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                    autoComplete="off"
                    initialValues={{
                        username: 'super_admin',
                        password: '123456',
                    }}
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
