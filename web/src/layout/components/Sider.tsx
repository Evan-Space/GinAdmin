import { Layout, Menu } from 'antd'
const { Sider } = Layout
import { useSideMenuStore } from '@src/store/sideMenu'
import { useNavigate, useRouterState } from '@tanstack/react-router'
import { MENU_LIST } from '@src/layout/constants'

export default function SiderComponent() {

    const { collapsed, setCollapsed } = useSideMenuStore()
    const pathname = useRouterState({ select: (state) => state.location.pathname })
    const navigate = useNavigate()





    return (
        <Sider
            theme="light"
            width={224}
            
            collapsible
            collapsed={collapsed}
            onCollapse={setCollapsed}
            trigger={null}
            style={{ borderInlineEnd: '1px solid #F0EEF8' }}
        >
            <Menu
                mode="inline"
                defaultSelectedKeys={['dashboard']}
                defaultOpenKeys={['system']}
                selectedKeys={[pathname]}
                style={{ borderInlineEnd: 'none' }}
                items={MENU_LIST}
                onClick={({key}) => {
                    navigate({ to: key })
                }}
            />
        </Sider>
    )
}
