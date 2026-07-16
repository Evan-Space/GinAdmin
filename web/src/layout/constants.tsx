import { HomeOutlined, UserOutlined } from "@ant-design/icons";
import type { MenuProps } from "antd";


export const MENU_LIST: MenuProps['items'] =  [
    { key: "/", label: "首页", icon: <HomeOutlined /> },
    { key: "/userList", label: "用户列表", icon: <HomeOutlined /> },

    { key: "/permission", label: "权限管理", icon: <UserOutlined /> },
    { key: "/log", label: "操作日志", icon: <UserOutlined />, children: [
        { key: "/log/requestLog", label: "请求日志", icon: <UserOutlined /> },
        { key: "/log/errorLog", label: "错误日志", icon: <UserOutlined /> },
    ] },
    { key: "/task", label: "任务中心", icon: <UserOutlined /> },
    { key: "/setting", label: "系统管理", icon: <UserOutlined /> },
]