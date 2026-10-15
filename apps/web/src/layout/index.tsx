import { Layout as AntdLayout } from 'antd'
import HeaderComponent from './components/Header'
import SiderComponent from './components/Sider'
import ContentComponent from './components/Content'
import FooterComponent from './components/Footer'

export default function Layout() {
    return (
        <AntdLayout style={{ height: '100%' }}>
            <HeaderComponent />
            <AntdLayout>
                <SiderComponent />
                <ContentComponent />
            </AntdLayout>
            <FooterComponent />
        </AntdLayout>
    )
}
