# 贡献指南

感谢你对 GoWind Content Hub（go-wind-cms）的关注！本文档说明如何向本项目提交代码、文档与问题反馈。在开始之前，请先阅读 [README.md](./README.md) 了解项目定位与技术栈。

参与本项目即代表你同意在讨论、Issue 与 PR 中保持友善与专业。

## 快速上手

1. **Fork 并克隆仓库**

   ```bash
   git clone https://github.com/<your-name>/go-wind-cms.git
   cd go-wind-cms
   git remote add upstream https://github.com/<org>/go-wind-cms.git
   ```

2. **搭建开发环境**

   - 后端环境与启动：[backend/README.md](./backend/README.md)、[backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md](./backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md)
   - 依赖服务一键起：`backend/scripts/docker/libs_only.ps1`（Windows）或 `libs_only.sh`（Linux / macOS）
   - 首次安装工具链：`cd backend && make init`（protoc 插件 + buf / ent / golangci-lint / gow）

3. **创建分支**

   请基于最新的 `main` 创建特性分支，不要直接在 `main` 上开发：

   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feat/your-feature
   ```

## 开发约定

### 代码生成

本项目大量依赖代码生成，**请不要手工编辑生成产物**：

| 产物 | 生成命令 | 说明 |
|------|---------|------|
| Go API / HTTP / gRPC / 校验代码 | `gow api` 或 `make api` | 由 `backend/api/protos/` 下的 proto 生成至 `backend/api/gen/go/` |
| OpenAPI 文档 | `make openapi` | 产物写到 `backend/app/admin/service/cmd/server/assets/openapi.yaml`，**打进二进制**（用途见下方提示） |
| TypeScript HTTP 客户端 | `make ts` | 一条命令依次跑五个 buf 模板：admin、react、taro、vue 各自的 `src/api/generated/`，外加 admin OpenAPI |
| Ent ORM 代码 | `gow ent admin` 或 `make ent` | 由 `app/core/service/internal/data/ent/schema/` 下的 schema 生成 |
| 新模块接入装配 | `make register ENTITY=xxx [SVC=core] [PKG=contentV1]` | 把仓储/服务构造行与路由注册注入手写的 `cmd/server/wiring.go` 和 `internal/server/rest_server.go`（靠锚点行定位，幂等可重跑） |

`gow` 与 `make` 产出完全一致，选 `gow` 的理由很实际：Windows 没有原生 make，且本仓 Makefile 是嵌套结构、需要跨目录 cd（`gow` 会自动发现 `app/*/service`，见根 [AGENTS.md](./AGENTS.md)）。`make ts` / `make openapi` 目前只有 Makefile 入口。

修改 `.proto` 或 `ent/schema` 后请重新生成并**连同产物一并提交**。`make ts` 的多端产物是同一条命令生成的、本来就同时更新，不存在"只提交一端"的选项。

> `make register` 只覆盖 CRUD 生成器的标准形态。依赖更多的模块请手工补那两行装配，别指望工具替你写完 wiring。

> 漏跑 `make openapi` 的后果很费解：管理页「接口同步」读的是**打进二进制的** `openapi.yaml`，所以它会告诉你同步成功了，新端点却依旧不在接口管理表里，而路由本身看起来完全正常。正确顺序是 `make api && make openapi` 后重启进程，再触发同步。

### 输入校验

字段约束请在 proto 中用 protoc-gen-validate 的 `(validate.rules)` 声明（配套 `import "validate/validate.proto";`），避免在 service 层手写判空。`audit` 目录下已有可参照的用例。

生成器会产出对应的 `*.pb.validate.go`。**注意**：本仓 REST 链路没有统一的 validate 中间件，规则主要在 gRPC 校验中间件（配置 `enable_validate`）生效——给关键入参加完约束后，请实测确认真的会拦，别以为写了规则就自动生效。

### 代码规范

- **Go**：遵循 [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)。提交前在 `backend/` 下跑 `make lint`（golangci-lint）和 `go vet ./...`。
- **前端**：各应用有独立的 eslint 配置，请遵守对应配置；提交前跑所在应用的 lint / typecheck（门禁命令见根 [AGENTS.md](./AGENTS.md)）。
- **命名与风格**：新代码应与所在文件既有风格保持一致（缩进、注释密度、命名习惯）。

### 提交规范

本项目采用 [Conventional Commits](https://www.conventionalcommits.org/zh-hans/)，历史提交已统一使用。提交信息格式：

```
<type>(<scope>): <subject>
```

| type | 用途 |
|------|------|
| `feat` | 新功能 |
| `fix` | 缺陷修复 |
| `refactor` | 重构（不改变外部行为） |
| `perf` | 性能优化 |
| `style` | 代码格式（不影响逻辑） |
| `docs` | 文档 |
| `test` | 补充测试 |
| `build` | 构建系统、依赖 |
| `chore` | 杂项、脚手架 |
| `ci` | CI 配置 |

- `scope` 可选，对应受影响的模块（如 `admin`、`app`、`core`、`content`、`upload`）
- `subject` 使用中文或英文均可，简明描述，结尾不加句号
- **示例**：`feat(admin): 文章编辑页参照重做并增加列表状态 Tab`

破坏性变更在 type 后加 `!`（如 `feat(admin)!:`）并在正文说明迁移方式。大改动请添加正文说明动机与影响范围，必要时附 Issue 链接。

## 提交 Pull Request

1. 确保本地通过后端 `make lint` + `go vet ./...` + `make test`，以及改动所在前端的 lint / typecheck
2. 若改动涉及 proto / schema，确认已重新生成并提交产物：`make api && make openapi && make ent`（`make ts` 在前端在范围内时跑）
3. PR 标题遵循上述提交规范
4. PR 描述请写清：类型、动机、改动内容、**影响范围**（落在哪一端/哪个服务）、自测情况
5. 一个 PR 只解决一个问题，便于评审与回退

> Windows 上 `make test` 偶发 `fork/exec %TEMP%\...test.exe: Access is denied.`——那是安全策略/杀软拦了刚链接好的测试二进制，**不是你的代码挂了**。按根 [AGENTS.md](./AGENTS.md) 的办法把测试二进制输出到工作区内再直接跑一次确认，别据此提 Issue。

### 评审标准

- 是否改动生成产物、是否漏跑生成
- 是否破坏多租户隔离（跨租户数据访问需带租户过滤）
- 是否引入新的硬编码密钥 / 凭据（应走配置或环境变量）
- 是否补齐了相应测试
- 文档是否同步更新

## 反馈问题

- 缺陷请开 Issue，写清：环境（OS / Go / Node 版本）、复现步骤、期望与实际行为、相关日志
- 新功能建议请开 Issue 并说明使用场景与预期效果
- 安全漏洞请按 [SECURITY.md](./SECURITY.md) 流程私下上报，**勿直接开公开 Issue**

## 致谢

你的每一次贡献——无论代码、文档还是问题反馈——都让 GoWind Content Hub 更好。感谢你的参与。
