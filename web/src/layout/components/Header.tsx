import { Layout, Space, Button, Avatar } from 'antd'
const { Header } = Layout

export default function HeaderComponent() {

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
                    {/* {collapsed ? '☰' : '⇤'} */}
                </Button>
                {/* <Text type="secondary">首页 / 仪表盘</Text> */}
            </Space>
            <Space align="center" size="middle">
                <Button type="default">帮助</Button>
                <Button type="primary">新建</Button>
                <Avatar style={{ background: '#7C6FF0' }}>用户名</Avatar>
            </Space>
        </Header>
    )
}
