# 备份恢复演练（Backup Restore Drill）

> 配套实现：备份导出（`TaskService.AsyncBackup` + `BackupRepo.ExportCoreTables`，
> core 服务 asynq `backup` 任务）与恢复（`BackupRepo.RestoreCoreTables`，
> 仅允许空库、保留原 ID、单事务）。本文档面向运维/交付：如何**安全地**验证
> "备份确实可用"，以及在灾难场景下如何执行恢复。

## 1. 能覆盖与不覆盖

| 覆盖 | 不覆盖 |
|---|---|
| 24 张核心身份/权限/组织配置表（见 `BackupRepo.coreTableExports`：tenants、users、user_credentials、roles、role_*、permission_*、user_*、membership*、org_units、positions、menus、apis、login_policies） | 审计日志等大体量表（备份本就不含） |
| 保留原 ID 与时间字段；`sys_user_credentials` 密码哈希随备份保真（恢复后可登录） | 内容业务表（posts/pages/media 等——内容数据请走数据库原生备份通道） |
| JSON 备份（`backups` 桶 `gowind-cms/YYYY/MM/DD/<name>-<time>.json.gz`） | 增量合并恢复（只支持空库整体恢复） |

> 等保视角：本工具链解决"核心配置表误删/迁移"场景；内容与审计数据的留存
> 依赖数据库层备份（pg_dump 等），两者互补，不互相替代。

## 2. 备份对象

- 桶：`backups`（MinIO，`UploadFile` 会自动建桶）
- 对象键：`gowind-cms/YYYY/MM/DD/<name>-<time>.json.gz` —— **带项目前缀**。
  dev/小部署的 MinIO 常多项目共用一个桶，清理只扫本前缀，不碰他项目对象。
- 触发：`sys_tasks` 表的 `type_name='backup'` 任务（demo 库自带每小时一行）；
  也可在任务管理页调整 cron 或手动控制（ControlTask）。
- 保留期：环境变量 `BACKUP_RETENTION_DAYS`（默认 30 天；**≤0 显式关闭清理**）。
  清理在每次备份上传成功后顺手执行（best-effort，失败只留日志）。

## 3. 恢复演练（标准流程）

**原则：演练永远不碰活库。** 用一次性全新库做导入验证。
恢复**没有**管理页/RPC 入口（有意设计：避免误触；恢复是运维动作，走一次性程序）。

### 3.1 前置

- 取到一份备份产物：MinIO `backups` 桶对象（`mc cp local/backups/gowind-cms/...`）；
- 一个**空的目标数据库**（ent schema 迁移完成后、未播种的库）。

### 3.2 步骤（一次性 Go 程序，参照 `backup_restore_test.go` 的往返测试）

```go
// 1) 读备份文件（.json.gz 先 gzip 解压）
raw, _ := os.ReadFile("backup.json.gz")
gz, _ := gzip.NewReader(bytes.NewReader(raw))
doc, _ := io.ReadAll(gz)

// 2) 解出 data 段
var doc0 map[string]json.RawMessage
_ = json.Unmarshal(doc, &doc0)

// 3) 连接空库（迁移 schema 后、未播种），导出/恢复都要用系统旁路上下文
ctx := appViewer.NewSystemViewerContext(context.Background())
repo := data.NewBackupRepo(bootstrapCtx, entCrudClient)
counts, err := repo.RestoreCoreTables(ctx, doc0["data"])
// 成功后 counts 即各表恢复行数；与备份 _meta.tables 对照
```

### 3.3 恢复后校验清单

1. `counts` 与备份 `_meta.tables` 一致；
2. 管理员账号能登录（密码哈希随 `sys_user_credentials` 保真）；
3. 抽查租户/角色/菜单绑定关系完整（`sys_user_roles`/`sys_role_permissions`）；
4. 平台超管权限码（`sys:platform_admin`）仍能拉到全部菜单（ID 保真的意义所在）；
5. PostgreSQL：确认各表 ID 序列已对齐（工具自动 setval，日志有 `sequence ... aligned`）。

### 3.4 频率建议

- 每季度一次演练；每次等保测评前必演；
- 演练记录留存（谁、何时、哪个备份、校验结论）——等保测评要查。

## 4. 灾难恢复（真实场景）

1. 定位故障（误删/迁移）→ 准备全新库（同版本 schema 迁移、未播种）；
2. 选最近一次成功备份（core 日志 `backup: completed successfully, object=...`）；
3. 按第 3 节流程恢复到新库 → 校验 → 把应用配置指到新库 → 启动；
4. 事后把误操作期间的业务数据人工补录（恢复点是备份时刻，之后的增量不在覆盖范围）。

## 5. 明确不做

- **不做活库合并恢复**：RestoreCoreTables 前置校验任一目标表非空即整体中止
  （这是安全设计：活库合并涉及 ID 冲突/外键顺序/序列回填，脚本化极易产出脏数据）；
- 不提供管理页/RPC 恢复入口（避免误触；恢复是运维动作，走一次性程序）；
- 不自动恢复（自动=把误删快速放大，人工确认是必要闸门）。
