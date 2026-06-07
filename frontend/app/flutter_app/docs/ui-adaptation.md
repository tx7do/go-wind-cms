# UI 适配规范

本文档描述 GoWind CMS Flutter 端的多端 UI 适配策略，涵盖断点体系、布局切换、尺寸适配、字体配置等内容。

---

## 1. 整体架构

项目采用 **双套布局 + 条件尺寸** 的混合适配方案：

```
┌──────────────────────────────────────────────┐
│  ScreenUtilInit (designSize: 375×812)        │
│  ┌──────────────────────────────────────────┐│
│  │  ResponsiveLayout                        ││
│  │  ├── mobileBody  (< 600dp)               ││
│  │  │   └── 使用 .sp .w .h .r 适配          ││
│  │  └── webBody    (>= 600dp)               ││
│  │       └── 使用固定 px 值                  ││
│  └──────────────────────────────────────────┘│
└──────────────────────────────────────────────┘
```

**核心原则：**
- **手机端**：使用 `flutter_screenutil` 的 `.sp/.w/.h/.r` 做等比缩放
- **Web/平板端**：使用固定像素值，内容区收窄居中
- **两套 UI 共享子组件**，通过 `isMobile` 参数切换尺寸策略

---

## 2. 断点体系

定义在 `lib/src/core/constants/breakpoints.dart`：

| 设备类型   | 宽度范围          | 布局策略      |
|--------|---------------|-----------|
| 手机     | < 600 dp      | 单栏瀑布流     |
| 平板     | 600 ~ 1024 dp | 双栏布局      |
| Web/桌面 | > 1024 dp     | 最大宽度居中/三栏 |

关键常量：

```dart
class Breakpoints {
  static const double mobile = 600;           // 手机端阈值
  static const double tablet = 1024;           // 平板端阈值
  static const double webContentMaxWidth = 1140; // Web 内容区最大宽度
  static const double webSidebarWidth = 260;    // Web 侧边栏宽度
  static const double webContentPadding = 32;   // Web 内容区间距

  static bool isMobile(double width) => width < mobile;
  static bool isTablet(double width) => width >= mobile && width < tablet;
  static bool isWeb(double width) => width >= tablet;
}
```

---

## 3. 布局切换组件

### 3.1 ResponsiveLayout

文件：`lib/src/core/widgets/responsive_layout.dart`

根据屏幕宽度自动选择渲染哪套布局：

```dart
ResponsiveLayout(
  mobileBody: HomeMobileView(),   // < 600dp
  tabletBody: SomeTabletView(),   // 600~1024dp（可选，缺省降级为 webBody）
  webBody: HomeWebView(),         // >= 600dp
);
```

内部通过 `LayoutBuilder` 获取 `constraints.maxWidth`，然后调用 `Breakpoints` 判断。

### 3.2 WebContentCenter

Web 端专用的居中内容容器，限制最大宽度并水平居中：

```dart
WebContentCenter(
  maxWidth: Breakpoints.webContentMaxWidth,
  padding: EdgeInsets.symmetric(horizontal: 24),
  child: content,
)
```

---

## 4. 尺寸适配策略

### 4.1 三元条件表达式（推荐）

项目中绝大多数组件采用 `isMobile ? mobileValue : webValue` 的模式：

```dart
// 字体
fontSize: isMobile ? 16.sp : 16,

// 间距
SizedBox(height: isMobile ? 14.h : 14),

// 内边距
padding: EdgeInsets.symmetric(horizontal: isMobile ? 20.w : 0),

// 圆角半径
BorderRadius.circular(isMobile ? 14.r : 14),

// 图标大小
size: isMobile ? 18.sp : 18,
```

**规律总结：**
- `isMobile` 为 `true` 时使用 `.sp/.w/.h/.r` 后缀（ScreenUtil 等比缩放）
- `isMobile` 为 `false` 时使用固定 `double` 值（Web 端无需缩放）

### 4.2 ResponsiveUtils 工具类

文件：`lib/src/core/utils/responsive_utils.dart`

提供统一的断点感知方法，无需手写三元表达式：

```dart
// 自适应内边距
ResponsiveUtils.padding(context, mobile: 16, tablet: 24, web: 32)

// 自适应水平内边距
ResponsiveUtils.horizontalPadding(context, mobile: 16, tablet: 24, web: 32)

// 自适应字体大小
ResponsiveUtils.fontSize(context, 14, tabletSize: 15, webSize: 16)

// 自适应间距
ResponsiveUtils.spacing(context, 12, tabletSize: 16, webSize: 20)

// 网格列数
ResponsiveUtils.gridColumns(context)         // 通用：1/2/3
ResponsiveUtils.postGridColumns(context)     // 文章：1/2/3
ResponsiveUtils.categoryGridColumns(context) // 分类：2/3/4

// 平台判断
ResponsiveUtils.isWideScreen(context)  // >= 1024dp
ResponsiveUtils.isMobile(context)      // < 600dp
ResponsiveUtils.isTablet(context)      // 600~1024dp
```

---

## 5. flutter_screenutil 配置

### 5.1 初始化

文件：`lib/src/app.dart`

```dart
ScreenUtilInit(
  designSize: const Size(375, 812),  // iPhone 6/7/8 设计稿尺寸
  minTextAdapt: true,
  splitScreenMode: true,
  ensureScreenSize: true,
  builder: (ctx, child) {
    // Web 端：动态重置 designSize 为当前视窗尺寸
    // 使 .w/.h/.sp 始终 1:1，字体不随窗口缩放
    if (kIsWeb) {
      ScreenUtil.init(
        ctx,
        designSize: Size(
          MediaQuery.of(ctx).size.width,
          MediaQuery.of(ctx).size.height,
        ),
        minTextAdapt: false,
      );
    }
    return _buildMaterialApp(context);
  },
);
```

**Web 端关键设计：** 将 `designSize` 动态设置为当前视窗尺寸，使得所有 `.w/.h/.sp` 值都等于原始 `dp` 值（1:1 映射）。用户缩小浏览器窗口时字体不会跟着缩小，符合 Web 阅读习惯。

### 5.2 ScreenUtil 后缀说明

| 后缀    | 含义     | 示例      |
|-------|--------|---------|
| `.sp` | 字体大小适配 | `14.sp` |
| `.w`  | 宽度方向适配 | `20.w`  |
| `.h`  | 高度方向适配 | `16.h`  |
| `.r`  | 圆角半径适配 | `14.r`  |

> **注意：** Web 端由于 designSize 被重置为视窗尺寸，这些后缀等同于原始值。因此 Web 端代码中直接使用固定值即可。

---

## 6. 字体配置

### 6.1 字体族

文件：`lib/src/core/themes/fonts.dart`、`lib/src/core/themes/const.dart`

```dart
// 主字体
const String kDefaultFontFamily = 'Noto Sans SC';

// 主题中配置
ThemeData(
  fontFamily: kDefaultFontFamily,
  fontFamilyFallback: AppFont.fontFamilyFallback,
);
```

### 6.2 平台字体回退链

`AppFont.fontFamilyFallback` 按平台提供不同回退字体：

| 平台      | 回退字体                               |
|---------|------------------------------------|
| Apple   | `.SF UI Text`, `PingFang SC/TC/HK` |
| Android | `miui`, `mipro`（小米）                |
| Windows | `Microsoft YaHei`                  |
| 通用      | `Noto Color Emoji`（Emoji 支持）       |

### 6.3 已注册字重

`pubspec.yaml` 中注册的 Noto Sans SC 字重：

| 字重            | 文件                       | 说明  |
|---------------|--------------------------|-----|
| 400 (Regular) | `NotoSansSC-Regular.ttf` | 正文  |
| 500 (Medium)  | `NotoSansSC-Medium.ttf`  | 副标题 |
| 700 (Bold)    | `NotoSansSC-Bold.ttf`    | 标题  |

使用方式：
```dart
fontWeight: FontWeight.w400    // → Regular
fontWeight: FontWeight.w500    // → Medium
fontWeight: FontWeight.bold   // → Bold (700)
```

---

## 7. 典型页面适配示例

### 7.1 首页（HomePage）

```dart
// 直接使用 ResponsiveLayout 分发
class HomePage extends StatelessWidget {
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: HomeMobileView(),  // NestedScrollView + TabBarView
      webBody: HomeWebView(),        // 顶部导航 + 左右双栏
    );
  }
}
```

### 7.2 文章详情页（PostDetailPage）

同一组件内部通过 `isMobile` 参数区分：

```dart
Widget _buildView(BuildContext context, {required bool isMobile}) {
  return Scaffold(
    body: Column(
      children: [
        Expanded(
          child: isMobile
              ? _buildMobileBody(context, post, comments)
              : _buildWebBody(context, post, comments),
        ),
        _CommentInputBar(isMobile: isMobile),
      ],
    ),
  );
}
```

### 7.3 内容卡片（PostCard）

```dart
// 三元条件贯穿每个尺寸属性
Container(
  padding: EdgeInsets.all(isMobile ? 16.w : 16),
  child: Column(
    children: [
      Image(height: isMobile ? 140.h : 140),
      SizedBox(height: isMobile ? 12.h : 12),
      Text(style: TextStyle(fontSize: isMobile ? 14.sp : 14)),
    ],
  ),
)
```

---

## 8. 新增页面适配 Checklist

新增页面或组件时，按以下步骤确保多端适配：

1. **布局层**：使用 `ResponsiveLayout` 包裹，分别提供 `mobileBody` 和 `webBody`
2. **尺寸层**：所有 `fontSize`、`padding`、`margin`、`size`、`radius` 使用 `isMobile ? mobileValue : webValue`
3. **手机端值**：带 ScreenUtil 后缀（`.sp`、`.w`、`.h`、`.r`）
4. **Web 端值**：使用固定 `dp` 值
5. **内容区宽度**：Web 端使用 `WebContentCenter` 或 `ConstrainedBox(maxWidth: Breakpoints.webContentMaxWidth)` 限制宽度
6. **网格列数**：使用 `ResponsiveUtils.postGridColumns()` 或 `ResponsiveUtils.categoryGridColumns()`
7. **测试**：分别在 < 600dp、600~1024dp、> 1024dp 三个宽度下验证

---

## 9. 相关文件索引

| 文件路径                                          | 说明             |
|-----------------------------------------------|----------------|
| `lib/src/core/constants/breakpoints.dart`     | 断点常量定义         |
| `lib/src/core/widgets/responsive_layout.dart` | 响应式布局组件        |
| `lib/src/core/utils/responsive_utils.dart`    | 响应式工具方法        |
| `lib/src/core/themes/const.dart`              | 默认字体、间距常量      |
| `lib/src/core/themes/fonts.dart`              | 字体族与回退链        |
| `lib/src/core/utilities/platform.dart`        | 平台判断工具         |
| `lib/src/app.dart`                            | ScreenUtil 初始化 |
| `pubspec.yaml`                                | 字体文件注册         |
