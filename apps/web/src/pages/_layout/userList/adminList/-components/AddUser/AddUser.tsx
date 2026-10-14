import { Drawer, Form, Input, InputNumber, Button, Space, Select } from 'antd'
import { useUserListStore } from '../../store'
import { useAddUser } from './hooks'
import { AddUserDrawerFormRules } from './constants'

export function AddUserDrawer() {
    const { addUserDrawerOpen } = useUserListStore()
    const { form, handleCancelBtn, handleAddBtn, handleCloseDrawer } = useAddUser()
    return (
        <Drawer
            title="添加账号"
            open={addUserDrawerOpen}
            size={800}
            onClose={async () => {
                handleCloseDrawer()
            }}
        >
            <Space orientation="vertical" size="medium" className="relative flex w-full h-full ">
                <Form
                    form={form}
                    wrapperCol={{ span: 8, style: { minWidth: '400px' } }}
                    labelCol={{ span: 4 }}
                >
                    <Form.Item label="昵称" name="nickname" rules={AddUserDrawerFormRules.nickname}>
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="用户名"
                        name="userName"
                        rules={AddUserDrawerFormRules.userName}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item label="密码" name="password" rules={AddUserDrawerFormRules.password}>
                        <Input />
                    </Form.Item>

                    <Form.Item label="年龄" name="age" rules={AddUserDrawerFormRules.age}>
                        <InputNumber />
                    </Form.Item>

                    <Form.Item label="状态" name="status" rules={AddUserDrawerFormRules.status}>
                        <Select
                            allowClear
                            options={[
                                { label: '启用', value: 1 },
                                { label: '禁用', value: 0 },
                            ]}
                        />
                    </Form.Item>
                </Form>

                <footer className="absolute bottom-0 left-0 right-0 flex justify-end gap-2">
                    <Button type="primary" onClick={handleAddBtn}>
                        添加
                    </Button>
                    <Button onClick={handleCancelBtn}>取消</Button>
                </footer>
            </Space>
        </Drawer>
    )
}
