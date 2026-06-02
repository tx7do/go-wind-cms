<script setup lang="ts">
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { createHighlighter, type Highlighter } from 'shiki'
import katex from 'katex'

export type ContentType = 'html' | 'markdown' | 'text'

const props = withDefaults(defineProps<{
  content?: string
  type?: ContentType
  class?: string
}>(), {
  content: '',
  type: 'markdown',
})

// ========== Shiki 单例（双主题 light / dark）==========
const BUNDLED_LANGUAGES = [
  'javascript', 'typescript', 'jsx', 'tsx',
  'go', 'python', 'bash', 'shell',
  'json', 'yaml', 'html', 'css',
  'sql', 'markdown', 'diff', 'xml',
] as const

const BUNDLED_THEMES = ['github-light', 'github-dark'] as const

let highlighterPromise: Promise<Highlighter> | null = null
let cachedHighlighter: Highlighter | null = null

function getHighlighter(): Promise<Highlighter> {
  if (!highlighterPromise) {
    highlighterPromise = createHighlighter({
      themes: [...BUNDLED_THEMES],
      langs: [...BUNDLED_LANGUAGES],
    }).then((h) => {
      cachedHighlighter = h
      return h
    }).catch((err) => {
      console.warn('Failed to initialize Shiki:', err)
      highlighterPromise = null
      throw err
    })
  }
  return highlighterPromise
}

// 语言别名
const LANG_ALIAS: Record<string, string> = {
  js: 'javascript',
  ts: 'typescript',
  py: 'python',
  sh: 'bash',
  shell: 'bash',
  yml: 'yaml',
  md: 'markdown',
}

const shikiReady = ref(false)

onMounted(async () => {
  try {
    await getHighlighter()
    shikiReady.value = true
  } catch {
    // Shiki 加载失败，静默处理
  }
})

// ========== HTML 转义 ==========
function escapeHtml(text: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  }
  return text.replace(/[&<>"']/g, (m) => map[m])
}

// ========== Mermaid 初始化（仅客户端）==========
onMounted(async () => {
  if (import.meta.client) {
    try {
      const mermaid = (await import('mermaid')).default
      mermaid.initialize({
        startOnLoad: true,
        theme: 'default',
        securityLevel: 'loose',
      })
    } catch {
      // mermaid 加载失败
    }
  }
})

// ========== 自定义 marked renderer ==========
const renderer = new marked.Renderer()

renderer.codespan = (code) => {
  return `<code class="inline-code">${escapeHtml(code.text)}</code>`
}

renderer.code = (code) => {
  const rawLang = (code.lang || 'plaintext').toLowerCase()
  const lang = LANG_ALIAS[rawLang] || rawLang

  // Mermaid 图表
  if (lang === 'mermaid') {
    try {
      return `<div class="mermaid">${escapeHtml(code.text)}</div>`
    } catch (error) {
      console.warn('Failed to render mermaid diagram:', error)
      return `<pre class="code-block mermaid-error" data-lang="mermaid"><code>${escapeHtml(code.text)}</code></pre>`
    }
  }

  // Shiki 同步高亮
  if (cachedHighlighter) {
    try {
      const loaded = cachedHighlighter.getLoadedLanguages()
      const actualLang = loaded.includes(lang) ? lang : 'plaintext'

      if (actualLang === 'plaintext' || actualLang === 'text') {
        return `<pre class="code-block shiki-container" data-lang="${rawLang || 'text'}"><code>${escapeHtml(code.text)}</code></pre>`
      }

      const html = cachedHighlighter.codeToHtml(code.text, {
        lang: actualLang,
        themes: {
          light: 'github-light',
          dark: 'github-dark',
        },
      })

      return html.replace(
        /<pre class="shiki([^"]*)"/,
        `<pre class="shiki$1 code-block" data-lang="${rawLang || actualLang}"`,
      )
    } catch (error) {
      console.warn(`Shiki highlight failed for ${lang}:`, error)
    }
  }

  // Fallback：未加载 Shiki 时
  return `<pre class="code-block shiki-pending" data-lang="${rawLang || 'text'}"><code>${escapeHtml(code.text)}</code></pre>`
}

renderer.heading = (heading) => {
  const inlineHtml = marked.parseInline(heading.text)
  return `<h${heading.depth} class="heading-anchor">${inlineHtml}</h${heading.depth}>`
}

function splitUrlAndText(content: string): string {
  return content.replace(/(https?:\/\/[^\s，]+)(，[^ \n]+)/g, (_match, url, desc) => {
    return `[${url}](${url})${desc}`
  })
}

renderer.link = (link) => {
  const isExternal = link.href.startsWith('http') || link.href.startsWith('//')
  if (link.href === link.text) {
    return `<a href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.text)}</a>`
  } else {
    return `<a href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.href)}</a>${escapeHtml(link.text.replace(link.href, ''))}`
  }
}

renderer.image = (image) => {
  return `<figure class="markdown-image">
    <img src="${image.href}" alt="${image.text}" class="md-img" />
    ${image.text ? `<figcaption>${image.text}</figcaption>` : ''}
  </figure>`
}

renderer.table = (token) => {
  let html = '<thead>\n<tr>\n'
  for (const header of token.header) {
    const align = header.align ? ` style="text-align:${header.align}"` : ''
    const cellHtml = marked.parseInline(header.text)
    html += `<th${align}>${cellHtml}</th>\n`
  }
  html += '</tr>\n</thead>\n<tbody>\n'
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

renderer.paragraph = (token) => {
  return `<p>${marked.parseInline(token.text)}</p>\n`
}

marked.setOptions({ renderer })

// ========== 数学公式处理 ==========
function processMathFormulas(html: string): string {
  // Block math: $$...$$
  html = html.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
    try {
      const rendered = katex.renderToString(formula, { displayMode: true })
      return `<div class="math-block">${rendered}</div>`
    } catch (error) {
      console.warn('Failed to render block math:', error)
      return `<div class="math-error"><code>${escapeHtml(formula)}</code></div>`
    }
  })

  // Inline math: $...$
  html = html.replace(/(?<!\$)\$([^$\n]+)\$(?!\$)/g, (_, formula) => {
    try {
      const rendered = katex.renderToString(formula, { displayMode: false })
      return `<span class="math-inline">${rendered}</span>`
    } catch (error) {
      console.warn('Failed to render inline math:', error)
      return `<code class="math-error">${escapeHtml(formula)}</code>`
    }
  })

  return html
}

// ========== HTML 消毒 ==========
function sanitizeHtml(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: [
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'p', 'br', 'hr', 'pre', 'code',
      'strong', 'b', 'em', 'i', 'u', 'del', 's',
      'a', 'img', 'blockquote', 'ul', 'ol', 'li',
      'table', 'thead', 'tbody', 'tfoot', 'tr', 'th', 'td',
      'video', 'iframe', 'figure', 'figcaption',
      'mark', 'sub', 'sup', 'span', 'div',
      'math', 'mrow', 'mi', 'mn', 'mo', 'mfrac', 'msup', 'msub', 'mover', 'munder',
      'svg', 'g', 'text', 'path', 'circle', 'line', 'polyline', 'polygon', 'rect', 'ellipse',
    ],
    ALLOWED_ATTR: [
      'href', 'title', 'target', 'rel',
      'src', 'alt', 'width', 'height',
      'class', 'id', 'style', 'data-lang',
      'data-*',
      'tabindex',
      'viewBox', 'xmlns', 'x', 'y', 'cx', 'cy', 'r', 'x1', 'y1', 'x2', 'y2',
      'd', 'fill', 'stroke', 'stroke-width', 'points', 'transform',
    ],
    KEEP_CONTENT: true,
    ALLOW_UNKNOWN_PROTOCOLS: true,
  }) as unknown as string
}

// ========== 渲染计算 ==========
const renderedContent = computed(() => {
  // shikiReady 作为隐式依赖
  void shikiReady.value

  if (!props.content) return ''

  try {
    switch (props.type) {
      case 'markdown': {
        let md = props.content
        md = splitUrlAndText(md)
        let html = marked.parse(md) as string
        html = processMathFormulas(html)
        return sanitizeHtml(html)
      }
      case 'html':
        return sanitizeHtml(props.content)
      case 'text': {
        return sanitizeHtml(`<pre class="plain-text-block">${escapeHtml(props.content)}</pre>`)
      }
      default:
        return props.content
    }
  } catch (error) {
    console.error('Error rendering content:', error)
    return '<p class="content-error">Failed to render content</p>'
  }
})

// Mermaid 内容渲染后初始化
const containerRef = ref<HTMLElement | null>(null)

watch([renderedContent, shikiReady], () => {
  if (import.meta.client && containerRef.value) {
    nextTick(async () => {
      try {
        const mermaid = (await import('mermaid')).default
        mermaid.contentLoaded()
      } catch {
        // mermaid 渲染失败
      }
    })
  }
})
</script>

<template>
  <div
    ref="containerRef"
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
  -moz-osx-font-smoothing: grayscale;
}

/* Headings */
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

.content-viewer :deep(h5),
.content-viewer :deep(h6) {
  font-size: 1rem;
}

/* Paragraphs */
.content-viewer :deep(p) {
  margin: 22px 0;
  line-height: 2.0;
  letter-spacing: 0.25px;
  word-spacing: 0.05em;
  text-align: justify;
}

.content-viewer :deep(p:first-child) {
  margin-top: 0;
}

.content-viewer :deep(p:last-child) {
  margin-bottom: 0;
}

/* Links */
.content-viewer :deep(a.markdown-link) {
  color: hsl(var(--primary));
  text-decoration: none;
  border-bottom: 1.5px solid hsl(var(--primary) / 0.3);
  transition: all 0.25s ease;
  font-weight: 500;
}

.content-viewer :deep(a.markdown-link:hover) {
  text-decoration: none;
  border-bottom-color: hsl(var(--primary));
  background: hsl(var(--primary) / 0.08);
  padding: 2px 4px;
  margin: -2px -4px;
  border-radius: 4px;
}

/* Strong and emphasis */
.content-viewer :deep(strong),
.content-viewer :deep(b) {
  font-weight: 700;
  color: hsl(var(--primary));
  letter-spacing: 0.2px;
}

.content-viewer :deep(em),
.content-viewer :deep(i) {
  font-style: italic;
  color: hsl(var(--muted-foreground));
  letter-spacing: 0.15px;
}

/* Inline code */
.content-viewer :deep(code.inline-code) {
  background: rgba(150, 150, 150, 0.12);
  padding: 3px 8px;
  border-radius: 5px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 0.88em;
  color: hsl(var(--primary));
  border: 1px solid hsl(var(--primary) / 0.18);
  letter-spacing: 0.1px;
  white-space: nowrap;
}

/* ========== Code blocks (Shiki + fallback) ========== */
.content-viewer :deep(pre.code-block) {
  padding: 32px 24px 20px;
  border-radius: 12px;
  overflow-x: auto;
  margin: 32px 0;
  position: relative;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 13.5px;
  line-height: 1.7;
  letter-spacing: 0.15px;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.16);
  white-space: pre;
  word-wrap: normal;
  overflow-wrap: normal;
  tab-size: 2;
}

.content-viewer :deep(pre.code-block::before) {
  content: attr(data-lang);
  position: absolute;
  top: 10px;
  right: 14px;
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  pointer-events: none;
  z-index: 1;
}

.content-viewer :deep(pre.code-block code) {
  background: none;
  padding: 0;
  border-radius: 0;
  font-size: inherit;
  line-height: inherit;
  color: inherit;
  font-family: inherit;
  white-space: inherit;
  word-wrap: inherit;
  overflow-wrap: inherit;
}

/* Plain text fallback */
.content-viewer :deep(pre.code-block.shiki-pending),
.content-viewer :deep(pre.code-block.shiki-container) {
  background: #282c34;
  color: #abb2bf;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

/* Shiki highlighted */
.content-viewer :deep(pre.shiki.code-block) {
  border: 1px solid hsl(var(--border));
  padding: 32px 24px 20px;
  margin: 32px 0;
  background-color: #ffffff;
}

.content-viewer :deep(pre.shiki.code-block code) {
  display: block;
  background: transparent !important;
}

.content-viewer :deep(pre.shiki.code-block .line) {
  display: block;
}

.content-viewer :deep(pre.shiki.code-block::before) {
  color: rgba(0, 0, 0, 0.35);
}

/* Dark mode: Shiki dark theme */
.content-viewer :deep(.dark pre.shiki.code-block),
.content-viewer :deep(.dark pre.shiki.code-block span) {
  color: var(--shiki-dark, inherit) !important;
  background-color: var(--shiki-dark-bg, transparent) !important;
  font-style: var(--shiki-dark-font-style, normal) !important;
  font-weight: var(--shiki-dark-font-weight, normal) !important;
  text-decoration: var(--shiki-dark-text-decoration, none) !important;
}

:global(.dark) .content-viewer :deep(pre.shiki.code-block),
:global(.dark) .content-viewer :deep(pre.shiki.code-block) :deep(span) {
  color: var(--shiki-dark, inherit) !important;
  background-color: var(--shiki-dark-bg, transparent) !important;
  font-style: var(--shiki-dark-font-style, normal) !important;
  font-weight: var(--shiki-dark-font-weight, normal) !important;
  text-decoration: var(--shiki-dark-text-decoration, none) !important;
}

:global(.dark) .content-viewer :deep(pre.shiki.code-block) {
  background-color: #0d1117 !important;
  border-color: rgba(255, 255, 255, 0.08);
}

:global(.dark) .content-viewer :deep(pre.shiki.code-block::before) {
  color: rgba(255, 255, 255, 0.4);
}

:global(.dark) .content-viewer :deep(pre.code-block.shiki-pending),
:global(.dark) .content-viewer :deep(pre.code-block.shiki-container) {
  background: #0d1117;
  border-color: rgba(255, 255, 255, 0.08);
}

/* Plain text blocks */
.content-viewer :deep(pre.plain-text-block) {
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

/* Blockquotes */
.content-viewer :deep(blockquote) {
  border-left: 4px solid hsl(var(--primary));
  padding: 20px 24px;
  margin: 32px 0;
  background: hsl(var(--primary) / 0.06);
  color: hsl(var(--muted-foreground));
  border-radius: 0 8px 8px 0;
  font-style: italic;
  position: relative;
  letter-spacing: 0.2px;
  line-height: 1.85;
}

.content-viewer :deep(blockquote::before) {
  content: '\201C';
  position: absolute;
  top: -10px;
  left: 12px;
  font-size: 56px;
  color: hsl(var(--primary));
  opacity: 0.2;
  font-family: Georgia, serif;
  line-height: 1;
}

.content-viewer :deep(blockquote p) {
  margin: 14px 0;
  text-indent: 0;
  text-align: left;
}

.content-viewer :deep(blockquote p:first-child) {
  margin-top: 0;
}

.content-viewer :deep(blockquote p:last-child) {
  margin-bottom: 0;
}

/* Lists */
.content-viewer :deep(ul),
.content-viewer :deep(ol) {
  padding-left: 2em;
  margin: 28px 0;
}

.content-viewer :deep(li) {
  margin: 12px 0;
  line-height: 1.95;
  letter-spacing: 0.2px;
  position: relative;
}

.content-viewer :deep(li p) {
  margin: 10px 0;
}

.content-viewer :deep(ul) {
  list-style-type: disc;
}

.content-viewer :deep(ul li::marker) {
  color: hsl(var(--primary));
  font-weight: 600;
}

.content-viewer :deep(ul ul) {
  list-style-type: circle;
}

.content-viewer :deep(ul ul ul) {
  list-style-type: square;
}

.content-viewer :deep(ol) {
  list-style-type: decimal;
}

.content-viewer :deep(ol li::marker) {
  color: hsl(var(--primary));
  font-weight: 600;
}

/* Images */
.content-viewer :deep(figure.markdown-image) {
  margin: 1.5em 0;
  text-align: center;
}

.content-viewer :deep(img.md-img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  display: inline-block;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.content-viewer :deep(figcaption) {
  margin-top: 0.5em;
  font-size: 0.875rem;
  color: hsl(var(--muted-foreground));
  text-align: center;
}

/* Tables */
.content-viewer :deep(.table-wrapper) {
  overflow-x: auto;
  margin: 32px 0;
  border-radius: 8px;
  border: 1px solid hsl(var(--border));
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.content-viewer :deep(table.markdown-table) {
  width: 100%;
  border-collapse: collapse;
  background: hsl(var(--card));
  font-size: 15px;
  letter-spacing: 0.15px;
}

.content-viewer :deep(table.markdown-table th),
.content-viewer :deep(table.markdown-table td) {
  padding: 14px 16px;
  border: 1px solid hsl(var(--border));
  text-align: left;
}

.content-viewer :deep(table.markdown-table th) {
  background: hsl(var(--primary));
  color: white;
  font-weight: 600;
  line-height: 1.5;
  letter-spacing: 0.3px;
}

.content-viewer :deep(table.markdown-table td) {
  line-height: 1.7;
}

.content-viewer :deep(table.markdown-table tr:nth-child(even)) {
  background: rgba(150, 150, 150, 0.05);
}

.content-viewer :deep(table.markdown-table tr:hover) {
  background: hsl(var(--primary) / 0.06);
}

/* HR */
.content-viewer :deep(hr) {
  border: none;
  border-top: 2px solid hsl(var(--border));
  margin: 48px 0;
}

/* Video / Iframe */
.content-viewer :deep(video) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 1.5em 0;
  display: block;
}

.content-viewer :deep(iframe) {
  max-width: 100%;
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
  margin: 1.5em 0;
}

/* Error state */
.content-viewer :deep(.content-error) {
  color: #ff6b6b;
  padding: 1em;
  background: rgba(255, 107, 107, 0.1);
  border: 1px solid #ff6b6b;
  border-radius: 8px;
  margin: 1.5em 0;
}

/* Math formulas */
.content-viewer :deep(.math-block) {
  display: flex;
  justify-content: center;
  margin: 1.5em 0;
  padding: 1em;
  background: rgba(150, 150, 150, 0.05);
  border-radius: 8px;
  overflow-x: auto;
}

.content-viewer :deep(.math-inline) {
  display: inline;
  font-size: 0.95em;
}

.content-viewer :deep(.math-error) {
  color: #ff6b6b;
  background: rgba(255, 107, 107, 0.1);
  padding: 0.25em 0.5em;
  border-radius: 4px;
}

/* Mermaid diagrams */
.content-viewer :deep(.mermaid) {
  display: flex;
  justify-content: center;
  margin: 1.5em 0;
  padding: 1em;
  background: rgba(150, 150, 150, 0.05);
  border-radius: 8px;
  border: 1px solid hsl(var(--border));
  overflow-x: auto;
}

.content-viewer :deep(.mermaid svg) {
  max-width: 100%;
  height: auto;
}

.content-viewer :deep(.mermaid-error) {
  color: #ff6b6b;
  background: rgba(255, 107, 107, 0.1);
  border-color: #ff6b6b;
}

/* KaTeX */
.content-viewer :deep(.katex) {
  font: normal 1.21em KaTeX_Main, 'Times New Roman', serif;
  line-height: 1.2;
  text-indent: 0;
}

.content-viewer :deep(.katex-display) {
  display: block;
  margin: 1em 0;
  text-align: center;
}

.content-viewer :deep(.katex-html) {
  display: none;
}

.content-viewer :deep(.katex-mathml) {
  position: absolute;
  clip: rect(1px, 1px, 1px, 1px);
  padding: 0;
  border: 0;
  height: 1px;
  width: 1px;
  overflow: hidden;
}
</style>
