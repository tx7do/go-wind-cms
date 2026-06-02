'use client';

import React, {useEffect, useRef, useState} from 'react';
import {marked} from 'marked';
import {createHighlighter, type Highlighter} from 'shiki';
import DOMPurify from 'dompurify';
import katex from 'katex';
import mermaid from 'mermaid';

import type {ContentViewerProps} from './types';
import styles from './ContentViewer.module.css';

// Initialize Mermaid
mermaid.initialize({
    startOnLoad: true,
    theme: 'default',
    securityLevel: 'loose',
});

/* ========== Shiki 单例（双主题 light / dark）========== */
const BUNDLED_LANGUAGES = [
    'javascript', 'typescript', 'jsx', 'tsx',
    'go', 'python', 'bash', 'shell',
    'json', 'yaml', 'html', 'css',
    'sql', 'markdown', 'diff', 'xml',
] as const;

const BUNDLED_THEMES = ['github-light', 'github-dark'] as const;

let highlighterPromise: Promise<Highlighter> | null = null;
let cachedHighlighter: Highlighter | null = null;

function getHighlighter(): Promise<Highlighter> {
    if (!highlighterPromise) {
        highlighterPromise = createHighlighter({
            themes: [...BUNDLED_THEMES],
            langs: [...BUNDLED_LANGUAGES],
        }).then((h) => {
            cachedHighlighter = h;
            return h;
        }).catch((err) => {
            console.warn('Failed to initialize Shiki:', err);
            highlighterPromise = null;
            throw err;
        });
    }
    return highlighterPromise;
}

// 语言别名映射（用户输入 → Shiki 接受的语言名）
const LANG_ALIAS: Record<string, string> = {
    js: 'javascript',
    ts: 'typescript',
    py: 'python',
    sh: 'bash',
    shell: 'bash',
    yml: 'yaml',
    md: 'markdown',
};

// Configure marked with custom renderer
const renderer = new marked.Renderer();

// Override code block rendering
renderer.codespan = (code) => {
    return `<code class="inline-code">${escapeHtml(code.text)}</code>`;
};

renderer.code = (code) => {
    const rawLang = (code.lang || 'plaintext').toLowerCase();
    const lang = LANG_ALIAS[rawLang] || rawLang;

    // Check if it's a mermaid diagram
    if (lang === 'mermaid') {
        try {
            return `<div class="mermaid">${escapeHtml(code.text)}</div>`;
        } catch (error) {
            console.warn('Failed to render mermaid diagram:', error);
            return `<pre class="code-block mermaid-error" data-lang="mermaid"><code>${escapeHtml(code.text)}</code></pre>`;
        }
    }

    // Shiki 同步高亮（如果 highlighter 已就绪）
    if (cachedHighlighter) {
        try {
            const loaded = cachedHighlighter.getLoadedLanguages();
            const actualLang = loaded.includes(lang) ? lang : 'plaintext';

            if (actualLang === 'plaintext' || actualLang === 'text') {
                // 纯文本：返回自定义结构以便保留容器样式
                return `<pre class="code-block shiki-container" data-lang="${rawLang || 'text'}"><code>${escapeHtml(code.text)}</code></pre>`;
            }

            // 使用双主题：light 为默认，dark 通过 CSS .dark 类切换
            // Shiki 会输出 --shiki-dark 等 CSS 变量到每个 <span>
            const html = cachedHighlighter.codeToHtml(code.text, {
                lang: actualLang,
                themes: {
                    light: 'github-light',
                    dark: 'github-dark',
                },
            });

            // 在 <pre> 上注入 data-lang 标签和 code-block 类名，以便 CSS 控制容器样式
            return html.replace(
                /<pre class="shiki([^"]*)"/,
                `<pre class="shiki$1 code-block" data-lang="${rawLang || actualLang}"`,
            );
        } catch (error) {
            console.warn(`Shiki highlight failed for ${lang}:`, error);
        }
    }

    // Fallback：未加载 Shiki 时显示纯文本（保留结构以便后续 hydrate）
    return `<pre class="code-block shiki-pending" data-lang="${rawLang || 'text'}"><code>${escapeHtml(code.text)}</code></pre>`;
};

// Override heading rendering
renderer.heading = (heading) => {
    const inlineHtml = marked.parseInline(heading.text);
    return `<h${heading.depth} class="heading-anchor">${inlineHtml}</h${heading.depth}>`;
};

function splitUrlAndText(content: string): string {
    return content.replace(/(https?:\/\/[^\s，]+)(，[^ \n]+)/g, (_match, url, desc) => {
        return `[${url}](${url})${desc}`;
    });
}

// Override link rendering
renderer.link = (link) => {
    const isExternal = link.href.startsWith('http') || link.href.startsWith('//');
    if (link.href === link.text) {
        return `<a href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.text)}</a>`;
    } else {
        return `<a href="${link.href}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.href)}</a>${escapeHtml(link.text.replace(link.href, ''))}`;
    }
};

// Override image rendering
renderer.image = (image) => {
    return `<figure class="markdown-image">
    <img src="${image.href}" alt="${image.text}" class="md-img" />
    ${image.text ? `<figcaption>${image.text}</figcaption>` : ''}
  </figure>`;
};

// Override table rendering
renderer.table = (token) => {
    let html = '<thead>\n<tr>\n';

    for (const header of token.header) {
        const align = header.align ? ` style="text-align:${header.align}"` : '';
        const cellHtml = marked.parseInline(header.text);
        html += `<th${align}>${cellHtml}</th>\n`;
    }

    html += '</tr>\n</thead>\n<tbody>\n';

    for (const row of token.rows) {
        html += '<tr>\n';
        for (const cell of row) {
            const align = cell.align ? ` style="text-align:${cell.align}"` : '';
            const cellHtml = marked.parseInline(cell.text);
            html += `<td${align}>${cellHtml}</td>\n`;
        }
        html += '</tr>\n';
    }

    html += '</tbody>';
    return `<div class="table-wrapper"><table class="markdown-table">${html}</table></div>`;
};

// Override paragraph rendering
renderer.paragraph = (token) => {
    return `<p>${marked.parseInline(token.text)}</p>\n`;
};

// Set custom renderer
marked.setOptions({renderer});

// Process math formulas
function processMathFormulas(html: string): string {
    // Block math: $$...$$
    html = html.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
        try {
            const rendered = katex.renderToString(formula, {displayMode: true});
            return `<div class="math-block">${rendered}</div>`;
        } catch (error) {
            console.warn('Failed to render block math:', error);
            return `<div class="math-error"><code>${escapeHtml(formula)}</code></div>`;
        }
    });

    // Inline math: $...$
    html = html.replace(/(?<!\$)\$([^$\n]+)\$(?!\$)/g, (_, formula) => {
        try {
            const rendered = katex.renderToString(formula, {displayMode: false});
            return `<span class="math-inline">${rendered}</span>`;
        } catch (error) {
            console.warn('Failed to render inline math:', error);
            return `<code class="math-error">${escapeHtml(formula)}</code>`;
        }
    });

    return html;
}

// Escape HTML for text type (preserve emoji)
function escapeHtml(text: string): string {
    const map: Record<string, string> = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;',
    };

    return text.replace(/[&<>"']/g, (m) => map[m]);
}

// Sanitize HTML
function sanitizeHtml(html: string): string {
    const config: any = {
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
        RETURN_DOM: false,
        RETURN_DOM_FRAGMENT: false,
        RETURN_DOM_IMPORT: false,
        FORCE_BODY: false,
        SANITIZE_DOM: true,
        IN_PLACE: false,
    };

    return DOMPurify.sanitize(html, config) as unknown as string;
}

const ContentViewer: React.FC<ContentViewerProps> = ({
                                                         content = '',
                                                         type = 'markdown',
                                                         className = ''
                                                     }) => {
    const containerRef = useRef<HTMLDivElement>(null);
    // Shiki 加载状态：首次加载后触发 re-render 以应用高亮
    const [shikiReady, setShikiReady] = useState<boolean>(() => cachedHighlighter !== null);

    useEffect(() => {
        if (!shikiReady) {
            let cancelled = false;
            getHighlighter().then(() => {
                if (!cancelled) setShikiReady(true);
            }).catch(() => {
                // Shiki 加载失败，静默处理（代码块保持纯文本回退）
            });
            return () => {
                cancelled = true;
            };
        }
    }, [shikiReady]);

    // Render content based on type
    // shikiReady 作为隐式依赖：当 Shiki 加载完成后，renderer.code 会使用同步高亮
    const getRenderedContent = () => {
        void shikiReady; // 触发依赖追踪
        if (!content) return '';

        try {
            let html = '';
            switch (type) {
                case 'markdown':
                    let md = content;
                    md = splitUrlAndText(md);
                    html = marked.parse(md) as string;
                    html = processMathFormulas(html);
                    return sanitizeHtml(html);
                case 'html':
                    return sanitizeHtml(content);
                case 'text':
                    return sanitizeHtml(`<pre class="plain-text-block">${escapeHtml(content)}</pre>`);
                default:
                    return content;
            }
        } catch (error) {
            console.error('Error rendering content:', error);
            return `<p class="content-error">Failed to render content</p>`;
        }
    };

    // Initialize mermaid diagrams after content is rendered
    useEffect(() => {
        if (containerRef.current) {
            try {
                mermaid.contentLoaded();
            } catch (error) {
                console.warn('Failed to initialize mermaid:', error);
            }
        }
    }, [shikiReady, content, type]);

    return (
        <div
            ref={containerRef}
            className={`${styles.contentViewer} ${className}`}
            dangerouslySetInnerHTML={{__html: getRenderedContent()}}
        />
    );
};

export default ContentViewer;
