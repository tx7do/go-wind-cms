# ContentEditor Component

A versatile content editor component that supports multiple editor types: Markdown, Rich Text, and HTML.

## Features

- **Multiple Editor Types**: Switch between Markdown, Rich Text, and HTML modes
- **Markdown Toolbar**: Quick access to common Markdown formatting tools
- **Live Preview**: Preview Markdown content in real-time (tab-based)
- **Customizable Height**: Set custom editor height
- **i18n Support**: Fully internationalized with Chinese and English support
- **Type-safe**: Written in TypeScript with full type definitions

## Usage

### Basic Usage

```vue
<script setup lang="ts">
import { ref } from 'vue'
import { ContentEditor } from '@/components/ContentEditor'

const content = ref('')

function handleChange(value: string) {
  console.log('Content changed:', value)
}
</script>

<template>
  <ContentEditor
    v-model="content"
    type="markdown"
    @change="handleChange"
  />
</template>
```

### With Custom Configuration

```vue
<template>
  <ContentEditor
    v-model="content"
    type="markdown"
    height="600px"
    :show-toolbar="true"
    :show-type-switch="true"
    :disabled="false"
    placeholder="Enter your content here..."
    @change="handleContentChange"
    @type-change="handleTypeChange"
  />
</template>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | `string` | `''` | The content value (v-model) |
| `type` | `'markdown' \| 'richtext' \| 'html'` | `'markdown'` | Editor type |
| `placeholder` | `string` | `''` | Placeholder text |
| `disabled` | `boolean` | `false` | Disable editor |
| `height` | `string \| number` | `'400px'` | Editor height |
| `showToolbar` | `boolean` | `true` | Show formatting toolbar |
| `showTypeSwitch` | `boolean` | `true` | Show type switcher |

## Events

| Event | Parameters | Description |
|-------|------------|-------------|
| `update:modelValue` | `(value: string)` | Emitted when content changes |
| `change` | `(value: string)` | Emitted when content changes |
| `typeChange` | `(type: EditorType)` | Emitted when editor type changes |

## Editor Types

### Markdown

The Markdown editor includes a toolbar with the following formatting options:

- **Bold** (`**text**`)
- *Italic* (`*text*`)
- Heading (`## text`)
- Link (`[text](url)`)
- Image (`![alt](url)`)
- Inline Code (`` `code` ``)
- Code Block (`` ```code``` ``)
- Quote (`> text`)
- Unordered List (`- item`)
- Ordered List (`1. item`)
- Horizontal Rule (`---`)

The Markdown editor also supports a preview mode that renders the Markdown content as HTML.

### Rich Text

The Rich Text editor provides a simple textarea interface. In future versions, this can be enhanced with a WYSIWYG editor like Quill or TipTap.

### HTML

The HTML editor provides a plain textarea for editing raw HTML content.

## Styling

The component uses CSS variables for theming and automatically adapts to light/dark mode:

- `--color-surface`: Background color
- `--color-bg`: Secondary background color
- `--color-border`: Border color
- `--color-text-primary`: Primary text color
- `--color-text-secondary`: Secondary text color
- `--color-brand`: Brand/accent color
- `--radius-md`: Border radius

## Demo

Visit `/editor-demo` to see a live demo of the ContentEditor component.

## Future Enhancements

- [ ] Integration with advanced Rich Text editors (Quill, TipTap)
- [ ] Image upload support
- [ ] Syntax highlighting for code blocks
- [ ] Full-screen mode
- [ ] Auto-save functionality
- [ ] Emoji picker
- [ ] Table support
- [ ] Export to various formats (PDF, DOCX, etc.)

## Dependencies

- Vue 3
- Naive UI
- vue-i18n

## License

MIT

