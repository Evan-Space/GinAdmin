import { Layout, Menu, MenuProps } from 'antd'
const { Sider } = Layout
import { useSideMenuStore } from '@src/store/sideMenu'

export default function SiderComponent() {

    const { collapsed, setCollapsed } = useSideMenuStore()

    const menuItems: MenuProps['items'] = [
        { key: 'dashboard', label: '仪表盘' },
        { key: 'user', label: '用户管理' },
        { key: 'role', label: '角色管理' },
        { key: 'menu', label: '菜单权限' },
        {
            key: 'system',
            label: '系统设置',
            children: [
                { key: 'config', label: '参数配置' },
                { key: 'log', label: '操作日志' },
            ],
        },
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
                items={menuItems}
                style={{ borderInlineEnd: 'none' }}
            />
        </Sider>
    )
}
