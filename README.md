<div align="center">

# GoWind Content Hub

### 风行 · 开箱即用的企业级前后端一体内容平台

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![React](https://img.shields.io/badge/React-19.x-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![Flutter](https://img.shields.io/badge/Flutter-3.x-02569B?logo=flutter&logoColor=white)](https://flutter.dev/)
[![Kratos](https://img.shields.io/badge/Kratos-2.9-00ADD8?logo=go&logoColor=white)](https://go-kratos.dev/)

**[English](./README.en-US.md)** · **中文** · **[日本語](./README.ja-JP.md)**

</div>

---

风行（GoWind HCH）是一款开箱即用的企业级 Golang 全栈 Headless 内容平台（HCH = Headless Content Hub，无头内容中枢），为企业提供灵活、可扩展的全域内容管理与分发解决方案。

**核心亮点：**

- **API 优先** — 完整的 RESTful / gRPC 双协议接口，OpenAPI 文档自动生成
- **多端适配** — 前台支持 Vue、React、Taro（小程序）、Flutter 四套前端
- **多租户架构** — 租户数据隔离，自动初始化部门、角色与管理员
- **精细化权限** — 菜单权限、接口权限、数据权限三级管控，支持按钮级控制
- **原生多语言** — 内容、菜单、提示信息统一翻译管理，支撑出海业务
- **微服务架构** — 基于 go-kratos，支持服务发现、链路追踪、分布式事务

## 演示地址

| 演示类型           | 访问地址                                                                                 |
|:---------------|:-------------------------------------------------------------------------------------|
| 管理后台           | [https://admin.cms.gowind.cloud](https://admin.cms.gowind.cloud)                     |
| 后台 API Swagger | [https://api.admin.cms.gowind.cloud/docs/](https://api.admin.cms.gowind.cloud/docs/) |
| 前台 API Swagger | [https://api.cms.gowind.cloud/docs/](https://api.cms.gowind.cloud/docs/)             |
| Vue 前台         | [https://vue.cms.gowind.cloud](https://vue.cms.gowind.cloud)                         |
| React 前台       | [https://react.cms.gowind.cloud](https://react.cms.gowind.cloud)                     |
| Taro 前台        | [https://taro.cms.gowind.cloud](https://taro.cms.gowind.cloud)                       |
| Flutter 前台     | [https://flutter.cms.gowind.cloud](https://flutter.cms.gowind.cloud)                 |

> 所有演示站点通用账号：`admin` / `admin`

## 技术栈

### 后端

| 层级     | 技术                                                                                   | 说明          |
|:-------|:-------------------------------------------------------------------------------------|:------------|
| 语言     | [Go 1.25+](https://go.dev/)                                                          | 高性能编译型语言    |
| 框架     | [go-kratos](https://go-kratos.dev/)                                                  | B站开源微服务框架   |
| 依赖注入   | [Wire](https://github.com/google/wire)                                               | 编译时依赖注入     |
| ORM    | [Ent](https://entgo.io/)                                                             | Go 实体框架     |
| 数据库    | [PostgreSQL](https://www.postgresql.org/) / [MySQL](https://www.mysql.com/)          | 关系型数据库      |
| 缓存     | [Redis](https://redis.io/)                                                           | 内存数据库       |
| 对象存储   | [MinIO](https://min.io/)                                                             | 兼容 S3 的对象存储 |
| 搜索引擎   | [OpenSearch](https://opensearch.org/)                                                | 全文检索        |
| 服务注册   | [Etcd](https://etcd.io/) / [Consul](https://www.consul.io/)                          | 服务发现与配置     |
| 链路追踪   | [Jaeger](https://www.jaegertracing.io/) + [OpenTelemetry](https://opentelemetry.io/) | 分布式可观测      |
| API 定义 | [Protobuf](https://protobuf.dev/) + [buf.build](https://buf.build/)                  | 接口契约优先      |
| 权限引擎   | [Casbin](https://casbin.org/) / [OPA](https://www.openpolicyagent.org/)              | 策略驱动鉴权      |

### 管理后台前端

| 技术                                            | 说明         |
|:----------------------------------------------|:-----------|
| [Vue 3](https://vuejs.org/)                   | 渐进式前端框架    |
| [TypeScript](https://www.typescriptlang.org/) | 类型安全       |
| [Ant Design Vue](https://antdv.com/)          | 企业级 UI 组件库 |
| [Vben Admin](https://doc.vben.pro/)           | 后台管理框架     |

### 前台展示前端

| 版本      | 技术栈                                                                                   | 适用场景         |
|:--------|:--------------------------------------------------------------------------------------|:-------------|
| Vue     | [Nuxt.js](https://nuxt.com/) + [shadcn-vue](https://www.shadcn-vue.com/)              | Web 应用 / SSR |
| React   | [Next.js](https://nextjs.org/) + [shadcn/ui](https://ui.shadcn.com/)                  | Web 应用 / SSR |
| Taro    | [Taro](https://docs.taro.zone/en/docs/) + React + [shadcn/ui](https://ui.shadcn.com/) | 微信小程序 / H5   |
| Flutter | [Flutter](https://flutter.dev/) + [BLoC](https://bloclibrary.dev/)                    | 跨平台原生应用      |

## 核心功能

### 组织与权限

| 功能    | 说明                                      |
|:------|:----------------------------------------|
| 多租户管理 | 租户新增、启用/禁用、套餐配置与数据隔离；新增租户自动初始化部门、角色与管理员 |
| 用户管理  | 用户全生命周期管理，支持多角色、多部门绑定，一键代登录，高级条件查询      |
| 角色管理  | 角色与角色分组管理，精细化配置菜单权限、接口权限与数据权限           |
| 权限管理  | 权限分组、菜单节点与权限点管理，支持按钮级细粒度控制，完整权限闭环       |
| 菜单管理  | 可视化菜单配置，支持目录、页面、按钮三级，按权限动态渲染            |
| 部门管理  | 多级部门树形管理，联动绑定用户，划分组织层级                  |

### 内容与站点

| 功能    | 说明                                          |
|:------|:--------------------------------------------|
| 内容建模  | 可视化自定义内容模型，支持文本、数字、富文本、图片、文件、关联等字段类型        |
| 内容管理  | 内容增删改查、发布/下架、置顶、排序、回收站与批量操作；富文本/Markdown 编辑 |
| 分类管理  | 多级分类树形管理，支持绑定内容模型，前台按分类快速筛选                 |
| 标签管理  | 标签增删改查与内容关联，支持按标签检索与聚合展示                    |
| 评论管理  | 评论审核、删除、回复、屏蔽，可按内容、用户、时间筛选                  |
| 多语言管理 | 语种管理、内容/菜单/提示统一翻译，原生国际化支持                   |
| 站点管理  | 多站点独立配置，独立域名、标题、Logo、SEO 与展示风格              |
| 站点配置  | 基础信息、SEO、上传限制、缓存策略、邮件/短信等全局参数配置             |

### 系统与运维

| 功能           | 说明                                    |
|:-------------|:--------------------------------------|
| 文件资源管理       | 图片/文档/视频统一管理，支持本地或 OSS 云存储，预览、下载、分组管理 |
| 字典管理         | 数据字典大类与子项管理，联动查询、排序、导入导出              |
| 接口管理         | 后端接口统一维护与自动同步，树形展示，可配置请求参数与响应结果       |
| 任务调度         | 定时任务管理，支持启动/暂停/立即执行，查看执行记录与运行日志       |
| 消息通知         | 多级消息分类，向指定用户发送消息，查看已读状态与已读时间          |
| 站内信          | 个人消息中心，支持查看、删除、单条/全部已读                |
| 缓存管理         | 缓存实时查询，按键精准清除或批量清理                    |
| 登录日志         | 登录成功/失败日志，含账号、IP、归属地、设备、时间，支持查询导出     |
| 操作日志         | 全链路操作日志，记录操作人、IP、请求参数与结果，支持详情追溯       |
| 个人中心         | 个人信息编辑、头像修改、密码重置、登录记录查看               |
| Headless API | 完整 OpenAPI 接口，支持内容 CRUD，适配多端调用        |

## 项目结构

```
go-wind-cms/
├── backend/                        # 后端服务
│   ├── api/                        # Protobuf API 定义 & 生成代码
│   │   ├── protos/                 # .proto 源文件
│   │   └── gen/                    # 生成代码 (Go / TypeScript / OpenAPI)
│   ├── app/
│   │   ├── admin/service/          # 管理后台服务 (HTTP/gRPC)
│   │   ├── app/service/            # 前台应用服务 (HTTP/gRPC)
│   │   └── core/service/           # 核心业务服务 (业务逻辑 + 数据层)
│   ├── pkg/                        # 公共包 (鉴权/加密/事件总线/JWT/中间件/OSS...)
│   └── scripts/                    # 部署脚本 (Docker/PM2/环境安装)
├── frontend/
│   ├── admin/                      # 管理后台前端 (Vue3 + Ant Design Vue + Vben Admin)
│   └── app/                        # 前台应用
│       ├── react/                  # React 前台 (Next.js)
│       ├── vue/                    # Vue 前台 (Nuxt.js)
│       ├── taro/                   # Taro 前台 (微信小程序/H5)
│       └── flutter_app/            # Flutter 前台 (跨平台原生)
└── docs/                           # 文档 & 截图
```

## 快速开始

### 环境要求

- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ & pnpm
- buf (Protobuf 工具链)

### 1. 启动依赖服务

```bash
cd backend

# Windows
.\scripts\docker\libs_only.ps1

# Linux / macOS
./scripts/docker/libs_only.sh
```

### 2. 启动后端服务

```bash
# 推荐方式：使用 gow CLI
gow run admin

# 或在 IDE 中直接调试运行
```

### 3. 启动前端

```bash
# 管理后台
cd frontend/admin
pnpm install
pnpm dev

# React 前台
cd frontend/app/react
pnpm install
pnpm dev
```

### 常用命令

```bash
cd backend

# 生成 Protobuf API 代码
make api

# 生成 OpenAPI 文档
make openapi

# 生成 TypeScript 代码
make ts

# 一键生成全部代码 (ent + wire + api + openapi)
make gen

# 构建所有服务
make build

# 运行测试
make test
```

> 更多开发工作流请参考 [后端文档](./backend/README.md) 和 [脚本指南](./backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md)。

## 后台截图

<table>
    <tr>
        <td><img src="./docs/images/admin_login.png" alt="后台用户登录界面"/></td>
        <td><img src="./docs/images/admin_register.png" alt="后台用户注册界面"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_post_list.png" alt="后台帖子列表"/></td>
        <td><img src="./docs/images/admin_post_edit.png" alt="后台帖子编辑"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_category_list.png" alt="后台分类列表"/></td>
        <td><img src="./docs/images/admin_category_edit.png" alt="后台分类编辑"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_tag_list.png" alt="后台标签列表"/></td>
        <td><img src="./docs/images/admin_tag_edit.png" alt="后台标签编辑"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_comment_list.png" alt="后台评论列表"/></td>
        <td><img src="./docs/images/admin_site_list.png" alt="后台站点列表"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_site_setting_list.png" alt="后台站点配置列表"/></td>
        <td><img src="./docs/images/admin_navigation_list.png" alt="后台导航列表"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_analysis.png" alt="后台数据分析"/></td>
        <td><img src="./docs/images/admin_media_asset_list.png" alt="后台媒体资源列表"/></td>
    </tr>
</table>

## 前台截图

<table>
    <tr>
        <td><img src="./docs/images/react_app_login.png" alt="React 前台登录界面"/></td>
        <td><img src="./docs/images/react_app_register.png" alt="React 前台注册界面"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_homepage.png" alt="React 前台主页"/></td>
        <td><img src="./docs/images/react_app_about.png" alt="React 前台关于页面"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_post_list.png" alt="React 前台帖子列表"/></td>
        <td><img src="./docs/images/react_app_post_detail.png" alt="React 前台帖子详情"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_category_list.png" alt="React 前台分类列表"/></td>
        <td><img src="./docs/images/react_app_category_detail.png" alt="React 前台分类详情"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_tag_list.png" alt="React 前台标签列表"/></td>
        <td><img src="./docs/images/react_app_tag_detail.png" alt="React 前台标签详情"/></td>
    </tr>
</table>

## 联系我们

- 微信个人号：`yang_lin_bo`（备注：`go-wind-cms`）
- 掘金专栏：[go-wind-cms](https://juejin.cn/column/7541283508041826367)

## 参与贡献

我们欢迎各种形式的贡献，包括但不限于：

- 提交 [Issue](../../issues) 报告 Bug 或提出功能建议
- 提交 [Pull Request](../../pulls) 修复问题或添加新功能
- 完善文档和翻译
- 分享使用经验

## 开源许可

本项目基于 [MIT License](./LICENSE) 开源。

## 致谢

[![JetBrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://jb.gg/OpenSource)

感谢 [JetBrains](https://jb.gg/OpenSource) 提供的免费 GoLand & WebStorm 开源许可。
