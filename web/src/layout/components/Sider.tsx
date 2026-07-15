import { Layout, Menu, MenuProps } from 'antd'
const { Sider } = Layout
import { useSideMenuStore } from '@src/store/sideMenu'
import { HomeOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate, useRouterState } from '@tanstack/react-router'


export default function SiderComponent() {

    const { collapsed, setCollapsed } = useSideMenuStore()
    const pathname = useRouterState({ select: (state) => state.location.pathname })
    const navigate = useNavigate()

    const menuItems: MenuProps['items'] = [
        { key: "/", label: "首页", icon: <HomeOutlined /> },
        { key: "/permission", label: "权限管理", icon: <UserOutlined /> },
        { key: "/log", label: "操作日志", icon: <UserOutlined />, children: [
            { key: "/log/requestLog", label: "请求日志", icon: <UserOutlined /> },
            { key: "/log/errorLog", label: "错误日志", icon: <UserOutlined /> },
        ] },
        { key: "/task", label: "任务中心", icon: <UserOutlined /> },
        { key: "/setting", label: "系统管理", icon: <UserOutlined /> },
    ]



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
                items={menuItems}
                onClick={({key}) => {
                    navigate({ to: key })
                }}
            />
        </Sider>
    )
}
