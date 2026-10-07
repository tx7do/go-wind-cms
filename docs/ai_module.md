# AI 模块（移植自 go-wind-admin 上游）

本仓 AI 模块自上游 `go-wind-admin` 移植并适配两服务拓扑（admin BFF :6600 + core）。
功能：模型提供商、流式对话（Function Call）、知识库 RAG、智能问数、用量记账/统计、
全局搜索语义化、Dashboard 安全洞察、审计/用量服务端导出、CMS 表单 AI 生成。

## 快速开始（本地开发）

1. **启动依赖**：Postgres（需可装 pgvector 扩展，`CREATE EXTENSION vector`）+ Redis + asynq（Redis）。
2. **启动 mock 模型端点**（无真实 API key 时）：

   ```bash
   node frontend/admin/mock/ai-server.cjs   # 监听 127.0.0.1:18080，OpenAI 兼容
   ```

3. **提供商配置**：管理后台「AI 助手 → 模型提供商」新建：
   - 本地模型：host=127.0.0.1，port=18080，模型名任意（如 mock-model）
   - 云端模型：OpenAI 兼容 base_url（如 `https://api.deepseek.com/v1`）+ API Key
   - 勾选「默认提供商」；对话/生成/问数/embedding 均走默认项
4. **知识库**：「AI 助手 → 知识库」建库（绑定 provider + embedding 模型名）→ 上传文档
   （粘贴文本或 txt/md/docx/pdf 文件）→ 对话页选择知识库提问。
5. **重建索引**：更换 embedding 模型/provider 后，知识库页「重建向量索引」按钮
   （入队 asynq 任务 `ai.doc.reindex`，异步执行）。

## 安全模型

| 能力 | 边界 |
|------|------|
| 智能问数 | 只读 SELECT/WITH + 表白名单（本仓 CMS 真实 schema）+ LIMIT≤100 + 只读事务；租户用户强制 `tenant_id = 自身` 谓词（fail-closed）；路由限 `sys:platform_admin` |
| 安全洞察 | 平台用户专属（审计明细为全平台数据且外发模型端点），路由级与卡片级双重隐藏 |
| 知识库 | 检索 SQL 按租户过滤；embedding 按知识库归属租户记账 |
| api_key | AES-GCM 加密落库（全局加密器未启用时明文回落，与上游同语义），读取恒为空只回脱敏 hint |
| 用量 | 全调用记账（对话/生成/embedding），usage-logs 支持服务端全量导出 |

## 两服务拓扑适配要点（与上游单体的差异）

- core 侧 Chat 为 **server-streaming**（chunk/tool/done 帧），admin BFF 消费后经
  SSE 网关（:6601）推 `ai_chat_chunk`/`ai_chat_tool` 事件；done 帧作同步响应。
- 流式 RPC 不走 unary 中间件：操作者身份由 BFF 显式写入请求体
  （`operator_user_id`/`operator_tenant_id`），core 信任内部跳。
- 审计/用量导出在 BFF 手动路由（二进制响应），逐页复用各域 List RPC 聚合。

## 已知边界

- 通知域（渠道/规则/模板/台账/偏好）未移植，`ai.query.run` 结果投递仅站内信。
- 问数定时任务（`ai.query.run`）与重索引（`ai.doc.reindex`）已注册 asynq worker，
  触发入口待接任务管理 UI（当前可经 asynq/任务工具直接入队）。
- 语义索引首次查询懒自建；菜单结构变更后需调
  `POST /admin/v1/ai/content/rebuild-index` 重建。
