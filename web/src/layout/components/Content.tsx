import { Layout } from 'antd'
import { Outlet } from '@tanstack/react-router'
const { Content } = Layout

export default function ContentComponent() {
    return (
        <Content style={{ padding: 24, overflow: 'auto' }}>
            <Outlet />
        </Content>
    )
}
