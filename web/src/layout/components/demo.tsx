import { useState } from 'react'
import {
    Layout as AntdLayout,
    Menu,
    Button,
    Card,
    Table,
    Tag,
    Avatar,
    Input,
    Space,
    Typography,
    Row,
    Col,
    Statistic,
} from 'antd'
import type { MenuProps, TableProps } from 'antd'

const { Header, Content, Sider } = AntdLayout
const { Title, Text } = Typography

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

interface UserRow {
    key: number
    name: string
    role: string
    status: 'active' | 'disabled'
    created: string
}

const data: UserRow[] = [
    { key: 1, name: '张三', role: '超级管理员', status: 'active', created: '2026-07-01' },
    { key: 2, name: '李四', role: '运营', status: 'active', created: '2026-07-03' },
    { key: 3, name: '王五', role: '访客', status: 'disabled', created: '2026-07-08' },
]

const columns: TableProps<UserRow>['columns'] = [
    { title: '用户名', dataIndex: 'name', key: 'name' },
    {
        title: '角色',
        dataIndex: 'role',
        key: 'role',
        render: (role: string) => <Tag color="purple">{role}</Tag>,
    },
    {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        render: (status: UserRow['status']) =>
            status === 'active' ? (
                <Tag color="success">启用</Tag>
            ) : (
                <Tag>禁用</Tag>
            ),
    },
    { title: '创建时间', dataIndex: 'created', key: 'created' },
    {
        title: '操作',
        key: 'action',
        render: () => (
            <Space size="small">
                <Button type="link" size="small">编辑</Button>
                <Button type="link" size="small" danger>删除</Button>
            </Space>
        ),
    },
]

export default function Layout() {
    const [collapsed, setCollapsed] = useState(false)

    return (
        <AntdLayout style={{ height: '100%' }}>
            <Sider
                theme="light"
                width={224}
                collapsible
                collapsed={collapsed}
                trigger={null}
                style={{ borderInlineEnd: '1px solid #F0EEF8' }}
            >
                <div
                    style={{
                        height: 64,
                        display: 'flex',
                        alignItems: 'center',
                        gap: 10,
                        paddingInline: 20,
                    }}
                >
                    <div
                        style={{
                            width: 30,
                            height: 30,
                            borderRadius: 8,
                            background: 'linear-gradient(135deg, #8B7CF0, #6A5CE8)',
                            color: '#fff',
                            fontWeight: 700,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            boxShadow: '0 4px 10px rgba(124,111,240,0.35)',
                        }}
                    >
                        G
                    </div>
                    {!collapsed && (
                        <Text strong style={{ fontSize: 16 }}>GinAdmin</Text>
                    )}
                </div>
                <Menu
                    mode="inline"
                    defaultSelectedKeys={['dashboard']}
                    defaultOpenKeys={['system']}
                    items={menuItems}
                    style={{ borderInlineEnd: 'none' }}
                />
            </Sider>

            <AntdLayout>
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
                            onClick={() => setCollapsed(!collapsed)}
                            style={{ fontSize: 16 }}
                        >
                            {collapsed ? '☰' : '⇤'}
                        </Button>
                        <Text type="secondary">首页 / 仪表盘</Text>
                    </Space>
                    <Space align="center" size="middle">
                        <Button type="default">帮助</Button>
                        <Button type="primary">新建</Button>
                        <Avatar style={{ background: '#7C6FF0' }}>管</Avatar>
                    </Space>
                </Header>

                <Content style={{ padding: 24, overflow: 'auto' }}>
                    <Row gutter={16} style={{ marginBottom: 16 }}>
                        <Col span={8}>
                            <Card>
                                <Statistic title="总用户" value={1280} valueStyle={{ color: '#7C6FF0' }} />
                            </Card>
                        </Col>
                        <Col span={8}>
                            <Card>
                                <Statistic title="今日活跃" value={326} valueStyle={{ color: '#22C55E' }} />
                            </Card>
                        </Col>
                        <Col span={8}>
                            <Card>
                                <Statistic title="待处理" value={12} valueStyle={{ color: '#F59E0B' }} />
                            </Card>
                        </Col>
                    </Row>

                    <Card>
                        <div
                            style={{
                                display: 'flex',
                                justifyContent: 'space-between',
                                alignItems: 'center',
                                marginBottom: 16,
                            }}
                        >
                            <Title level={5} style={{ margin: 0 }}>用户列表</Title>
                            <Space>
                                <Input.Search placeholder="搜索用户" style={{ width: 220 }} allowClear />
                                <Button type="primary">新增用户</Button>
                            </Space>
                        </div>
                        <Table columns={columns} dataSource={data} pagination={{ pageSize: 5 }} />
                    </Card>
                </Content>
            </AntdLayout>
        </AntdLayout>
    )
}
