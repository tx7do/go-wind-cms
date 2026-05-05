<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'

import {marked} from 'marked'
import hljs from 'highlight.js'
import DOMPurify from 'dompurify'
import katex from 'katex'
import mermaid from 'mermaid'

import type {Props} from './types'

// Initialize Mermaid
mermaid.initialize({
  startOnLoad: true,
  theme: 'default',
  securityLevel: 'loose',
})

const containerRef = ref<HTMLDivElement>()
const mermaidReady = ref(false)

const props = withDefaults(defineProps<Props>(), {
  content: '',
  type: 'markdown',
  class: '',
})

// Configure marked with custom renderer for better code handling
marked.setOptions({
  breaks: true,
  gfm: true,
  pedantic: false,
})

// Custom renderer for better code block handling
const renderer = new marked.Renderer()

// Override code block rendering
renderer.codespan = (code) => {
  return `<code class="inline-code">${escapeHtml(code.text)}</code>`
}

renderer.code = (code) => {
  const lang = code.lang || 'plaintext'
  let highlighted: string

  // Check if it's a mermaid diagram
  if (lang === 'mermaid') {
    try {
      return `<div class="mermaid">${escapeHtml(code.text)}</div>`
    } catch (error) {
      console.warn('Failed to render mermaid diagram:', error)
      return `<pre class="code-block mermaid-error" data-lang="mermaid"><code>${escapeHtml(code.text)}</code></pre>`
    }
  }

  // Try to highlight the code
  try {
    if (hljs.getLanguage(lang)) {
      highlighted = hljs.highlight(code.text, {language: lang, ignoreIllegals: true}).value
    } else {
      highlighted = escapeHtml(code.text)
    }
  } catch (error) {
    console.warn(`Failed to highlight code with language ${lang}:`, error)
    highlighted = escapeHtml(code.text)
  }

  return `<pre class="code-block" data-lang="${lang}"><code class="hljs language-${lang}">${highlighted}</code></pre>`
}

// Override heading rendering - 不添加 id，让父组件处理
renderer.heading = (heading) => {
  const inlineHtml = marked.parseInline(heading.text)
  // 不添加 id，由父组件的 generateTableOfContents 函数动态添加
  return `<h${heading.depth} class="heading-anchor">${inlineHtml}</h${heading.depth}>`
}

function splitUrlAndText(content: string): string {
  // 匹配 URL 后跟中文逗号和描述
  return content.replace(/(https?:\/\/[^\s，]+)(，[^ \n]+)/g, (_match, url, desc) => {
    return `[${url}](${url})${desc}`
  })
}

// Override link rendering
renderer.link = (link) => {
  const isExternal = link.href.startsWith('http') || link.href.startsWith('//')
  // 只渲染链接部分，后续文本不包裹
  if (link.href === link.text) {
    return `<a href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.text)}</a>`
  } else {
    // 链接文本和 href 不一致，分开渲染
    return `<a href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.href)}</a>${escapeHtml(link.text.replace(link.href, ''))}`
  }
}

// Override image rendering
renderer.image = (image) => {
  return `<figure class="markdown-image">
    <img src="${image.href}" alt="${image.text}" class="md-img" />
    ${image.text ? `<figcaption>${image.text}</figcaption>` : ''}
  </figure>`
}

// Override table rendering
renderer.table = (token) => {
  let html = '<thead>\n<tr>\n'

  // Render table header
  for (const header of token.header) {
    const align = header.align ? ` style="text-align:${header.align}"` : ''
    const cellHtml = marked.parseInline(header.text)
    html += `<th${align}>${cellHtml}</th>\n`
  }

  html += '</tr>\n</thead>\n<tbody>\n'

  // Render table rows
  for (const row of token.rows) {
    html += '<tr>\n'
    for (const cell of row) {
      const align = cell.align ? ` style="text-align:${cell.align}"` : ''
      const cellHtml = marked.parseInline(cell.text)
      html += `<td${align}>${cellHtml}</td>\n`
    }
    html += '</tr>\n'
  }

  html += '</tbody>'
  return `<div class="table-wrapper"><table class="markdown-table">${html}</table></div>`
}

// Override paragraph rendering
renderer.paragraph = (token) => {
  return `<p>${marked.parseInline(token.text)}</p>\n`
}


// Set custom renderer
marked.setOptions({renderer})

// Render content based on type
const renderedContent = computed(() => {
  if (!props.content) return ''

  try {
    let html = ''
    switch (props.type) {
      case 'markdown':
        let md = props.content
        md = splitUrlAndText(md)
        html = marked.parse(md) as string
        // Process math formulas
        html = processMathFormulas(html)
        return sanitizeHtml(html)
      case 'html':
        return sanitizeHtml(props.content)
      case 'text':
        return sanitizeHtml(`<pre class="plain-text-block">${escapeHtml(props.content)}</pre>`)
      default:
        return props.content
    }
  } catch (error) {
    console.error('Error rendering content:', error)
    return `<p class="content-error">Failed to render content</p>`
  }
})

// Process math formulas (both inline and block)
function processMathFormulas(html: string): string {
  // Block math: $$...$$
  html = html.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
    try {
      const rendered = katex.renderToString(formula, {displayMode: true})
      return `<div class="math-block">${rendered}</div>`
    } catch (error) {
      console.warn('Failed to render block math:', error)
      return `<div class="math-error"><code>${escapeHtml(formula)}</code></div>`
    }
  })

  // Inline math: $...$
  html = html.replace(/(?<!\$)\$([^$\n]+)\$(?!\$)/g, (_, formula) => {
    try {
      const rendered = katex.renderToString(formula, {displayMode: false})
      return `<span class="math-inline">${rendered}</span>`
    } catch (error) {
      console.warn('Failed to render inline math:', error)
      return `<code class="math-error">${escapeHtml(formula)}</code>`
    }
  })

  return html
}

// Escape HTML for text type (preserve emoji)
function sanitizeHtml(html: string): string {
  // Configure DOMPurify to allow emoji and unicode
  const config: any = {
    ALLOWED_TAGS: [
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'p', 'br', 'hr', 'pre', 'code',
      'strong', 'b', 'em', 'i', 'u', 'del', 's',
      'a', 'img', 'blockquote', 'ul', 'ol', 'li',
      'table', 'thead', 'tbody', 'tfoot', 'tr', 'th', 'td',
      'video', 'iframe', 'figure', 'figcaption',
      'mark', 'sub', 'sup', 'span', 'div',
      // Math formula tags
      'math', 'mrow', 'mi', 'mn', 'mo', 'mfrac', 'msup', 'msub', 'mover', 'munder',
      // Mermaid diagram
      'svg', 'g', 'text', 'path', 'circle', 'line', 'polyline', 'polygon', 'rect', 'ellipse',
    ],
    ALLOWED_ATTR: [
      'href', 'title', 'target', 'rel',
      'src', 'alt', 'width', 'height',
      'class', 'id', 'style', 'data-lang',
      'data-*',
      // SVG attributes
      'viewBox', 'xmlns', 'x', 'y', 'cx', 'cy', 'r', 'x1', 'y1', 'x2', 'y2',
      'd', 'fill', 'stroke', 'stroke-width', 'points', 'transform',
    ],
    KEEP_CONTENT: true,
    // Don't strip whitespace or special characters that might be emoji
    ALLOW_UNKNOWN_PROTOCOLS: true,
    // Use HTML parser instead of DOM parser to preserve text nodes
    RETURN_DOM: false,
    RETURN_DOM_FRAGMENT: false,
    RETURN_DOM_IMPORT: false,
    FORCE_BODY: false,
    SANITIZE_DOM: true,
    IN_PLACE: false,
  }

  // Use DOMPurify's default sanitization with our config
  return DOMPurify.sanitize(html, config) as unknown as string
}

// Escape HTML for text type (preserve emoji)
function escapeHtml(text: string): string {
  // Create a temporary element and use textContent to properly escape HTML
  // while preserving emoji and unicode characters
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  }

  // Only escape HTML special characters, NOT unicode/emoji
  return text.replace(/[&<>"']/g, (m) => map[m])
}

// Initialize mermaid diagrams after content is rendered
onMounted(() => {
  if (containerRef.value) {
    try {
      mermaid.contentLoaded()
      mermaidReady.value = true
    } catch (error) {
      console.warn('Failed to initialize mermaid:', error)
    }
  }
})
</script>

<template>
  <div ref="containerRef" :class="['content-viewer', props.class]" v-html="renderedContent"/>
</template>

<style scoped lang="less">
.content-viewer {
  color: var(--color-text-primary);
  line-height: 1.9;
  word-wrap: break-word;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;

  // Typography - Headings
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin-top: 1.5em;
    margin-bottom: 0.75em;
    font-weight: 600;
    line-height: 1.35;
    color: var(--color-text-primary);

    &:first-child {
      margin-top: 0;
    }
  }

  :deep(h1) {
    font-size: 2rem;
    border-bottom: 2px solid var(--color-border);
    padding-bottom: 0.5rem;
  }

  :deep(h2) {
    font-size: 1.625rem;
    border-bottom: 2px solid var(--color-border);
    padding-bottom: 0.5rem;
    margin: 48px 0 24px;
  }

  :deep(h3) {
    font-size: 1.375rem;
    margin: 36px 0 20px;
  }

  :deep(h4) {
    font-size: 1.125rem;
  }

  :deep(h5), :deep(h6) {
    font-size: 1rem;
  }

  // Paragraphs
  :deep(p) {
    margin: 22px 0;
    line-height: 2.0;
    letter-spacing: 0.25px;
    word-spacing: 0.05em;
    text-align: justify;

    &:first-child {
      margin-top: 0;
    }

    &:last-child {
      margin-bottom: 0;
    }

    & + p {
      margin-top: 20px;
    }
  }

  // Links
  :deep(a.markdown-link) {
    color: var(--color-brand);
    text-decoration: none;
    border-bottom: 1.5px solid rgba(168, 85, 247, 0.3);
    transition: all 0.25s ease;
    font-weight: 500;

    &:hover {
      text-decoration: none;
      border-bottom-color: var(--color-brand);
      background: rgba(168, 85, 247, 0.08);
      padding: 2px 4px;
      margin: -2px -4px;
      border-radius: 4px;
    }
  }

  // Strong and emphasis
  :deep(strong), :deep(b) {
    font-weight: 700;
    color: var(--color-brand);
    letter-spacing: 0.2px;
  }

  :deep(em), :deep(i) {
    font-style: italic;
    color: var(--color-text-secondary);
    letter-spacing: 0.15px;
  }

  // Inline code
  :deep(code.inline-code) {
    background: rgba(150, 150, 150, 0.12);
    padding: 3px 8px;
    border-radius: 5px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    font-size: 0.88em;
    color: var(--color-brand);
    border: 1px solid rgba(168, 85, 247, 0.18);
    letter-spacing: 0.1px;
    white-space: nowrap;
  }

  // Code blocks with syntax highlighting
  :deep(pre.code-block) {
    background: #282c34;
    color: #abb2bf;
    padding: 20px;
    border-radius: 10px;
    overflow-x: auto;
    margin: 32px 0;
    border: 1px solid rgba(255, 255, 255, 0.1);
    position: relative;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    line-height: 1.65;
    letter-spacing: 0.15px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    white-space: pre;
    word-wrap: normal;
    overflow-wrap: normal;

    &::before {
      content: attr(data-lang);
      position: absolute;
      top: 8px;
      right: 12px;
      font-size: 12px;
      color: rgba(255, 255, 255, 0.4);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    code {
      background: none;
      padding: 0;
      border-radius: 0;
      font-size: 13.5px;
      line-height: 1.65;
      color: inherit;
      font-family: inherit;
      white-space: inherit;
      word-wrap: inherit;
      overflow-wrap: inherit;

      // Syntax highlighting colors (One Dark theme)
      .hljs-attr {
        color: #e06c75;
      }

      .hljs-attr-value {
        color: #98c379;
      }

      .hljs-number {
        color: #d19a66;
      }

      .hljs-literal {
        color: #56b6c2;
      }

      .hljs-string {
        color: #98c379;
      }

      .hljs-string .hljs-emphasis {
        color: #e06c75;
      }

      .hljs-section {
        color: #e06c75;
      }

      .hljs-link {
        color: #d19a66;
      }

      .hljs-operator {
        color: #56b6c2;
      }

      .hljs-keyword {
        color: #c678dd;
      }

      .hljs-type {
        color: #56b6c2;
      }

      .hljs-name {
        color: #e06c75;
      }

      .hljs-selector-class {
        color: #e06c75;
      }

      .hljs-selector-id {
        color: #e06c75;
      }

      .hljs-variable {
        color: #e06c75;
      }

      .hljs-template-variable {
        color: #e06c75;
      }

      .hljs-title {
        color: #61afef;
      }

      .hljs-title.class_ {
        color: #e5c07b;
      }

      .hljs-title.class_.inherited__ {
        color: #98c379;
      }

      .hljs-title.function_ {
        color: #61afef;
      }

      .hljs-params {
        color: #abb2bf;
      }

      .hljs-class {
        color: #e5c07b;
      }

      .hljs-function {
        color: #61afef;
      }

      .hljs-comment {
        color: #5c6370;
      }

      .hljs-doctag {
        color: #5c6370;
      }

      .hljs-meta {
        color: #5c6370;
      }

      .hljs-meta .hljs-keyword {
        color: #5c6370;
      }

      .hljs-meta .hljs-string {
        color: #98c379;
      }

      .hljs-builtin-name {
        color: #e5c07b;
      }

      .hljs-code {
        color: #98c379;
      }

      .hljs-tag {
        color: #e06c75;
      }

      .hljs-tag .hljs-attr {
        color: #e5c07b;
      }

      .hljs-tag .hljs-title {
        color: #e06c75;
      }

      .hljs-tag .hljs-literal {
        color: #abb2bf;
      }

      .hljs-name {
        color: #e06c75;
      }

      .hljs-attr {
        color: #e5c07b;
      }

      .hljs-bullet {
        color: #56b6c2;
      }

      .hljs-quote {
        color: #5c6370;
      }

      .hljs-punctuation {
        color: #abb2bf;
      }
    }
  }

  // Plain text blocks
  :deep(pre.plain-text-block) {
    background: rgba(150, 150, 150, 0.1);
    color: var(--color-text-primary);
    padding: 1rem;
    border-radius: 8px;
    overflow-x: auto;
    margin: 1.5em 0;
    border: 1px solid var(--color-border);
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    font-size: 0.875rem;
    line-height: 1.6;
  }

  // Blockquotes
  :deep(blockquote) {
    border-left: 4px solid var(--color-brand);
    padding: 20px 24px;
    margin: 32px 0;
    background: rgba(168, 85, 247, 0.06);
    color: var(--color-text-secondary);
    border-radius: 0 8px 8px 0;
    font-style: italic;
    position: relative;
    letter-spacing: 0.2px;
    line-height: 1.85;

    &::before {
      content: '"';
      position: absolute;
      top: -10px;
      left: 12px;
      font-size: 56px;
      color: var(--color-brand);
      opacity: 0.2;
      font-family: Georgia, serif;
      line-height: 1;
    }

    p {
      margin: 14px 0;
      text-indent: 0;
      text-align: left;

      &:first-child {
        margin-top: 0;
      }

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  // Lists
  :deep(ul), :deep(ol) {
    padding-left: 2em;
    margin: 28px 0;

    li {
      margin: 12px 0;
      line-height: 1.95;
      letter-spacing: 0.2px;
      position: relative;

      p {
        margin: 10px 0;
      }
    }
  }

  :deep(ul) {
    list-style-type: disc;

    li::marker {
      color: var(--color-brand);
      font-weight: 600;
    }

    ul {
      list-style-type: circle;

      ul {
        list-style-type: square;
      }
    }
  }

  :deep(ol) {
    list-style-type: decimal;

    li::marker {
      color: var(--color-brand);
      font-weight: 600;
    }
  }

  // Images
  :deep(figure.markdown-image) {
    margin: 1.5em 0;
    text-align: center;

    :deep(img.md-img) {
      max-width: 100%;
      height: auto;
      border-radius: 8px;
      display: inline-block;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    :deep(figcaption) {
      margin-top: 0.5em;
      font-size: 0.875rem;
      color: var(--color-text-secondary);
      text-align: center;
    }
  }

  // Tables
  :deep(.table-wrapper) {
    overflow-x: auto;
    margin: 32px 0;
    border-radius: 8px;
    border: 1px solid var(--color-border);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  }

  :deep(table.markdown-table) {
    width: 100%;
    border-collapse: collapse;
    background: var(--color-surface);
    font-size: 15px;
    letter-spacing: 0.15px;

    th, td {
      padding: 14px 16px;
      border: 1px solid var(--color-border);
      text-align: left;
    }

    th {
      background: var(--color-brand);
      color: white;
      font-weight: 600;
      line-height: 1.5;
      letter-spacing: 0.3px;
    }

    td {
      line-height: 1.7;
    }

    tr:nth-child(even) {
      background: rgba(150, 150, 150, 0.05);
    }

    tr:hover {
      background: rgba(168, 85, 247, 0.06);
    }
  }

  // Horizontal rule
  :deep(hr) {
    border: none;
    border-top: 2px solid var(--color-border);
    margin: 48px 0;
  }

  // Video
  :deep(video) {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
    margin: 1.5em 0;
    display: block;
  }

  // Iframe (for embedded content)
  :deep(iframe) {
    max-width: 100%;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    margin: 1.5em 0;
  }

  // Error state
  :deep(.content-error) {
    color: #ff6b6b;
    padding: 1em;
    background: rgba(255, 107, 107, 0.1);
    border: 1px solid #ff6b6b;
    border-radius: 8px;
    margin: 1.5em 0;
  }

  // Math formulas
  :deep(.math-block) {
    display: flex;
    justify-content: center;
    margin: 1.5em 0;
    padding: 1em;
    background: rgba(150, 150, 150, 0.05);
    border-radius: 8px;
    overflow-x: auto;
  }

  :deep(.math-inline) {
    display: inline;
    font-size: 0.95em;
  }

  :deep(.math-error) {
    color: #ff6b6b;
    background: rgba(255, 107, 107, 0.1);
    padding: 0.25em 0.5em;
    border-radius: 4px;
  }

  // Mermaid diagrams
  :deep(.mermaid) {
    display: flex;
    justify-content: center;
    margin: 1.5em 0;
    padding: 1em;
    background: rgba(150, 150, 150, 0.05);
    border-radius: 8px;
    border: 1px solid var(--color-border);
    overflow-x: auto;

    svg {
      max-width: 100%;
      height: auto;
    }
  }

  :deep(.mermaid-error) {
    color: #ff6b6b;
    background: rgba(255, 107, 107, 0.1);
    border-color: #ff6b6b;
  }

  // KaTeX CSS (imported from katex package)
  :deep(.katex) {
    font: normal 1.21em KaTeX_Main, 'Times New Roman', serif;
    line-height: 1.2;
    text-indent: 0;
  }

  :deep(.katex-display) {
    display: block;
    margin: 1em 0;
    text-align: center;
  }

  :deep(.katex-html) {
    display: none;
  }

  :deep(.katex-mathml) {
    position: absolute;
    clip: rect(1px, 1px, 1px, 1px);
    padding: 0;
    border: 0;
    height: 1px;
    width: 1px;
    overflow: hidden;
  }
}
</style>

