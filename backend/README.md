# GoWind Content Hub Backend

GoWind（风行内容中台）— 企业级无头内容中枢（Headless Content Hub），基于 Go + Vue3 + TypeScript 构建，采用 API
优先、前后端分离架构，提供多语言原生支持、多端内容分发、精细化权限管控，为出海企业、多终端项目提供高可用、可扩展的解决方案。

- 后端基于 [golang](https://go.dev/) + [微服务框架 go-kratos](https://go-kratos.dev/)
- 前端基于 [Vue.js 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/)

## 技术栈

- [Kratos](https://go-kratos.dev/) -- B站微服务框架
- [Consul](https://www.consul.io/) / [Etcd](https://etcd.io/) -- 服务发现和配置管理
- [OpenTelemetry](https://opentelemetry.io/) -- 分布式可观察系统
- 手写依赖装配（wiring.go，无 DI 框架）
- [OpenAPI](https://www.openapis.org/) -- RESTful API 文档
- [MinIO](https://min.io/) -- 对象存储服务器
- [Redis](https://redis.io/) -- 非关系型数据库
- [PostgreSQL](https://www.postgresql.org/) / [MySQL](https://www.mysql.com/) -- 关系型数据库
- [Ent](https://entgo.io/) / [GORM](https://gorm.io/index.html) -- Golang ORM框架
- [ElasticSearch](https://www.elastic.co/cn/elasticsearch) -- 搜索引擎

## API文档

### Swagger UI

- [Admin Swagger UI](http://localhost:6600/docs/)
- [App Swagger UI](http://localhost:6700/docs/)

### openapi.yaml

- [Admin openapi.yaml](http://localhost:6600/docs/openapi.yaml)
- [Front openapi.yaml](http://localhost:6700/docs/openapi.yaml)

## 生成Protobuf API

本项目使用[buf.build](https://buf.build/)进行Protobuf API构建。

相关命令行工具和插件的具体安装方法请参见：[Kratos微服务框架API工程化指南](https://juejin.cn/post/7191095845096259641)

本项目提供了两种生成API代码的方式：`Makefile`和`gow cli`。

一键生成API的所有代码，该命令可以在后端项目的任何位置执行：

```bash
gow api
```

### 生成GO代码

```bash
cd `{项目根目录}/backend`
make api
```

### 生成TypeScript代码

```bash
cd `{项目根目录}/backend`
make ts
```

### 生成OpenAPI v3文档

```bash
cd `{项目根目录}/backend`
make openapi
```

## 其他代码生成

### 生成ent代码

```bash
cd `{项目根目录}/backend/app/{服务名}/service`
make ent
```

### 构建Docker镜像

```bash
make docker
```

## 构建程序

```bash
cd `{项目根目录}/backend/app/{服务名}/service`
make build
```

### 调试运行

```bash
cd `{项目根目录}/backend/app/{服务名}/service`
make run
```

## Docker部署
