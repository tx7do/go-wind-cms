# GoWind 风行 · 设计语言（Design Language）

> 本文档是**前台四端**（react / vue / taro / flutter）与 admin 共用的视觉设计唯一事实源（Single Source of Truth）。
> 品牌色板与 [docs/brand/README.md](./brand/README.md) 同源；品牌图形标为方案 B「旋涡 Vortex」。
> 2026-09-16 制定并落地 C 端。改任何颜色/尺寸前先改这里，再同步各端 token 文件。

---

## 1. 品牌与设计概念

品牌主名 **GoWind**，中文名 **风行**。一切视觉决策围绕「风」的三个意象：

| 意象 | 视觉转译 | 落点 |
|---|---|---|
| 速度 | 高饱和单一主色、干脆的动效（150–300ms）、无重装饰 | 按钮/链接/进度 |
| 流动 | 冷灰蓝中性底、圆角与留白、渐变只用于大面积氛围（hero/logo），禁止小元素渐变 | 背景/卡片/hero |
| 汇聚 | 内容向卡片聚合、层级靠「底色深浅」而非重边框、主色只出现在行动点 | 布局/层级 |

**三条铁律**：
1. **主色只有一个**：品牌蓝 `hsl(212 100% 45%)`。禁止再出现第二品牌色（历史遗留的绿 #16A34A、钢蓝 #3A7CA5 已废止）。
2. **改色先改 hue 锚**：所有派生色（hover/tint/暗色变体）只允许在 hue 212 ± 8 的蓝域内调明度饱和度；语义色（success/warning/destructive）除外。
3. **语义色不许挪用**：success 绿、warning 橙、destructive 红只表达状态，不参与装饰。

---

## 2. 色彩系统

### 2.1 品牌色板（权威值，来自 docs/brand）

| 用途 | 值 |
|---|---|
| 品牌主色 | `hsl(212 100% 45%)` ≈ `#006BE6` |
| 渐变亮端 | `#70B3FF` |
| 渐变中段 | `#0D74F2` |
| 渐变暗端 | `#005AC2` |
| 应用图标底渐变 | `#1E7BFF` → `#0053BD`（135°） |

### 2.2 语义 Token（shadcn 约定，HSL 三元组）

Web 端（react/vue/taro-H5）统一以 CSS 变量承载，组件只允许引用 token，不允许写死色值。

#### 浅色模式 —「晴风」

| Token | 值 | 说明 |
|---|---|---|
| `--background` | `210 40% 98%` | 冷灰蓝底 `#F4F7F9`，风感基底 |
| `--foreground` | `224 71.4% 4.1%` | slate-950 正文 |
| `--card` | `0 0% 100%` | 纯白浮空卡 |
| `--card-foreground` | 同 foreground | |
| `--popover` | `0 0% 100%` | |
| `--primary` | **`212 100% 45%`** | 品牌蓝，唯一主色 |
| `--primary-foreground` | `0 0% 100%` | |
| `--secondary` | `210 40% 93.1%` | 次级按钮/徽底 |
| `--secondary-foreground` | `224 71.4% 4.1%` | |
| `--muted` | `210 40% 93.1%` | 弱底 |
| `--muted-foreground` | `215.4 16.3% 46.9%` | slate-500 辅助文字 |
| `--accent` | `212 100% 95%` | **品牌蓝浅底**：hover/焦点底色（旧版饱和蓝废止） |
| `--accent-foreground` | `212 100% 38%` | 浅蓝底上的深蓝文字 |
| `--destructive` | `0 84.2% 60.2%` | |
| `--success` | `142.1 76.2% 36.3%` | 语义绿（仅状态） |
| `--warning` | `38 92% 50%` | 语义橙（仅状态） |
| `--border` | `214.3 32% 88%` | slate-200 |
| `--input` | `214.3 32% 88%` | |
| `--ring` | `212 100% 45%` | 焦点环随主色 |
| `--radius` | `0.6rem` | 全局圆角基值 |

#### 深色模式 —「夜风」

| Token | 值 | 说明 |
|---|---|---|
| `--background` | `224 45% 6%` | 深空夜蓝 `#050B14`，非死黑 |
| `--foreground` | `210 20% 98%` | |
| `--card` | `222.2 47.4% 11%` | 比背景亮一阶的磨砂蓝黑 |
| `--popover` | 同 card | |
| `--primary` | **`212 100% 58%`** | 品牌蓝提亮变体（hue 不变，暗底可读）≈ `#2E96FF` |
| `--primary-foreground` | `0 0% 100%` | |
| `--secondary` / `--muted` | `217.2 32.6% 13.5%` | |
| `--muted-foreground` | `215 20.2% 65.1%` | |
| `--accent` | `212 50% 18%` | 品牌蓝暗底 tint |
| `--accent-foreground` | `212 100% 78%` | |
| `--destructive` | `0 62.8% 50%` | |
| `--success` | `142.1 86.2% 50.3%` | |
| `--warning` | `38 92% 55%` | |
| `--border` / `--input` | `217.2 32.6% 17.5%` | |
| `--ring` | `212 100% 58%` | |

### 2.3 用色纪律

- **60-30-10**：约 60% 中性底（background/card）、30% 文字与分隔、10% 主色。主色只落在：主按钮、链接、激活态、焦点环、品牌图形。
- 主色透明度阶梯做层次：`primary/10` 徽底、`primary/30` 描边、`primary` 实色，禁止再引入新蓝。
- 渐变只允许两处：hero 氛围层（`#70B3FF→#0D74F2→#005AC2` 低透明度）与品牌图形本体。
- 正文对比度 ≥ 4.5:1；辅助文字（muted-foreground）只用于元信息，不用于行动点。

---

## 3. 字体排印

| 层级 | 规格 | 用途 |
|---|---|---|
| Display | 30–36px / 800 / -0.02em | 首页 hero 标题 |
| H1 | 24–28px / 700 | 页面标题、文章题图标题 |
| H2 | 20px / 700 | 区块标题（配 3px 主色左条或图标） |
| H3 | 16px / 600 | 卡片标题，两行截断 |
| Body | 14–15px / 400 / 1.7 | 正文；文章正文 16–17px / 1.8 |
| Caption | 12–13px / 400 / muted-foreground | 时间、计数、元信息 |

- 字体栈：系统栈优先（`system-ui, -apple-system, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', 'Noto Sans SC', sans-serif`），不引入外部字体文件（品牌字标已转路径，不依赖字体）。
- 中文正文不使用斜体；强调用字重与主色。

---

## 4. 空间、圆角、阴影

- **间距**：4px 基数。卡片内边距 20px（p-5），区块纵向节奏 48/64px（`py-12/py-16`），页面容器 `max-w-[1200px]` + `px-4`（移动 `px-4`、平板 `px-6`）。
- **圆角**：全局 `--radius: 0.6rem`；卡片 `rounded-2xl`（1rem）、按钮/输入 `rounded-lg`（0.6rem）、徽章/标签全圆 `rounded-full`。
- **阴影**：三级
  - `shadow-sm`：静态卡片（近乎不可见，靠 border 分层）
  - `shadow-[0_12px_30px_rgba(0,107,230,0.08)]`：hover 浮起（**主色色相投影**，禁止沿用旧绿影）
  - `shadow-lg`：弹层（popover/dropdown）
- **边框**：`border-border` 一律 1px；hover 时可过渡到 `border-primary/30`。

---

## 5. 布局

| 项 | 值 |
|---|---|
| Header | 高 64px（移动 56px），白底/卡底 + 细分割线，sticky；logo 32–36px + 站点名 |
| 内容容器 | `max-w-[1200px]` 居中 |
| Footer | 深色（`foreground` 反白底或 card 底 + 上边框），四栏链接，移动单列 |
| 首页 | Hero（题图轮播/品牌渐变层）→ 最新文章（两栏卡片流 + 侧栏热标签）→ 分类宫格 → 推荐阅读 |
| 详情页 | 正文卡 `max-w-[760px]` 居中，侧边操作浮层（点赞/收藏/分享）桌面右浮动、移动底栏 |
| 移动断点 | `<768px` 单列；`768–1024` 双列卡片；`>1024` 桌面栅格 + 侧栏 |
| 触达 | 移动端遵守安全区（`env(safe-area-inset-*)`），底部操作栏高 56px + safe-bottom |

---

## 6. 动效

| 类型 | 规格 |
|---|---|
| 微交互 | 150–200ms `ease-out`；hover 位移 ≤ 4px（`-translate-y-1`） |
| 卡片浮起 | 300ms `ease-out`：上浮 4px + 主色柔影 + 边框转 `primary/30` |
| 页面过渡 | fade 300ms（现有 `.fade-*` 约定保留） |
| 滚动入场 | `.scroll-reveal-item` 600ms 逐项 80ms 延迟（保留） |
| 品牌动效 | 仅登录页 `vortex-login.svg`（24s/圈缓旋）；尊重 `prefers-reduced-motion` |

---

## 7. 组件约定

- **按钮**：primary 实色蓝底白字；secondary `secondary` 底；ghost hover `accent` 底。圆角 `rounded-lg`，高度 40px（移动 44px）。
- **卡片**：`bg-card border border-border rounded-2xl`，静态 `shadow-sm`，hover 浮起（见 §6）。卡片内层级：题图 16:9 → 标题 H3 → 摘要 2 行截断 → 元信息行（头像+作者+时间）。
- **占位图形**（无图文章/分类）：禁止灰块白图标；统一「品牌蓝 tint 底（`primary/8`）+ 主色线性图标 + 右下角淡色风纹曲线」的组合，颜色全部跟随 token（换主题自动适配）。
- **标签/徽章**：`primary/10` 底 + `primary` 文字 + 全圆；热门标签可加 `Warning`/火焰色但仅限一个强调位。
- **表单**：输入框 `border-input rounded-lg`，focus 时 `ring-2 ring-ring ring-offset-1`；错误态 `destructive`。
- **空态/加载**：骨架屏用 `muted`；空态配线性插画（风纹母题）+ muted-foreground 文案。

---

## 8. 品牌资产使用

| 资产 | 用途 | 位置 |
|---|---|---|
| `logo.png`（vortex tile 200×200） | Header/侧栏 logo | 各端 `public/logo.png` |
| `favicon.ico` | 浏览器图标 | 各端 `public/favicon.ico` |
| `vortex-lockup.svg` | 需要横版组合的场景（footer 关于、关于页） | docs/brand/ |
| `vortex-login.svg` | 登录页品牌插画 | docs/brand/（admin 三端已入） |
| 旧橘猫吉祥物 | **废止**，从 public/ 移除 | — |

- 净空：logo 四周 ≥ 1/4 高度；禁拉伸、禁改色相、禁加投影描边。
- 产品命名「GoWind + 产品名」；C 端站点名以后台站点配置为准（demo 为「风行内容中台」）。

---

## 9. 各端实现映射（落地清单）

| 端 | Token 文件 | 说明 |
|---|---|---|
| react (Next.js) | `frontend/app/react/src/app/globals.css` | Tailwind v4 `@theme inline` + shadcn 变量，照 §2.2 置换 |
| vue (Nuxt) | `frontend/app/vue/app/assets/css/main.css` | 同上；另有 `--hero-gradient-*` 变量转蓝 |
| taro | `frontend/app/taro/src/app.css` | `page` 级 CSS 变量，中性色对齐 slate、新增 `--color-primary` 系 |
| flutter | `frontend/app/flutter_app/lib/src/core/themes/{light,dark}_theme.dart` | `primaryColor #006BE6`，container/tint 按 Material 由 seed 派生，seed 同步 |
| admin（参照） | `frontend/admin/packages/@core/base/design/src/design-tokens/` | 已落地品牌蓝，作为跨端一致性基准 |

> 验收方式：四端亮/暗两态截图比对；grep 四端源码不得再出现 `16A34A / 22C55E / 34,197,94 / 3A7CA5` 等废止色值。
