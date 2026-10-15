import type { ThemeConfig } from 'antd'

// shadcn/ui 色系（zinc）· 近黑主色 + 全中性灰阶 · 极简高级
export const appTheme: ThemeConfig = {
  token: {
    // —— 品牌 / 功能色 ——
    colorPrimary: '#18181B',
    colorInfo: '#18181B',
    colorLink: '#18181B',
    colorLinkHover: '#3F3F46',
    colorLinkActive: '#09090B',
    colorSuccess: '#16A34A',
    colorWarning: '#D97706',
    colorError: '#DC2626',

    // —— 文字（zinc 灰阶）——
    colorText: '#09090B',
    colorTextSecondary: '#52525B',
    colorTextTertiary: '#71717A',
    colorTextQuaternary: '#A1A1AA',

    // —— 边框 / 分割线 ——
    colorBorder: '#E4E4E7',
    colorBorderSecondary: '#F3F3F4',
    colorSplit: '#F1F1F3',

    // —— 背景 ——
    colorBgLayout: '#FAFAFA',
    colorBgContainer: '#FFFFFF',
    colorBgElevated: '#FFFFFF',
    colorBgSpotlight: '#18181B',

    // —— 填充 ——
    colorFill: '#E4E4E7',
    colorFillSecondary: '#F4F4F5',
    colorFillTertiary: '#FAFAFA',
    colorFillQuaternary: '#FAFAFA',

    // —— 圆角 ——
    borderRadius: 8,
    borderRadiusSM: 6,
    borderRadiusLG: 12,
    borderRadiusXS: 4,

    // —— 阴影（克制）——
    boxShadow: '0 1px 2px 0 rgba(9,9,11,0.05)',
    boxShadowSecondary: '0 4px 6px -1px rgba(9,9,11,0.07), 0 2px 4px -2px rgba(9,9,11,0.05)',
    boxShadowTertiary: '0 1px 2px 0 rgba(9,9,11,0.04)',

    // —— 字体 / 尺寸 ——
    fontSize: 14,
    controlHeight: 36,
    wireframe: false,
  },

  components: {
    Layout: {
      siderBg: '#FAFAFA',
      headerBg: '#FFFFFF',
      bodyBg: '#FAFAFA',
      footerBg: '#FFFFFF',
      headerHeight: 64,
      headerPadding: '0 24px',
    },

    Menu: {
      itemBg: 'transparent',
      subMenuItemBg: 'transparent',
      itemColor: '#3F3F46',
      itemHoverColor: '#18181B',
      itemHoverBg: '#F4F4F5',
      itemActiveBg: '#F4F4F5',
      itemSelectedColor: '#FAFAFA',
      itemSelectedBg: '#27272A',
      itemBorderRadius: 6,
      itemMarginInline: 8,
      itemHeight: 40,
      iconSize: 16,
    },

    Button: {
      primaryShadow: '0 1px 2px 0 rgba(9,9,11,0.06)',
      defaultShadow: '0 1px 2px 0 rgba(9,9,11,0.04)',
      fontWeight: 500,
    },

    Card: {
      boxShadowTertiary: '0 1px 2px 0 rgba(9,9,11,0.04)',
      headerBg: 'transparent',
      paddingLG: 24,
    },

    Table: {
      headerBg: '#FAFAFA',
      headerColor: '#52525B',
      rowHoverBg: '#FAFAFA',
      borderColor: '#F1F1F3',
      headerSplitColor: '#E4E4E7',
    },

    Input: {
      hoverBorderColor: '#D4D4D8',
      activeShadow: '0 0 0 3px rgba(9,9,11,0.08)',
    },

    Select: {
      optionSelectedBg: '#F4F4F5',
      optionActiveBg: '#FAFAFA',
    },

    Tag: {
      defaultBg: '#F4F4F5',
      defaultColor: '#3F3F46',
    },

    Tabs: {
      inkBarColor: '#18181B',
      itemSelectedColor: '#18181B',
      itemHoverColor: '#18181B',
    },
  },
}
