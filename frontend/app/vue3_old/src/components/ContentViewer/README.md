# ContentViewer Component

A versatile content viewer component for rendering Markdown, HTML, and plain text content with beautiful styling.

## Features

- **Multiple Content Types**: Supports Markdown, HTML, and plain text
- **Markdown Rendering**: Uses `marked` library for robust Markdown parsing
- **Beautiful Typography**: Professional styling for all content elements
- **Responsive Images**: Automatic image sizing and styling
- **Code Highlighting**: Styled code blocks and inline code
- **Table Support**: Fully styled tables with alternating rows
- **GFM Support**: GitHub Flavored Markdown including task lists
- **Theme Compatible**: Automatically adapts to light/dark themes
- **Safe HTML**: Sanitizes and renders HTML content safely

## Usage

### Basic Usage

```vue
<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'

const content = `
# Hello World

This is a **Markdown** content viewer.

- Item 1
- Item 2
- Item 3
`
</script>

<template>
  <ContentViewer :content="content" type="markdown" />
</template>
```

### With HTML Content

```vue
<template>
  <ContentViewer 
    :content="htmlContent" 
    type="html" 
  />
</template>

<script setup lang="ts">
const htmlContent = '<h1>Hello</h1><p>This is <strong>HTML</strong> content.</p>'
</script>
```

### With Plain Text

```vue
<template>
  <ContentViewer 
    :content="textContent" 
    type="text" 
  />
</template>

<script setup lang="ts">
const textContent = 'This is plain text content.\nIt will be displayed as-is.'
</script>
```

### With Custom Styling

```vue
<template>
  <ContentViewer 
    :content="content" 
    type="markdown"
    class="custom-viewer"
  />
</template>

<style>
.custom-viewer {
  padding: 2rem;
  background: var(--color-surface);
  border-radius: 8px;
}
</style>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `content` | `string` | `''` | The content to display |
| `type` | `'markdown' \| 'html' \| 'text'` | `'markdown'` | Content type |
| `class` | `string` | `''` | Custom CSS class |

## Content Types

### Markdown

Supports full GitHub Flavored Markdown (GFM) including:

- **Headings** (H1-H6)
- **Text formatting** (bold, italic, strikethrough)
- **Links** and **images**
- **Code blocks** with language support
- **Inline code**
- **Blockquotes**
- **Lists** (ordered and unordered)
- **Tables**
- **Task lists** (- [ ] and - [x])
- **Horizontal rules**
- **Line breaks**

### HTML

Renders HTML content with full styling support. All standard HTML elements are styled consistently.

### Text

Displays plain text in a preformatted block, preserving whitespace and line breaks.

## Styled Elements

The ContentViewer includes beautiful styling for:

- **Typography**: Headings, paragraphs, text formatting
- **Code**: Inline code and code blocks with syntax highlighting support
- **Quotes**: Styled blockquotes with accent border
- **Lists**: Bulleted and numbered lists with proper indentation
- **Tables**: Responsive tables with alternating row colors
- **Images**: Responsive images with rounded corners and shadows
- **Links**: Styled links with hover effects
- **Horizontal Rules**: Styled dividers

## Styling

The component uses CSS variables for consistent theming:

- `--color-text-primary`: Primary text color
- `--color-text-secondary`: Secondary text color
- `--color-brand`: Accent/brand color
- `--color-border`: Border color
- `--color-surface`: Surface background color

These variables automatically adapt to your application's theme (light/dark mode).

## Example: Blog Post Viewer

```vue
<template>
  <div class="blog-post">
    <h1>{{ post.title }}</h1>
    <div class="post-meta">
      <span>{{ post.author }}</span>
      <span>{{ formatDate(post.createdAt) }}</span>
    </div>
    <ContentViewer 
      :content="post.content" 
      type="markdown"
      class="post-content"
    />
  </div>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'
import type { Post } from '@/typings'

interface Props {
  post: Post
}

const props = defineProps<Props>()

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}
</script>

<style scoped>
.blog-post {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.post-meta {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.post-content {
  margin-top: 2rem;
}
</style>
```

## Example: Comment Viewer

```vue
<template>
  <div class="comment">
    <div class="comment-header">
      <strong>{{ comment.author }}</strong>
      <span>{{ timeAgo(comment.createdAt) }}</span>
    </div>
    <ContentViewer 
      :content="comment.content" 
      type="markdown"
      class="comment-content"
    />
  </div>
</template>
```

## Security

- HTML content is rendered using Vue's `v-html` directive
- Plain text content is automatically escaped to prevent XSS attacks
- For user-generated content, consider additional sanitization

## Dependencies

- `marked` - Markdown parsing and rendering
- Vue 3 - Component framework

## Browser Support

Works in all modern browsers that support:
- ES6+
- CSS Custom Properties (CSS Variables)
- Flexbox

## License

MIT

