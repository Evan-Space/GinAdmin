import { create } from 'zustand'
import type { MenuProps } from 'antd'

interface SideMenuStore {
    collapsed: boolean
    setCollapsed: (collapsed: boolean) => void
    menus: MenuProps
}

export const useSideMenuStore = create<SideMenuStore>((set) => ({
    collapsed: false,
    setCollapsed: (collapsed) => set({ collapsed }),
    menus: [],
    setMenus: (menus: MenuProps) => set({ menus }),
}))
