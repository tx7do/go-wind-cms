# AGENTS.md — 脚手架项目二次开发指南

> 本文件是前端 admin 子项目的 AI 编码规范单一事实源，适用于所有支持 AGENTS.md 的 AI 编码工具（ZCode、GitHub Copilot、Cursor、Codex、Gemini CLI 等）。Claude Code 通过 `CLAUDE.md` 中的 `@AGENTS.md` 引用加载。

## 项目概览

基于 Vue 3 + Vite + TypeScript 的中后台管理系统脚手架（Vben Admin 配置型框架），面向二开场景。

**核心技术栈**: Vue 3.5, Ant Design Vue 4.2, Tailwind CSS, Shadcn-ui, Pinia, Vue Router, Vue Query, VxeTable, i18n, Axios

**应用入口**: `apps/admin/src/`

## Vben 框架核心机制

本项目基于 Vben Admin **配置型框架**，通过配置而非编码来完成大部分开发工作。

### 组件注册机制（不要绕过）

框架有 **两套组件注册体系**：

**1. VbenForm 组件（schema 中使用）** — 定义在 `adapter/component/index.ts`，通过 `globalShareState.setComponents()` 注册。`schema` 的 `component` 字段 **只能使用这些注册过的名称**：

```
Input, InputNumber, InputPassword, Select, ApiSelect, TreeSelect,
ApiTreeSelect, RadioGroup, Checkbox, CheckboxGroup, Switch, DatePicker,
RangePicker, TimePicker, Textarea, Upload, Editor, IconPicker, AutoComplete,
Mentions, Rate, Divider, Space, DefaultButton, PrimaryButton, ApiTree
```

**2. Template 全局组件（template 中使用）** — 定义在 `registerGlobComp.ts`，通过 `app.use()` 全局注册。在 template 中 **直接用 `a-*` 前缀**：

```
<a-button>, <a-tag>, <a-popconfirm>, <a-input>, <a-select>,
<a-tree>, <a-table>, <a-dropdown>, <a-menu>, <a-card>, <a-space>,
<a-switch>, <a-tabs>, <a-divider>, <a-layout>
```

> **禁止**: `import { Tag, Button } from 'ant-design-vue'` 然后在 template 中用 `<Tag>` 或 `<Button>`。

### 配置驱动模式（不要自己写逻辑）

| 场景 | 配置方式 | 不要做 |
|---|---|---|
| 表格列日期格式化 | `formatter: 'formatDateTime'` | Slot 中用 dayjs 格式化 |
| 表单必填校验 | `rules: 'required'` / `'selectRequired'` | 用 Zod 或自定义校验函数 |
| 列表数据加载 | `proxyConfig.ajax.query` | 手动 watch + ref + async function |
| 分页 | `pagerConfig: {}` | 手动管理分页状态 |
| 表格刷新 | `gridApi.reload()` | 手动重新请求数据 |
| 表单赋值/取值 | `baseFormApi.setValues()` / `baseFormApi.getValues()` | 直接操作 DOM 或 ref |
| 表单校验 | `baseFormApi.validate()` | 手动检查每个字段 |

## 目录结构

```
apps/admin/src/
├── api/                  # API 层（两层架构）
│   ├── generated/        # ← protobuf 自动生成，禁止手动编辑
│   ├── client.ts         # ← ApiClient 单例（ClientTransport 适配器）
│   └── composables/      # ← Vue Query hooks 层：use*/fetch*/枚举工具
├── adapter/              # VbenForm + VxeTable 适配器配置
├── router/routes/modules/# ← 路由模块（按功能拆分）
├── stores/               # Pinia 状态管理
├── views/app/            # 业务页面（按功能模块组织）
├── locales/langs/        # i18n 国际化文件（zh-CN/en-US: enum.json, menu.json, page.json, ui.json）
└── transport/rest/       # HTTP 传输层（PaginationQuery, requestApi）
```

## 关键约定（必须遵守）

### 数据层约定

1. **禁止直接引用 `#/api/generated/` 路径** — 业务层通过 `#/api` 统一入口导入
2. **composables 直接使用 `apiClient`** — 导入 `apiClient` from `#/api/client`，调用 `apiClient.xxxService.Method()`
3. **组件内用 `use*` hooks，组件外（Store/路由守卫）用 `fetch*` 函数**
4. **更新操作只传变化字段** — `useUpdate*` 内部自动生成 `updateMask`
5. **Pinia Store 中不可依赖 `useRouter()`** — Store 初始化时路由可能未就绪
6. **所有列表查询统一使用 `PaginationQuery`**

### Vben 框架强规约

7. **表单组件必须使用注册名** — `schema` 中 `component` 只能用 `adapter/component/index.ts` 中注册的名称（如 `Input`、`Select`、`ApiSelect`），不要用 `AInput`、`ASelect` 或原生 HTML 标签
8. **Template 中使用 `a-*` 前缀** — Ant Design Vue 组件已全局注册，直接用 `<a-button>`、`<a-tag>`、`<a-popconfirm>` 等
9. **图标在 `:icon` prop 中必须用 `h()` 渲染** — `:icon="h(LucideFilePenLine)"`，图标从 `@vben/icons` 导入
10. **日期列用 `formatter: 'formatDateTime'`** — 不要在 Slot 中用 dayjs 手动格式化
11. **表单校验用内置规则名** — `rules: 'required'` 或 `rules: 'selectRequired'`，不要用 Zod 表达式
12. **删除操作必须二次确认** — 使用 `<a-popconfirm>`，不要用 `window.confirm` 或直接删除
13. **页面必须用 `<Page auto-content-height>` 包裹** — 不要用 `<div>` 替代
14. **消息提示用 `notification`** — 从 `ant-design-vue` 导入，不要用 `alert()` 或 `ElMessage`

### 通用约定

15. **国际化文本不硬编码** — 使用 `$t()` 引用 locales 文件中的 key
16. **使用严格相等运算符 `===`** — 禁止使用 `==`

## 导入路径约定

```typescript
// API 导入 — 统一通过 #/api 入口
import { useListUsers, fetchListUsers, PaginationQuery } from '#/api';
import { type identityservicev1_User as User } from '#/api';

// ✅ 正确 — composables 内部导入 apiClient
import { apiClient } from '#/api/client';

// 适配器导入
import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';

// 布局与通用组件
import { Page, useVbenDrawer } from '@vben/common-ui';

// 国际化
import { $t } from '@vben/locales';

// 图标（lucide）
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

// 消息提示
import { notification } from 'ant-design-vue';
```

```typescript
// ❌ 错误 — 禁止直接引用 generated 路径
import { useListUsers } from '#/api/composables/user';
import type { User } from '#/api/generated/admin/service/v1';
```

## 新建业务模块 Checklist

```
- [ ] Step 1: 确认 generated 层已有类型（protobuf 已生成）
- [ ] Step 2: 创建 composables 层（src/api/composables/xxx.ts）
- [ ] Step 3: 注册导出（composables/index.ts）
- [ ] Step 4: 添加 i18n 翻译（zh-CN + en-US 的 enum.json, menu.json, page.json）
- [ ] Step 5: 创建路由模块（router/routes/modules/app/xxx.ts）
- [ ] Step 6: 创建视图页面（views/app/xxx/）
```

### Step 1: 确认 generated 类型

确保 `src/api/generated/admin/service/v1` 中包含目标模块的：
- 类型定义: `xxxservicev1_EntityName`
- 请求/响应类型
- ApiClient 中已有 `xxxService` getter

### Step 2-3: API 层开发

见下方 [API 两层架构](#api-两层架构) 章节，包含完整的 composable 模板代码。

### Step 4: i18n 国际化

在四个 JSON 文件中添加对应 key：

| 文件 | 用途 | key 前缀示例 |
|---|---|---|
| `zh-CN/enum.json` | 枚举翻译 | `"enum.xxx.status.NORMAL": "正常"` |
| `zh-CN/menu.json` | 菜单名称 | `"menu.xxx.moduleName": "模块名"` |
| `zh-CN/page.json` | 页面文本 | `"page.xxx.name": "名称"` |
| `zh-CN/ui.json` | 通用 UI 文本 | 已有通用 key，通常无需修改 |

`en-US/` 目录下同步添加英文翻译。

### Step 5: 路由配置

新建 `router/routes/modules/app/xxx.ts`：

```typescript
import type { RouteRecordRaw } from 'vue-router';
import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    path: '/xxx',
    name: 'XxxManagement',
    component: BasicLayout,
    redirect: '/xxx/list',
    meta: {
      order: 3001,
      icon: 'lucide:xxx',
      title: $t('menu.xxx.moduleName'),
      authority: ['sys:platform_admin', 'sys:tenant_manager'],
    },
    children: [
      {
        path: 'list',
        name: 'XxxList',
        meta: { title: $t('menu.xxx.list'), authority: ['sys:platform_admin'] },
        component: () => import('#/views/app/xxx/index.vue'),
      },
    ],
  },
];
export default routes;
```

### Step 6: 视图页面开发

见下方 [视图/页面开发](#视图页面开发) 章节，包含列表页、Drawer 表单页的完整模板。

## API 两层架构

```
generated/  +  client.ts  →  composables/  →  views/stores
(自动生成)     (ApiClient 单例)   (Vue Query hooks)
```

**依赖方向**: `views/stores → composables → apiClient.xxxService → transport.unary → requestApi`

### client.ts — ApiClient 单例

文件路径: `src/api/client.ts`。生成的 `ApiClient` 通过 `ClientTransport` 接口发送请求，`client.ts` 将已有的 `requestApi`（基于 axios 的 `RequestClient`）适配为 `ClientTransport`，保留 token 注入、错误拦截、自动刷新等全部已有逻辑：

```typescript
import { type ClientTransport, createApiClient } from '#/api/generated/admin/service/v1';
import { requestApi } from '#/transport/rest';

const transport: ClientTransport = {
  unary(path, method, body, _meta) { return requestApi({ body, method, path }); },
  serverStream(path, _meta) { throw new Error(`serverStream not supported via ApiClient: ${path}`); },
  duplexStream(path, _meta) { throw new Error(`duplexStream not supported via ApiClient: ${path}`); },
};
export const apiClient = createApiClient(transport);
```

**ApiClient 提供的全部 Service Client getter**：

```typescript
apiClient.adminPortalService              // 管理门户
apiClient.apiService                      // API 管理
apiClient.authenticationService           // 认证服务
apiClient.dictEntryService                // 字典条目
apiClient.dictTypeService                 // 字典类型
apiClient.fileService                     // 文件管理
apiClient.fileTransferService             // 文件传输
apiClient.internalMessageService          // 内部消息
apiClient.internalMessageCategoryService  // 内部消息分类
apiClient.internalMessageRecipientService // 内部消息接收
apiClient.languageService                 // 语言管理
apiClient.loginPolicyService              // 登录策略
apiClient.menuService                     // 菜单管理
apiClient.orgUnitService                  // 组织单元
apiClient.permissionService               // 权限管理
apiClient.permissionGroupService          // 权限组管理
apiClient.positionService                 // 岗位管理
apiClient.roleService                     // 角色管理
apiClient.taskService                     // 异步任务
apiClient.tenantService                   // 租户管理
apiClient.userProfileService              // 用户资料
apiClient.userService                     // 用户管理
apiClient.apiAuditLogService              // API 审计日志
apiClient.dataAccessAuditLogService       // 数据访问审计日志
apiClient.loginAuditLogService            // 登录审计日志
apiClient.operationAuditLogService        // 操作审计日志
apiClient.permissionAuditLogService       // 权限审计日志
apiClient.policyEvaluationLogService      // 策略评估日志
```

### composables 层模板

文件路径: `src/api/composables/{module-name}.ts`。每个 composable 文件直接导入 `apiClient` 并调用对应的 Service Client：

```typescript
import type {
  xxxservicev1_GetXxxRequest,
  xxxservicev1_ListXxxResponse,
  xxxservicev1_Xxx,
  xxxservicev1_Xxx_Status as Xxx_Status,
} from '#/api/generated/admin/service/v1';

import { computed } from 'vue';
import { i18n } from '@vben/locales';
import { useMutation, type UseMutationOptions, useQuery, type UseQueryOptions } from '@tanstack/vue-query';
import { apiClient } from '#/api/client';
import { queryClient } from '#/plugins/vue-query';
import { makeUpdateMask, type PaginationQuery } from '#/transport/rest';

const t = i18n.global.t;

// 列表 — 组件内 hook
export function useListXxxs(query: PaginationQuery, options?: UseQueryOptions<xxxservicev1_ListXxxResponse, Error>) {
  return useQuery({ queryKey: ['listXxxs', query], queryFn: () => apiClient.xxxService.List(query.toRawParams()), ...options });
}

// 列表 — Store / 外部调用
export async function fetchListXxxs(params: PaginationQuery) {
  return queryClient.fetchQuery({ queryKey: ['listXxxs', params], queryFn: () => apiClient.xxxService.List(params.toRawParams()), retry: 0 });
}

// 详情 — 组件内 hook
export function useGetXxx(req: xxxservicev1_GetXxxRequest, options?: UseQueryOptions<xxxservicev1_Xxx, Error>) {
  return useQuery({ queryKey: ['getXxx', req], queryFn: () => apiClient.xxxService.Get(req), ...options });
}

// 详情 — Store / 外部调用
export async function fetchXxx(params: xxxservicev1_GetXxxRequest) {
  return queryClient.fetchQuery({ queryKey: ['getXxx', params], queryFn: () => apiClient.xxxService.Get(params), retry: 0 });
}

// 创建
export function useCreateXxx(options?: UseMutationOptions<object, Error, { data: xxxservicev1_Xxx }>) {
  return useMutation({ mutationFn: ({ data }) => apiClient.xxxService.Create({ data }), ...options });
}

// 更新（自动生成 updateMask）
export function useUpdateXxx(options?: UseMutationOptions<object, Error, { id: number; values: Record<string, any> }>) {
  return useMutation({
    mutationFn: ({ id, values }: { id: number; values: Record<string, any> }) =>
      apiClient.xxxService.Update({ id, data: { ...values } as any, updateMask: makeUpdateMask(Object.keys(values ?? {})) }),
    ...options,
  });
}

// 删除
export function useDeleteXxx(options?: UseMutationOptions<object, Error, number>) {
  return useMutation({ mutationFn: (id) => apiClient.xxxService.Delete({ id }), ...options });
}

// 枚举与工具函数
export const xxxStatusList = computed(() => [
  { value: 'ON', label: t('enum.xxx.status.ON') },
  { value: 'OFF', label: t('enum.xxx.status.OFF') },
]);

const XXX_STATUS_COLOR_MAP: Record<string, string> = { ON: '#52C41A', OFF: '#909399', DEFAULT: '#86909C' };

export function xxxStatusToColor(status: Xxx_Status) {
  return XXX_STATUS_COLOR_MAP[status as string] ?? XXX_STATUS_COLOR_MAP.DEFAULT ?? '#86909C';
}

export function xxxStatusToName(status?: Xxx_Status) {
  const map: Record<string, string> = { ON: t('enum.xxx.status.ON'), OFF: t('enum.xxx.status.OFF') };
  return map[status as string] ?? '';
}
```

**注册导出** — 在 `src/api/composables/index.ts` 中添加：`export * from './xxx';`

### 特殊模式：Store 直接调用的纯异步函数

某些函数需要同时供 Pinia Store 和 Mutation builder 使用（如 auth 的 `login`、`logout`），先定义为纯异步函数：

```typescript
// composables/auth.ts
import { apiClient } from '#/api/client';

// 纯异步函数 — 供 Store 和 Mutation builder 使用
export async function login(request: authenticationservicev1_LoginRequest) {
  return apiClient.authenticationService.Login(request);
}

export async function logout() {
  return apiClient.authenticationService.Logout({});
}

// Mutation builder — 供 Store 在非 vue 上下文中使用
export const loginMutation = queryClient.getMutationCache().build(queryClient, {
  mutationKey: ['login'],
  mutationFn: login,
  retry: 0,
});
```

### Query Key 命名约定

| 操作 | Query Key 格式 | 示例 |
|---|---|---|
| 列表 | `["list{Entity}s", query]` | `["listUsers", query]` |
| 详情 | `["get{Entity}", req]` | `["getUser", { id: 1 }]` |
| 创建 / 更新 / 删除 | mutation（无 query key） | — |

### PaginationQuery 使用模式

```typescript
import { PaginationQuery } from '#/api';

new PaginationQuery({ paging: { page: 1, pageSize: 20 } });                           // 基础分页
new PaginationQuery({ paging: { page: 1, pageSize: 20 }, formValues: { status: 'ON' } }); // 带搜索
new PaginationQuery({ orderBy: ['-created_at', 'name'] });                              // 带排序
new PaginationQuery({ formValues: { status: 'ON' } });                                  // 不分页（获取全量）
new PaginationQuery({ fieldMask: 'id,name,status' });                                   // 只返回指定字段
```

### Vue Query 全局配置

| 配置 | 值 | 说明 |
|---|---|---|
| `staleTime` | 60s | 数据在 60 秒内视为新鲜 |
| `retry` | false | 失败不自动重试 |
| `refetchOnWindowFocus` | false | 窗口聚焦不刷新 |
| `refetchOnReconnect` | false | 网络重连不刷新 |

## 视图/页面开发

### 标准模块文件结构

```
views/app/{module}/
├── index.vue              # 列表页（VxeTable + 搜索表单）
└── {entity}-drawer.vue    # 新建/编辑 Drawer 表单
```

如需详情页可增加 `detail/` 子目录。

### 核心概念：Grid 的两套 Schema

| 配置项 | 作用 |
|---|---|
| `formOptions.schema` | **搜索表单**（使用 VbenForm 组件） |
| `gridOptions.columns` | **表格列**（使用 VxeTable 列配置） |

两者通过 `useVbenVxeGrid({ gridOptions, formOptions })` 合并，搜索表单的值自动传入 `proxyConfig.ajax.query` 的 `formValues`。

### 搜索表单组件 (formOptions.schema)

搜索表单使用 VbenForm 体系，`component` 字段指定组件类型。**必须使用以下注册过的组件名，不要使用原生 HTML 或 Ant Design Vue 的原始组件。**

| component | 用途 | 关键 props |
|---|---|---|
| `Input` | 文本输入 | `placeholder`, `allowClear` |
| `InputNumber` | 数字输入 | `placeholder`, `allowClear`, `defaultValue` |
| `Select` | 下拉选择 | `options`, `showSearch`, `allowClear`, `filterOption` |
| `ApiSelect` | 远程下拉 | `api`, `afterFetch`, `showSearch`, `allowClear`, `filterOption`, `immediate`, `alwaysLoad` |
| `ApiTreeSelect` | 远程树选择 | `api`, `treeDefaultExpandAll`, `labelField`, `valueField`, `childrenField`, `numberToString` |
| `RadioGroup` | 单选按钮组 | `optionType: 'button'`, `buttonStyle: 'solid'`, `options` |
| `RangePicker` | 日期范围 | `showTime`, `allowClear`, `presets` |
| `DatePicker` | 日期选择 | `placeholder`, `allowClear` |
| `Switch` | 开关 | — |
| `Textarea` | 多行文本 | `placeholder`, `allowClear` |

**Input — 文本搜索**

```typescript
{
  component: 'Input',
  fieldName: 'name',
  label: $t('page.xxx.name'),
  componentProps: { placeholder: $t('ui.placeholder.input'), allowClear: true },
}
```

**Select — 下拉筛选（本地选项）**

```typescript
{
  component: 'Select',
  fieldName: 'status',
  label: $t('ui.table.status'),
  componentProps: {
    options: xxxStatusList,         // 从 composables 导入的 computed 选项列表
    placeholder: $t('ui.placeholder.select'),
    allowClear: true,
    showSearch: true,
    filterOption: (input: string, option: any) =>
      option.label.toLowerCase().includes(input.toLowerCase()),
  },
}
```

> `options` 必须是 `{ value, label }[]` 格式，通常从 `composables` 层的 computed 导出。`filterOption` 是搜索过滤函数，**必须提供**才能支持输入搜索。

**ApiSelect — 远程数据下拉**

```typescript
{
  component: 'ApiSelect',
  fieldName: 'roleId',
  label: $t('page.user.form.role'),
  componentProps: {
    allowClear: true,
    showSearch: true,
    placeholder: $t('ui.placeholder.select'),
    filterOption: (input: string, option: any) =>
      option.label.toLowerCase().includes(input.toLowerCase()),
    afterFetch: (data: any[]) => data.map((item: any) => ({ label: item.name, value: item.id })),
    api: async () => {
      const result = await fetchListRoles(new PaginationQuery({ formValues: { status: 'ON' } }));
      return result.items;
    },
  },
}
```

> **关键**: `afterFetch` 必须将数据转为 `{ label, value }[]` 格式。`api` 返回的是原始 API 数据。

**ApiTreeSelect — 远程树形下拉**

```typescript
{
  component: 'ApiTreeSelect',
  fieldName: 'orgUnitId',
  label: $t('page.position.orgUnit'),
  componentProps: {
    placeholder: $t('ui.placeholder.select'),
    numberToString: true,            // 数字 ID 转字符串（TreeSelect 需要）
    showSearch: true,
    treeDefaultExpandAll: true,
    allowClear: true,
    childrenField: 'children',       // 子节点字段名
    labelField: 'name',              // 显示文本字段
    valueField: 'id',                // 值字段
    treeNodeFilterProp: 'label',     // 搜索过滤属性
    api: async () => {
      const result = await fetchListOrgUnits(new PaginationQuery({ formValues: { status: 'ON' } }));
      return result.items;
    },
  },
}
```

**RangePicker — 日期范围搜索**

```typescript
import dayjs from 'dayjs';

{
  component: 'RangePicker',
  fieldName: 'createdAt',
  label: $t('page.xxx.createdAt'),
  componentProps: {
    showTime: true,
    allowClear: true,
    presets: [
      { label: $t('ui.dateRange.today'), value: [dayjs().startOf('day'), dayjs().endOf('day')] },
      { label: $t('ui.dateRange.yesterday'), value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')] },
      { label: $t('ui.dateRange.thisWeek'), value: [dayjs().startOf('week'), dayjs().endOf('week')] },
      { label: $t('ui.dateRange.lastWeek'), value: [dayjs().subtract(1, 'week').startOf('week'), dayjs().subtract(1, 'week').endOf('week')] },
      { label: $t('ui.dateRange.thisMonth'), value: [dayjs().startOf('month'), dayjs().endOf('month')] },
      { label: $t('ui.dateRange.lastMonth'), value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')] },
    ],
  },
}
```

> **query 中使用**: RangePicker 的值是 `[dayjs, dayjs]` 数组，在 `proxyConfig.ajax.query` 中需手动拆分为 `created_at__gte` 和 `created_at__lte` 参数。

**搜索表单完整配置**

```typescript
const formOptions: VbenFormProps = {
  collapsed: false,               // 默认是否折叠
  showCollapseButton: false,      // 是否显示展开/折叠按钮（字段多时设为 true）
  submitOnEnter: true,            // 回车提交搜索
  schema: [ /* ... */ ],
};
```

### 表格列配置 (gridOptions.columns)

| 用法 | 配置方式 | 说明 |
|---|---|---|
| 普通文本列 | `{ title, field }` | 直接显示字段值 |
| 序号列 | `{ type: 'seq' }` | 自动生成行号 |
| 日期列 | `{ title, field, formatter: 'formatDateTime' }` | 自动格式化时间（全局已注册） |
| 状态/枚举列 | `{ title, field, slots: { default: 'slotName' } }` | 通过 Slot 渲染 Tag |
| 操作列 | `{ title, field: 'action', fixed: 'right', slots: { default: 'action' } }` | 编辑/删除按钮 |
| 树节点列 | `{ title, field, treeNode: true }` | 树形表格的展开列 |

普通文本列支持可选参数：`width`（固定列宽）、`align`（对齐，长文本建议 `'left'`）、`showOverflow: 'tooltip'`（超长显示 tooltip）。

**嵌套对象字段** — VxeTable 支持点号路径访问：

```typescript
{ title: '操作系统', field: 'deviceInfo.osName' }
{ title: '菜单标题', field: 'meta.title' }
{ title: '排序', field: 'meta.order' }
```

**树形表格配置**（如组织架构、菜单管理）需要额外配置 `treeConfig`：

```typescript
const gridOptions: VxeGridProps<OrgUnit> = {
  pagerConfig: { enabled: false },    // 树形表格通常不分页
  treeConfig: {
    parentField: 'parentId',   // 方式一：通过 parentField 自动构建树
    rowField: 'id',
    transform: true,
  },
  // 或者（后端已返回树形结构时）：
  // treeConfig: { childrenField: 'children', rowField: 'id' },
  columns: [
    { title: $t('page.xxx.name'), field: 'name', treeNode: true },  // 树节点列必须标记 treeNode: true
    // ...其他列
  ],
};
```

> **注意**: 树形表格必须将 `pagerConfig.enabled` 设为 `false`，否则会冲突。

### Slot 用法详解

在 `columns` 中通过 `slots: { default: 'slotName' }` 声明，然后在 `<template #slotName="{ row }">` 中实现渲染。**必须使用 Ant Design Vue 的 `a-tag`、`a-button`、`a-popconfirm` 等组件，不要使用原生 HTML 元素。**

**状态 Tag Slot（最常用）**

```html
<template #status="{ row }">
  <a-tag :color="xxxStatusToColor(row.status)">
    {{ xxxStatusToName(row.status) }}
  </a-tag>
</template>
```

> 颜色和名称映射函数在 `composables` 层定义（如 `userStatusToColor`、`statusToName`），从 `#/api` 导入。**不要**在 Slot 中硬编码颜色值或用 if/else 判断。

**多值 Tag Slot（数组字段，如角色列表）**

```html
<template #role="{ row }">
  <div>
    <a-tag
      v-for="role in row.roleNames"
      :key="role"
      class="mb-1 mr-1"
      :style="{ backgroundColor: getRandomColor(role), color: '#333', border: 'none' }"
    >
      {{ role }}
    </a-tag>
  </div>
</template>
```

> 使用 `getRandomColor()` 工具函数（从 `#/utils/color` 导入）生成随机背景色。

**布尔值 Tag Slot（启用/禁用）**

```html
<template #isEnabled="{ row }">
  <a-tag :color="enableBoolToColor(row.isEnabled)">
    {{ enableBoolToName(row.isEnabled) }}
  </a-tag>
</template>
```

> `enableBoolToColor` 和 `enableBoolToName` 从 `#/api` 的 `shared.ts` 导出。

**嵌套对象 Slot**

```html
<template #geoLocation="{ row }">
  {{ row.geoLocation.province }} {{ row.geoLocation.city }}
</template>
```

**带图标的树节点 Slot（菜单管理）**

```html
<template #title="{ row }">
  <div class="flex w-full items-center gap-1">
    <div class="size-5 flex-shrink-0">
      <IconifyIcon v-if="row.type === 'button'" icon="carbon:security" class="size-full" />
      <IconifyIcon
        v-else-if="row.meta?.icon"
        :icon="row.meta?.icon || 'carbon:circle-dash'"
        class="size-full"
      />
    </div>
    <span class="flex-auto">{{ $t(row.meta?.title) }}</span>
  </div>
</template>
```

> `IconifyIcon` 从 `@vben/icons` 导入。

**操作列 Slot（固定模式）**

```html
<template #action="{ row }">
  <a-button type="link" :icon="h(LucideFilePenLine)" @click.stop="handleEdit(row)" />
  <a-popconfirm
    :cancel-text="$t('ui.button.cancel')"
    :ok-text="$t('ui.button.ok')"
    :title="$t('ui.text.do_you_want_delete', { moduleName: $t('page.xxx.moduleName') })"
    @confirm="handleDelete(row)"
  >
    <a-button danger type="link" :icon="h(LucideTrash2)" />
  </a-popconfirm>
</template>
```

> 图标使用 `h()` 渲染函数（`from 'vue'` 导入 `h`），图标组件从 `@vben/icons` 导入。删除前必须用 `a-popconfirm` 二次确认，**不要**用 `window.confirm`。

带「查看详情」按钮的操作列，在最前面加一个 `<a-button type="link" :icon="h(LucideInfo)" @click.stop="handleDetail(row)" />` 即可。

**工具栏 Slot (toolbar-tools)**

```html
<template #toolbar-tools>
  <a-button class="mr-2" type="primary" @click="handleCreate">
    {{ $t('page.xxx.button.create') }}
  </a-button>
</template>
```

树形表格的展开/折叠按钮：

```html
<template #toolbar-tools>
  <a-button class="mr-2" type="primary" @click="handleCreate">{{ $t('page.xxx.button.create') }}</a-button>
  <a-button class="mr-2" @click="expandAll">{{ $t('ui.tree.expand_all') }}</a-button>
  <a-button class="mr-2" @click="collapseAll">{{ $t('ui.tree.collapse_all') }}</a-button>
</template>
```

```typescript
const expandAll = () => gridApi.grid?.setAllTreeExpand(true);
const collapseAll = () => gridApi.grid?.setAllTreeExpand(false);
```

### VxeGrid 全局配置

**基础列表页配置（带分页）**

```typescript
const gridOptions: VxeGridProps<Xxx> = {
  toolbarConfig: { custom: true, export: true, refresh: true, zoom: true },
  exportConfig: {},
  pagerConfig: {},              // 启用分页（默认）
  rowConfig: { isHover: true }, // 行悬停高亮
  height: 'auto',               // 自适应高度
  stripe: true,                 // 斑马纹
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        return await fetchListXxxs(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues,
          }),
        );
      },
    },
  },
  columns: [ /* ... */ ],
};
```

不分页列表：设 `pagerConfig: { enabled: false }`。

**带 Tooltip 的配置**

```typescript
tooltipConfig: {
  showAll: true,
  enterable: true,
  contentMethod: ({ column, row }) => {
    if (column.field === 'roleNames') return `${row[column.field]}`;
    return null; // 其余列使用默认行为
  },
},
```

列上配合 `showOverflow: 'tooltip'`。

**Grid 事件监听**

```typescript
import { type VxeGridListeners } from '#/adapter/vxe-table';

const gridEvents: VxeGridListeners<User> = {
  cellDblclick: ({ row }) => { handleDetail(row); },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions, gridEvents });
```

### Drawer 表单组件 (useVbenForm)

Drawer 中的表单同样使用 VbenForm 体系，组件类型与搜索表单相同。表单 Schema 可用组件：`Input`, `InputNumber`, `Select`, `ApiSelect`, `ApiTreeSelect`, `RadioGroup`, `Textarea`, `Switch`, `DatePicker`。

**表单校验规则**（已在 `adapter/form.ts` 全局注册，自动国际化提示）：

```typescript
rules: 'required'          // 必填（输入框类）
rules: 'selectRequired'    // 必选（下拉/单选类）
```

> **不要**自行写 `rules: z.string().min(1)` 之类的 Zod 校验。

**表单默认值** — 通过 `defaultValue` 设置：

```typescript
{ component: 'RadioGroup', fieldName: 'status', defaultValue: 'ON', /* ... */ },
{ component: 'InputNumber', fieldName: 'sortOrder', defaultValue: 1, /* ... */ },
```

**表单公共配置**

```typescript
const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,   // 隐藏默认的提交/重置按钮
  commonConfig: { componentProps: { class: 'w-full' } },  // 所有表单项宽度 100%
  schema: [ /* ... */ ],
});
```

### Drawer 表单模板 (xxx-drawer.vue)

```vue
<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { PaginationQuery, xxxStatusList, useCreateXxx, useUpdateXxx } from '#/api';

const { mutateAsync: createXxx } = useCreateXxx();
const { mutateAsync: updateXxx } = useUpdateXxx();

const data = ref();

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('page.xxx.moduleName') })
    : $t('ui.modal.update', { moduleName: $t('page.xxx.moduleName') }),
);

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: { componentProps: { class: 'w-full' } },
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.xxx.name'),
      componentProps: { placeholder: $t('ui.placeholder.input'), allowClear: true },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: $t('page.xxx.code'),
      componentProps: { placeholder: $t('ui.placeholder.input'), allowClear: true },
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: $t('ui.table.status'),
      defaultValue: 'ON',
      rules: 'selectRequired',
      componentProps: { optionType: 'button', buttonStyle: 'solid', class: 'flex flex-wrap', options: xxxStatusList },
    },
    { component: 'Textarea', fieldName: 'remark', label: $t('ui.table.remark') },
  ],
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() { drawerApi.close(); },

  async onConfirm() {
    const validate = await baseFormApi.validate();
    if (!validate.valid) return;

    setLoading(true);
    const values = await baseFormApi.getValues();

    try {
      await (data.value?.create ? createXxx(values) : updateXxx({ id: data.value.row.id, values }));
      notification.success({
        message: data.value?.create ? $t('ui.notification.create_success') : $t('ui.notification.update_success'),
      });
    } catch {
      notification.error({
        message: data.value?.create ? $t('ui.notification.create_failed') : $t('ui.notification.update_failed'),
      });
    } finally {
      drawerApi.close();
      setLoading(false);
    }
  },

  onOpenChange(isOpen) {
    if (isOpen) {
      data.value = drawerApi.getData<Record<string, any>>();
      baseFormApi.setValues(data.value?.row);
      setLoading(false);
    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ loading });
}
</script>

<template>
  <Drawer :title="getTitle">
    <BaseForm />
  </Drawer>
</template>
```

### 完整列表页模板 (index.vue)

```vue
<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { Page, useVbenDrawer, type VbenFormProps } from '@vben/common-ui';
import { LucideFilePenLine, LucideTrash2 } from '@vben/icons';

import { notification } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  PaginationQuery,
  xxxStatusList,
  xxxStatusToColor,
  xxxStatusToName,
  useDeleteXxx,
  fetchListXxxs,
} from '#/api';
import { type xxxservicev1_Xxx as Xxx } from '#/api';
import { $t } from '@vben/locales';

import XxxDrawer from './xxx-drawer.vue';

const { mutateAsync: deleteXxx } = useDeleteXxx();

// 搜索表单
const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.xxx.name'),
      componentProps: { placeholder: $t('ui.placeholder.input'), allowClear: true },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('ui.table.status'),
      componentProps: {
        options: xxxStatusList,
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
        showSearch: true,
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
      },
    },
  ],
};

// 表格配置
const gridOptions: VxeGridProps<Xxx> = {
  toolbarConfig: { custom: true, export: true, refresh: true, zoom: true },
  exportConfig: {},
  pagerConfig: {},
  rowConfig: { isHover: true },
  height: 'auto',
  stripe: true,
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        return await fetchListXxxs(
          new PaginationQuery({
            paging: { page: page.currentPage, pageSize: page.pageSize },
            formValues,
          }),
        );
      },
    },
  },
  columns: [
    { title: $t('page.xxx.name'), field: 'name' },
    { title: $t('page.xxx.code'), field: 'code' },
    { title: $t('ui.table.status'), field: 'status', slots: { default: 'status' }, width: 95 },
    { title: $t('ui.table.createdAt'), field: 'createdAt', formatter: 'formatDateTime', width: 140 },
    { title: $t('ui.table.remark'), field: 'remark' },
    { title: $t('ui.table.action'), field: 'action', fixed: 'right', slots: { default: 'action' }, width: 90 },
  ],
};

// Grid + Drawer 初始化
const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: XxxDrawer,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) gridApi.reload();
  },
});

function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({ create, row });
  drawerApi.open();
}

function handleCreate() { openDrawer(true); }
function handleEdit(row: any) { openDrawer(false, row); }

async function handleDelete(row: any) {
  try {
    await deleteXxx(row.id);
    notification.success({ message: $t('ui.notification.delete_success') });
    await gridApi.reload();
  } catch {
    notification.error({ message: $t('ui.notification.delete_failed') });
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid :table-title="$t('menu.xxx.moduleName')">
      <template #toolbar-tools>
        <a-button class="mr-2" type="primary" @click="handleCreate">
          {{ $t('page.xxx.button.create') }}
        </a-button>
      </template>
      <template #status="{ row }">
        <a-tag :color="xxxStatusToColor(row.status)">
          {{ xxxStatusToName(row.status) }}
        </a-tag>
      </template>
      <template #action="{ row }">
        <a-button type="link" :icon="h(LucideFilePenLine)" @click.stop="handleEdit(row)" />
        <a-popconfirm
          :cancel-text="$t('ui.button.cancel')"
          :ok-text="$t('ui.button.ok')"
          :title="$t('ui.text.do_you_want_delete', { moduleName: $t('page.xxx.moduleName') })"
          @confirm="handleDelete(row)"
        >
          <a-button danger type="link" :icon="h(LucideTrash2)" />
        </a-popconfirm>
      </template>
    </Grid>
    <Drawer />
  </Page>
</template>
```

## Vben 框架 API 速查

### Page 组件（必须包裹）

所有页面 **必须** 用 `<Page>` 组件包裹，`auto-content-height` 自适应内容高度。从 `@vben/common-ui` 导入。**禁止**用 `<div>` 或原生元素替代。

### Grid API 方法

`useVbenVxeGrid` 返回 `[Grid, gridApi]`：

| 方法 | 用途 | 常见场景 |
|---|---|---|
| `gridApi.reload()` | 重新加载数据 | 删除/编辑后刷新表格 |
| `gridApi.grid` | VxeTable 实例 | 调用原生 VxeTable API |

原生方法：`gridApi.grid?.setAllTreeExpand(true/false)` — 展开/折叠所有树节点。

### Drawer API 方法

`useVbenDrawer` 返回 `[Drawer, drawerApi]`：

| 方法 | 用途 | 说明 |
|---|---|---|
| `drawerApi.open()` / `close()` | 打开 / 关闭抽屉 | — |
| `drawerApi.setData(data)` | 传入数据 | 通常传 `{ create, row }` |
| `drawerApi.getData<T>()` | 获取传入数据 | 在 `onOpenChange` 中使用 |
| `drawerApi.setState({ loading })` | 设置加载状态 | 显示/隐藏提交按钮的 loading |

### Form API 方法

`useVbenForm` 返回 `[Form, formApi]`：

| 方法 | 用途 | 说明 |
|---|---|---|
| `formApi.getValues()` | 获取表单值 | 返回 `Promise<Record>` |
| `formApi.setValues(obj)` | 设置表单值 | 编辑时回填数据 |
| `formApi.validate()` | 校验表单 | 返回 `{ valid: boolean }` |
| `formApi.resetForm()` | 重置表单 | 清空所有字段（关闭 Drawer 时框架自动重置，无需手动清空） |

### Notification（消息提示）

```typescript
import { notification } from 'ant-design-vue';

notification.success({ message: $t('ui.notification.delete_success') });
notification.error({ message: $t('ui.notification.delete_failed') });
```

> **禁止**使用 `alert()`、`console.log` 向用户展示操作结果，必须使用 `notification`。项目已预定义通用 i18n key：`ui.notification.create_success`、`update_success`、`delete_success` 及对应 `_failed` 版本。

## 常见错误与纠正

| 错误做法 | 正确做法 |
|---|---|
| 原生 `<span style="color:red">正常</span>` | `<a-tag :color="statusToColor(row.status)">` |
| `window.confirm('确定删除?')` | `<a-popconfirm>` |
| 原生 `<button>` | `<a-button type="primary/link">` |
| Slot 中 `dayjs(row.createdAt).format(...)` | `formatter: 'formatDateTime'` |
| `rules: z.string().min(1, '不能为空')` | `rules: 'required'` 或 `'selectRequired'` |
| Slot 中 `v-if` 硬编码颜色 | composables 层的 `xxxToColor()` |
| `import { Tag } from 'ant-design-vue'` 后用 `<Tag>` | 直接用 `<a-tag>`（全局注册） |
| `<div>` 包裹页面 | `<Page auto-content-height>` |
| `ref(false)` 管理 loading | `drawerApi.setState({ loading })` |
| 表单提交后手动清空字段 | 关闭 Drawer 时框架自动重置 |
| `alert()` / `ElMessage` 提示 | `notification` from `ant-design-vue` |
| `component: 'AInput'` 或 `'a-input'` | `component: 'Input'`（注册名） |
| 手动 watch 搜索条件重新请求 | 搜索表单与 Grid 自动关联，通过 `proxyConfig.ajax.query` 的 `formValues` 获取 |
| `:icon="LucideTrash2"` 直接传组件 | `:icon="h(LucideTrash2)"` 用 h() 包裹 |
| `ref` + `watch` 管理列表数据 | `proxyConfig.ajax.query` + `PaginationQuery`，由框架自动管理 |
