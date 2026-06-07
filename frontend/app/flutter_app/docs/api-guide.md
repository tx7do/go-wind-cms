# API 层开发指南

本文档面向二开人员，介绍 GoWind CMS Flutter 端的 API 架构设计、代码生成流程、Service 编写规范和分页查询机制。

---

## 1. 整体架构

```
┌─────────────────────────────────────────────────────────┐
│  UI 层（Page / Widget）                                  │
│  调用 Service 方法，处理返回数据                            │
├─────────────────────────────────────────────────────────┤
│  Service 层（lib/src/features/cms/services/）             │
│  业务封装，异常处理，typedef 短类名                         │
│  继承 BaseService → handleDioError()                     │
├─────────────────────────────────────────────────────────┤
│  生成层（lib/generated/api/）                             │
│  RestClient → XxxServiceClient（Retrofit 接口）           │
│  models/（JSON 序列化模型）                                │
├─────────────────────────────────────────────────────────┤
│  传输层（lib/src/core/transport/）                        │
│  Dio 实例 + 拦截器（认证、日志）                            │
│  环境配置（.env → Environments.apiBaseUrl）                │
└─────────────────────────────────────────────────────────┘
```

**调用链路：** `Page → Service → RestClient.XxxService → Dio → HTTP`

---

## 2. 代码生成流程

后端 OpenAPI 规范文件自动生成 Dart API 客户端代码，无需手写网络请求。

### 2.1 入口命令

```bash
# 方式一：swagger_parser（推荐，Dart 原生，无需 Java）
dart run bin/clean_and_gen.dart

# 方式二：openapi-generator-cli（需要 Java）
dart run bin/clean_and_gen.dart --openapi
```

### 2.2 生成流程

```
后端 openapi.yaml
  ↓ YAML → JSON 清洗（递归删除 null，修复 swagger_parser bug）
  ↓ swagger_parser 解析
  ↓ 生成 RestClient + ServiceClient + Models
  ↓ build_runner 编译 .g.dart（JSON 序列化）
lib/generated/api/
```

### 2.3 生成产物

```
lib/generated/api/
├── rest_client.dart                    # 聚合入口，持有所有 ServiceClient
├── post_service/
│   ├── post_service_client.dart        # Retrofit 接口（@GET/@POST/@PUT/@DELETE）
│   └── post_service_client.g.dart      # Retrofit 生成实现
├── models/
│   ├── content_service_v1_post.dart    # 数据模型
│   ├── content_service_v1_post.g.dart  # JSON 序列化
│   └── ...                             # 其他模型
└── ...其他 service 目录
```

### 2.4 配置文件

| 文件                       | 作用                                           |
|--------------------------|----------------------------------------------|
| `swagger_parser.yaml`    | swagger_parser 配置（源文件路径、输出目录、语言）             |
| `generator_config.yaml`  | openapi-generator-cli 配置（仅 `--openapi` 模式使用） |
| `bin/clean_and_gen.dart` | 一键清洗 + 生成脚本                                  |

### 2.5 重要原则

- **禁止手动编辑** `lib/generated/api/` 下的任何文件
- 修改后端 API 后，重新运行 `dart run bin/clean_and_gen.dart`
- 生成的模型类名较长（如 `ContentServiceV1Post`），通过 typedef 映射短类名

---

## 3. 传输层

### 3.1 Dio 初始化

文件：`lib/src/core/transport/http/http_client.dart`

```dart
Dio createDio() {
  final dio = Dio();
  dio.options.baseUrl = Environments.apiBaseUrl;    // 从 .env 读取
  dio.options.connectTimeout = Environments.connectionTimeout;
  dio.options.receiveTimeout = Environments.receiveTimeout;
  dio.options.responseType = ResponseType.json;
  // ...
  return dio;
}
```

Dio 实例通过 GetIt 注册为全局单例：

```dart
// lib/src/core/transport/init.dart
getIt.registerLazySingleton<Dio>(() => createDio());
```

### 3.2 环境配置

文件：`lib/src/core/config/environments.dart`

通过 `flutter_dotenv` 从 `.env`（生产）或 `.dev.env`（开发）加载：

```env
# .dev.env
API_BASE_URL=https://api.cms.gowind.cloud
CONNECTION_TIMEOUT=3000
RECEIVE_TIMEOUT=3000
```

### 3.3 RestClient 聚合入口

文件：`lib/generated/api/rest_client.dart`

所有 ServiceClient 的聚合入口，通过 GetIt 全局获取：

```dart
// 获取 RestClient 单例
final api = GetIt.instance<RestClient>();

// 懒加载各 ServiceClient
api.postService        // PostServiceClient
api.categoryService    // CategoryServiceClient
api.commentService     // CommentServiceClient
api.tagService         // TagServiceClient
api.pageService        // PageServiceClient
api.navigationService  // NavigationServiceClient
api.fileTransferService // FileTransferServiceClient
api.authenticationService // AuthenticationServiceClient
api.userProfileService // UserProfileServiceClient
```

---

## 4. Service 层编写规范

### 4.1 标准 Service 结构

以 `PostService` 为例：

```dart
// lib/src/features/cms/services/post_service.dart

// ① import 生成模型，用 typedef 定义短类名
typedef Post = ContentServiceV1Post;
typedef ListPostResponse = ContentServiceV1ListPostResponse;
typedef CreatePostRequest = ContentServiceV1CreatePostRequest;
typedef UpdatePostRequest = ContentServiceV1UpdatePostRequest;

// ② 继承 BaseService
class PostService extends BaseService {
  PostService() : super(tag: 'PostService');

  // ③ 通过 GetIt 获取 RestClient，访问对应的 ServiceClient
  PostServiceClient get _api => GetIt.instance<RestClient>().postService;

  // ④ 直接调用方法（try-catch + handleDioError）
  Future<dynamic> list([PaginationQuery? query]) async {
    final q = query ?? const PaginationQuery();
    try {
      return await _api.postServiceList(
        page: q.page,
        pageSize: q.pageSize,
        noPaging: q.noPaging,
        orderBy: q.orderByString,
        query: q.queryString,
      );
    } on DioException catch (e) {
      return handleDioError(e);   // Dio 异常 → Status 对象
    }
  }
}
```

### 4.2 typedef 短类名规范

生成模型的完整类名较长（如 `ContentServiceV1Category`），在 Service 文件顶部统一用 typedef 映射：

```dart
// 在 service 文件顶部
typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Comment = CommentServiceV1Comment;
```

外部使用时直接 `import` service 文件，通过 typedef 使用短类名。

### 4.3 异常处理模式

所有 Service 方法统一使用 `try-catch + handleDioError`：

```dart
Future<dynamic> get(int id) async {
  try {
    return await _api.postServiceGet(id: id);
  } on DioException catch (e) {
    return handleDioError(e);  // → Status(code, reason, message)
  }
}
```

`handleDioError` 定义在 `BaseService` 中，将 Dio 异常转为统一的 `Status` 对象：

```dart
// lib/src/core/services/base_service.dart
Status handleDioError(DioException e) {
  final data = e.response?.data;
  if (data is Map<String, dynamic>) {
    return Status(
      code: e.response?.statusCode,
      reason: data['reason'],
      message: data['message'],
    );
  }
  return Status(code: e.response?.statusCode, message: e.message);
}
```

### 4.4 返回值约定

| 方法类型         | 成功返回                     | 失败返回        |
|--------------|--------------------------|-------------|
| 查询（list/get） | `XxxResponse` 或 `Xxx` 模型 | `Status` 对象 |
| 创建（create）   | 创建的资源模型                  | `Status` 对象 |
| 更新（update）   | 更新后的资源模型                 | `Status` 对象 |
| 删除（delete）   | `null`                   | `Status` 对象 |

**调用方需要判断返回类型：**

```dart
final result = await _postService.list(query);
if (result is ListPostResponse) {
  // 成功，使用 result.items
} else if (result is Status) {
  // 失败，显示 result.message
}
```

---

## 5. 分页查询（PaginationQuery）

### 5.1 概述

文件：`lib/src/core/services/pagination_query.dart`

统一封装所有 List API 的查询参数，参考 TS 端的 `PaginationQuery` 设计。

### 5.2 核心参数

```dart
PaginationQuery(
  page: 1,                    // 页码（从 1 开始）
  pageSize: 10,               // 每页条数
  formValues: {               // 过滤条件 → 序列化为 JSON query 字符串
    'category_id': 5,
    'status': 'published',
  },
  orderBy: ['-created_at'],   // 排序（默认 ['-created_at']）
  fieldMask: 'id,title',      // 字段白名单（SELECT）
  locale: 'zh-CN',            // 语言（不传则自动从用户偏好获取）
  skipLocale: false,           // 是否跳过 locale 注入
  isTenantUser: false,         // 是否清理租户字段
);
```

### 5.3 计算属性

| 属性                   | 说明                                               |
|----------------------|--------------------------------------------------|
| `noPaging`           | `page == null && pageSize == null` 时为 true（全量加载） |
| `queryString`        | `formValues` + `locale` 合并为 JSON 字符串，自动清理 null   |
| `orderByString`      | 排序列表 → JSON 数组字符串                                |
| `formattedFieldMask` | `List<String>` → 逗号分隔，`String` → 原样              |

### 5.4 常用场景

```dart
// 全量加载（不分页）
_postService.list(const PaginationQuery());

// 分页加载
_postService.list(PaginationQuery(page: 1, pageSize: 10));

// 带过滤条件
_postService.list(PaginationQuery(
  formValues: {'category_id': 5},
));

// 评论按 objectId 过滤（跳过 locale）
_commentService.list(PaginationQuery(
  skipLocale: true,
  formValues: {'object_id': postId},
));

// 下一页
query = query.nextPage();
```

### 5.5 locale 自动注入

`PaginationQuery` 会自动从 `UserPreferenceCache` 读取用户语言偏好，注入到 `query` 的 `locale` 字段。

**何时需要 `skipLocale: true`：**
- Comment API 不支持 locale 参数
- 其他不涉及多语言内容的 API

---

## 6. 认证拦截器

文件：`lib/src/core/transport/http/interceptors/authentication_interceptor.dart`

```dart
// 注册认证拦截器（可选）
registerInterceptor(AuthenticationInterceptor(
  authService: myAuthService,
  autoRefreshToken: true,
));
```

功能：
- 自动在请求头添加 `Authorization: Bearer <token>`
- 401 响应时自动刷新 token
- 刷新失败触发 `authenticationFailed()` 回调

---

## 7. 文件传输

文件：`lib/src/features/cms/services/file_transfer.dart`

支持两种上传模式：

### 7.1 普通上传

```dart
final fileService = FileTransferService();
final response = await fileService.uploadFile(
  fileBytes: bytes,
  fileName: 'photo.jpg',
  bucketName: 'uploads',
);
```

### 7.2 预签名直传（推荐，支持进度回调）

```dart
final response = await fileService.uploadFilePresigned(
  fileBytes: bytes,
  fileName: 'photo.jpg',
  presignMethod: 'PUT',
  presignExpireSeconds: 3600,
  onSendProgress: (sent, total) {
    print('进度: ${(sent / total * 100).toStringAsFixed(0)}%');
  },
);
```

流程：`获取预签名 URL → 直接 PUT 到存储服务 → 返回文件信息`

---

## 8. 新增 API 接入 Checklist

当后端新增或修改 API 后，按以下步骤操作：

1. **重新生成代码**
   ```bash
   dart run bin/clean_and_gen.dart
   ```

2. **检查生成产物** — 确认 `lib/generated/api/` 下有对应的 ServiceClient 和 Model

3. **新建或更新 Service**
   - 在 `lib/src/features/cms/services/` 下新建或编辑 Service 文件
   - 继承 `BaseService`
   - 顶部添加 `typedef` 短类名
   - 实现 CRUD 方法，统一 `try-catch + handleDioError`

4. **使用 Service**
   - 在 Page/Widget 中 import Service
   - 调用方法，判断返回类型（模型 vs Status）

---

## 9. 现有 Service 清单

| Service             | 文件                          | 说明         |
|---------------------|-----------------------------|------------|
| PostService         | `post_service.dart`         | 帖子/文章 CRUD |
| CategoryService     | `category_service.dart`     | 分类 CRUD    |
| TagService          | `tag_service.dart`          | 标签 CRUD    |
| CommentService      | `comment_service.dart`      | 评论查询       |
| PageService         | `page_service.dart`         | 页面 CRUD    |
| NavigationService   | `navigation_service.dart`   | 导航管理       |
| FileTransferService | `file_transfer.dart`        | 文件上传/下载    |
| UserProfileService  | `user_profile_service.dart` | 用户资料管理     |

---

## 10. 相关文件索引

| 文件路径                                           | 说明                |
|------------------------------------------------|-------------------|
| `bin/clean_and_gen.dart`                       | API 代码生成脚本        |
| `swagger_parser.yaml`                          | swagger_parser 配置 |
| `lib/generated/api/rest_client.dart`           | API 聚合入口          |
| `lib/generated/api/models/`                    | 数据模型目录            |
| `lib/src/core/services/base_service.dart`      | Service 基类（异常处理）  |
| `lib/src/core/services/pagination_query.dart`  | 分页查询封装            |
| `lib/src/core/transport/http/http_client.dart` | Dio 初始化           |
| `lib/src/core/transport/init.dart`             | 传输层注册             |
| `lib/src/core/config/environments.dart`        | 环境变量              |
| `.dev.env` / `.env`                            | 环境配置文件            |
