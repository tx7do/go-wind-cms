<script setup lang="ts">
import { marked } from 'marked'
import DOMPurify from 'dompurify'

export type ContentType = 'html' | 'markdown' | 'text'

const props = withDefaults(defineProps<{
  content?: string
  type?: ContentType
  class?: string
}>(), {
  content: '',
  type: 'markdown',
})

const renderedContent = computed(() => {
  if (!props.content) return ''

  try {
    switch (props.type) {
      case 'markdown': {
        const html = marked.parse(props.content) as string
        return DOMPurify.sanitize(html, {
          ALLOWED_TAGS: [
            'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
            'p', 'br', 'hr', 'pre', 'code',
            'strong', 'b', 'em', 'i', 'u', 'del', 's',
            'a', 'img', 'blockquote', 'ul', 'ol', 'li',
            'table', 'thead', 'tbody', 'tfoot', 'tr', 'th', 'td',
            'video', 'iframe', 'figure', 'figcaption',
            'mark', 'sub', 'sup', 'span', 'div',
          ],
          ALLOWED_ATTR: [
            'href', 'title', 'target', 'rel',
            'src', 'alt', 'width', 'height',
            'class', 'id', 'style',
          ],
        })
      }
      case 'html':
        return DOMPurify.sanitize(props.content)
      case 'text': {
        const escaped = props.content
          .replace(/&/g, '&amp;')
          .replace(/</g, '&lt;')
          .replace(/>/g, '&gt;')
          .replace(/"/g, '&quot;')
          .replace(/'/g, '&#039;')
        return `<pre class="plain-text-block">${escaped}</pre>`
      }
      default:
        return props.content
    }
  } catch (error) {
    console.error('Error rendering content:', error)
    return '<p class="content-error">Failed to render content</p>'
  }
})
</script>

<template>
  <div
    class="content-viewer"
    :class="props.class"
    v-html="renderedContent"
  />
</template>

<style scoped>
.content-viewer {
  color: hsl(var(--foreground));
  line-height: 1.9;
  word-wrap: break-word;
  -webkit-font-smoothing: antialiased;
}

.content-viewer :deep(h1),
.content-viewer :deep(h2),
.content-viewer :deep(h3),
.content-viewer :deep(h4),
.content-viewer :deep(h5),
.content-viewer :deep(h6) {
  margin-top: 1.5em;
  margin-bottom: 0.75em;
  font-weight: 600;
  line-height: 1.35;
  color: hsl(var(--foreground));
}

.content-viewer :deep(h1:first-child),
.content-viewer :deep(h2:first-child) {
  margin-top: 0;
}

.content-viewer :deep(h1) {
  font-size: 2rem;
  border-bottom: 2px solid hsl(var(--border));
  padding-bottom: 0.5rem;
}

.content-viewer :deep(h2) {
  font-size: 1.625rem;
  border-bottom: 2px solid hsl(var(--border));
  padding-bottom: 0.5rem;
  margin: 48px 0 24px;
}

.content-viewer :deep(h3) {
  font-size: 1.375rem;
  margin: 36px 0 20px;
}

.content-viewer :deep(h4) {
  font-size: 1.125rem;
}

.content-viewer :deep(p) {
  margin: 22px 0;
  line-height: 2.0;
  letter-spacing: 0.25px;
  text-align: justify;
}

.content-viewer :deep(p:first-child) {
  margin-top: 0;
}

.content-viewer :deep(p:last-child) {
  margin-bottom: 0;
}

.content-viewer :deep(a) {
  color: hsl(var(--primary));
  text-decoration: none;
  border-bottom: 1.5px solid hsl(var(--primary) / 0.3);
  transition: all 0.25s ease;
  font-weight: 500;
}

.content-viewer :deep(a:hover) {
  text-decoration: none;
  border-bottom-color: hsl(var(--primary));
  background: hsl(var(--primary) / 0.08);
  padding: 2px 4px;
  margin: -2px -4px;
  border-radius: 4px;
}

.content-viewer :deep(strong),
.content-viewer :deep(b) {
  font-weight: 700;
  color: hsl(var(--primary));
}

.content-viewer :deep(em),
.content-viewer :deep(i) {
  font-style: italic;
  color: hsl(var(--muted-foreground));
}

.content-viewer :deep(code) {
  background: rgba(150, 150, 150, 0.12);
  padding: 3px 8px;
  border-radius: 5px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 0.88em;
  color: hsl(var(--primary));
  border: 1px solid hsl(var(--primary) / 0.18);
}

.content-viewer :deep(pre) {
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1.5em 0;
  background: rgba(150, 150, 150, 0.1);
  border: 1px solid hsl(var(--border));
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 0.875rem;
  line-height: 1.6;
}

.content-viewer :deep(pre code) {
  background: none;
  padding: 0;
  border-radius: 0;
  font-size: inherit;
  color: inherit;
  border: none;
}

.content-viewer :deep(blockquote) {
  border-left: 4px solid hsl(var(--primary));
  padding: 20px 24px;
  margin: 32px 0;
  background: hsl(var(--primary) / 0.06);
  color: hsl(var(--muted-foreground));
  border-radius: 0 8px 8px 0;
  font-style: italic;
}

.content-viewer :deep(ul),
.content-viewer :deep(ol) {
  padding-left: 2em;
  margin: 28px 0;
}

.content-viewer :deep(li) {
  margin: 12px 0;
  line-height: 1.95;
}

.content-viewer :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  display: inline-block;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.content-viewer :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 32px 0;
  font-size: 15px;
}

.content-viewer :deep(th),
.content-viewer :deep(td) {
  padding: 14px 16px;
  border: 1px solid hsl(var(--border));
  text-align: left;
}

.content-viewer :deep(th) {
  background: hsl(var(--primary));
  color: white;
  font-weight: 600;
}

.content-viewer :deep(hr) {
  border: none;
  border-top: 2px solid hsl(var(--border));
  margin: 48px 0;
}

.content-viewer :deep(.content-error) {
  color: #ff6b6b;
  padding: 1em;
  background: rgba(255, 107, 107, 0.1);
  border: 1px solid #ff6b6b;
  border-radius: 8px;
  margin: 1.5em 0;
}

.content-viewer :deep(.plain-text-block) {
  background: rgba(150, 150, 150, 0.1);
  color: hsl(var(--foreground));
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1.5em 0;
  border: 1px solid hsl(var(--border));
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 0.875rem;
  line-height: 1.6;
}
</style>
