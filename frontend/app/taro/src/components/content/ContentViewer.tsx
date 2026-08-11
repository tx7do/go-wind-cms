import React, {useMemo} from 'react';
import {View} from '@tarojs/components';
import {marked} from 'marked';
import DOMPurify from 'dompurify';
import type {ContentViewerProps} from './types';

// 语言别名
const LANG_ALIAS: Record<string, string> = {
    js: 'javascript',
    ts: 'typescript',
    py: 'python',
    sh: 'bash',
    shell: 'bash',
    yml: 'yaml',
    md: 'markdown',
};

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

// 简化版 marked 渲染器
const renderer = new marked.Renderer();

renderer.codespan = (code) => {
    return `<code class="inline-code">${escapeHtml(code.text)}</code>`;
};

renderer.code = (code) => {
    const rawLang = (code.lang || 'plaintext').toLowerCase();
    const lang = LANG_ALIAS[rawLang] || rawLang;
    return `<pre class="code-block" data-lang="${rawLang || 'text'}"><code class="language-${lang}">${escapeHtml(code.text)}</code></pre>`;
};

renderer.heading = (heading) => {
    const inlineHtml = marked.parseInline(heading.text);
    return `<h${heading.depth} class="heading-anchor">${inlineHtml}</h${heading.depth}>`;
};

renderer.link = (link) => {
    const isExternal = link.href.startsWith('http') || link.href.startsWith('//');
    // 转义 href 防止属性逃逸；DOMPurify 仍会随后过滤 javascript: 等危险协议
    const safeHref = escapeHtml(link.href);
    return `<a href="${safeHref}" ${isExternal ? 'target="_blank" rel="noopener noreferrer"' : ''} class="markdown-link">${escapeHtml(link.text)}</a>`;
};

renderer.image = (image) => {
    const safeSrc = escapeHtml(image.href);
    const safeAlt = escapeHtml(image.text);
    return `<figure class="markdown-image"><img src="${safeSrc}" alt="${safeAlt}" class="md-img" />${image.text ? `<figcaption>${escapeHtml(image.text)}</figcaption>` : ''}</figure>`;
};

renderer.paragraph = (token) => {
    return `<p>${marked.parseInline(token.text)}</p>\n`;
};

marked.setOptions({renderer});

// ========== HTML 消毒 ==========
// 仅允许 http(s)、mailto、tel 以及相对路径/锚点的 URI，拒绝 javascript:/data: 等危险协议。
// 对齐 Vue 端 Viewer.vue 的安全基线，修复存储型 XSS（AUD9-C1）。
const SAFE_URI_REGEXP = /^(?:(?:https?|mailto|tel):|[/#.])/i;

function sanitizeHtml(html: string): string {
    return DOMPurify.sanitize(html, {
        ALLOWED_TAGS: [
            'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
            'p', 'br', 'hr', 'pre', 'code',
            'strong', 'b', 'em', 'i', 'u', 'del', 's',
            'a', 'img', 'blockquote', 'ul', 'ol', 'li',
            'table', 'thead', 'tbody', 'tfoot', 'tr', 'th', 'td',
            'video', 'figure', 'figcaption',
            'mark', 'sub', 'sup', 'span', 'div',
        ],
        ALLOWED_ATTR: [
            'href', 'title', 'target', 'rel',
            'src', 'alt', 'width', 'height',
            'class', 'id', 'data-lang',
            'data-*',
            'tabindex',
        ],
        KEEP_CONTENT: true,
        ALLOWED_URI_REGEXP: SAFE_URI_REGEXP,
    }) as unknown as string;
}

const ContentViewer: React.FC<ContentViewerProps> = ({
    content = '',
    type = 'markdown',
    className = '',
}) => {
    const renderedContent = useMemo(() => {
        if (!content) return '';

        try {
            let html = '';
            switch (type) {
                case 'markdown':
                    html = marked.parse(content) as string;
                    return sanitizeHtml(html);
                case 'html':
                    // html 类型同样必须经过 sanitize —— 原文不可信（来自 admin 富文本）
                    return sanitizeHtml(content);
                case 'text':
                    return `<pre class="plain-text-block">${escapeHtml(content)}</pre>`;
                default:
                    return sanitizeHtml(content);
            }
        } catch (error) {
            console.error('Error rendering content:', error);
            return `<p class="content-error">Failed to render content</p>`;
        }
    }, [content, type]);

    return (
        <View
          className={`content-viewer ${className}`}
        >
            <View dangerouslySetInnerHTML={{__html: renderedContent}} />
        </View>
    );
};

export default ContentViewer;
