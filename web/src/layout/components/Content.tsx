import { Layout, Row, Col, Card, Statistic, Space, Input, Button, Table } from 'antd'
const { Content } = Layout

export default function ContentComponent() {
    return (
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
                    {/* <Title level={5} style={{ margin: 0 }}>
                        用户列表
                    </Title> */}
                    <Space>
                        <Input.Search placeholder="搜索用户" style={{ width: 220 }} allowClear />
                        <Button type="primary">新增用户</Button>
                    </Space>
                </div>
                {/* <Table columns={columns} dataSource={data} pagination={{ pageSize: 5 }} /> */}
            </Card>
        </Content>
    )
}
