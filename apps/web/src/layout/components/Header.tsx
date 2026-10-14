import { Layout, Space, Button, message } from 'antd'
const { Header } = Layout
import { useNavigate } from '@tanstack/react-router'

export default function HeaderComponent() {

    const navigate = useNavigate()
    const handleLogout = () => {
        localStorage.removeItem('auth_token')
        message.success('退出成功')
        navigate({ to: '/login' })
    }
    return (
        <Header
            style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                paddingInline: 20,
                boxShadow: '0 1px 4px rgba(33,30,51,0.06)',
                position: 'relative',
                zIndex: 10,
            }}
        >
            <Space align="center" size="middle">
                <Button
                    type="text"
                    style={{ fontSize: 16 }}
                >
                    GinAdmin
                </Button>
            </Space>
            <Space align="center" size="middle">
                <Button type="default">帮助</Button>
                <Button type="primary">新建</Button>
                <Button type="primary" onClick={handleLogout}>退出</Button>
            </Space>
        </Header>
    )
}
