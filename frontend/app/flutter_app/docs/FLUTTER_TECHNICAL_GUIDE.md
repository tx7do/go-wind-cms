# 基于 Flutter 的 Headless CMS 全平台前端架构：技术解析与二次开发导引

> 本文面向希望基于此项目进行二次开发的 Flutter 工程师，从技术栈选型、核心架构设计、关键模块实现到二开实践路径，提供一份完整的技术地图。

---

## 一、技术栈总览

本项目是一个 **Flutter 全平台** CMS 内容展示前端，一套 Dart 代码同时编译为 iOS、Android、Web、macOS、Windows 等多端应用：

| 层面       | 技术                       | 版本               | 用途                      |
|----------|--------------------------|------------------|-------------------------|
| 框架       | Flutter                  | 3.x (Dart 3.12+) | 全平台 UI 框架               |
| 语言       | Dart                     | 3.x              | 类型安全 + 空安全              |
| 状态管理     | flutter_bloc + Cubit     | 9.x              | 响应式状态（Cubit 模式）         |
| 服务定位     | GetIt                    | 9.x              | 轻量 IoC 容器（单例管理）         |
| 路由       | GoRouter                 | 17.x             | 声明式路由 + ShellRoute      |
| 国际化      | flutter_intl (intl)      | 0.21.x           | ARB 翻译文件 + 代码生成         |
| HTTP 客户端 | Dio + Retrofit           | 5.x / 4.x        | REST 通信 + 类型安全客户端       |
| API 代码生成 | swagger_parser           | 1.43.x           | OpenAPI → Dart 模型 + 客户端 |
| 数据缓存     | cached_query             | 3.x              | Query/Mutation 缓存管理     |
| 响应式适配    | flutter_screenutil       | 5.x              | 手机端设计稿适配（Web 端禁用）       |
| Markdown | flutter_markdown         | 0.7.x            | 内容解析                    |
| HTML 渲染  | flutter_widget_from_html | 0.17.x           | 富文本渲染                   |
| 加密       | encrypt + crypto         | 5.x / 3.x        | AES 加密（Token 持久化）       |
| JWT      | jose                     | 0.3.x            | JWT Token 解析            |
| 图片缓存     | cached_network_image     | 3.x              | 网络图片缓存                  |
| 骨架屏      | shimmer                  | 3.x              | 加载占位动画                  |
| 环境变量     | flutter_dotenv           | 6.x              | .env 文件管理               |
| 日志       | logger                   | 2.x              | 结构化日志                   |

**代码生成工具链：** `swagger_parser`（API 模型） + `intl_utils`（国际化） + `build_runner` + `freezed` + `json_serializable` + `retrofit_generator`

---

## 二、核心架构设计

### 2.1 全平台编译模式

项目基于 Flutter 3.x，通过 `flutter` CLI 将同一份 Dart 代码编译为多端产物：

```bash
# Web 开发
flutter run -d chrome

# iOS 开发
flutter run -d ios

# Android 开发
flutter run -d android

# 构建产物
flutter build web          # → build/web/ (SPA)
flutter build apk          # → build/app/outputs/ (APK)
flutter build ios          # → build/ios/ (IPA)
flutter build macos        # → build/macos/ (macOS App)
flutter build windows      # → build/windows/ (Windows App)
```

**环境配置**（`.dev.env` / `.env`）：

```env
API_BASE_URL="https://api.cms.gowind.cloud"    # 后端 API 地址
SSE_URL="https://sse.cms.gowind.cloud/events"   # SSE 推送地址
CONNECTION_TIMEOUT=3000                          # 连接超时（毫秒）
RECEIVE_TIMEOUT=3000                             # 接收超时（毫秒）
AES_KEY="f51d66a73d8a0927"                       # AES 加密密钥
SENTRY_DSN="https://ingest.sentry.io/"           # Sentry 异常监控
```

> 环境变量通过 `flutter_dotenv` 加载，在 `Environments` 类中统一导出。Debug 模式加载 `.dev.env`，Release 模式加载 `.env`。

### 2.2 全平台路由架构

使用 **GoRouter** 实现声明式路由，支持嵌套路由、重定向和 ShellRoute 持久化布局：

```dart
// app_router.dart
static final router = GoRouter(
  initialLocation: '/',
  routes: [
    // Web 端持久化顶部导航栏 Shell
    ShellRoute(
      builder: (context, state, child) {
        if (ResponsiveUtils.isMobile(context)) {
          return child;  // 移动端直接透传
        }
        return WebShellLayout(  // Web 端包裹持久化导航
          currentPath: state.uri.path,
          child: child,
        );
      },
      routes: [
        GoRoute(path: '/', name: 'home', builder: ...),
        GoRoute(path: '/post/:id', name: 'postDetail', builder: ...),
        GoRoute(path: '/explore', name: 'explore', builder: ...),
        GoRoute(path: '/settings', name: 'settings', builder: ...),
        // ... 更多路由
      ],
    ),
    // 登录页（不在 Shell 内）
    GoRoute(path: '/login', ...),
  ],
);
```

**路由常量集中管理**：

```dart
// router_paths.dart
class AppRoutePath {
  static const initial = '/';
  static const login = '/login';
  static const explore = '/explore';
  static const profile = '/profile';
  static const settings = '/settings';
  // ...
}
```

**路由目录结构**：

```
lib/src/app_router/
├── app_router.dart          # GoRouter 配置（ShellRoute + 所有路由定义）
└── route_names.dart         # 路由名称常量
```

### 2.3 三层 API 架构

API 层遵循**生成层 → 服务层 → 调用层**的三层分离架构：

```
lib/generated/api/           # [自动生成] swagger_parser 产出
├── rest_client.dart         #   Retrofit 总客户端
├── post_service/            #   文章服务客户端
├── models/                  #   所有请求/响应模型
└── ...

lib/src/features/cms/services/ # [服务封装] 业务逻辑、错误处理、Query/Mutation
├── post_service.dart        #   文章服务
├── category_service.dart    #   分类服务
├── tag_service.dart         #   标签服务
├── comment_service.dart     #   评论服务
├── navigation_service.dart  #   导航服务
└── ...
```

**第一层 — 自动生成的客户端**（`generated/api/`）：由 OpenAPI 定义通过 `swagger_parser` 自动生成的 Retrofit 客户端代码，不应手动编辑。

**第二层 — 服务封装**（`services/`）：基于生成的客户端封装业务方法，统一错误处理和分页参数：

```dart
// services/post_service.dart
class PostService extends BaseService {
  PostServiceClient get _api => GetIt.instance<RestClient>().postService;

  // 直接调用（适合页面简单场景）
  Future<dynamic> list([PaginationQuery? query]) async {
    try {
      return await _api.postServiceList(
        page: q.page,
        pageSize: q.pageSize,
        query: q.queryString,
      );
    } on DioException catch (e) {
      return handleDioError(e);  // 统一错误转换
    }
  }

  // Query 缓存（适合需要缓存/刷新的场景）
  Query<ListPostResponse> listQuery([PaginationQuery? query]) {
    return Query<ListPostResponse>(
      key: 'posts',
      queryFn: () => _api.postServiceList(...),
    );
  }

  // Mutation（适合写操作，自动失效缓存）
  Mutation<Post, Post> createMutation() {
    return Mutation<Post, Post>(
      mutationFn: (post) => _api.postServiceCreate(body: ...),
      invalidateQueries: ['posts'],
    );
  }
}
```

**第三层 — 页面调用**：页面组件直接实例化 Service 并调用方法，不额外封装 Hook。

### 2.4 Dio + Retrofit — HTTP 通信内核

基于 Dio 封装的全局单例（通过 GetIt 注册），通过拦截器链实现完整的认证生命周期：

```
请求拦截链：Token 注入 → Locale 注入 → 日志
响应拦截链：数据解构 → 401 认证处理 → 错误消息提取
```

**初始化流程**：

```dart
// init.dart
Future<void> init() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Environments.init();           // 加载 .env
  await initThirdPartyPlugins();       // 初始化 SharedPreferences 等
  _initTransport();                    // 注册 Dio + RestClient 到 GetIt
  await repos.init();                  // 初始化缓存仓库
}

void _initTransport() {
  transport.init();
  final getIt = GetIt.instance;
  getIt.registerLazySingleton<RestClient>(
    () => RestClient(GetIt.instance<Dio>()),
  );
}
```

**BaseService 错误处理**：

```dart
// base_service.dart
abstract class BaseService {
  /// 将 DioException 统一转换为 Status 对象
  Status handleDioError(DioException e) {
    final data = e.response?.data;
    if (data is Map<String, dynamic>) {
      return Status(
        code: e.response?.statusCode,
        reason: data['reason'],
        message: data['message'] ?? e.message,
      );
    }
    return Status(code: e.response?.statusCode, message: e.message);
  }
}
```

二开时如需对接不同后端，只需修改 `.env` 中的 `API_BASE_URL` 和 `generated/api/` 下的生成代码。

### 2.5 状态管理 — BLoC / Cubit 模式

采用 `flutter_bloc` 的 Cubit 模式（Bloc 的简化版），用于主题、语言等全局状态：

```dart
// main.dart
void main() async {
  await init();
  runApp(
    MultiBlocProvider(
      providers: [BlocProvider(create: (_) => AppThemeCubit())],
      child: const CMSApp(),
    ),
  );
}
```

**AppThemeCubit** 同时管理主题模式（亮色/暗色/跟随系统）、主题色（seedColor）和语言切换：

```dart
// app_theme_cubit.dart
class AppThemeCubit extends Cubit<AppThemeState> {
  // 主题模式
  modify(ThemeMode newMode) async { ... }
  // 主题色
  void modifySeedColor(Color color) async { ... }
  // 语言切换
  void modifyLocale(Locale newLocale) async { ... }
}
```

**页面级状态**：页面内的局部状态（如列表数据、加载状态）使用 `StatefulWidget` 的 `setState` 管理，不引入额外的状态管理库。

### 2.6 偏好持久化 — UserPreferenceCache

`UserPreferenceCache` 基于 `SharedPreferences` 封装，管理用户偏好设置：

```dart
class UserPreferenceCache {
  // 主题模式
  ThemeMode get themeMode;
  Future<void> setMaterialThemeMode(ThemeMode mode);
  // 主题色
  Color? get seedColorValue;
  Future<void> setSeedColor(Color color);
  // 语言
  String get language;        // "zh_CN" / "en_US"
  Future<void> setLanguage(String lang, bool notify);
}
```

主题支持三种模式：`light`（亮色）、`dark`（暗色）、`system`（跟随系统），通过 Material 3 的 `ColorScheme.fromSeed` 动态生成色板。

### 2.7 国际化体系

使用 **flutter_intl** (intl + ARB 文件) 实现国际化：

```
lib/l10n/
├── intl_zh_CN.arb       # 中文翻译（7.4KB）
└── intl_en_US.arb       # 英文翻译（7.8KB）

lib/generated/           # [自动生成] intl_utils 产出
├── l10n.dart            # S 类（统一导出所有翻译方法）
└── intl/
    ├── messages_zh_CN.dart
    └── messages_en_US.dart
```

**代码生成命令**：

```bash
flutter pub run intl_utils:generate
```

**使用方式**：

```dart
// 在 Widget 中
Text(S.of(context).appName)           // 获取当前语言的翻译
Text(S.of(context).postsCount(5))     // 带参数的翻译
```

**多语言内容获取**：后端返回的实体（Post、Category 等）携带 `translations[]` 数组，前端通过辅助函数按当前 locale 提取对应翻译：

```dart
// translation_helpers.dart
String getPostTitle(Post post) {
  final translation = getTranslation(post.translations ?? []);
  return translation?.title ?? '';
}
```

---

## 三、关键模块深度解析

### 3.1 响应式布局系统

项目支持全平台响应式，采用断点 + 双视图模式：


| 设备类型      | 屏幕宽度        | 布局策略               |
|-----------|-------------|--------------------|
| 手机 Mobile | < 600 dp    | 纵向单栏瀑布流 + 底部导航栏    |
| 平板 Tablet | 600~1024 dp | 双栏布局               |
| 网页 Web    | > 1024 dp   | 三栏/居中布局 + 持久化顶部导航栏 |


**核心组件**：

```dart
// ResponsiveLayout — 双视图切换
ResponsiveLayout(
  mobileBody: _buildMobileView(),
  webBody: _buildWebView(),
)

// ResponsiveUtils — 响应式工具
ResponsiveUtils.isMobile(context)     // 是否手机
ResponsiveUtils.postGridColumns(ctx)  // 文章网格列数（1/2/3）
```

**Web 端 ShellRoute 持久化导航**：

Web 端通过 `ShellRoute` + `WebShellLayout` 实现贯穿所有页面的顶部导航栏（与移动端底部导航栏效果一致）：

```dart
ShellRoute(
  builder: (context, state, child) {
    if (ResponsiveUtils.isMobile(context)) return child;
    return WebShellLayout(currentPath: state.uri.path, child: child);
  },
  routes: [/* 所有页面路由 */],
)
```

**ScreenUtil 策略**：
- 手机端：使用 `.w` / `.h` / `.sp` 适配不同手机分辨率
- Web 端：动态将 `designSize` 设为当前视窗尺寸，使 `.w/.h/.sp` 始终 1:1，字体不随窗口缩放

### 3.2 主题系统

采用 **Material 3 + ColorScheme.fromSeed** 方案，支持动态主题色切换：

```dart
// light_theme.dart
ThemeData getLightTheme({Color? seedColor}) {
  final colorScheme = ColorScheme.fromSeed(
    seedColor: seedColor ?? kDefaultSeedColor,
    brightness: Brightness.light,
  );
  return ThemeData(
    colorScheme: colorScheme,
    useMaterial3: true,
    // ...
  );
}
```

**预设主题色**（设置页可切换）：

```dart
static const List<Color> _presetColors = [
  Color(0xFF3A7CA5),  // 默认蓝
  Color(0xFF6750A4),  // 紫色
  Color(0xFF006B5E),  // 墨绿
  Color(0xFFB7262E),  // 中国红
  Color(0xFFF57C00),  // 橙色
  // ...
];
```

**设计要点：**
- **全平台统一**：Material 3 在 iOS/Android/Web 上表现一致
- **暗色模式**：通过 `themeMode` 切换 `theme` / `darkTheme`
- **动态主题色**：Cubit 管理状态，SharedPreferences 持久化

### 3.3 内容渲染管线

`ContentViewer` 组件实现 HTML/Markdown 内容渲染：

```
后端 HTML 内容
  ↓ flutter_widget_from_html（HTML 渲染器）
  ├── 图片 → cached_network_image 缓存加载
  ├── 链接 → url_launcher 外部跳转
  └── 自定义样式 → 响应式排版
  ↓ 渲染输出
```

### 3.4 认证流程

```
用户登录 → 存储 accessToken + refreshToken（SharedPreferences + AES 加密）
    ↓
每次请求 → Dio 拦截器注入 Authorization 头
    ↓
401 响应 → 自动调用 refreshToken 接口
    ↓ 成功               ↓ 失败
更新 Token 继续请求    清除凭证 → 跳转登录页
```

**登录状态管理**：通过 `UserAuthCache`（基于 GetIt 单例）+ `ValueNotifier` 实现响应式登录状态：

```dart
// 全局监听登录状态
ValueListenableBuilder<bool>(
  valueListenable: GetIt.instance<UserAuthCache>().loginStateNotifier,
  builder: (context, hasLogin, _) {
    return Text(hasLogin ? S.of(context).me : S.of(context).login);
  },
)
```

### 3.5 评论系统

评论模块支持多级嵌套回复，包含：

- **评论列表**（`CommentSection`）：树形结构展示，`buildCommentTree()` 递归构建
- **评论输入**（`CommentInputBar`）：支持回复指定评论
- **评论服务**（`CommentService`）：CRUD + 服务端过滤（`objectId`）

### 3.6 导航系统

后端动态配置导航项，前端根据位置渲染：

```dart
// NavigationService 加载导航数据
final result = await _navService.list();
final headerItems = getFlatNavItems(navigations, NavigationLocation.header);
final footerItems = getFlatNavItems(navigations, NavigationLocation.footer);
```

导航链接支持内部路由、外部链接和新标签页打开，通过 `NavBarLink` 组件统一渲染并高亮当前路由。

---

## 四、项目目录结构与职责

```
lib/
├── main.dart                    # 程序入口（init + MultiBlocProvider）
├── src/
│   ├── app.dart                 # CMSApp（ScreenUtilInit + MaterialApp.router）
│   ├── init.dart                # 应用初始化（环境变量、传输层、仓库）
│   ├── init_thirdparty_plugins.dart # 第三方插件初始化
│   ├── app_router/
│   │   ├── app_router.dart      #   GoRouter 路由配置
│   │   └── route_names.dart     #   路由名称常量
│   ├── core/
│   │   ├── config/
│   │   │   └── environments.dart    # 环境变量管理
│   │   ├── constants/
│   │   │   ├── breakpoints.dart     #   设备断点常量
│   │   │   ├── router_paths.dart    #   路由路径常量
│   │   │   └── index.dart           #   统一导出
│   │   ├── extensions/              # Dart 扩展方法
│   │   ├── logic/                   # 业务逻辑层
│   │   ├── preference/
│   │   │   ├── user_preference.dart      # 用户偏好定义
│   │   │   └── user_preference_cache.dart # SharedPreferences 封装
│   │   ├── repositories/            # 数据仓库（认证缓存等）
│   │   │   ├── user_auth_cache.dart     #   登录状态 + Token 管理
│   │   │   └── init.dart               #   仓库初始化
│   │   ├── services/
│   │   │   ├── base_service.dart        #   服务基类（统一错误处理）
│   │   │   └── pagination_query.dart    #   通用分页查询参数
│   │   ├── themes/
│   │   │   ├── cubit/
│   │   │   │   ├── app_theme_cubit.dart #   主题/语言 Cubit
│   │   │   │   └── app_theme_state.dart #   主题状态
│   │   │   ├── light_theme.dart        #   亮色主题配置
│   │   │   ├── dark_theme.dart         #   暗色主题配置
│   │   │   └── theme.dart              #   主题工具
│   │   ├── transport/
│   │   │   ├── http/
│   │   │   │   ├── base_api_client.dart #   API 客户端基类
│   │   │   │   ├── http_client.dart     #   Dio 初始化
│   │   │   │   ├── dio_extension.dart   #   Dio 扩展
│   │   │   │   ├── interceptors/        #   拦截器（认证、白名单等）
│   │   │   │   └── status.dart          #   统一错误状态
│   │   │   └── init.dart               #   传输层初始化
│   │   ├── utils/
│   │   │   └── responsive_utils.dart   #   响应式工具类
│   │   ├── utilities/                # 通用工具函数
│   │   └── widgets/                  # 核心共享组件
│   │       ├── responsive_layout.dart #   响应式布局切换
│   │       ├── web_shell_layout.dart  #   Web 端持久化导航
│   │       ├── app_back_button.dart   #   智能返回按钮
│   │       ├── app_bottom_nav_bar.dart #   移动端底部导航栏
│   │       ├── not_found_page.dart    #   404 页面
│   │       └── error_page.dart        #   错误页面
│   └── features/
│       ├── auth/
│       │   ├── pages/
│       │   │   └── login_page.dart    #   登录页
│       │   └── services/
│       │       └── authentication_service.dart # 认证服务
│       └── cms/
│           ├── pages/                 # CMS 页面
│           │   ├── home/              #   首页（Web 三栏 + Mobile 单栏）
│           │   ├── explore/           #   发现/分类页
│           │   ├── post_detail/       #   文章详情页
│           │   ├── post_list/         #   文章列表页
│           │   ├── category_list/     #   分类列表页
│           │   ├── tag_feed/          #   标签文章页
│           │   ├── tag_list/          #   标签列表页
│           │   ├── search/            #   搜索页
│           │   ├── profile/           #   个人中心
│           │   ├── bookmarks/         #   收藏页
│           │   ├── settings/          #   设置页
│           │   ├── my_comments/       #   我的评论
│           │   ├── about/             #   关于页
│           │   └── legal/             #   法律/信息页
│           ├── services/              # CMS 服务层
│           │   ├── post_service.dart  #   文章服务
│           │   ├── category_service.dart # 分类服务
│           │   ├── tag_service.dart    #   标签服务
│           │   ├── comment_service.dart # 评论服务
│           │   └── navigation_service.dart # 导航服务
│           └── widgets/               # CMS 共享组件
│               ├── post_card.dart     #   文章卡片
│               ├── featured_carousel.dart # 轮播图
│               ├── content_viewer.dart #  内容渲染器
│               ├── tag_chip.dart      #   标签芯片
│               └── tag_cloud.dart     #   标签云
├── generated/                        # [自动生成]
│   ├── l10n.dart                     #   国际化 S 类
│   ├── api/                          #   API 客户端 + 模型
│   └── intl/                         #   翻译消息
└── l10n/                             # 国际化 ARB 文件
    ├── intl_zh_CN.arb
    └── intl_en_US.arb
```

---

## 五、二次开发导引

### 5.1 环境搭建

```bash
# 安装依赖
flutter pub get

# 生成代码（API 模型 + 国际化）
dart run build_runner build --delete-conflicting-outputs
flutter pub run intl_utils:generate

# 启动 Web 开发
flutter run -d chrome

# 启动 iOS 开发
flutter run -d ios

# 启动 Android 开发
flutter run -d android

# 构建生产产物
flutter build web
flutter build apk
flutter build ios

# 代码分析
flutter analyze

# 运行测试
flutter test
```

**环境变量配置**（`.dev.env`）：

```env
API_BASE_URL="https://api.cms.gowind.cloud"
SSE_URL="https://sse.cms.gowind.cloud/events"
CONNECTION_TIMEOUT=3000
RECEIVE_TIMEOUT=3000
AES_KEY="f51d66a73d8a0927"
```

> 修改 `.env` 文件后需重启应用。Debug 模式加载 `.dev.env`，Release 模式加载 `.env`。

### 5.2 新增一个业务页面

以「产品」模块为例：

**Step 1 — 生成 API 客户端**

将后端 OpenAPI 定义更新后，运行代码生成：

```bash
dart run build_runner build --delete-conflicting-outputs
```

这会在 `lib/generated/api/` 下自动生成 `ProductServiceClient` 和相关模型。

**Step 2 — 封装服务层**

创建 `lib/src/features/cms/services/product_service.dart`：

```dart
class ProductService extends BaseService {
  ProductServiceClient get _api => GetIt.instance<RestClient>().productService;

  Future<dynamic> list([PaginationQuery? query]) async {
    final q = query ?? const PaginationQuery();
    try {
      return await _api.productServiceList(
        page: q.page, pageSize: q.pageSize, query: q.queryString,
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}
```

**Step 3 — 创建页面**

创建 `lib/src/features/cms/pages/product/product_page.dart`：

```dart
class ProductPage extends StatefulWidget {
  @override
  State<ProductPage> createState() => _ProductPageState();
}

class _ProductPageState extends State<ProductPage> {
  final _service = ProductService();
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _service.list();
    if (!mounted) return;
    setState(() { /* 处理数据 */ _isLoading = false; });
  }

  @override
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: _buildMobileView(),
      webBody: _buildWebView(),
    );
  }
}
```

**Step 4 — 注册路由**

在 `app_router.dart` 的 ShellRoute.routes 中添加：

```dart
GoRoute(
  name: RouteNames.product,
  path: '/product',
  builder: (context, state) => const ProductPage(),
),
```

在 `route_names.dart` 和 `router_paths.dart` 中添加对应常量。

**Step 5 — 添加导航入口**

在 `WebShellLayout` 的导航链接中添加产品页入口（如果后端导航配置了），或在移动端 `AppBottomNavBar` 中添加。

### 5.3 新增一种语言

**Step 1 — 创建翻译文件**

在 `lib/l10n/` 下创建新的 ARB 文件（如 `intl_ja_JP.arb`），复制现有文件并翻译。

**Step 2 — 生成代码**

```bash
flutter pub run intl_utils:generate
```

这会自动在 `lib/generated/` 下生成对应的翻译代码。

**Step 3 — 更新设置页**

在 `settings_page.dart` 的 `localeLabels` 中添加新语言标签：

```dart
final localeLabels = <Locale, String>{
  const Locale('zh', 'CN'): '中文(中国)',
  const Locale('en', 'US'): 'English',
  const Locale('ja', 'JP'): '日本語',  // 新增
};
```

### 5.4 自定义主题配色

修改 `settings_page.dart` 中的预设主题色列表：

```dart
static const List<Color> _presetColors = [
  Color(0xFF3A7CA5),  // 默认蓝
  Color(0xFF6750A4),  // 紫色
  // 添加更多主题色
  Color(0xFFE91E63),  // 粉红
];
```

所有使用 `theme.colorScheme.*` 的组件会自动跟随变化，因为色板由 `ColorScheme.fromSeed` 动态生成。

### 5.5 对接不同后端

本项目前端与后端通过 REST API 通信，对接不同后端的核心修改点：

1. **`.dev.env` / `.env`** — 修改 `API_BASE_URL`
2. **`swagger_parser.yaml`** — 更新 OpenAPI 定义文件 URL
3. **重新生成 API 代码** — `dart run build_runner build --delete-conflicting-outputs`
4. **`services/*.dart`** — 调整请求参数格式和响应结构

认证流程可通过修改 `interceptors/` 下的拦截器来自定义。

### 5.6 部署

**Web 部署：**

```bash
flutter build web   # 输出到 build/web/
```

产物为 SPA 静态文件，需配置服务器 fallback 到 `index.html`。

**Nginx 配置示例：**

```nginx
server {
    listen 80;
    root /var/www/cms;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

**Android 部署：**

```bash
flutter build apk --release     # 输出 APK
flutter build appbundle --release # 输出 AAB（推荐用于 Play Store）
```

**iOS 部署：**

```bash
flutter build ios --release   # 输出到 build/ios/
```

通过 Xcode 打开 `Runner.xcworkspace`，签名后上传至 App Store Connect。

---

## 六、开发规范与注意事项

### 6.1 响应式适配规范

- 使用 `ResponsiveLayout` 组件区分 mobile/web 视图，不要在一个 build 方法中混用
- 手机端可使用 `flutter_screenutil` 的 `.w` / `.h` / `.sp` 适配
- **Web 端禁止使用 `.w` / `.h` / `.sp`**（Web 端 ScreenUtil designSize 被设为视窗尺寸，1:1 映射，使用固定值）
- 使用 `ResponsiveUtils.isMobile(context)` 判断设备类型
- 使用 `Breakpoints` 常量定义断点值，不要硬编码

### 6.2 API 服务规范

- 所有 Service 必须继承 `BaseService`，使用 `handleDioError` 统一处理 DioException
- 使用 `PaginationQuery` 封装分页参数，不要手动拼接 query 字符串
- 分类数据获取后按 `sortOrder` 排序：`_categories.sort((a, b) => (a.sortOrder ?? 0).compareTo(b.sortOrder ?? 0))`

### 6.3 路由注意事项

- 使用 `context.go()` 进行顶级路由切换（如首页 → 设置页）
- 使用 `context.push()` 进行子页面导航（如首页 → 文章详情）
- 返回按钮使用 `AppBackButton`（内置 canPop 检查 + fallback）
- 路由路径常量集中管理在 `router_paths.dart` 和 `route_names.dart`

### 6.4 Web 端特殊注意事项

- **GridView 限制**：Web 端避免在 `CustomScrollView` 的 `SliverToBoxAdapter` 中嵌套 `GridView` + `NeverScrollableScrollPhysics`，会触发 viewport hitTest null 错误。改用 `LayoutBuilder` + `Row`/`Column` 手动网格布局
- **嵌套 Scaffold**：Web 端 `ShellRoute` 已提供外层 Scaffold，顶级页面的 Scaffold 在 Web 端会形成嵌套，需注意滚动冲突
- **导航栏**：Web 端通过 `WebShellLayout` 提供持久化顶部导航栏，页面无需再显示 AppBar

### 6.5 多语言内容

获取后端实体的多语言字段时，使用 `translation_helpers.dart` 中的辅助函数：

```dart
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

final title = getPostTitle(post);         // 自动匹配当前语言
final summary = getPostSummary(post);
final categoryName = getCategoryName(cat);
final tagName = getTagName(tag);
```

### 6.6 环境变量

环境变量通过 `flutter_dotenv` 加载，在 `Environments` 类中统一导出。`.env` 文件中键名自由命名（不要求前缀），通过 `dotenv.get('KEY')` 读取。

### 6.7 代码生成

修改以下文件后需重新生成代码：

| 修改内容          | 生成命令                                                       |
|---------------|------------------------------------------------------------|
| OpenAPI 定义更新  | `dart run build_runner build --delete-conflicting-outputs` |
| ARB 翻译文件更新    | `flutter pub run intl_utils:generate`                      |
| Freezed 模型变更  | `dart run build_runner build --delete-conflicting-outputs` |
| Retrofit 接口变更 | `dart run build_runner build --delete-conflicting-outputs` |

---

## 七、技术亮点总结

1. **一套代码全平台运行**：Flutter 框架实现 iOS + Android + Web + macOS + Windows 的统一代码库
2. **类型安全的 API 层**：OpenAPI 自动生成 → Retrofit 类型安全客户端 → BaseService 统一错误处理
3. **响应式断点设计**：三级断点（手机/平板/Web）+ `ResponsiveLayout` 双视图组件
4. **ShellRoute 持久化导航**：Web 端顶部导航栏贯穿所有页面，与移动端底部导航栏效果一致
5. **Material 3 动态主题**：`ColorScheme.fromSeed` 实现动态主题色，8 种预设色一键切换
6. **Token AES 加密持久化**：SharedPreferences + AES 对称加密，安全存储认证凭证
7. **flutter_bloc Cubit 模式**：轻量级状态管理，主题/语言状态全局共享
8. **PaginationQuery 统一分页**：参考 TS 端设计，封装 locale 自动注入、query JSON 序列化

---

> **快速开始**：`flutter pub get && flutter pub run intl_utils:generate && flutter run -d chrome`，打开 `http://localhost:xxxx` 即可运行。
