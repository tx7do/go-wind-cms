# AGENTS.md — Flutter 全平台前端开发指南

> 本文件是 `frontend/app/flutter_app` 子项目的 AI 编码规范单一事实源，适用于所有支持 AGENTS.md 的 AI 编码工具。Claude Code 通过同级 `CLAUDE.md` 中的 `@AGENTS.md` 引用加载。

## 项目概览

基于 **Flutter** 的全平台 CMS 内容展示前端，一套 Dart 代码编译为 iOS / Android / Web / macOS / Windows / Linux。

**核心技术栈**：Flutter 3.x (Dart 3.12+) + flutter_bloc/Cubit（状态管理）+ GoRouter（路由）+ GetIt（IoC）+ Dio + Retrofit（HTTP）+ swagger_parser（API 生成）+ cached_query（缓存）+ flutter_screenutil（响应式）+ flutter_intl（i18n）+ flutter_widget_from_html/flutter_markdown（内容渲染）+ Material 3

**代码生成工具链**：swagger_parser（API 模型）+ intl_utils（i18n）+ build_runner + freezed + json_serializable + retrofit_generator

## 关键架构认知

### Feature-First 模块化架构

```
lib/
├── main.dart                    # 入口（init + MultiBlocProvider）
├── src/
│   ├── app.dart                 # CMSApp（ScreenUtilInit + MaterialApp.router）
│   ├── init.dart                # 应用初始化（环境变量、传输层、仓库）
│   ├── app_router/              # GoRouter 路由配置 + 路由名称常量
│   ├── core/                    # 核心基础设施
│   │   ├── config/              #   environments.dart（环境变量）
│   │   ├── constants/           #   breakpoints / router_paths
│   │   ├── preference/          #   UserPreferenceCache（SharedPreferences）
│   │   ├── repositories/        #   user_auth_cache（登录态 + Token）
│   │   ├── services/            #   base_service（统一错误处理）+ pagination_query
│   │   ├── themes/              #   cubit/（AppThemeCubit）+ light/dark_theme
│   │   ├── transport/http/      #   Dio + 拦截器 + status
│   │   ├── utils/               #   responsive_utils（响应式）
│   │   └── widgets/             #   responsive_layout / web_shell_layout / 底部导航栏
│   └── features/                # ★ Feature-First 业务模块
│       ├── auth/                #   pages/ + services/
│       └── cms/
│           ├── pages/           #   home/explore/post_detail/post_list/category_list/tag_feed/search/profile/bookmarks/settings/...
│           ├── services/        #   post/category/tag/comment/navigation_service
│           └── widgets/         #   post_card/featured_carousel/content_viewer/tag_chip
├── generated/                   # [自动生成] l10n.dart + api/ + intl/
└── l10n/                        # i18n ARB 文件（intl_zh_CN.arb / intl_en_US.arb）
```

### 响应式布局（三级断点）

| 设备 | 屏宽 | 布局策略 |
|------|------|----------|
| 手机 Mobile | < 600 dp | 纵向单栏瀑布流 + 底部导航栏 |
| 平板 Tablet | 600~1024 dp | 双栏布局 |
| 网页 Web | > 1024 dp | 三栏/居中 + 持久化顶部导航栏 |

```dart
ResponsiveLayout(
  mobileBody: _buildMobileView(),
  webBody: _buildWebView(),
)
ResponsiveUtils.isMobile(context)        // 判断设备
ResponsiveUtils.postGridColumns(ctx)     // 网格列数（1/2/3）
```

**Web 端 ShellRoute 持久化导航**：Web 端通过 `ShellRoute` + `WebShellLayout` 实现贯穿所有页面的顶部导航栏。

### 三层 API 架构

```
lib/generated/api/               # [自动生成] swagger_parser 产出（RestClient + 各 ServiceClient + models）
lib/src/features/cms/services/   # [服务封装] 继承 BaseService，封装业务 + Query/Mutation
```

```dart
class PostService extends BaseService {
  PostServiceClient get _api => GetIt.instance<RestClient>().postService;

  Future<dynamic> list([PaginationQuery? query]) async {
    try {
      return await _api.postServiceList(page: q.page, pageSize: q.pageSize, query: q.queryString);
    } on DioException catch (e) {
      return handleDioError(e);  // 统一错误转换
    }
  }

  Query<ListPostResponse> listQuery([PaginationQuery? query]) { /* 缓存查询 */ }
  Mutation<Post, Post> createMutation() { /* 写操作 + 自动失效缓存 */ }
}
```

页面直接实例化 Service 调用，不额外封装 Hook。

### Dio + Retrofit + GetIt（HTTP 通信）

Dio 全局单例（通过 GetIt 注册），拦截器链：Token 注入 → Locale → 日志 / 响应：数据解构 → 401 认证 → 错误消息。`BaseService.handleDioError` 统一把 `DioException` 转 `Status`。

### 状态管理 — BLoC / Cubit

`AppThemeCubit` 管理全局状态（主题模式 / 主题色 / 语言），页面局部状态用 `StatefulWidget` + `setState`。登录状态通过 `UserAuthCache`（GetIt 单例）+ `ValueNotifier` 响应式。

### 主题系统（Material 3 + ColorScheme.fromSeed）

```dart
ThemeData getLightTheme({Color? seedColor}) {
  final colorScheme = ColorScheme.fromSeed(seedColor: seedColor ?? kDefaultSeedColor, brightness: Brightness.light);
  return ThemeData(colorScheme: colorScheme, useMaterial3: true);
}
```

支持 `light`/`dark`/`system` 三种模式 + 8 种预设主题色，Cubit 管理状态，SharedPreferences 持久化。

### 国际化（flutter_intl + ARB）

```
lib/l10n/intl_zh_CN.arb / intl_en_US.arb   # 翻译源
lib/generated/l10n.dart                     # [生成] S 类
```

```dart
Text(S.of(context).appName)              // 获取翻译
Text(S.of(context).postsCount(5))        // 带参数
```

多语言内容获取用 `translation_helpers.dart` 辅助函数（`getPostTitle(post)` 等）。

## 关键约定（必须遵守）

1. **Service 必须继承 `BaseService`** — 用 `handleDioError` 统一处理 `DioException`
2. **分页用 `PaginationQuery`** — 不要手动拼接 query 字符串
3. **禁止手改 `lib/generated/`** — swagger_parser / intl_utils / build_runner 自动生成
4. **响应式用 `ResponsiveLayout`** — 不要在一个 build 方法混用 mobile/web 视图
5. **Web 端禁止 `.w`/`.h`/`.sp`** — Web 端 ScreenUtil designSize 设为视窗尺寸（1:1），用固定值；手机端可用
6. **断点用 `Breakpoints` 常量** — 不要硬编码屏宽数值
7. **路由用 `context.go()`（顶级切换）/ `context.push()`（子页面）** — 返回用 `AppBackButton`（内置 canPop 检查）
8. **路由路径集中管理** — `router_paths.dart` + `route_names.dart`
9. **多语言内容用辅助函数** — `getPostTitle(post)` 等，不直接访问 `translations`

## 代码生成（改后必须重新生成）

| 修改内容 | 命令 |
|---|---|
| OpenAPI 定义 / Freezed 模型 / Retrofit 接口 | `dart run build_runner build --delete-conflicting-outputs` |
| ARB 翻译文件 | `flutter pub run intl_utils:generate` |

## 开发命令

```bash
flutter pub get                              # 安装依赖
flutter pub run intl_utils:generate          # 生成 i18n
dart run build_runner build --delete-conflicting-outputs  # 生成 API/模型
flutter run -d chrome                        # Web 开发
flutter run -d ios / android                 # 移动端开发
flutter build web / apk / ios / macos / windows  # 构建生产产物
flutter analyze                              # 代码分析
flutter test                                 # 测试
```

**环境变量**（`.dev.env` Debug / `.env` Release，通过 flutter_dotenv 加载）：

```env
API_BASE_URL="https://api.cms.gowind.cloud"
SSE_URL="https://sse.cms.gowind.cloud/events"
CONNECTION_TIMEOUT=3000
RECEIVE_TIMEOUT=3000
AES_KEY="f51d66a73d8a0927"
```

## 新增业务模块 Checklist（以"产品"为例）

```
- [ ] Step 1: 生成 API 客户端（更新 OpenAPI 后 dart run build_runner build --delete-conflicting-outputs）
- [ ] Step 2: 封装服务层（lib/src/features/cms/services/product_service.dart，继承 BaseService）
- [ ] Step 3: 创建页面（lib/src/features/cms/pages/product/product_page.dart，StatefulWidget + ResponsiveLayout）
- [ ] Step 4: 注册路由（app_router.dart 的 ShellRoute.routes + route_names.dart + router_paths.dart）
- [ ] Step 5: 添加导航入口（WebShellLayout 或 AppBottomNavBar）
```

## 新增语言

1. 在 `lib/l10n/` 创建 `intl_ja_JP.arb`（复制翻译）
2. `flutter pub run intl_utils:generate` 生成代码
3. 在 `settings_page.dart` 的 `localeLabels` 添加语言标签

## Web 端特殊注意事项

- **GridView 限制**：避免在 `CustomScrollView` 的 `SliverToBoxAdapter` 中嵌套 `GridView` + `NeverScrollableScrollPhysics`（触发 viewport hitTest null 错误），改用 `LayoutBuilder` + `Row`/`Column`
- **嵌套 Scaffold**：Web 端 `ShellRoute` 已提供外层 Scaffold，顶级页面需注意滚动冲突
- **导航栏**：Web 端 `WebShellLayout` 提供持久化顶部导航栏，页面无需再显示 AppBar

## 对接不同后端

1. `.dev.env` / `.env` — 修改 `API_BASE_URL`
2. `swagger_parser.yaml` — 更新 OpenAPI 定义 URL
3. 重新生成 API 代码 — `dart run build_runner build --delete-conflicting-outputs`
4. `services/*.dart` — 调整请求参数/响应结构
5. `interceptors/` — 自定义认证流程

## 常见错误与纠正

| 错误做法 | 正确做法 |
|---|---|
| 手改 `lib/generated/` | 改 OpenAPI/ARB 源后重新生成 |
| Service 不继承 `BaseService` | 继承并用 `handleDioError` 处理错误 |
| 手动拼接分页 query 字符串 | 用 `PaginationQuery` 封装 |
| 一个 build 方法混用 mobile/web | 用 `ResponsiveLayout` 双视图 |
| Web 端用 `.w`/`.h`/`.sp` | Web 端用固定值（手机端才用 ScreenUtil） |
| 硬编码屏宽数值 | 用 `Breakpoints` 常量 |
| 直接访问 `post.translations` | 用 `getPostTitle(post)` 辅助函数 |
| Web 端 SliverToBoxAdapter 嵌套 GridView | 用 `LayoutBuilder` + `Row`/`Column` |
| 返回按钮不检查 canPop | 用 `AppBackButton` 组件 |
| 路由路径硬编码 | 集中管理在 `router_paths.dart` |
