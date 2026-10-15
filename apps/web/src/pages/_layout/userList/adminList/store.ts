import { create } from 'zustand'


interface UserListStore {
    addUserDrawerOpen: boolean // 添加成员账号抽屉是否打开
    setAddUserDrawerOpen: (addUserDrawerOpen: boolean) => void // 设置添加成员账号抽屉是否打开
}

export const useUserListStore = create<UserListStore>((set) => ({
    addUserDrawerOpen: false,
    setAddUserDrawerOpen: (addUserDrawerOpen) => set({ addUserDrawerOpen }),
}))
