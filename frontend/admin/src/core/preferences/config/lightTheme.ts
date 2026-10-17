import type { ThemeConfig } from 'antd';

/**
 * 浅色模式扩展色板 —— docs/design-language.md §2.2 的落地值
 *
 *   canvas #F1F3F6  页面画布（卡片在其上"浮起"）
 *   surface #FFFFFF 卡片 / 表格容器 / 顶栏 / 侧栏
 *   border  #E4E4E7 通用边框（≈ zinc-200）
 *
 * 沿革：本文件 2026-10-16 之前不存在——react 端此前只有 darkTheme.ts 一套覆盖，浅色全部
 * 走 antd defaultAlgorithm 派生值（实测画布 #F5F5F5、卡片边框 #F0F0F0、卡片投影 none），
 * 白卡压在 #F5F5F5 上分离度只有 1.049、边框对白底 1.045，等于"没有层级系统"。
 * §2.2 开头那条"react 半 ✓ … 要严格对齐需新增一套 light tokens，未做，属新工作"记录的就是这里。
 */
export const LIGHT_PALETTE = {
  bgLayout: '#F1F3F6',
  bgSurface: '#FFFFFF',
  borderCard: '#E4E4E7',
} as const;

/**
 * 浅色模式扩展 token —— 叠加在 theme.defaultAlgorithm 之上
 *
 * 只补"层级"这一条线（画布/容器/边框/圆角），文字档不动：
 * 主/常规/次要三档已由 src/styles/semantic-text.css 与 colorTextDescription
 * （constants.ts 的 SECONDARY_TEXT_LIGHT）接上，再在 token 层覆盖一遍会派生出
 * 第二套真相。
 *
 * colorTextPlaceholder 同样**故意不动**：§2.2 记的浅色档是 #9CA3AF，但 antd
 * 默认的 rgba(0,0,0,.25)（合成 #BFBFBF）与已输入值 #1F2937 的分离度 6.06 倍，
 * 比 #9CA3AF 的 5.78 倍更"退后"。本轮的诉求是去掉生硬、不是加重提示语，
 * 钉成 #9CA3AF 会把暗色那轮刚压下去的问题在浅色侧重犯一遍。
 */
export const lightThemeTokens: ThemeConfig['token'] = {
  // 画布：§2.2 的 #F1F3F6（冷调微蓝），替代 antd 派生的中性 #F5F5F5
  colorBgLayout: LIGHT_PALETTE.bgLayout,
  colorBgContainer: LIGHT_PALETTE.bgSurface,
  // 边框：antd 默认 colorBorderSecondary 是 #F0F0F0，在白底上 1.045 倍、几乎隐形，
  // 卡片只剩一条看不见的线 ⇒ "平"。提到 §2.2 档后 1.135 倍，卡片边界可辨。
  // 注意连带效应（antd/es/theme/util/alias.js:51）：
  // colorSplit = getAlphaColor(colorBorderSecondary, colorBgContainer)，
  // 所以这一条同时把浅色分隔线从 rgba(5,5,5,.06) 抬到等效 #E4E4E7 ——
  // 正好补齐 §2.2 的边框档（此前浅色抽屉/表头分隔线实测只有 1.140，低于规范 1.269）。
  colorBorderSecondary: LIGHT_PALETTE.borderCard,
  // 大圆角：§2.4 卡片/抽屉/弹窗 = 12px。不设时浅色拿的是 borderRadius 派生值
  // （偏好 radius=8 ⇒ 实测 10px），而暗色在 darkThemeTokens 里显式钉了 12，两态形状对不上。
  borderRadiusLG: 12,
};

// 浅色投影两档（静止 0 1px 2px / 悬浮 0 4px 12px）不在这里：antd 没有可覆写的
// "带边框卡片静止投影"token——boxShadowCard 是内部派生值，写进 ThemeConfig.token
// 会被 TS 挡掉（AliasToken 上不存在），而它只服务 hoverable 卡片的 hover。
// 因此两档统一落在 styles/global.css 的 [data-theme='light'] .ant-card 上，
// 那里静止/悬浮都能覆盖，不留第二处真相。

// 浅色没有组件级 token：Card 描边已由全局 colorBorderSecondary 覆盖（组件级默认取它），
// 输入类边框刻意留在 antd 的 colorBorder #D9D9D9 —— §2.2 的 #E4E4E7 比 #D9D9D9 更浅，
// 套到输入框上会让控件描边倒退一档。
