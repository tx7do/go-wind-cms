# 更新日志

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/) 规范，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

> 以下内容归入 [1.0.0]（首个正式版本）。
> 条目基于 Git 提交记录回溯整理（自 2026-08 起），仅记录对使用者有影响的功能性变更，
> 日常重构、代码风格、CI 调整等不在逐一列举；更早历史请查阅 Git 提交记录。

## [1.0.0] - 2026-10-08

### 新增

- AI 模块（自 go-wind-admin 上游移植并适配两服务拓扑）：模型提供商管理、
  SSE 流式对话（Function Call 工具循环）、知识库 RAG（pgvector 切片向量化
  检索，支持 txt/md/docx/pdf 上传入库）、智能问数（NL→只读 SQL→结构化结果，
  表白名单与租户谓词 fail-closed）、用量记账与统计页、审计/用量服务端全量
  导出（XLSX/CSV）。
- 全局搜索语义化：菜单与帖子标题/摘要向量化索引（pgvector），搜索直达帖子
  编辑页；平台用户全量、租户用户按租户过滤，菜单项按角色可见集裁剪；
  索引空时懒自建，支持手动全量重建。
- Dashboard 安全洞察：规则预筛近 24h 审计明细行为模式（疑似口令尝试/深夜
  操作/失败集中/敏感操作），结构化事实 + LLM 总评按界面语言生成；平台用户专属。
- AI 内容生成接入 CMS 表单：帖子标题/摘要/正文、页面标题、分类/标签描述、
  评论回复（携带评论原文上下文）、角色描述、站内信正文（公告场景）。
- 对话助手挂载 CMS 域工具：帖子列表查询、状态分布统计、内容概览、评论统计、
  互动 TOP 帖——模型在对话中自主调用只读工具后作答（工具调用过程可见）。
- 异步任务：知识库切片重索引（更换 embedding 模型后全量重算）、定时问数
  （执行问数并以站内信投递结果给平台用户）。
- 管理后台前端目录由 `frontend/admin-react` 更名为 `frontend/admin`。

- 内容建模子系统：可视化自定义内容模型与字段定义（文本/数字/富文本/图片/文件/关联等类型），
  后端写入侧与管理侧全链路落地；Section 随后从顶层独立模块降级为 Page 嵌套子部件。
- 管理后台分析页接入真实统计：新建 stats 域统计 RPC；全局补时间默认值钩子
  （ApplyTimeDefaultHooks，杜绝 created_at 全 NULL 一类问题）。
- 翻译子系统增强：翻译引擎改环境变量配置化并新增百度引擎；一键翻译增加 Markdown 保护
  （代码块与链接不被译坏）。
- 文章编辑页重做：右栏设置面板、草稿落库、定时发布、未保存离开拦截；编辑器类型收窄为
  富文本 / Markdown；帖子列表新增状态 Tab（草稿箱，带计数）与按分类过滤。
- 登录支持邮箱身份凭证。
- API 错误信息五端统一改用 reason 查各端 i18n 字典（回退 DEFAULT），不再直接展示后端 message。
- C 端品牌蓝设计语言统一（react / taro / vue / flutter）并沉淀设计规范文档
  （`docs/design-language.md`）；GoWind 风行品牌标识与三语 README 品牌头部同步落地。
- thumbnail 字段从翻译表上移主表（post / page / category 全链路；存量库需手工搬数据）。

### 变更

- **管理后台切换至 React 版**（React 19 + antd 6 + ProComponents + TanStack Query），
  移除 Vben 前端（破坏性变更；dev 端口 15999）。
- 移除 google/wire，依赖装配改为手写 `wiring.go`，新增 `make register` 辅助登记 CRUD 模块。
- 四端 dev 端口统一升至 15xxx 段（admin 15999 / react 15001 / vue 15011 / taro 15021），
  避开姐妹仓默认段。
- 上传链路加固：流式上传消除整文件 OOM（minio-go 自动分片），新增上传对象大小硬限制；
  清理文件传输 / 头像上传 / OssService 死代码与预签名 URL 调试残留。
- demo 数据 SQL 对齐 ent schema 并统一 pg / mysql 两版。
- engine/opa 升至 v1.1.16（SetPolicies 竞态修复版）。
- 管理后台侧边栏菜单文案统一为四字；编辑器暗色主题对齐应用色板。

### 修复

- 全量 GUI 审计发现的 15 项增删改查缺陷（租户创建被静默阻断、站点更新 mask 500、
  字典增删全断、序列脱节、slug 覆盖等），及遗留审计问题（任务控制 / 操作审计日志接线 / 评论树形）。
- 翻译 Update 改按语言 upsert，消除重复行累积与跨语言误删；GetTranslation 防重复 500。
- C 端应用全链路修复（注册 / 游客评论 / 搜索 / 令牌刷新 / 角色细分）；文章作者名称读时回填；
  文章列表统一最新在前。
- 管理后台登录后地址栏与路由渲染"劈叉"（store 整页跳转与 router navigate 竞态）；
  暗色模式下列表 tag 颜色错误；编辑器高度未撑满编辑页剩余空间。
- 管理后台 drawer 表单补齐缺失输入项（页面别名与状态/类型选择、站点配置 value/options/isRequired），
  消除 Update 路径静默丢数据。
- 关闭 TanStack Query 缓存，列表与单查每次走实时请求，消除编辑回显脏数据。
- demo 数据 SQL 语法碎片（多余逗号行致语句断裂）。

[Unreleased]: https://github.com/tx7do/go-wind-cms/commits/main
