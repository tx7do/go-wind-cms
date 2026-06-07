<div align="center">

# GoWind Content Hub

### FengXing · Out-of-the-Box Enterprise Integrated Frontend-Backend Content Platform

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![React](https://img.shields.io/badge/React-19.x-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![Flutter](https://img.shields.io/badge/Flutter-3.x-02569B?logo=flutter&logoColor=white)](https://flutter.dev/)
[![Kratos](https://img.shields.io/badge/Kratos-2.9-00ADD8?logo=go&logoColor=white)](https://go-kratos.dev/)

**English** · **[中文](./README.md)** · **[日本語](./README.ja-JP.md)**

</div>

---

FengXing (GoWind HCH) is an out-of-the-box enterprise-grade Golang full-stack Headless Content Platform (HCH = Headless Content Hub), providing enterprises with flexible and scalable full-domain content management and distribution solutions.

**Key Highlights:**

- **API First** — Complete RESTful / gRPC dual-protocol interfaces with auto-generated OpenAPI docs
- **Multi-platform** — Frontend support for Vue, React, Taro (Mini Program), and Flutter
- **Multi-tenant** — Tenant data isolation with auto-initialized departments, roles, and admins
- **Fine-grained Permissions** — Three-level access control: menu, API, and data permissions with button-level granularity
- **Native i18n** — Unified translation management for content, menus, and UI text
- **Microservices** — Built on go-kratos with service discovery, distributed tracing, and distributed transactions

## Demo

| Demo Type | Access URL |
|:---|:---|
| Admin Panel | [https://admin.cms.gowind.cloud](https://admin.cms.gowind.cloud) |
| Admin API Swagger | [https://api.admin.cms.gowind.cloud/docs/](https://api.admin.cms.gowind.cloud/docs/) |
| App API Swagger | [https://api.cms.gowind.cloud/docs/](https://api.cms.gowind.cloud/docs/) |
| Vue Frontend | [https://cms.gowind.cloud](https://cms.gowind.cloud) |
| React Frontend | [https://react.cms.gowind.cloud](https://react.cms.gowind.cloud) |
| Taro Frontend | [https://taro.cms.gowind.cloud](https://taro.cms.gowind.cloud) |
| Flutter Frontend | [https://flutter.cms.gowind.cloud](https://flutter.cms.gowind.cloud) |

> Default credentials for all demo sites: `admin` / `admin`

## Tech Stack

### Backend

| Layer | Technology | Description |
|:---|:---|:---|
| Language | [Go 1.25+](https://go.dev/) | High-performance compiled language |
| Framework | [go-kratos](https://go-kratos.dev/) | Microservice framework by Bilibili |
| DI | [Wire](https://github.com/google/wire) | Compile-time dependency injection |
| ORM | [Ent](https://entgo.io/) | Go entity framework |
| Database | [PostgreSQL](https://www.postgresql.org/) / [MySQL](https://www.mysql.com/) | Relational database |
| Cache | [Redis](https://redis.io/) | In-memory data store |
| Object Storage | [MinIO](https://min.io/) | S3-compatible object storage |
| Search Engine | [OpenSearch](https://opensearch.org/) | Full-text search |
| Service Registry | [Etcd](https://etcd.io/) / [Consul](https://www.consul.io/) | Service discovery & configuration |
| Tracing | [Jaeger](https://www.jaegertracing.io/) + [OpenTelemetry](https://opentelemetry.io/) | Distributed observability |
| API Definition | [Protobuf](https://protobuf.dev/) + [buf.build](https://buf.build/) | Contract-first API design |
| AuthZ Engine | [Casbin](https://casbin.org/) / [OPA](https://www.openpolicyagent.org/) | Policy-driven authorization |

### Admin Frontend

| Technology | Description |
|:---|:---|
| [Vue 3](https://vuejs.org/) | Progressive frontend framework |
| [TypeScript](https://www.typescriptlang.org/) | Type-safe development |
| [Ant Design Vue](https://antdv.com/) | Enterprise UI component library |
| [Vben Admin](https://doc.vben.pro/) | Admin management framework |

### App Frontend

| Version | Tech Stack | Use Case |
|:---|:---|:---|
| Vue | [Nuxt.js](https://nuxt.com/) + [shadcn-vue](https://www.shadcn-vue.com/) | Web App / SSR |
| React | [Next.js](https://nextjs.org/) + [shadcn/ui](https://ui.shadcn.com/) | Web App / SSR |
| Taro | [Taro](https://docs.taro.zone/en/docs/) + React + [shadcn/ui](https://ui.shadcn.com/) | WeChat Mini Program / H5 |
| Flutter | [Flutter](https://flutter.dev/) + [BLoC](https://bloclibrary.dev/) | Cross-platform native app |

## Core Features

### Organization & Permissions

| Feature | Description |
|:---|:---|
| Multi-tenant Management | Tenant creation, enable/disable, package config & data isolation; auto-init departments, roles & admin |
| User Management | Full lifecycle management, multi-role/dept binding, proxy login, advanced queries |
| Role Management | Role & role group management with fine-grained menu, API & data permissions |
| Permission Management | Permission groups, menu nodes & permission points with button-level control |
| Menu Management | Visual menu config with directory/page/button levels, dynamic rendering by permission |
| Department Management | Multi-level department tree management with user binding |

### Content & Site

| Feature | Description |
|:---|:---|
| Content Modeling | Visual custom content models with text, number, rich text, image, file, relation fields |
| Content Management | CRUD, publish/unpublish, pin, sort, recycle bin & batch operations; rich text / Markdown |
| Category Management | Multi-level category tree, content model binding, frontend filtering |
| Tag Management | Tag CRUD with content association, search & aggregation |
| Comment Management | Comment audit, delete, reply, block; filter by content, user, time |
| Multi-language Management | Language management, unified translation for content, menus & UI text |
| Site Management | Multi-site config with independent domain, title, logo, SEO & style |
| Site Configuration | Basic info, SEO, upload limits, cache policy, email/SMS global params |

### System & Operations

| Feature | Description |
|:---|:---|
| File Management | Image/document/video management, local or OSS storage, preview, download, grouping |
| Dictionary Management | Data dictionary categories & items, linked queries, sorting, import/export |
| API Management | Backend API maintenance & auto-sync, tree view, request/response config |
| Task Scheduling | Cron job management, start/pause/execute, execution history & logs |
| Message Notification | Multi-level message categories, targeted messaging, read status tracking |
| Internal Messages | Personal inbox, view/delete, single/batch mark as read |
| Cache Management | Real-time cache query, precise key-based or batch cleanup |
| Login Logs | Login success/failure logs with account, IP, location, device, time |
| Operation Logs | Full-link audit trail with operator, IP, params & results |
| Personal Center | Profile editing, avatar, password reset, login history |
| Headless API | Complete OpenAPI for content CRUD, multi-terminal support |

## Project Structure

```
go-wind-cms/
├── backend/                        # Backend services
│   ├── api/                        # Protobuf API definitions & generated code
│   │   ├── protos/                 # .proto source files
│   │   └── gen/                    # Generated code (Go / TypeScript / OpenAPI)
│   ├── app/
│   │   ├── admin/service/          # Admin service (HTTP/gRPC)
│   │   ├── app/service/            # App service (HTTP/gRPC)
│   │   └── core/service/           # Core service (business logic + data layer)
│   ├── pkg/                        # Shared packages (auth/crypto/eventbus/JWT/middleware/OSS...)
│   └── scripts/                    # Deploy scripts (Docker/PM2/env setup)
├── frontend/
│   ├── admin/                      # Admin frontend (Vue3 + Ant Design Vue + Vben Admin)
│   └── app/                        # App frontends
│       ├── react/                  # React app (Next.js)
│       ├── vue/                    # Vue app (Nuxt.js)
│       ├── taro/                   # Taro app (WeChat Mini Program / H5)
│       └── flutter_app/            # Flutter app (cross-platform native)
└── docs/                           # Documentation & screenshots
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ & pnpm
- buf (Protobuf toolchain)

### 1. Start Dependencies

```bash
cd backend

# Windows
.\scripts\docker\libs_only.ps1

# Linux / macOS
./scripts/docker/libs_only.sh
```

### 2. Start Backend

```bash
# Recommended: use gow CLI
gow run admin

# Or debug directly in your IDE
```

### 3. Start Frontend

```bash
# Admin panel
cd frontend/admin
pnpm install
pnpm dev

# React app
cd frontend/app/react
pnpm install
pnpm dev
```

### Common Commands

```bash
cd backend

# Generate Protobuf API code
make api

# Generate OpenAPI docs
make openapi

# Generate TypeScript code
make ts

# Generate all code (ent + wire + api + openapi)
make gen

# Build all services
make build

# Run tests
make test
```

> For more development workflows, see [Backend Docs](./backend/README.md) and [Scripts Guide](./backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md).

## Backend Screenshots

<table>
    <tr>
        <td><img src="./docs/images/admin_login.png" alt="Backend login"/></td>
        <td><img src="./docs/images/admin_register.png" alt="Backend registration"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_post_list.png" alt="Backend post list"/></td>
        <td><img src="./docs/images/admin_post_edit.png" alt="Backend post editing"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_category_list.png" alt="Backend category list"/></td>
        <td><img src="./docs/images/admin_category_edit.png" alt="Backend category editing"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_tag_list.png" alt="Backend tag list"/></td>
        <td><img src="./docs/images/admin_tag_edit.png" alt="Backend tag editing"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_comment_list.png" alt="Backend comment list"/></td>
        <td><img src="./docs/images/admin_site_list.png" alt="Backend site list"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_site_setting_list.png" alt="Backend site settings"/></td>
        <td><img src="./docs/images/admin_navigation_list.png" alt="Backend navigation list"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_analysis.png" alt="Backend data analysis"/></td>
        <td><img src="./docs/images/admin_media_asset_list.png" alt="Backend media assets"/></td>
    </tr>
</table>

## Frontend Screenshots

<table>
    <tr>
        <td><img src="./docs/images/react_app_login.png" alt="React frontend login"/></td>
        <td><img src="./docs/images/react_app_register.png" alt="React frontend registration"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_homepage.png" alt="React frontend homepage"/></td>
        <td><img src="./docs/images/react_app_about.png" alt="React frontend about"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_post_list.png" alt="React frontend post list"/></td>
        <td><img src="./docs/images/react_app_post_detail.png" alt="React frontend post detail"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_category_list.png" alt="React frontend category list"/></td>
        <td><img src="./docs/images/react_app_category_detail.png" alt="React frontend category detail"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_tag_list.png" alt="React frontend tag list"/></td>
        <td><img src="./docs/images/react_app_tag_detail.png" alt="React frontend tag detail"/></td>
    </tr>
</table>

## Contact

- WeChat: `yang_lin_bo` (Note: `go-wind-cms`)
- Juejin Column: [go-wind-cms](https://juejin.cn/column/7541283508041826367)

## Contributing

We welcome all forms of contribution, including but not limited to:

- Submit [Issues](../../issues) to report bugs or suggest features
- Submit [Pull Requests](../../pulls) to fix issues or add features
- Improve documentation and translations
- Share your experience

## License

This project is licensed under the [MIT License](./LICENSE).

## Acknowledgements

[![JetBrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://jb.gg/OpenSource)

Thanks to [JetBrains](https://jb.gg/OpenSource) for providing free GoLand & WebStorm open source licenses.
