import { Form, App } from 'antd'
import { useUserListStore } from '../../store'
import { createUserAPI } from '@src/request/userList'

export const useAddUser = () => {
    const [form] = Form.useForm()
    const { setAddUserDrawerOpen } = useUserListStore()
    const { message: messageApi } = App.useApp()



    /**
     * 添加按钮
    */
    const handleAddBtn = async () => {
        const values = await form.getFieldsValue(true)
        const { code, msg } = await createUserAPI(values)
        if (code !== 0) {
            messageApi.error(msg || '添加失败，请重试')
            return
        }
        messageApi.success('添加成功')
        handleCloseDrawer()
    }

    /**
     * 取消
     */
    const handleCancelBtn = async () => {
        await form.resetFields()
        setAddUserDrawerOpen(false)
    }

    /**
     * 关闭抽屉
    */
    const handleCloseDrawer = async () => {
        await form.resetFields()
        setAddUserDrawerOpen(false)
    }

    return {
        form,
        handleCancelBtn, // 取消按钮
        handleAddBtn, // 添加按钮
        handleCloseDrawer, // 关闭抽屉
    }
}
