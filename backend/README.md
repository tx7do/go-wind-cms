# GoWind CMS Backend

一个Golang的Headless CMS。

- 后端基于 golang微服务框架 [go-kratos](https://go-kratos.dev/)
- 前端基于 [VUE3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/)

## 技术栈

* [Kratos](https://go-kratos.dev/) -- B站微服务框架
* [Consul](https://www.consul.io/) -- 服务发现和配置管理
* [OpenTelemetry](https://opentelemetry.io/) -- 分布式可观察系统
* [Wire](https://github.com/google/wire) -- 依赖注入框架
* [OpenAPI](https://www.openapis.org/) -- RESTful API 文档
* [Ent](https://entgo.io/) -- Facebook ORM 数据库实体框架
* [Redis](https://redis.io/) -- 非关系型数据库
* [PostgreSQL](https://www.postgresql.org/) -- 关系型数据库
* [MinIO](https://min.io/) -- 对象存储服务器

## API文档

### Swagger UI

- [Admin Swagger UI](http://localhost:8800/docs/)
- [Front Swagger UI](http://localhost:9800/docs/)

### openapi.yaml

- [Admin openapi.yaml](http://localhost:8800/docs/openapi.yaml)
- [Front openapi.yaml](http://localhost:9800/docs/openapi.yaml)

## 生成Protobuf API

使用[buf.build](https://buf.build/)进行Protobuf API的构建。

相关命令行工具和插件的具体安装方法请参见：[Kratos微服务框架API工程化指南](https://juejin.cn/post/7191095845096259641)

在`{项目根目录}/backend/api`下执行命令：

- 更新buf.lock

    ```bash
    buf dep update
    ```

- 生成GO代码

    ```bash
    buf generate
    ```

- 生成OpenAPI v3文档

    ```bash
    buf generate --template buf.admin.openapi.gen.yaml
    buf generate --template buf.front.openapi.gen.yaml
    ```

## Make构建

请在`app/{服务名}/service`下执行：

- 初始化开发环境

   ```bash
   make init
   ```

- 生成API的go代码

   ```bash
   make api
   ```

- 生成API的OpenAPI v3 文档

   ```bash
   make openapi
   ```

- 生成ent代码

   ```bash
   make ent
   ```

- 生成wire代码

   ```bash
   make wire
   ```

- 构建程序

   ```bash
   make build
   ```

- 调试运行

   ```bash
   make run
   ```

- 构建Docker镜像

   ```bash
   make docker
   ```

## Bazel构建

使用[bazel.build](https://bazel.build/)进行服务器程序的构建。

如何安装bazel.build的文档，请参考官方文档：<https://bazel.build/install>。

在`blog-backend`根目录下执行命令：

- 更新GO依赖库引入的构建配置文件repos.bzl

   ```bash
   bazel run //:gazelle-update-repos
   ```

- 拉取依赖项，生成配置文件BUILD.bazel

   ```bash
   bazel run //:gazelle
   ```

- 构建单个服务

  ```bash
  bazel build //app/admin/service/cmd/server:server
  bazel build //app/comment/service/cmd/server:server
  bazel build //app/content/service/cmd/server:server
  bazel build //app/file/service/cmd/server:server
  bazel build //app/user/service/cmd/server:server
  ```

  或者

  ```bash
  bazel build //:admin-service
  bazel build //:comment-service
  bazel build //:content-service
  bazel build //:file-service
  bazel build //:user-service
  ```

- 运行单个服务

  ```bash
  bazel run //app/admin/service/cmd/server:server
  bazel run //app/comment/service/cmd/server:server
  bazel run //app/content/service/cmd/server:server
  bazel run //app/file/service/cmd/server:server
  bazel run //app/user/service/cmd/server:server
  ```

  或者

  ```bash
  bazel run //:admin-service
  bazel run //:comment-service
  bazel run //:content-service
  bazel run //:file-service
  bazel run //:user-service
  ```

- 单个服务生成Docker镜像tar文件

  ```bash
  bazel build //:admin-service-image
  bazel build //:comment-service-image
  bazel build //:content-service-image
  bazel build //:file-service-image
  bazel build //:user-service-image
  ```

- 单个服务生成本地Docker镜像

  ```bash
  bazel run //:admin-service-image
  bazel run //:comment-service-image
  bazel run //:content-service-image
  bazel run //:file-service-image
  bazel run //:user-service-image
  ```

- 单个服务推送Docker镜像到DockerHub

  ```bash
  bazel run //:admin-service-image-push
  bazel run //:comment-service-image-push
  bazel run //:content-service-image-push
  bazel run //:file-service-image-push
  bazel run //:user-service-image-push
  ```

- 构建全部服务

  ```bash
  bazel build //...
  ```

## Docker部署
