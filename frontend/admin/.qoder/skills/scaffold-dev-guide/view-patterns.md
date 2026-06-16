# 视图/页面开发详细模板

## 标准模块文件结构

```
views/app/{module}/
├── index.vue              # 列表页（VxeTable + 搜索表单）
└── {entity}-drawer.vue    # 新建/编辑 Drawer 表单
```

如需详情页可增加 `detail/` 子目录。

---

## 核心概念：Grid 的两套 Schema

本项目中每个列表页都涉及 **两套独立的 Schema 配置**，不要混淆：

| 配置项                   | 作用                  | 位置                                |
|-----------------------|---------------------|-----------------------------------|
| `formOptions.schema`  | **搜索表单**（表格上方的筛选区域） | 使用 `#/adapter/form` 的 VbenForm 组件 |
| `gridOptions.columns` | **表格列**（数据展示区域）     | 使用 VxeTable 的列配置                  |

两者通过 `useVbenVxeGrid({ gridOptions, formOptions })` 合并为一个 Grid 组件，搜索表单的值会自动传入 `proxyConfig.ajax.query` 的第二个参数 `formValues`。

---

## 搜索表单组件 (formOptions.schema)

搜索表单使用 VbenForm 体系，`component` 字段指定组件类型。**必须使用以下注册过的组件名，不要使用原生 HTML 或 Ant Design Vue 的原始组件。**

### 可用组件列表

| component       | 用途    | 关键 props                                                                                     |
|-----------------|-------|----------------------------------------------------------------------------------------------|
| `Input`         | 文本输入  | `placeholder`, `allowClear`                                                                  |
| `InputNumber`   | 数字输入  | `placeholder`, `allowClear`, `defaultValue`                                                  |
| `Select`        | 下拉选择  | `options`, `showSearch`, `allowClear`, `filterOption`                                        |
| `ApiSelect`     | 远程下拉  | `api`, `afterFetch`, `showSearch`, `allowClear`, `filterOption`, `immediate`, `alwaysLoad`   |
| `ApiTreeSelect` | 远程树选择 | `api`, `treeDefaultExpandAll`, `labelField`, `valueField`, `childrenField`, `numberToString` |
| `RadioGroup`    | 单选按钮组 | `optionType: 'button'`, `buttonStyle: 'solid'`, `options`                                    |
| `Textarea`      | 多行文本  | `placeholder`, `allowClear`                                                                  |
| `Switch`        | 开关    | —                                                                                            |
| `DatePicker`    | 日期选择  | `placeholder`, `allowClear`                                                                  |
| `RangePicker`   | 日期范围  | `showTime`, `allowClear`, `presets`                                                          |

### Input — 文本搜索

```typescript
{
  component: 'Input',
  fieldName: 'name',
  label: $t('page.xxx.name'),
  componentProps: {
    placeholder: $t('ui.placeholder.input'),
    allowClear: true,
  },
}
```

### Select — 下拉筛选（本地选项）

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

> **注意**: `options` 必须是 `{ value, label }[]` 格式的数组，通常从 `composables` 层的 computed 导出（如 `statusList`、`userStatusList`）。
> `filterOption` 是搜索过滤函数，**必须提供**才能支持输入搜索。

### ApiSelect — 远程数据下拉

用于从 API 获取选项数据的下拉选择（如角色列表、岗位列表）：

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
    afterFetch: (data: any[]) => {
      // 将 API 返回数据转为 { label, value } 格式
      return data.map((item: any) => ({
        label: item.name,
        value: item.id,
      }));
    },
    api: async () => {
      const result = await fetchListRoles(
        new PaginationQuery({
          formValues: { status: 'ON' },
        }),
      );
      return result.items;
    },
  },
}
```

> **关键**: `afterFetch` 必须将数据转为 `{ label, value }[]` 格式。`api` 返回的是原始 API 数据。

### ApiTreeSelect — 远程树形下拉

用于选择树形结构数据（如组织架构、菜单分类）：

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
      const result = await fetchListOrgUnits(
        new PaginationQuery({
          formValues: { status: 'ON' },
        }),
      );
      return result.items;
    },
  },
}
```

### RangePicker — 日期范围搜索

用于按时间范围筛选（审计日志等场景）：

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

### 搜索表单完整配置

```typescript
const formOptions: VbenFormProps = {
  collapsed: false,               // 默认是否折叠
  showCollapseButton: false,      // 是否显示展开/折叠按钮（字段多时设为 true）
  submitOnEnter: true,            // 回车提交搜索
  schema: [ /* ... */ ],
};
```

---

## 表格列配置 (gridOptions.columns)

表格使用 VxeTable，通过 `columns` 数组定义每一列。

### 列类型速查

| 用法     | 配置方式                                                                       | 说明             |
|--------|----------------------------------------------------------------------------|----------------|
| 普通文本列  | `{ title, field }`                                                         | 直接显示字段值        |
| 序号列    | `{ type: 'seq' }`                                                          | 自动生成行号         |
| 状态/枚举列 | `{ title, field, slots: { default: 'slotName' } }`                         | 通过 Slot 渲染 Tag |
| 日期列    | `{ title, field, formatter: 'formatDateTime' }`                            | 自动格式化时间        |
| 操作列    | `{ title, field: 'action', fixed: 'right', slots: { default: 'action' } }` | 编辑/删除按钮        |
| 树节点列   | `{ title, field, treeNode: true }`                                         | 树形表格的展开列       |

### 普通文本列

```typescript
{ title: $t('page.xxx.name'), field: 'name' }
```

支持可选参数：
- `width: 120` — 固定列宽
- `align: 'left'` — 对齐方式（默认居中，长文本建议左对齐）
- `showOverflow: 'tooltip'` — 超长文本显示 tooltip

### 序号列

```typescript
{ title: $t('ui.table.seq'), type: 'seq', width: 50 }
```

### 日期时间列

```typescript
{
  title: $t('ui.table.createdAt'),
  field: 'createdAt',
  formatter: 'formatDateTime',
  width: 140,
}
```

> **注意**: `formatDateTime` 已在 `adapter/vxe-table.ts` 中全局注册，直接使用即可。**不要**自行用 Slot 手动格式化日期。

### 嵌套对象字段

VxeTable 支持点号路径访问嵌套字段：

```typescript
{ title: '操作系统', field: 'deviceInfo.osName' }
{ title: '菜单标题', field: 'meta.title' }
{ title: '排序', field: 'meta.order' }
```

### 树形表格配置

树形表格（如组织架构、菜单管理）需要额外配置 `treeConfig`：

```typescript
const gridOptions: VxeGridProps<OrgUnit> = {
  // ...其他配置
  pagerConfig: {
    enabled: false,    // 树形表格通常不分页
  },
  treeConfig: {
    parentField: 'parentId',   // 方式一：通过 parentField 自动构建树
    rowField: 'id',            // 行唯一标识
    transform: true,           // 自动将扁平数据转为树形
  },
  // 或者（后端已返回树形结构时）：
  treeConfig: {
    childrenField: 'children', // 方式二：直接使用后端返回的 children 字段
    rowField: 'id',
  },
  columns: [
    // 树节点列必须标记 treeNode: true
    { title: $t('page.xxx.name'), field: 'name', treeNode: true },
    // ...其他列
  ],
};
```

> **注意**: 树形表格必须将 `pagerConfig.enabled` 设为 `false`，否则会冲突。

---

## Slot 用法详解

Slot 是 VxeTable 中自定义单元格渲染的核心机制。在 `columns` 中通过 `slots: { default: 'slotName' }` 声明，然后在 `<template #slotName="{ row }">` 中实现渲染。

> **重要**: 本项目使用 **Ant Design Vue** 的 `a-tag`、`a-button`、`a-popconfirm` 等组件来渲染 Slot 内容。不要使用原生 HTML 元素（如 `<span style="color:red">`、`<div class="badge">`）来替代。

### Slot 的工作流程

```
1. columns 中声明 slot
   → { field: 'status', slots: { default: 'status' } }

2. template 中实现 slot
   → <template #status="{ row }">
        <a-tag :color="statusToColor(row.status)">
          {{ statusToName(row.status) }}
        </a-tag>
      </template>
```

### 状态 Tag Slot（最常用）

用于将枚举字段渲染为彩色标签：

```typescript
// columns 中
{
  title: $t('ui.table.status'),
  field: 'status',
  slots: { default: 'status' },
  width: 95,
}
```

```html
<!-- template 中 -->
<template #status="{ row }">
  <a-tag :color="xxxStatusToColor(row.status)">
    {{ xxxStatusToName(row.status) }}
  </a-tag>
</template>
```

> **规范**: 颜色和名称映射函数在 `composables` 层定义（如 `userStatusToColor`、`statusToName`），从 `#/api` 导入。**不要**在 Slot 中硬编码颜色值或使用 if/else 判断。

### 多值 Tag Slot（数组字段）

用于渲染数组类型的字段（如角色列表、标签列表）：

```html
<template #role="{ row }">
  <div>
    <a-tag
      v-for="role in row.roleNames"
      :key="role"
      class="mb-1 mr-1"
      :style="{
        backgroundColor: getRandomColor(role),
        color: '#333',
        border: 'none',
      }"
    >
      {{ role }}
    </a-tag>
  </div>
</template>
```

> 使用 `getRandomColor()` 工具函数（从 `#/utils/color` 导入）生成随机背景色。

### 布尔值 Tag Slot

用于启用/禁用等布尔字段：

```html
<template #isEnabled="{ row }">
  <a-tag :color="enableBoolToColor(row.isEnabled)">
    {{ enableBoolToName(row.isEnabled) }}
  </a-tag>
</template>
```

> `enableBoolToColor` 和 `enableBoolToName` 从 `#/api` 的 `shared.ts` 导出。

### 嵌套对象 Slot

用于显示嵌套对象中的多个字段：

```html
<template #geoLocation="{ row }">
  {{ row.geoLocation.province }} {{ row.geoLocation.city }}
</template>

<template #platform="{ row }">
  {{ row.deviceInfo.osName }} {{ row.deviceInfo.browserName }}
</template>
```

### 带图标的树节点 Slot

用于菜单管理等需要显示图标的树形表格：

```html
<template #title="{ row }">
  <div class="flex w-full items-center gap-1">
    <div class="size-5 flex-shrink-0">
      <IconifyIcon
        v-if="row.type === 'button'"
        icon="carbon:security"
        class="size-full"
      />
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

### 权限标签 Slot

用于显示多个权限标签（带随机颜色）：

```html
<template #authority="{ row }">
  <a-tag
    v-for="auth in normalizeAuthority(row.meta?.authority)"
    :key="auth"
    class="mb-1 mr-1"
    :style="{
      backgroundColor: getRandomColor(auth),
      color: '#333',
      border: 'none',
    }"
  >
    {{ auth }}
  </a-tag>
</template>
```

### 操作列 Slot（固定模式）

操作列使用统一的编辑+删除按钮模式：

```typescript
// columns 中
{
  title: $t('ui.table.action'),
  field: 'action',
  fixed: 'right',
  slots: { default: 'action' },
  width: 90,
}
```

```html
<template #action="{ row }">
  <a-button
    type="link"
    :icon="h(LucideFilePenLine)"
    @click.stop="handleEdit(row)"
  />
  <a-popconfirm
    :cancel-text="$t('ui.button.cancel')"
    :ok-text="$t('ui.button.ok')"
    :title="
      $t('ui.text.do_you_want_delete', {
        moduleName: $t('page.xxx.moduleName'),
      })
    "
    @confirm="handleDelete(row)"
  >
    <a-button danger type="link" :icon="h(LucideTrash2)" />
  </a-popconfirm>
</template>
```

> **注意**: 图标使用 `h()` 渲染函数方式传入（`from 'vue'` 导入 `h`），图标组件从 `@vben/icons` 导入（如 `LucideFilePenLine`、`LucideTrash2`、`LucideInfo`）。
> 删除前必须使用 `a-popconfirm` 二次确认。**不要**使用原生的 `window.confirm` 或直接删除。

### 带详情按钮的操作列

如需在操作列中增加「查看详情」按钮：

```html
<template #action="{ row }">
  <a-button
    type="link"
    :icon="h(LucideInfo)"
    @click.stop="handleDetail(row)"
  />
  <a-button
    type="link"
    :icon="h(LucideFilePenLine)"
    @click.stop="handleEdit(row)"
  />
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

### 工具栏 Slot (toolbar-tools)

用于在表格标题右侧放置操作按钮（新建、同步、展开/折叠等）：

```html
<template #toolbar-tools>
  <a-button class="mr-2" type="primary" @click="handleCreate">
    {{ $t('page.xxx.button.create') }}
  </a-button>
</template>
```

带二次确认的操作按钮（如同步）：

```html
<template #toolbar-tools>
  <a-button type="primary" class="mr-2" @click="handleCreate">
    {{ $t('page.xxx.button.create') }}
  </a-button>
  <a-popconfirm
    :cancel-text="$t('ui.button.cancel')"
    :ok-text="$t('ui.button.ok')"
    :title="$t('ui.text.do_you_want_sync_permissions')"
    @confirm="handleSync"
  >
    <a-button type="primary" danger class="mr-2">
      {{ $t('page.xxx.button.sync') }}
    </a-button>
  </a-popconfirm>
</template>
```

树形表格的展开/折叠按钮：

```html
<template #toolbar-tools>
  <a-button class="mr-2" type="primary" @click="handleCreate">
    {{ $t('page.xxx.button.create') }}
  </a-button>
  <a-button class="mr-2" @click="expandAll">
    {{ $t('ui.tree.expand_all') }}
  </a-button>
  <a-button class="mr-2" @click="collapseAll">
    {{ $t('ui.tree.collapse_all') }}
  </a-button>
</template>
```

```typescript
const expandAll = () => {
  gridApi.grid?.setAllTreeExpand(true);
};
const collapseAll = () => {
  gridApi.grid?.setAllTreeExpand(false);
};
```

---

## VxeGrid 全局配置参考

### 基础列表页配置（带分页）

```typescript
const gridOptions: VxeGridProps<Xxx> = {
  toolbarConfig: {
    custom: true,     // 自定义列
    export: true,     // 导出
    refresh: true,    // 刷新
    zoom: true,       // 全屏
  },
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

### 不分页列表配置

```typescript
const gridOptions: VxeGridProps<Xxx> = {
  // ...
  pagerConfig: {
    enabled: false,   // 禁用分页
  },
  // ...
};
```

### 带 Tooltip 的配置

```typescript
const gridOptions: VxeGridProps<User> = {
  // ...
  tooltipConfig: {
    showAll: true,
    enterable: true,
    contentMethod: ({ column, row }) => {
      const { field } = column;
      if (field === 'roleNames') {
        return `${row[field]}`;
      }
      return null; // 其余列使用默认行为
    },
  },
  columns: [
    // ...
    {
      title: '角色',
      field: 'roleNames',
      slots: { default: 'role' },
      showOverflow: 'tooltip',  // 超长显示 tooltip
    },
  ],
};
```

### Grid 事件监听

```typescript
import { type VxeGridListeners } from '#/adapter/vxe-table';

const gridEvents: VxeGridListeners<User> = {
  cellDblclick: ({ row }) => {
    handleDetail(row);
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions,
  formOptions,
  gridEvents,     // 传入事件监听
});
```

---

## Drawer 表单组件 (useVbenForm)

Drawer 中的表单同样使用 VbenForm 体系，组件类型与搜索表单相同。

### 表单 Schema 常用组件

| component       | 用途    | 关键 props                                                                                     |
|-----------------|-------|----------------------------------------------------------------------------------------------|
| `Input`         | 文本输入  | `placeholder`, `allowClear`                                                                  |
| `InputNumber`   | 数字输入  | `placeholder`, `allowClear`, `defaultValue`                                                  |
| `Select`        | 下拉选择  | `options`, `showSearch`, `allowClear`, `filterOption`                                        |
| `ApiSelect`     | 远程下拉  | `api`, `afterFetch`, `showSearch`, `allowClear`, `filterOption`                              |
| `ApiTreeSelect` | 远程树选择 | `api`, `treeDefaultExpandAll`, `labelField`, `valueField`, `childrenField`, `numberToString` |
| `RadioGroup`    | 单选按钮组 | `optionType: 'button'`, `buttonStyle: 'solid'`, `options`                                    |
| `Textarea`      | 多行文本  | `placeholder`, `allowClear`                                                                  |
| `Switch`        | 开关    | —                                                                                            |
| `DatePicker`    | 日期选择  | `placeholder`, `allowClear`                                                                  |

### 表单校验规则

```typescript
// 必填（输入框类）
rules: 'required'

// 必选（下拉/单选类）
rules: 'selectRequired'
```

这两个规则已在 `adapter/form.ts` 中全局注册，会自动进行国际化提示（显示「请输入xxx」或「请选择xxx」）。**不要**自行写 `rules: z.string().min(1)` 之类的 Zod 校验。

### 表单默认值

通过 `defaultValue` 设置：

```typescript
{
  component: 'RadioGroup',
  fieldName: 'status',
  defaultValue: 'ON',
  // ...
},
{
  component: 'InputNumber',
  fieldName: 'sortOrder',
  defaultValue: 1,
  // ...
},
```

### 表单公共配置

```typescript
const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,   // 隐藏默认的提交/重置按钮
  commonConfig: {
    componentProps: {
      class: 'w-full',         // 所有表单项宽度 100%
    },
  },
  schema: [ /* ... */ ],
});
```

---

## Drawer 表单模板 (xxx-drawer.vue)

```vue
<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  PaginationQuery,
  xxxStatusList,
  useCreateXxx,
  useUpdateXxx,
} from '#/api';

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
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.xxx.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: $t('page.xxx.code'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: $t('ui.table.status'),
      defaultValue: 'ON',
      rules: 'selectRequired',
      componentProps: {
        optionType: 'button',
        buttonStyle: 'solid',
        class: 'flex flex-wrap',
        options: xxxStatusList,
      },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: $t('ui.table.remark'),
    },
  ],
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },

  async onConfirm() {
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);
    const values = await baseFormApi.getValues();

    try {
      await (data.value?.create
        ? createXxx(values)
        : updateXxx({ id: data.value.row.id, values }));

      notification.success({
        message: data.value?.create
          ? $t('ui.notification.create_success')
          : $t('ui.notification.update_success'),
      });
    } catch {
      notification.error({
        message: data.value?.create
          ? $t('ui.notification.create_failed')
          : $t('ui.notification.update_failed'),
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

---

## 完整列表页模板 (index.vue)

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

// ==============================
// 搜索表单
// ==============================
const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.xxx.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
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

// ==============================
// 表格配置
// ==============================
const gridOptions: VxeGridProps<Xxx> = {
  toolbarConfig: {
    custom: true,
    export: true,
    refresh: true,
    zoom: true,
  },
  exportConfig: {},
  pagerConfig: {},
  rowConfig: {
    isHover: true,
  },
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
    {
      title: $t('ui.table.status'),
      field: 'status',
      slots: { default: 'status' },
      width: 95,
    },
    {
      title: $t('ui.table.createdAt'),
      field: 'createdAt',
      formatter: 'formatDateTime',
      width: 140,
    },
    { title: $t('ui.table.remark'), field: 'remark' },
    {
      title: $t('ui.table.action'),
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      width: 90,
    },
  ],
};

// ==============================
// Grid + Drawer 初始化
// ==============================
const [Grid, gridApi] = useVbenVxeGrid({ gridOptions, formOptions });

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: XxxDrawer,
  onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      gridApi.reload();
    }
  },
});

function openDrawer(create: boolean, row?: any) {
  drawerApi.setData({ create, row });
  drawerApi.open();
}

function handleCreate() {
  openDrawer(true);
}

function handleEdit(row: any) {
  openDrawer(false, row);
}

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
        <a-button
          type="link"
          :icon="h(LucideFilePenLine)"
          @click.stop="handleEdit(row)"
        />
        <a-popconfirm
          :cancel-text="$t('ui.button.cancel')"
          :ok-text="$t('ui.button.ok')"
          :title="
            $t('ui.text.do_you_want_delete', {
              moduleName: $t('page.xxx.moduleName'),
            })
          "
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

---

## Vben 框架 API 速查

### Page 组件（必须包裹）

所有页面 **必须** 用 `<Page>` 组件包裹，它是 Vben 提供的页面容器：

```html
<Page auto-content-height>
  <Grid :table-title="...">
    <!-- ... -->
  </Grid>
  <Drawer />
</Page>
```

- `auto-content-height` — 自适应内容高度（推荐始终添加）
- 从 `@vben/common-ui` 导入：`import { Page } from '@vben/common-ui'`
- **禁止**用 `<div>` 或原生元素替代 `<Page>`

### Grid API 方法

`useVbenVxeGrid` 返回 `[Grid, gridApi]`，`gridApi` 提供以下方法：

| 方法                | 用途       | 常见场景                 |
|-------------------|----------|----------------------|
| `gridApi.reload()` | 重新加载数据   | 删除/编辑后刷新表格           |
| `gridApi.grid`    | VxeTable 实例 | 调用原生 VxeTable API   |

常用 grid 原生方法：
- `gridApi.grid?.setAllTreeExpand(true)` — 展开所有树节点
- `gridApi.grid?.setAllTreeExpand(false)` — 折叠所有树节点

### Drawer API 方法

`useVbenDrawer` 返回 `[Drawer, drawerApi]`，`drawerApi` 提供以下方法：

| 方法                           | 用途         | 说明                     |
|------------------------------|------------|------------------------|
| `drawerApi.open()`           | 打开抽屉       |                        |
| `drawerApi.close()`          | 关闭抽屉       |                        |
| `drawerApi.setData(data)`    | 传入数据       | 通常传 `{ create, row }`   |
| `drawerApi.getData<T>()`     | 获取传入数据     | 在 `onOpenChange` 中使用   |
| `drawerApi.setState({ loading })` | 设置加载状态 | 显示/隐藏提交按钮的 loading |

### Form API 方法

`useVbenForm` 返回 `[Form, formApi]`，`formApi` 提供以下方法：

| 方法                       | 用途     | 说明                      |
|--------------------------|--------|-------------------------|
| `formApi.getValues()`    | 获取表单值  | 返回 `Promise<Record>`    |
| `formApi.setValues(obj)` | 设置表单值  | 编辑时回填数据                 |
| `formApi.validate()`     | 校验表单   | 返回 `{ valid: boolean }` |
| `formApi.resetForm()`    | 重置表单   | 清空所有字段                  |

### Notification（消息提示）

使用 Ant Design Vue 的 `notification`，从 `ant-design-vue` 导入：

```typescript
import { notification } from 'ant-design-vue';

// 成功
notification.success({ message: $t('ui.notification.delete_success') });

// 失败
notification.error({ message: $t('ui.notification.delete_failed') });
```

> **禁止**使用 `alert()`、`console.log` 向用户展示操作结果，必须使用 `notification`。
> 项目中已预定义了通用 i18n key：`ui.notification.create_success`、`ui.notification.update_success`、
> `ui.notification.delete_success` 以及对应的 `_failed` 版本。

---

## 常见错误与纠正

| 错误做法                                                    | 正确做法                                                                                  |
|---------------------------------------------------------|---------------------------------------------------------------------------------------|
| 使用原生 `<span style="color:red">正常</span>` 渲染状态           | 使用 `<a-tag :color="statusToColor(row.status)">{{ statusToName(row.status) }}</a-tag>` |
| 使用 `window.confirm('确定删除?')` 确认删除                       | 使用 `<a-popconfirm>` 组件                                                                |
| 使用原生 `<button>` 或 `<button class="ant-btn)">`           | 使用 `<a-button type="primary/link">`                                                   |
| 日期格式化在 Slot 中用 `dayjs(row.createdAt).format(...)`       | 使用 `formatter: 'formatDateTime'` 列配置                                                  |
| 表单校验使用 `rules: z.string().min(1, '不能为空')`               | 使用 `rules: 'required'` 或 `rules: 'selectRequired'`                                    |
| 在 Slot 中用 `v-if` 判断状态值硬编码颜色                             | 使用 composables 层的 `xxxToColor()` 映射函数                                                 |
| 直接 `import { Tag } from 'ant-design-vue'` 然后用 `<Tag>`   | 直接在 template 中使用 `<a-tag>`（全局注册）                                                      |
| 用 `<div>` 包裹页面内容                                        | 必须用 `<Page auto-content-height>` 包裹                                                   |
| 手动管理 loading 状态 `ref(false)`                            | 使用 `drawerApi.setState({ loading })` 控制提交按钮 loading                                   |
| 表单提交后手动清空字段                                             | 关闭 Drawer 时框架自动重置，无需手动清空                                                              |
| 用 `alert()` 或 `ElMessage` 提示操作结果                        | 使用 `notification` from `ant-design-vue`                                               |
| schema 中 `component: 'AInput'` 或 `component: 'a-input'` | 使用注册名 `component: 'Input'`（见 adapter/component/index.ts）                              |
| 手动 watch 搜索条件变化重新请求                                     | 搜索表单与 Grid 自动关联，通过 `proxyConfig.ajax.query` 的 `formValues` 参数获取                       |
| 图标用 `:icon="LucideTrash2"` 直接传组件                        | 图标必须用 `h()` 包裹：`:icon="h(LucideTrash2)"`                                              |
| 列表数据用 `ref` + `watch` 手动管理                              | 使用 `proxyConfig.ajax.query` + `PaginationQuery`，由框架自动管理                               |
