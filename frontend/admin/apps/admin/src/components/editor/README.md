# 编辑器组件使用说明

## 概述

已成功整合三种编辑器：
1. **UEditor** - 富文本编辑器（使用 vue-ueditor-wrap）
2. **Markdown Editor** - Markdown 编辑器（临时使用 textarea，可升级为 @toast-ui/editor）
3. **JSON Editor** - JSON 编辑器（临时使用 textarea + 验证，可升级为 jsoneditor）

## 文件结构

```
apps/admin/src/components/editor/
├── editor.vue          # 主编辑器组件（根据类型动态加载）
├── UEditor.vue         # UEditor 富文本编辑器
├── MarkdownEditor.vue  # Markdown 编辑器
├── JsonEditor.vue      # JSON 编辑器
├── types.ts            # TypeScript 类型定义
└── index.ts            # 导出所有组件
```

## 使用方法

### 1. 基本使用

```vue
<script setup lang="ts">
import { ref } from 'vue';
import { Editor, EditorType } from '#/components/editor';

const content = ref('');
const editorType = ref(EditorType.MARKDOWN);
</script>

<template>
  <Editor
    v-model="content"
    :editor-type="editorType"
    :height="600"
    placeholder="Enter your content..."
  />
</template>
```

### 2. 编辑器类型

```typescript
export enum EditorType {
  MARKDOWN = 'EDITOR_TYPE_MARKDOWN',        // Markdown 编辑器
  RICH_TEXT = 'EDITOR_TYPE_RICH_TEXT',      // UEditor 富文本
  HTML = 'EDITOR_TYPE_HTML',                // UEditor（HTML模式）
  JSON = 'EDITOR_TYPE_JSON_BLOCK',          // JSON 编辑器
  WYSIWYG = 'EDITOR_TYPE_WYSIWYG',          // UEditor（所见即所得）
  // 以下类型暂时映射到 Markdown 编辑器
  PLAIN_TEXT = 'EDITOR_TYPE_PLAIN_TEXT',
  CODE = 'EDITOR_TYPE_CODE',
}
```

### 3. 动态切换编辑器

```vue
<template>
  <!-- 编辑器类型选择 -->
  <a-select v-model:value="editorType">
    <a-select-option :value="EditorType.MARKDOWN">
      Markdown Editor
    </a-select-option>
    <a-select-option :value="EditorType.RICH_TEXT">
      Rich Text Editor (UEditor)
    </a-select-option>
    <a-select-option :value="EditorType.JSON">
      JSON Editor
    </a-select-option>
  </a-select>

  <!-- 编辑器 -->
  <Editor
    v-model="content"
    :editor-type="editorType"
    :height="600"
  />
</template>
```

### 4. 单独使用特定编辑器

```vue
<script setup lang="ts">
import { UEditor, MarkdownEditor, JsonEditor } from '#/components/editor';
</script>

<template>
  <!-- 直接使用 UEditor -->
  <UEditor
    v-model="content"
    :height="500"
    :config="{ /* UEditor 配置 */ }"
  />

  <!-- 直接使用 Markdown 编辑器 -->
  <MarkdownEditor
    v-model="markdown"
    :height="500"
  />

  <!-- 直接使用 JSON 编辑器 -->
  <JsonEditor
    v-model="jsonData"
    :height="500"
  />
</template>
```

## Props 说明

### Editor 组件

| Prop | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | string | - | 编辑器内容（v-model） |
| editorType | EditorType \| string | MARKDOWN | 编辑器类型 |
| height | string \| number | 500 | 编辑器高度 |
| disabled | boolean | false | 是否禁用 |
| placeholder | string | 'Please enter content...' | 占位符文本 |
| ueditorConfig | object | {} | UEditor 特定配置 |
| markdownOptions | object | {} | Markdown 编辑器配置 |
| jsonOptions | object | {} | JSON 编辑器配置 |

### UEditor 特定配置

```typescript
{
  UEDITOR_HOME_URL: '/UEditor/',
  serverUrl: '/api/ueditor/controller',
  initialFrameHeight: 500,
  readonly: false,
  // ... 更多 UEditor 配置
}
```

### Markdown 编辑器配置

```typescript
{
  initialEditType: 'markdown' | 'wysiwyg',
  previewStyle: 'tab' | 'vertical',
  hideModeSwitch: boolean,
  toolbarItems: string[][]
}
```

### JSON 编辑器配置

```typescript
{
  mode: 'tree' | 'code' | 'form' | 'text' | 'view',
  modes: ['tree', 'code', 'form', 'text', 'view'],
  search: boolean
}
```

## Events

| Event | 参数 | 说明 |
|-------|------|------|
| update:modelValue | value: string | 内容更新 |
| change | value: string | 内容变化 |
| ready | - | 编辑器加载完成 |
| error | error: Error | 错误（仅 JSON 编辑器） |

## 已知限制

### 当前版本（临时实现）

1. **Markdown Editor**: 使用简单的 textarea，缺少语法高亮和预览功能
2. **JSON Editor**: 使用简单的 textarea + JSON 验证，缺少树形视图

### 升级建议

如需完整功能，请安装以下依赖：

```bash
# Markdown 编辑器
pnpm add @toast-ui/editor @toast-ui/vue-editor

# JSON 编辑器
pnpm add jsoneditor vue-jsoneditor
```

然后更新对应的编辑器组件文件即可。

## 在 Post 编辑页面中的使用

已经在 `/content/posts/create` 和 `/content/posts/edit/:id` 页面中集成了编辑器：

```vue
<!-- 编辑器类型选择 -->
<a-form-item label="Editor Type" name="editorType">
  <a-select v-model:value="formData.editorType">
    <a-select-option :value="EditorType.MARKDOWN">
      Markdown Editor
    </a-select-option>
    <a-select-option :value="EditorType.RICH_TEXT">
      Rich Text Editor (UEditor)
    </a-select-option>
    <a-select-option :value="EditorType.JSON">
      JSON Editor
    </a-select-option>
  </a-select>
</a-form-item>

<!-- 内容编辑器 -->
<a-form-item label="Content" name="content" required>
  <Editor
    v-model="formData.content"
    :editor-type="formData.editorType"
    :height="600"
    placeholder="Enter your content here..."
  />
</a-form-item>
```

## UEditor 后端配置

UEditor 需要后端支持图片上传等功能。请确保：

1. 将 UEditor 资源文件放在 `public/UEditor/` 目录
2. 配置后端接口 `/api/ueditor/controller` 处理上传请求
3. 参考 UEditor 官方文档配置服务端

## 注意事项

1. **UEditor 依赖**: 项目已安装 `vue-ueditor-wrap`，可直接使用
2. **资源路径**: UEditor 静态资源需要放在 public 目录
3. **类型安全**: 使用 TypeScript 类型定义确保类型安全
4. **懒加载**: 编辑器组件采用懒加载方式，提高性能
5. **样式隔离**: 使用 scoped 样式避免样式冲突

## 扩展开发

### 添加新的编辑器

1. 在 `components/editor/` 下创建新编辑器组件
2. 在 `types.ts` 中添加新的编辑器类型
3. 在 `editor.vue` 的 switch 语句中添加新的 case
4. 在 `index.ts` 中导出新组件

### 示例：添加代码编辑器

```vue
<!-- CodeEditor.vue -->
<script setup lang="ts">
import { ref } from 'vue';
// import Monaco Editor or CodeMirror
</script>

<template>
  <div class="code-editor">
    <!-- 编辑器实现 -->
  </div>
</template>
```

然后在 `editor.vue` 中添加：

```typescript
case EditorType.CODE: {
  return CodeEditor;
}
```

## 总结

编辑器组件已经成功集成到项目中，可以在 Post 编辑页面中使用。当前版本使用简化实现，可根据需求安装完整的编辑器库进行升级。

